package util

import "testing"

var memG = []int{1, 2, 3, 5, 7, 10, 15, 20, 24, 30, 32}

func TestSigmoid(t *testing.T) {
	for _, v := range memG {
		t.Logf("%d, %v", v, sigmoid(float64(v)))
	}

}
