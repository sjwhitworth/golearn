package optimisation

import (
	"testing"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/matrix/mat64"
)

func init() {
	mat64.Register(cblas.Blas{})
}

// f(x) = 2x + 2y
// Parameters should be really, really close to 2.
func TestBGD(t *testing.T) {
	x := mat64.NewDense(2, 2, []float64{1, 3, 5, 8})
	y := mat64.NewDense(2, 1, []float64{8, 26})
	theta := mat64.NewDense(2, 1, []float64{0, 0})
	results := BatchGradientDescent(x, y, theta, 0.005, 10000)
	if results.At(0, 0) <= 1.99 || results.At(0, 0) >= 2.01 {
		t.Error("Innaccurate convergence of batch gradient descent")
	}
}

// f(x) = 2x + 2y
// Parameters should be really, really close to 2.
func TestSGD(t *testing.T) {
	x := mat64.NewDense(2, 2, []float64{1, 3, 5, 8})
	y := mat64.NewDense(2, 1, []float64{8, 26})
	theta := mat64.NewDense(2, 1, []float64{0, 0})
	results := StochasticGradientDescent(x, y, theta, 0.005, 10000, 30)
	if results.At(0, 0) <= 1.99 || results.At(0, 0) >= 2.01 {
		t.Error("Innaccurate convergence of batch gradient descent")
	}
}
