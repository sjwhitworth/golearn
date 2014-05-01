package main

import (
	"fmt"

	data "github.com/sjwhitworth/golearn/data"
	knnclass "github.com/sjwhitworth/golearn/knn"
	util "github.com/sjwhitworth/golearn/utilities"
	mat "github.com/skelterjohn/go.matrix"
)

func main() {
	//Parses the infamous Iris data.
	cols, rows, _, labels, data := data.ParseCsv("datasets/iris.csv", 4, []int{0, 1, 2})

	//Initialises a new KNN classifier
	knn := knnclass.KNNClassifier{}
	knn.New(labels, data, rows, cols, "euclidean")

	for {
		//Creates a random array of N float64s between 0 and 7
		randArray := util.RandomArray(3, 7)

		//Initialises a vector with this array
		random := mat.MakeDenseMatrix(randArray, 1, 3)

		//Calculates the Euclidean distance and returns the most popular label
		labels, _ := knn.Predict(random, 3)
		fmt.Println(labels)
	}
}
