// Demonstrates decision tree classification

package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/amclay/golearn/base"
	"github.com/amclay/golearn/ensemble"
	"github.com/amclay/golearn/evaluation"
)

func main() {

	var tree base.Classifier

	// Load in the iris dataset
	iris, err := base.ParseCSVToInstances("../datasets/iris_headers.csv", true)
	if err != nil {
		panic(err)
	}

	for i := 1; i < 60; i += 2 {
		// Demonstrate the effect of adding more trees to the forest
		// and also how much better it is without discretisation.
		rand.Seed(44111342)

		tree = ensemble.NewRandomForest(i, 4)
		cfs, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(iris, tree, 5)
		if err != nil {
			panic(err)
		}

		mean, variance := evaluation.GetCrossValidatedMetric(cfs, evaluation.GetAccuracy)
		stdev := math.Sqrt(variance)

		fmt.Printf("%d\t%.2f\t(+/- %.2f)\n", i, mean, stdev*2)
	}
}
