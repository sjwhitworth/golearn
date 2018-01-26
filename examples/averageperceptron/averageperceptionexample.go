package main

import (
	"fmt"
	"math/rand"

	base "github.com/amclay/golearn/base"
	evaluation "github.com/amclay/golearn/evaluation"
	perceptron "github.com/amclay/golearn/perceptron"
)

func main() {

	rand.Seed(4402201)

	rawData, err := base.ParseCSVToInstances("../datasets/house-votes-84.csv", true)
	if err != nil {
		panic(err)
	}

	//Initialises a new AveragePerceptron classifier
	cls := perceptron.NewAveragePerceptron(10, 1.2, 0.5, 0.3)

	//Do a training-test split
	trainData, testData := base.InstancesTrainTestSplit(rawData, 0.50)
	fmt.Println(trainData)
	fmt.Println(testData)
	cls.Fit(trainData)

	predictions := cls.Predict(testData)

	// Prints precision/recall metrics
	confusionMat, _ := evaluation.GetConfusionMatrix(testData, predictions)
	fmt.Println(evaluation.GetSummary(confusionMat))
}
