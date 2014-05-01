package utilities

import (
	"fmt"
	"math"

	mat "github.com/skelterjohn/go.matrix"
)

// Computes the 'distance' between two vectors, where the distance is one of the following methods -
// euclidean (more to come)
func ComputeDistance(metric string, vector *mat.DenseMatrix, testrow *mat.DenseMatrix) (float64, error) {
	var sum float64

	switch metric {
	case "euclidean":
		{
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
			return eucdistance, nil
		}
	default:
		return 0.0, fmt.Errorf("ValueError: %s is not an implemented distance method", metric)
	}
}
