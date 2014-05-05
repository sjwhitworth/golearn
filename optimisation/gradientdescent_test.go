package optimisation

import (
	"testing"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/matrix/mat64"
)

func init() {
	mat64.Register(cblas.Blas{})
}

func TestBGD(t *testing.T) {
	x := mat64.NewDense(2, 3, []float64{1, 1, 2, 1, 2, 3})
	y := mat64.NewDense(2, 1, []float64{3, 4})
	theta := mat64.NewDense(3, 1, []float64{1, 1, 1})
	results := BatchGradientDescent(x, y, theta, 0.0001, 10000)
	if results.At(0, 0) < 0.880 || results.At(0, 0) > 0.881 {
		t.Error("Innaccurate convergence of batch gradient descent")
	}
}

func TestSGD(t *testing.T) {
	x := mat64.NewDense(2, 3, []float64{1, 1, 2, 1, 2, 3})
	y := mat64.NewDense(2, 1, []float64{3, 4})
	theta := mat64.NewDense(3, 1, []float64{1, 1, 1})
	results := StochasticGradientDescent(x, y, theta, 0.0001, 10000)
	if results.At(0, 0) < 0.880 || results.At(0, 0) > 0.881 {
		t.Error("Innaccurate convergence of batch gradient descent")
	}
}
