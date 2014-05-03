package distance

import (
	"testing"
	"math"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEucledianDistance(t *testing.T) {

	Convey("Given two vectors", t, func() {
		distance := eucledianDistance([]float64 { 1, -2, 3, 4 }, []float64 { -5, -6, 7, 8 })

		Convey("The result should be a square root of 84", func() {
			So(distance, ShouldEqual, math.Sqrt(84))
		})

	})

}
