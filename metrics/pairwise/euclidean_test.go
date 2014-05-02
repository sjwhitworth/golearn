package pairwise

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"

	mat "github.com/skelterjohn/go.matrix"
)

func TestEuclideanInnerProduct(t *testing.T) {
	euclidean := NewEuclidean()

	Convey("Given two vectors", t, func() {
		vectorX := mat.MakeDenseMatrix([]float64{1, 2, 3}, 3, 1)
		vectorY := mat.MakeDenseMatrix([]float64{3, 4, 5}, 3, 1)

		Convey("When doing inner product", func() {
			result := euclidean.InnerProduct(vectorX, vectorY)

			Convey("The result should be 26", func() {
				So(result, ShouldEqual, 26)
			})

		})

		Convey("When dimension not match", func() {
			vectorY = mat.MakeDenseMatrix([]float64{3, 4, 5}, 1, 3)

			Convey("It should panic with Dimension mismatch", func() {
				So(func() { euclidean.InnerProduct(vectorX, vectorY) }, ShouldPanicWith, "Dimension mismatch")
			})

		})

	})
}
