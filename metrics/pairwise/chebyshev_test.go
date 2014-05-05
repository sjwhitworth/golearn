package pairwise

import (
	"testing"

	"github.com/gonum/matrix/mat64"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChebyshev(t *testing.T) {
	var vectorX, vectorY *mat64.Dense
	chebyshev := NewChebyshev()

	Convey("Given two vectors", t, func() {
		vectorX = mat64.NewDense(4, 1, []float64{1, 2, 3, 4})
		vectorY = mat64.NewDense(4, 1, []float64{-5, -6, 7, 8})

		Convey("When calculating distance with two vectors", func() {
			result := chebyshev.Distance(vectorX, vectorY)

			Convey("The result should be 8", func() {
				So(result, ShouldEqual, 8)
			})
		})

		Convey("When calculating distance with row vectors", func() {
			vectorX.TCopy(vectorX)
			vectorY.TCopy(vectorY)
			result := chebyshev.Distance(vectorX, vectorY)

			Convey("The result should be 8", func() {
				So(result, ShouldEqual, 8)
			})
		})

		Convey("When calculating distance with different dimention matrices", func() {
			vectorX.TCopy(vectorX)
			So(func() { chebyshev.Distance(vectorX, vectorY) }, ShouldPanicWith, mat64.ErrShape)
		})

	})
}
