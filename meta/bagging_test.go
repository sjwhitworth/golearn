package meta

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"github.com/sjwhitworth/golearn/trees"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkBaggingRandomForestFit(testEnv *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
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
	testEnv.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Fit(instf)
	}
}

func BenchmarkBaggingRandomForestPredict(testEnv *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
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
	testEnv.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Predict(instf)
	}
}

func TestRandomForest1(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.6)
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
	fmt.Println(rf)
	predictions := rf.Predict(testDataf)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testDataf, predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}
