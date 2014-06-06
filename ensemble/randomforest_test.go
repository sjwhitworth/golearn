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
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.60)
	filt := filters.NewChiMergeFilter(trainData, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(testData)
	filt.Run(trainData)
	rf := NewRandomForest(10, 3)
	rf.Fit(trainData)
	predictions := rf.Predict(testData)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetSummary(confusionMat))
}
