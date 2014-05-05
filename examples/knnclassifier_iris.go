package main

import (
	"fmt"

	data "github.com/sjwhitworth/golearn/data"
	knn "github.com/sjwhitworth/golearn/knn"
	util "github.com/sjwhitworth/golearn/utilities"
)

func main() {
	//Parses the infamous Iris data.
	cols, rows, _, labels, data := data.ParseCsv("datasets/iris.csv", 4, []int{0, 1, 2})

	//Initialises a new KNN classifier
	cls := knn.NewKnnClassifier("euclidean")
	cls.Fit(labels, data, rows, cols)

	for {
		//Creates a random array of N float64s between 0 and 7
		randArray := util.RandomArray(3, 7)

		//Calculates the Euclidean distance and returns the most popular label
		labels := cls.Predict(randArray, 3)
		fmt.Println(labels)
	}
}
