package meta

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	filters "github.com/sjwhitworth/golearn/filters"
	trees "github.com/sjwhitworth/golearn/trees"
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
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(inst)
	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}
	testEnv.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Fit(inst)
	}
}

func BenchmarkBaggingRandomForestPredict(testEnv *testing.B) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(inst)
	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}
	rf.Fit(inst)
	testEnv.ResetTimer()
	for i := 0; i < 20; i++ {
		rf.Predict(inst)
	}
}

func TestRandomForest1(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	insts := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(insts[1])
	filt.Run(insts[0])
	rf := new(BaggedModel)
	for i := 0; i < 10; i++ {
		rf.AddModel(trees.NewRandomTree(2))
	}
	rf.Fit(insts[0])
	fmt.Println(rf)
	predictions := rf.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}
