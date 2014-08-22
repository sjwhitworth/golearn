package trees

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	eval "github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"math"
	"testing"
)

func TestRandomTree(testEnv *testing.T) {
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

	r := new(RandomTreeRuleGenerator)
	r.Attributes = 2
	fmt.Println(instf)
	root := InferID3Tree(instf, r)
	fmt.Println(root)
}

func TestRandomTreeClassification(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	trainDataF := base.NewLazilyFilteredInstances(trainData, filt)
	testDataF := base.NewLazilyFilteredInstances(testData, filt)

	r := new(RandomTreeRuleGenerator)
	r.Attributes = 2
	root := InferID3Tree(trainDataF, r)
	fmt.Println(root)
	predictions := root.Predict(testDataF)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testDataF, predictions)
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
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.4)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	trainDataF := base.NewLazilyFilteredInstances(trainData, filt)
	testDataF := base.NewLazilyFilteredInstances(testData, filt)

	root := NewRandomTree(2)
	root.Fit(trainDataF)
	fmt.Println(root)
	predictions := root.Predict(testDataF)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testDataF, predictions)
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
	trainData, testData := base.InstancesTrainTestSplit(inst, 0.6)
	filt := filters.NewChiMergeFilter(inst, 0.90)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	trainDataF := base.NewLazilyFilteredInstances(trainData, filt)
	testDataF := base.NewLazilyFilteredInstances(testData, filt)

	root := NewRandomTree(2)
	fittrainData, fittestData := base.InstancesTrainTestSplit(trainDataF, 0.6)
	root.Fit(fittrainData)
	root.Prune(fittestData)
	fmt.Println(root)
	predictions := root.Predict(testDataF)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testDataF, predictions)
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
	inst, err := base.ParseCSVToInstances("../examples/datasets/tennis.csv", true)
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
		fmt.Println(rainyChild.SplitAttr)
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

func TestID3Classification(testEnv *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(inst)
	filt := filters.NewBinningFilter(inst, 10)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	fmt.Println(filt)
	instf := base.NewLazilyFilteredInstances(inst, filt)
	fmt.Println("INSTFA", instf.AllAttributes())
	fmt.Println("INSTF", instf)
	trainData, testData := base.InstancesTrainTestSplit(instf, 0.70)

	// Build the decision tree
	rule := new(InformationGainRuleGenerator)
	root := InferID3Tree(trainData, rule)
	fmt.Println(root)
	predictions := root.Predict(testData)
	fmt.Println(predictions)
	confusionMat := eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(confusionMat)
	fmt.Println(eval.GetMacroPrecision(confusionMat))
	fmt.Println(eval.GetMacroRecall(confusionMat))
	fmt.Println(eval.GetSummary(confusionMat))
}

func TestID3(testEnv *testing.T) {

	// Import the "PlayTennis" dataset
	inst, err := base.ParseCSVToInstances("../examples/datasets/tennis.csv", true)
	fmt.Println(inst)
	if err != nil {
		panic(err)
	}

	// Build the decision tree
	tree := NewID3DecisionTree(0.0)
	tree.Fit(inst)
	root := tree.Root

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
