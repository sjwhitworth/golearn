package ensemble

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	filters "github.com/sjwhitworth/golearn/filters"
	"testing"
)

func TestRandomForest1(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	insts := base.InstancesTrainTestSplit(inst, 0.60)
	filt := filters.NewChiMergeFilter(insts[0], 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(insts[1])
	filt.Run(insts[0])
	rf := NewRandomForest(10, 3)
	rf.Fit(insts[0])
	predictions := rf.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetSummary(confusionMat))
}
