package main

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	evaluation "github.com/sjwhitworth/golearn/evaluation"
	knn "github.com/sjwhitworth/golearn/knn"
)

func main() {
	rawData, err := base.ParseCSVToInstances("datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}
	rawData.Shuffle()
	//Initialises a new KNN classifier
	cls := knn.NewKnnClassifier("euclidean", 2)

	//Do a training-test split
	trainTest := base.InstancesTrainTestSplit(rawData, 0.50)
	trainData := trainTest[0]
	testData := trainTest[1]
	cls.Fit(trainData)

	//Calculates the Euclidean distance and returns the most popular label
	predictions := cls.Predict(testData)
	fmt.Println(predictions)

	// Prints precision/recall metrics
	confusionMat := evaluation.GetConfusionMatrix(testData, predictions)
	fmt.Println(evaluation.GetSummary(confusionMat))
}
