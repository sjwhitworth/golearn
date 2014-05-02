package pairwise

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type RBFKernel struct {
	gamma float64
}

func NewRBFKernel(gamma float64) *RBFKernel {
	return &RBFKernel{gamma: gamma}
}

func (self *RBFKernel) InnerProduct(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	euclidean := NewEuclidean()
	distance := euclidean.Distance(vectorX, vectorY)

	result := math.Exp(self.gamma * math.Pow(distance, 2))

	return result
}
