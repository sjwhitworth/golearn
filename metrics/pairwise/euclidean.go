package pairwise

import (
	"fmt"
	"math"

	mat "github.com/skelterjohn/go.matrix"
)

// We may need to create Metrics / Vector interface for this
func Euclidean(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) (float64, error) {
	var sum float64

	difference, err := vectorY.MinusDense(vectorX)

	flat := difference.Array()

	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	for _, i := range flat {
		squared := math.Pow(i, 2)
		sum += squared
	}

	distance := math.Sqrt(sum)
	return distance, nil
}
