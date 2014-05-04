package pairwise

import (
	"testing"

	"github.com/gonum/matrix/mat64"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManhattan(t *testing.T) {
	var vectorX, vectorY *mat64.Dense
	manhattan := NewManhattan()

	Convey("Given two vectors", t, func() {
		vectorX = mat64.NewDense(3, 1, []float64{2, 2, 3})
		vectorY = mat64.NewDense(3, 1, []float64{1, 4, 5})

		Convey("When calculating distance with column vectors", func() {
			result := manhattan.Distance(vectorX, vectorY)

			Convey("The result should be 5", func() {
				So(result, ShouldEqual, 5)
			})
		})

		Convey("When calculating distance with row vectors", func() {
			vectorX.TCopy(vectorX)
			vectorY.TCopy(vectorY)
			result := manhattan.Distance(vectorX, vectorY)

			Convey("The result should be 5", func() {
				So(result, ShouldEqual, 5)
			})
		})

		Convey("When calculating distance with row and column vectors", func() {
			vectorX.TCopy(vectorX)
			So(func() { manhattan.Distance(vectorX, vectorY) }, ShouldPanicWith, mat64.ErrShape)
		})

	})
}
