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
	root := InferDecisionTree(inst, r)
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
	root := InferDecisionTree(insts[0], r)
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
