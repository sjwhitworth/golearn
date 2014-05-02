package pairwise

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	mat "github.com/skelterjohn/go.matrix"
)

func TestEuclidean(t *testing.T) {
	euclidean := NewEuclidean()

	Convey("Given two vectors", t, func() {
		vectorX := mat.MakeDenseMatrix([]float64{1, 2, 3}, 3, 1)
		vectorY := mat.MakeDenseMatrix([]float64{2, 4, 5}, 3, 1)

		Convey("When doing inner product", func() {
			result := euclidean.InnerProduct(vectorX, vectorY)

			Convey("The result should be 25", func() {
				So(result, ShouldEqual, 25)
			})

			Convey("When dimension not match", func() {
				vectorZ := mat.MakeDenseMatrix([]float64{3, 4, 5}, 1, 3)

				Convey("It should panic with Dimension mismatch", func() {
					So(func() { euclidean.InnerProduct(vectorX, vectorZ) }, ShouldPanicWith, "Dimension mismatch")
				})

			})

		})

		Convey("When calculating distance", func() {
			result, err := euclidean.Distance(vectorX, vectorY)

			Convey("The err should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("The result should be 3", func() {
				So(result, ShouldEqual, 3)
			})

		})

	})
}
