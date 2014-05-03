package pairwise

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type Manhattan struct{}

func NewManhattan() *Manhattan {
	return &Manhattan{}
}

func (self *Manhattan) Distance(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	subVector := mat64.NewDense(0, 0, nil)
	subVector.Sub(vectorX, vectorY)

	r, _ := subVector.Dims()

	result := .0
	for i := 0; i < r; i++ {
		result += math.Abs(subVector.At(i, 0))
	}

	return result
}
