package ensemble

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"testing"
)

func TestRandomForest1(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	trainData, testData := base.InstancesTrainTestSplit(instf, 0.60)

	rf := NewRandomForest(10, 3)
	rf.Fit(trainData)
	predictions := rf.Predict(testData)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetSummary(confusionMat))
}
