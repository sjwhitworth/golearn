package main

import (
		mat "github.com/skelterjohn/go.matrix"
		base "golearn/base"
		util "golearn/utilities"
		knnclass "golearn/knn"
		"fmt"
		)

func main(){
	//Parses the infamous Iris data.
	cols, rows, _, labels, data := base.ParseCsv("datasets/iris.csv", 4, []int{0,1,2})

	//Initialises a new KNN classifier
	knn := knnclass.KNNClassifier{}
	knn.New("Testing", labels, data, rows, cols)
	
	for {
		//Creates a random array of N float64s between 0 and 7
		randArray := util.RandomArray(3)

		//Initialises a vector with this array
		random := mat.MakeDenseMatrix(randArray,1,3)

		//Calculates the Euclidean distance and returns the most popular label
		labels, _ := knn.Predict(random, 3)
		fmt.Println(labels)
	}
}