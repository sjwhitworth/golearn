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
	insts := base.InstancesTrainTestSplit(inst, 0.4)
	filt := filters.NewBinningFilter(insts[0], 10)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(insts[1])
	filt.Run(insts[0])
	rf := NewRandomForest(10, 2)
	rf.Fit(insts[0])
	predictions := rf.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}
