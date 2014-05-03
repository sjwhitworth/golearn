package pairwise

import (
	"testing"

	"github.com/gonum/matrix/mat64"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRBFKernel(t *testing.T) {
	rbfKernel := NewRBFKernel(0.1)

	Convey("Given two vectors", t, func() {
		vectorX := mat64.NewDense(3, 1, []float64{1, 2, 3})
		vectorY := mat64.NewDense(3, 1, []float64{2, 4, 5})

		Convey("When doing inner product", func() {
			result := rbfKernel.InnerProduct(vectorX, vectorY)

			Convey("The result should be 2.45960311115695", func() {
				So(result, ShouldEqual, 2.45960311115695)

			})
		})

	})
}
