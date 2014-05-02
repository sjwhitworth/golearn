package pairwise

import (
	"errors"
	"math"

	mat "github.com/skelterjohn/go.matrix"
)

type PolyKernel struct {
	degree int
}

func NewPolyKernel(degree int) *PolyKernel {
	return &PolyKernel{degree: degree}
}

func (self *PolyKernel) InnerProduct(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) (float64, error) {
	if !CheckDimMatch(vectorX, vectorY) {
		return 0, errors.New("Dimension mismatch")
	}

	result := mat.Product(mat.Transpose(vectorX), vectorY).Get(0, 0)
	result = math.Pow(result+1, float64(self.degree))

	return result, nil
}

func (self *PolyKernel) Distance(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) (float64, error) {
	difference, err := vectorY.MinusDense(vectorX)
	result, err := self.InnerProduct(difference, difference)

	return math.Sqrt(result), err
}
