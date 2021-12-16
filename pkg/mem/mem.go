package mem

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

const Name = "NodeAllocatableMemory"

type NodeMemoryPlugin struct {
	handle framework.Handle
}

func (n NodeMemoryPlugin) Name() string {
	return Name
}

func (n NodeMemoryPlugin) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := n.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.AsStatus(fmt.Errorf("getting node %q from Snapshot: %w", nodeName, err))
	}
	klog.Infoln("node-memory-plugin: ", nodeInfo, pod)
	klog.Infoln("node-memory-plugin: use allocatable mem as score: ", nodeInfo.Allocatable.Memory)
	return nodeInfo.Allocatable.Memory, nil
}

func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	return &NodeMemoryPlugin{
		handle: f,
	}, nil
}
