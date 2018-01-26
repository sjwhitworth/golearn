// Demonstrates decision tree classification

package main

import (
	"fmt"
	"math/rand"

	"github.com/amclay/golearn/base"
	"github.com/amclay/golearn/ensemble"
	"github.com/amclay/golearn/evaluation"
	"github.com/amclay/golearn/filters"
	"github.com/amclay/golearn/trees"
)

func main() {

	var tree base.Classifier

	rand.Seed(44111342)

	// Load in the iris dataset
	iris, err := base.ParseCSVToInstances("../datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	// Discretise the iris dataset with Chi-Merge
	filt := filters.NewChiMergeFilter(iris, 0.999)
	for _, a := range base.NonClassFloatAttributes(iris) {
		filt.AddAttribute(a)
	}
	filt.Train()
	irisf := base.NewLazilyFilteredInstances(iris, filt)

	// Create a 60-40 training-test split
	trainData, testData := base.InstancesTrainTestSplit(irisf, 0.60)

	//
	// First up, use ID3
	//
	tree = trees.NewID3DecisionTree(0.6)
	// (Parameter controls train-prune split.)

	// Train the ID3 tree
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// Generate predictions
	predictions, err := tree.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Evaluate
	fmt.Println("ID3 Performance (information gain)")
	cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))

	tree = trees.NewID3DecisionTreeFromRule(0.6, new(trees.InformationGainRatioRuleGenerator))
	// (Parameter controls train-prune split.)

	// Train the ID3 tree
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// Generate predictions
	predictions, err = tree.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Evaluate
	fmt.Println("ID3 Performance (information gain ratio)")
	cf, err = evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))

	tree = trees.NewID3DecisionTreeFromRule(0.6, new(trees.GiniCoefficientRuleGenerator))
	// (Parameter controls train-prune split.)

	// Train the ID3 tree
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}

	// Generate predictions
	predictions, err = tree.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Evaluate
	fmt.Println("ID3 Performance (gini index generator)")
	cf, err = evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))
	//
	// Next up, Random Trees
	//

	// Consider two randomly-chosen attributes
	tree = trees.NewRandomTree(2)
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}
	predictions, err = tree.Predict(testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("RandomTree Performance")
	cf, err = evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))

	//
	// Finally, Random Forests
	//
	tree = ensemble.NewRandomForest(70, 3)
	err = tree.Fit(trainData)
	if err != nil {
		panic(err)
	}
	predictions, err = tree.Predict(testData)
	if err != nil {
		panic(err)
	}
	fmt.Println("RandomForest Performance")
	cf, err = evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))
}
