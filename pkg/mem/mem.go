package mem

import (
	"context"
	"fmt"
	"github.com/lwabish/k8s-scheduler/pkg/utils"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	runtime2 "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
	"math"
	"net/http"
)

const Name = "NodeAvailableMemory"

type NodeAvailableMemoryPluginArg struct {
	PrometheusEndpoint string `json:"prometheus_endpoint,omitempty"`
	MaxMemory          int    `json:"max_memory,omitempty"`
}

type NodeAvailableMemoryPlugin struct {
	handle framework.Handle
	args   NodeAvailableMemoryPluginArg
}

func (n NodeAvailableMemoryPlugin) Name() string {
	return Name
}

func (n NodeAvailableMemoryPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	queryString := fmt.Sprintf("node_memory_MemAvailable_bytes{kubernetes_node=\"%s\"}", nodeName)
	r, _ := http.Get(fmt.Sprintf("http://%s/api/v1/query?query=%s", n.args.PrometheusEndpoint, queryString))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(r.Body)
	jsonString, _ := ioutil.ReadAll(r.Body)
	nodeMemory, _ := utils.ParseNodeMemory(string(jsonString))
	normalized := utils.NormalizationMem(int64(n.args.MaxMemory*1024*1024*1024), nodeMemory)
	sigmoid := utils.Sigmoid(normalized)
	score := int64(math.Round(sigmoid * 100))
	klog.Infof("node %s counting detail:available %v normalized %v, sigmoid %v, score %v", nodeMemory, nodeName, normalized, sigmoid, score)
	return score, nil
}

func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	args := &NodeAvailableMemoryPluginArg{}
	err := runtime2.DecodeInto(configuration, args)
	if err != nil {
		return nil, err
	}
	return &NodeAvailableMemoryPlugin{
		handle: f,
		args:   *args,
	}, nil
}

func (n *NodeAvailableMemoryPlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
