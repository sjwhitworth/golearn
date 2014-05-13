package filters

import (
	base "github.com/sjwhitworth/golearn/base"
	"math"
	"testing"
)

func TestBinning(testEnv *testing.T) {
	inst1, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	inst2, err := base.ParseCSVToInstances("../examples/datasets/iris_binned.csv", true)
	inst3, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	filt := NewBinningFilter(inst1, 10)
	filt.AddAttribute(inst1.GetAttr(0))
	filt.Build()
	filt.Run(inst1)
	for i := 0; i < inst1.Rows; i++ {
		val1 := inst1.Get(i, 0)
		val2 := inst2.Get(i, 0)
		val3 := inst3.Get(i, 0)
		if math.Abs(val1-val2) >= 1 {
			testEnv.Error(val1, val2, val3, i)
		}
	}
}
