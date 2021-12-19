package utils

import (
	"math"
	"testing"
)

var caseSigmoid = []int{1, 2, 3, 5, 7, 10, 15, 20, 24, 30, 32}
var caseSigmoidMemory = []int{5446642176}

func TestSigmoid(t *testing.T) {
	for _, v := range caseSigmoid {
		s := Sigmoid(float64(v) / 10)
		t.Logf("%d, %v, %v", v, s, int64(math.Round(s*100)))
	}
	for _, v := range caseSigmoidMemory {
		memG := v / 1024 / 1024 / 1024
		t.Log(memG)
		countBase := float64(memG) / 10
		t.Log(countBase)
		t.Logf("%v", Sigmoid(countBase))
	}

}

var caseParseNodeMemory = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"10.233.92.93:9100","job":"node-exporter","kubernetes_namespace":"lens-metrics","kubernetes_node":"node3"},"value":[1639801850.034,"3909283840"]}]}}`

func TestParseNodeMemory(t *testing.T) {
	mem, _ := ParseNodeMemory(caseParseNodeMemory)
	t.Logf("result: %v bytes", mem)
}
