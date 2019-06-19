package pairwise

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
)

func TestChebyshev(t *testing.T) {
	var vectorX, vectorY *mat.Dense
	chebyshev := NewChebyshev()

	Convey("Given two vectors", t, func() {
		vectorX = mat.NewDense(4, 1, []float64{1, 2, 3, 4})
		vectorY = mat.NewDense(4, 1, []float64{-5, -6, 7, 8})

		Convey("When calculating distance with two vectors", func() {
			result := chebyshev.Distance(vectorX, vectorY)

			Convey("The result should be 8", func() {
				So(result, ShouldEqual, 8)
			})
		})

		Convey("When calculating distance with row vectors", func() {
			vectorX.Copy(vectorX.T())
			vectorY.Copy(vectorY.T())
			result := chebyshev.Distance(vectorX, vectorY)

			Convey("The result should be 8", func() {
				So(result, ShouldEqual, 8)
			})
		})

		Convey("When calculating distance with different dimension matrices", func() {
			vectorX.CloneFrom(vectorX.T())
			So(func() { chebyshev.Distance(vectorX, vectorY) }, ShouldPanic)
		})

	})
}
