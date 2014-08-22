package ensemble

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"testing"
)

func TestRandomForest1(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
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
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}
