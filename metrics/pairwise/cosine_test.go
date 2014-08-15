package pairwise

import (
	"testing"

	"github.com/gonum/matrix/mat64"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManCosine(t *testing.T) {
	var vectorX, vectorY *mat64.Dense
	cosine := NewCosine()

	Convey("Given two vectors that are same", t, func() {
		vec := mat64.NewDense(7, 1, []float64{0, 1, -2, 3.4, 5, -6.7, 89})
		distance := cosine.Distance(vec, vec)

		Convey("The result should be 1", func() {
			So(distance, ShouldAlmostEqual, 1.00)
		})
	})

	Convey("Given two vectors", t, func() {
		vectorX = mat64.NewDense(3, 1, []float64{2, 2, 3})
		vectorY = mat64.NewDense(3, 1, []float64{1, 4, 5})

		Convey("When calculating distance with column vectors", func() {
			result := cosine.Distance(vectorX, vectorY)

			Convey("The result should be 0.9356", func() {
				So(result, ShouldAlmostEqual, 0.9356, 0.0001)
			})
		})

		Convey("When calculating distance with row vectors", func() {
			vectorX.TCopy(vectorX)
			vectorY.TCopy(vectorY)
			result := cosine.Distance(vectorX, vectorY)

			Convey("The result should be 0.9356", func() {
				So(result, ShouldAlmostEqual, 0.9356, 0.0001)
			})
		})

		Convey("When calculating distance with different dimention matrices", func() {
			vectorX.TCopy(vectorX)
			So(func() { cosine.Distance(vectorX, vectorY) }, ShouldPanicWith, mat64.ErrShape)
		})

	})
}
