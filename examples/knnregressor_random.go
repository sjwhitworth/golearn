package main

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
	data "github.com/sjwhitworth/golearn/data"
	knn "github.com/sjwhitworth/golearn/knn"
	util "github.com/sjwhitworth/golearn/utilities"
)

func main() {
	//Parses the infamous Iris data.
	cols, rows, _, labels, data := data.ParseCsv("datasets/randomdata.csv", 2, []int{0, 1})
	newlabels := util.ConvertLabelsToFloat(labels)

	//Initialises a new KNN classifier
	cls := knn.NewKnnRegressor("euclidean")
	cls.Fit(newlabels, data, rows, cols)

	for {
		//Creates a random array of N float64s between 0 and Y
		randArray := util.RandomArray(2, 100)

		//Initialises a vector with this array
		random := mat64.NewDense(1, 2, randArray)

		//Calculates the Euclidean distance and returns the most popular label
		outcome := cls.Predict(random, 3)
		fmt.Println(outcome)
	}
}
