package pairwise

import (
	"math"

	"github.com/gonum/matrix/mat64"
)

type PolyKernel struct {
	degree int
}

func NewPolyKernel(degree int) *PolyKernel {
	return &PolyKernel{degree: degree}
}

func (self *PolyKernel) InnerProduct(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	result := vectorX.Dot(vectorY)
	result = math.Pow(result+1, float64(self.degree))

	return result
}

func (self *PolyKernel) Distance(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {
	subVector := mat64.NewDense(0, 0, nil)
	subVector.Sub(vectorX, vectorY)
	result := self.InnerProduct(subVector, subVector)

	return math.Sqrt(result)
}
