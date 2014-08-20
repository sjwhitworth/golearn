package trees

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
	"math"
	"testing"
)

func TestRandomTree(t *testing.T) {
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

	r := new(RandomTreeRuleGenerator)
	r.Attributes = 2

	_ = InferID3Tree(instf, r)
}

func TestRandomTreeClassification(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
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
	predictions, err := root.Predict(testDataF)
	if err != nil {
		t.Fatalf("Predicting failed: %s", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(testDataF, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}

func TestRandomTreeClassification2(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
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
	err = root.Fit(trainDataF)
	if err != nil {
		t.Fatalf("Fitting failed: %s", err.Error())
	}

	predictions, err := root.Predict(testDataF)
	if err != nil {
		t.Fatalf("Predicting failed: %s", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(testDataF, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}

func TestPruning(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
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

	err = root.Fit(fittrainData)
	if err != nil {
		t.Fatalf("Fitting failed: %s", err.Error())
	}

	root.Prune(fittestData)
	predictions, err := root.Predict(testDataF)
	if err != nil {
		t.Fatalf("Predicting failed: %s", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(testDataF, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}

func TestInformationGain(t *testing.T) {
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
		t.Error(entropy)
	}
}

func TestID3Inference(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/tennis.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	// Build the decision tree
	rule := new(InformationGainRuleGenerator)
	root := InferID3Tree(inst, rule)

	// Verify the tree
	// First attribute should be "outlook"
	if root.SplitAttr.GetName() != "outlook" {
		t.Error(root)
	}
	sunnyChild := root.Children["sunny"]
	overcastChild := root.Children["overcast"]
	rainyChild := root.Children["rainy"]
	if sunnyChild.SplitAttr.GetName() != "humidity" {
		t.Error(sunnyChild)
	}
	if rainyChild.SplitAttr.GetName() != "windy" {
		t.Error(rainyChild)
	}
	if overcastChild.SplitAttr != nil {
		t.Error(overcastChild)
	}

	sunnyLeafHigh := sunnyChild.Children["high"]
	sunnyLeafNormal := sunnyChild.Children["normal"]
	if sunnyLeafHigh.Class != "no" {
		t.Error(sunnyLeafHigh)
	}
	if sunnyLeafNormal.Class != "yes" {
		t.Error(sunnyLeafNormal)
	}
	windyLeafFalse := rainyChild.Children["false"]
	windyLeafTrue := rainyChild.Children["true"]
	if windyLeafFalse.Class != "yes" {
		t.Error(windyLeafFalse)
	}
	if windyLeafTrue.Class != "no" {
		t.Error(windyLeafTrue)
	}

	if overcastChild.Class != "yes" {
		t.Error(overcastChild)
	}
}

func TestID3Classification(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	filt := filters.NewBinningFilter(inst, 10)
	for _, a := range base.NonClassFloatAttributes(inst) {
		filt.AddAttribute(a)
	}
	filt.Train()
	instf := base.NewLazilyFilteredInstances(inst, filt)

	trainData, testData := base.InstancesTrainTestSplit(instf, 0.70)

	// Build the decision tree
	rule := new(InformationGainRuleGenerator)
	root := InferID3Tree(trainData, rule)

	predictions, err := root.Predict(testData)
	if err != nil {
		t.Fatalf("Predicting failed: %s", err.Error())
	}

	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		t.Fatalf("Unable to get confusion matrix: %s", err.Error())
	}
	_ = evaluation.GetSummary(confusionMat)
}

func TestID3(t *testing.T) {
	inst, err := base.ParseCSVToInstances("../examples/datasets/tennis.csv", true)
	if err != nil {
		t.Fatal("Unable to parse CSV to instances: %s", err.Error())
	}

	// Build the decision tree
	tree := NewID3DecisionTree(0.0)
	tree.Fit(inst)
	root := tree.Root

	// Verify the tree
	// First attribute should be "outlook"
	if root.SplitAttr.GetName() != "outlook" {
		t.Error(root)
	}
	sunnyChild := root.Children["sunny"]
	overcastChild := root.Children["overcast"]
	rainyChild := root.Children["rainy"]
	if sunnyChild.SplitAttr.GetName() != "humidity" {
		t.Error(sunnyChild)
	}
	if rainyChild.SplitAttr.GetName() != "windy" {
		t.Error(rainyChild)
	}
	if overcastChild.SplitAttr != nil {
		t.Error(overcastChild)
	}

	sunnyLeafHigh := sunnyChild.Children["high"]
	sunnyLeafNormal := sunnyChild.Children["normal"]
	if sunnyLeafHigh.Class != "no" {
		t.Error(sunnyLeafHigh)
	}
	if sunnyLeafNormal.Class != "yes" {
		t.Error(sunnyLeafNormal)
	}

	windyLeafFalse := rainyChild.Children["false"]
	windyLeafTrue := rainyChild.Children["true"]
	if windyLeafFalse.Class != "yes" {
		t.Error(windyLeafFalse)
	}
	if windyLeafTrue.Class != "no" {
		t.Error(windyLeafTrue)
	}

	if overcastChild.Class != "yes" {
		t.Error(overcastChild)
	}
}
