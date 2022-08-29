package mem

import (
	"context"
	"fmt"
	"github.com/lwabish/k8s-scheduler/pkg/utils"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	runtime2 "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	metricsClientSet "k8s.io/metrics/pkg/client/clientset/versioned"
	"math"
	"net/http"
)

const Name = "NodeAvailableMemory"

type NodeAvailableMemoryPluginArg struct {
	PrometheusEndpoint string `json:"prometheus_endpoint,omitempty"`
	MaxMemory          int    `json:"max_memory,omitempty"`
	MetricsClientSet   *metricsClientSet.Clientset
}

type NodeAvailableMemoryPlugin struct {
	handle framework.Handle
	args   NodeAvailableMemoryPluginArg
}

func (n *NodeAvailableMemoryPlugin) Name() string {
	return Name
}

func (n *NodeAvailableMemoryPlugin) Score(_ context.Context, _ *framework.CycleState, _ *v1.Pod, nodeName string) (int64, *framework.Status) {
	var nodeMemory int64
	if n.args.PrometheusEndpoint == "" {
		used, err := n.getNodeUsedMemory(nodeName)
		if err != nil {
			return 0, framework.AsStatus(err)
		}
		nodeInfo, err := n.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
		if err != nil {
			return 0, framework.AsStatus(err)
		}
		nodeMemory = nodeInfo.Allocatable.Memory - used
	} else {
		nm, err := n.getNodeAvailableMemory(nodeName)
		nodeMemory = nm
		if err != nil {
			return 0, framework.AsStatus(err)
		}
	}
	normalized := utils.NormalizationMem(int64(n.args.MaxMemory*1024*1024*1024), nodeMemory)
	sigmoid := utils.Sigmoid(normalized)
	score := int64(math.Round(sigmoid * 100))
	klog.Infof("node %s counting detail:available %v normalized %v, sigmoid %v, score %v", nodeName, nodeMemory, normalized, sigmoid, score)
	return score, nil
}

func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	args := &NodeAvailableMemoryPluginArg{}
	err := runtime2.DecodeInto(configuration, args)

	config, err := utils.GetClientConfig()
	if err != nil {
		return nil, err
	}
	mcs, err := metricsClientSet.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	args.MetricsClientSet = mcs

	return &NodeAvailableMemoryPlugin{
		handle: f,
		args:   *args,
	}, nil
}

func (n *NodeAvailableMemoryPlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}

// getNodeUsedMemory 输入节点名称，输出节点已用内存量
func (n *NodeAvailableMemoryPlugin) getNodeUsedMemory(nodeName string) (int64, error) {
	nodeMetrics, err := n.args.MetricsClientSet.MetricsV1beta1().NodeMetricses().Get(context.TODO(), nodeName, metaV1.GetOptions{})
	if err != nil {
		return 0, err
	}
	return nodeMetrics.Usage.Memory().Value(), nil
}

// GetNodeAvailableMemory 输入节点名称，获取该节点的可用内存量（从prometheus获取）
func (n *NodeAvailableMemoryPlugin) getNodeAvailableMemory(nodeName string) (int64, error) {
	queryString := fmt.Sprintf("node_memory_MemAvailable_bytes{kubernetes_node=\"%s\"}", nodeName)
	r, err := http.Get(fmt.Sprintf("http://%s/api/v1/query?query=%s", n.args.PrometheusEndpoint, queryString))
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)
	jsonString, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	nodeMemory, err := utils.ParseNodeMemory(string(jsonString))
	if err != nil {
		return 0, err
	}
	return nodeMemory, nil
}
