package optimisation

import (
	"testing"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/matrix/mat64"
	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	mat64.Register(cblas.Blas{})
}

func TestGradientDescent(t *testing.T) {
	Convey("When y = 2x_0 + 2x_1", t, func() {
		x := mat64.NewDense(2, 2, []float64{1, 3, 5, 8})
		y := mat64.NewDense(2, 1, []float64{8, 26})

		Convey("When estimating the parameters with Batch Gradient Descent", func() {
			theta := mat64.NewDense(2, 1, []float64{0, 0})
			results := BatchGradientDescent(x, y, theta, 0.005, 10000)

			Convey("The estimated parameters should be really close to 2, 2", func() {
				So(results.At(0, 0), ShouldAlmostEqual, 2.0, 0.01)
			})
		})

		Convey("When estimating the parameters with Stochastic Gradient Descent", func() {
			theta := mat64.NewDense(2, 1, []float64{0, 0})
			results := StochasticGradientDescent(x, y, theta, 0.005, 10000, 30)

			Convey("The estimated parameters should be really close to 2, 2", func() {
				So(results.At(0, 0), ShouldAlmostEqual, 2.0, 0.01)
			})
		})
	})
}
