// Example of how to use CART trees for both Classification and Regression

package main

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/trees"
)

func main() {
	/* Performance of CART Algorithm:

		Training Time for Titanic Dataset ≈ 611 µs
		Prediction Time for Titanic Datset ≈ 101 µs

		Complexity Analysis:
			1x Dataset --   x ms
			2x Dataset --   1.7x ms
			128x Dataset -- 74x ms

			Complexity is sub linear

		Sklearn:
		Training Time for Titanic Dataset ≈ 8.8 µs
		Prediction Time for Titanic Datset ≈ 7.87 µs


		This implementation and sci-kit learn produce the exact same tree for the exact same dataset.
		Predictions on the same test set also yield the exact same accuracy.

		This implementation is optimized to prevent redundant iterations over the dataset, but it is not completely optimized. Also, sklearn makes use of numpy to access column easily, whereas here a complete iteration is required.
	 	In terms of Hyperparameters, this implmentation gives you the ability to choose the impurity function and the maxDepth.
		Many of the other hyperparameters used in sklearn are not here, but pruning and impurity is included.
	*/

	// Load Titanic Data For classification
	classificationData, err := base.ParseCSVToInstances("../../datasets/titanic.csv", false)
	if err != nil {
		panic(err)
	}
	trainData, testData := base.InstancesTrainTestSplit(classificationData, 0.5)

	// Create New Classification Tree
	// Hyperparameters - loss function, max Depth (-1 will split until pure), list of unique labels
	decTree := trees.NewDecisionTreeClassifier("entropy", -1, []int64{0, 1})

	// Train Tree
	err = decTree.Fit(trainData)
	if err != nil {
		panic(err)
	}
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
	regTree := trees.NewDecisionTreeRegressor("mse", -1)

	// Train Tree
	err = regTree.Fit(trainRegData)
	if err != nil {
		panic(err)
	}

	// Print out tree for visualization
	fmt.Println(regTree.String())

	// Access Predictions
	regressionPreds := regTree.Predict(testRegData)

	fmt.Println("Boston House Price Predictions")
	fmt.Println(regressionPreds)

}
