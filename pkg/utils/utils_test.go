package utils

import (
	"math"
	"testing"
)

// 对归一化后的结果放大一定的倍数，使结果更适合作为sigmoid函数的参数
const multiplicationFactor = 4

var caseAvailableMemBytes []int64
var normalized []float64

func TestNormalizationMem(t *testing.T) {
	for i := 0.0; i < 32.0; i += 0.5 {
		caseAvailableMemBytes = append(caseAvailableMemBytes, int64(i*1024*1024*1024))
	}
	for _, v := range caseAvailableMemBytes {
		result := NormalizationMem(32*1024*1024*1024, v)
		normalized = append(normalized, result)
		t.Logf("mem: %v, normalized: %v", v, result)
	}
}

func TestSigmoid(t *testing.T) {
	TestNormalizationMem(t)
	for i, v := range normalized {
		s := Sigmoid(v * multiplicationFactor)
		score := int64(math.Round(s * 100))
		t.Logf("memBytes: %v, normalized: %v, sigmoid: %v, score: %v", caseAvailableMemBytes[i], v, s, score)
	}
}

var casePrometheusNodeMemoryResult = `{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"node_memory_MemAvailable_bytes","instance":"10.233.92.93:9100","job":"node-exporter","kubernetes_namespace":"lens-metrics","kubernetes_node":"node3"},"value":[1639801850.034,"3909283840"]}]}}`

func TestParseNodeMemory(t *testing.T) {
	mem, _ := ParseNodeMemory(casePrometheusNodeMemoryResult)
	t.Logf("result: %v bytes", mem)
}
