package distance

import (
	"testing"
	"math"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEucledianDistance(t *testing.T) {

	Convey("Given two vectors", t, func() {
		distance := EucledianDistance([]float64 { 1, -2, 3, 4 }, []float64 { -5, -6, 7, 8 })

		Convey("The result should be a square root of 84", func() {
			So(distance, ShouldEqual, math.Sqrt(84))
		})

	})

}

func TestCranberraDistance(t *testing.T) {

	Convey("Given two vectors that are same", t, func() {
		vec := []float64 { 0, 1, -2, 3.4, 5, -6.7, 89 }
		distance := CranberraDistance(vec, vec)

		Convey("The result should be 0", func() {
			So(distance, ShouldEqual, 0)
		})

	})

	Convey("Given two vectors", t, func() {
		p1 := []float64 { 1, 2, 3, 4, 9 }
		p2 := []float64 { -5, -6, 7, 4, 3 };
		distance := CranberraDistance(p1, p2)

		Convey("The result should be 0", func() {
			So(distance, ShouldEqual, 2.9)
		})

	})

}


func TestChebyshevDistance(t *testing.T) {

	Convey("Given two vectors", t, func() {
		p1 := []float64 { 1, 2, 3, 4 }
		p2 := []float64 { -5, -6, 7, 8 }

		distance := ChebyshevDistance(p1, p2)

		Convey("The result should be a square root of 8", func() {
			So(distance, ShouldEqual, 8)
		})

	})

}
