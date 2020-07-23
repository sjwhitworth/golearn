// Example of how to use CART trees for both Classification and Regression

package main

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
)

func main() {

	// Load Titanic Data For classification
	classificationData, err := base.ParseCSVToInstances("../datasets/titanic.csv", false)
	if err != nil {
		panic(err)
	}
	trainData, testData := base.InstancesTrainTestSplit(classificationData, 0.5)

	// Create New Classification Tree
	// Hyperparameters - loss function, max Depth (-1 will split until pure), list of unique labels
	decTree = NewDecisionTreeClassifier("entropy", -1, []int64{0, 1})

	// Train Tree
	decTree.Fit(trainData)
	// Print out tree for visualization - shows splits and feature and predictions
	fmt.Println(decTree.String())

	// Access Predictions
	classificationPreds := decTree.Predict(testData)

	fmt.Println("Titanic Predictions")
	fmt.Println(classificationPreds)

	// Evaluate Accuracy on Test Data
	fmt.Println(decTree.Evaluate(testData))

	// Load House Price Data For Regression

	regressionData, err := base.ParseCSVToInstances("../datasets/boston_house_prices.csv", false)
	if err != nil {
		panic(err)
	}
	trainRegData, testRegData := base.InstancesTrainTestSplit(regressionData, 0.5)

	// Hyperparameters - Loss function, max Depth (-1 will split until pure)
	regTree := NewDecisionTreeRegressor("mse", -1)

	// Train Tree
	regTree.Fit(trainRegData)

	// Print out tree for visualization
	fmt.Println(regTree.String())

	// Access Predictions
	regressionPreds := regTree.Predict(testRegData)

	fmt.Println("Boston House Price Predictions")
	fmt.Println(regressionPreds)

}
