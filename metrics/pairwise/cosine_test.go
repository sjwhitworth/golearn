package pairwise

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"gonum.org/v1/gonum/mat"
)

func TestCosine(t *testing.T) {
	var vectorX, vectorY *mat.Dense
	cosine := NewCosine()

	Convey("Given two vectors", t, func() {
		vectorX = mat.NewDense(3, 1, []float64{1, 2, 3})
		vectorY = mat.NewDense(3, 1, []float64{2, 4, 6})

		Convey("When doing inner Dot", func() {
			result := cosine.Dot(vectorX, vectorY)

			Convey("The result should be 25", func() {
				So(result, ShouldEqual, 28)
			})
		})

		Convey("When calculating distance", func() {
			result := cosine.Distance(vectorX, vectorY)

			Convey("The result should be 0", func() {
				So(result, ShouldEqual, 0)
			})

		})

	})
}
