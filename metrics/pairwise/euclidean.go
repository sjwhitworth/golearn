package pairwise

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type Euclidean struct{}

func NewEuclidean() *Euclidean {
	return &Euclidean{}
}

func (self *Euclidean) InnerProduct(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	result := vectorX.Dot(vectorY)

	return result
}

// We may need to create Metrics / Vector interface for this
func (self *Euclidean) Distance(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	subVector := mat64.NewDense(0, 0, nil)
	subVector.Sub(vectorX, vectorY)

	result := self.InnerProduct(subVector, subVector)

	return math.Sqrt(result)
}
