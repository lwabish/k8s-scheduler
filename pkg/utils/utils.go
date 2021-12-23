package utils

import (
	"github.com/thedevsaddam/gojsonq/v2"
	"math"
	"strconv"
)

func Sigmoid(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func NormalizationMem(maxBytes int64, valueBytes int64) float64 {
	return float64(valueBytes) / float64(maxBytes)
}

const nodeMemoryJsonPath = "data.result.[0].value.[1]"

func ParseNodeMemory(responseString string) (int64, error) {
	r := gojsonq.New().FromString(responseString).Find(nodeMemoryJsonPath)
	return strconv.ParseInt(r.(string), 10, 64)
}
