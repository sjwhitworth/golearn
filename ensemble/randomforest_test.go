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
	err = rf.Fit(trainData)
	if err != nil {
		t.Fatalf("Fitting failed: %s", err.Error())
	}
	predictions, err := rf.Predict(testData)
	if err != nil {
		t.Fatalf("Predicting failed: %s", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}

func TestRandomForestFitErrorWithNotEnoughFeatures(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	rf := NewRandomForest(10, len(base.NonClassAttributes(inst))+1)
	err = rf.Fit(inst)
	if err == nil {
		t.Fatalf("Fitting failed: %s", err.Error())
	}
}
