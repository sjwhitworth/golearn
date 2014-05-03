package pairwise

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type Manhattan struct{}

func NewManhattan() *Manhattan {
	return &Manhattan{}
}

// Manhattan distance, also known as L1 distance.
// Compute sum of absolute values of elements.
func (self *Manhattan) Distance(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	var length int
	subVector := mat64.NewDense(0, 0, nil)
	subVector.Sub(vectorX, vectorY)

	r, c := subVector.Dims()

	if r == 1 {
		// Force transpose to column vector
		subVector.TCopy(subVector)
		length = c
	} else if c == 1 {
		length = r
	} else {
		panic(mat64.ErrShape)
	}

	result := .0
	for i := 0; i < length; i++ {
		result += math.Abs(subVector.At(i, 0))
	}

	return result
}
