// Example of how to use Isolation Forest for outlier detection

package main

import (
	"fmt"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/trees"
)

func main() {
	/* Isolation Forest is used for outlier detection
		 The algorithm works by randomly splitting the data, so results won't be exactly reproducible
	 	 but generally outliers will still be classified as outliers. */

	// Load data for outlier detection - includes gaussian distribution, and ten outliers at the end
	// Dataset has 1000 normal datapoints, and 10 outliers at the ned
	csvData, err := base.ParseCSVToInstances("../../datasets/gaussian_outliers.csv", true)
	if err != nil {
		panic(err)
	}

	// Create New Isolation Forest with 100 trees, max depth 100, and each tree will use 850 datapoints
	forest := trees.NewIsolationForest(100, 100, 850)

	// fit the isolation forest to the data. Note that all class attributes are also used during training.
	// Remove all class attributes you don't want to use before calling fit.
	forest.Fit(csvData)

	// Make predictions. Generally, IsolationForest is used for Interpolation, not Extrapolation.
	// Predictions are returned as Anomaly Scores from 0 to 1. close to 0 - not outlier, close to 1 - outlier
	preds := forest.Predict(csvData)

	// Let's find the average and minimum Anomaly Score for normal data
	var avgScore float64
	var min float64
	min = 1

	for i := 0; i < 1000; i++ {
		temp := preds[i]
		avgScore += temp
		if temp < min {
			min = temp
		}
	}
	fmt.Println(avgScore / 1000)
	fmt.Println(min)

	// Now let's print the anomaly scores for the outliers.
	// You should find that these values are much higher (around 0.7) as comapred to the scores for normal data.
	fmt.Println("Anomaly Scores for outliers are ")
	for i := 1000; i < 1010; i++ {
		fmt.Print("      ")
		fmt.Println(preds[i])
	}
}
