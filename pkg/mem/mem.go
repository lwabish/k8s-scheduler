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
	"math"
	"net/http"
)

const Name = "NodeAvailableMemory"

type NodeAvailableMemoryPlugin struct {
	handle framework.Handle
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
	return score, nil
}

func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	return &NodeAvailableMemoryPlugin{
		handle: f,
	}, nil
}

func (n *NodeAvailableMemoryPlugin) ScoreExtensions() framework.ScoreExtensions {
	return nil
}
