package mem

import (
	"context"
	"fmt"
	"github.com/lwabish/k8s-scheduler/pkg/utils"
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
	MinMemory          int    `json:"min_memory,omitempty"`
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
	r, _ := http.Get(fmt.Sprintf("http://localhost:62222/api/v1/query?query=%s", queryString))
	defer r.Body.Close()
	jsonString, _ := ioutil.ReadAll(r.Body)
	nodeMemory, _ := utils.ParseNodeMemory(string(jsonString))
	klog.Infof("node %s mem available is: %v bytes", nodeName, nodeMemory)
	countBase := float64(nodeMemory/1024/1024/1024) / 10
	score := int64(math.Round(utils.Sigmoid(countBase) * 100))
	klog.Infof("node %s counting detail is: %v %v", nodeName, countBase, score)
	klog.Infoln(n.args.PrometheusEndpoint, n.args.MaxMemory, n.args.MinMemory)
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
