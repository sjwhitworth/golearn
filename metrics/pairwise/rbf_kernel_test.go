package pairwise

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
)

func TestRBFKernel(t *testing.T) {
	var vectorX, vectorY *mat.Dense
	rbfKernel := NewRBFKernel(0.1)

	Convey("Given two vectors", t, func() {
		vectorX = mat.NewDense(3, 1, []float64{1, 2, 3})
		vectorY = mat.NewDense(3, 1, []float64{2, 4, 5})

		Convey("When doing inner product", func() {
			result := rbfKernel.InnerProduct(vectorX, vectorY)

			Convey("The result should almost equal 0.4065696597405991", func() {
				So(result, ShouldAlmostEqual, 0.4065696597405991)

			})
		})

	})
}
