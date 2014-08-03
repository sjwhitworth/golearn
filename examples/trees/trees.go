// Demonstrates decision tree classification

package main

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	ensemble "github.com/sjwhitworth/golearn/ensemble"
	eval "github.com/sjwhitworth/golearn/evaluation"
	filters "github.com/sjwhitworth/golearn/filters"
	trees "github.com/sjwhitworth/golearn/trees"
	"math/rand"
	"time"
)

func main() {

	var tree base.Classifier

	rand.Seed(time.Now().UTC().UnixNano())

	// Load in the iris dataset
	iris, err := base.ParseCSVToInstances("../datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	// Discretise the iris dataset with Chi-Merge
	filt := filters.NewChiMergeFilter(iris, 0.99)
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
	tree.Fit(trainData)

	// Generate predictions
	predictions := tree.Predict(testData)

	// Evaluate
	fmt.Println("ID3 Performance")
	cf := eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(eval.GetSummary(cf))

	//
	// Next up, Random Trees
	//

	// Consider two randomly-chosen attributes
	tree = trees.NewRandomTree(2)
	tree.Fit(testData)
	predictions = tree.Predict(testData)
	fmt.Println("RandomTree Performance")
	cf = eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(eval.GetSummary(cf))

	//
	// Finally, Random Forests
	//
	tree = ensemble.NewRandomForest(100, 3)
	tree.Fit(trainData)
	predictions = tree.Predict(testData)
	fmt.Println("RandomForest Performance")
	cf = eval.GetConfusionMatrix(testData, predictions)
	fmt.Println(eval.GetSummary(cf))
}
