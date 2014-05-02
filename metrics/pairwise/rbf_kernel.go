package pairwise

import (
	"math"

	mat "github.com/skelterjohn/go.matrix"
)

type RBFKernel struct {
	gamma float64
}

func NewRBFKernel(gamma float64) *RBFKernel {
	return &RBFKernel{gamma: gamma}
}

func (self *RBFKernel) InnerProduct(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) (float64, error) {
	CheckDimMatch(vectorX, vectorY)

	euclidean := NewEuclidean()
	distance, err := euclidean.Distance(vectorX, vectorY)

	if err != nil {
		return 0, err
	}

	result := math.Exp(self.gamma * math.Pow(distance, 2))

	return result, nil
}
