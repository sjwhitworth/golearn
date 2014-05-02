package pairwise

import (
	"math"

	mat "github.com/skelterjohn/go.matrix"
)

type Euclidean struct{}

func NewEuclidean() *Euclidean {
	return &Euclidean{}
}

func (self *Euclidean) InnerProduct(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) float64 {
	CheckDimMatch(vectorX, vectorY)

	result := mat.Product(mat.Transpose(vectorX), vectorY).Get(0, 0)

	return result
}

// We may need to create Metrics / Vector interface for this
func (self *Euclidean) Distance(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) (float64, error) {
	difference, err := vectorY.MinusDense(vectorX)
	result := self.InnerProduct(difference, difference)

	return math.Sqrt(result), err
}
