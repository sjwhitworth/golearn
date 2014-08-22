package meta

import (
	"github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBaggingRandomForestFit(t *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}

	t.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Fit(instf)
	}
}

func BenchmarkBaggingRandomForestPredict(t *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}

	rf.Fit(instf)
	t.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Predict(instf)
	}
}

func TestRandomForest1(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.6)

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	trainDataf := base.NewLazilyFilteredInstances(trainData, filt)
	testDataf := base.NewLazilyFilteredInstances(testData, filt)

	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}

	rf.Fit(trainDataf)
	predictions := rf.Predict(testDataf)
	confusionMat, err := eval.GetConfusionMatrix(testDataf, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = eval.GetSummary(confusionMat)
}
