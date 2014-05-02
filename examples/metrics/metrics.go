package main

import (
	"fmt"

	pariwiseMetrics "github.com/sjwhitworth/golearn/metrics/pairwise"
	"github.com/sjwhitworth/golearn/utilities"
	mat "github.com/skelterjohn/go.matrix"
)

func main() {
	randArray := utilities.RandomArray(3, 7)
	vectorX := mat.MakeDenseMatrix(randArray, 1, 3)
	randArray = utilities.RandomArray(3, 7)
	vectorY := mat.MakeDenseMatrix(randArray, 1, 3)

	distance, err := pariwiseMetrics.Euclidean(vectorX, vectorY)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Euclidean distance of " + vectorX.String() + " and " + vectorY.String() + " is: ")
	fmt.Println(distance)
}
