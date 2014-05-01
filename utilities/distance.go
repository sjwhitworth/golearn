package utilities

import (
	util "../utilities"
	"fmt"
	mat "github.com/skelterjohn/go.matrix"
	"math"
)

//Computes the Euclidean distance between two vectors.
func ComputeDistance(metric string, vector *mat.DenseMatrix, testrow *mat.DenseMatrix) float64 {
	var sum float64

	// Compute a variety of distance metrics
	switch metric:
	case "euclidean": {
		difference, err := testrow.MinusDense(vector)
		flat := difference.Array()

		if err != nil {
			fmt.Println(err)
		}

		for _, i := range flat {
			squared := math.Pow(i, 2)
			sum += squared
		}

		eucdistance := math.Sqrt(sum)
		return eucdistance
	}