package trees

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	filters "github.com/sjwhitworth/golearn/filters"
	"math"
	"testing"
)

func TestRandomTree(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(inst)
	fmt.Println(inst)
	r := new(RandomTreeRuleGenerator)
	r.Attributes = 2
	root := InferID3Tree(inst, r)
	fmt.Println(root)
}

func TestRandomTreeClassification(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	insts := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	filt.Run(insts[0])
	filt.Run(insts[1])
	fmt.Println(inst)
	r := new(RandomTreeRuleGenerator)
	r.Attributes = 2
	root := InferID3Tree(insts[0], r)
	fmt.Println(root)
	predictions := root.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}

func TestRandomTreeClassification2(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	insts := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	fmt.Println(insts[1])
	filt.Run(insts[1])
	filt.Run(insts[0])
	root := NewRandomTree(2)
	root.Fit(insts[0])
	fmt.Println(root)
	predictions := root.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}

func TestPruning(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	insts := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	filt.AddAllNumericAttributes()
	filt.Build()
	fmt.Println(insts[1])
	filt.Run(insts[1])
	filt.Run(insts[0])
	root := NewRandomTree(2)
	fitInsts := base.InstancesTrainTestSplit(insts[0], 0.6)
	root.Fit(fitInsts[0])
	root.Prune(fitInsts[1])
	fmt.Println(root)
	predictions := root.Predict(insts[1])
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(insts[1], predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}

func TestInformationGain(testEnv *testing.T) {
	outlook := make(map[string]map[string]int)
	outlook["sunny"] = make(map[string]int)
	outlook["overcast"] = make(map[string]int)
	outlook["rain"] = make(map[string]int)
	outlook["sunny"]["play"] = 2
	outlook["sunny"]["noplay"] = 3
	outlook["overcast"]["play"] = 4
	outlook["rain"]["play"] = 3
	outlook["rain"]["noplay"] = 2

	entropy := getSplitEntropy(outlook)
	if math.Abs(entropy-0.694) > 0.001 {
		testEnv.Error(entropy)
	}
}

func TestID3Inference(testEnv *testing.T) {

	// Import the "PlayTennis" dataset
	inst, err := base.ParseCSVToInstances("./tennis.csv", true)
	if err != nil {
		panic(err)
	}

	// Build the decision tree
	rule := new(InformationGainRuleGenerator)
	root := InferID3Tree(inst, rule)

	// Verify the tree
	// First attribute should be "outlook"
	if root.SplitAttr.GetName() != "outlook" {
		testEnv.Error(root)
	}
	sunnyChild := root.Children["sunny"]
	overcastChild := root.Children["overcast"]
	rainyChild := root.Children["rainy"]
	if sunnyChild.SplitAttr.GetName() != "humidity" {
		testEnv.Error(sunnyChild)
	}
	if rainyChild.SplitAttr.GetName() != "windy" {
		testEnv.Error(rainyChild)
	}
	if overcastChild.SplitAttr != nil {
		testEnv.Error(overcastChild)
	}

	sunnyLeafHigh := sunnyChild.Children["high"]
	sunnyLeafNormal := sunnyChild.Children["normal"]
	if sunnyLeafHigh.Class != "no" {
		testEnv.Error(sunnyLeafHigh)
	}
	if sunnyLeafNormal.Class != "yes" {
		testEnv.Error(sunnyLeafNormal)
	}

	windyLeafFalse := rainyChild.Children["false"]
	windyLeafTrue := rainyChild.Children["true"]
	if windyLeafFalse.Class != "yes" {
		testEnv.Error(windyLeafFalse)
	}
	if windyLeafTrue.Class != "no" {
		testEnv.Error(windyLeafTrue)
	}

	if overcastChild.Class != "yes" {
		testEnv.Error(overcastChild)
	}
}
