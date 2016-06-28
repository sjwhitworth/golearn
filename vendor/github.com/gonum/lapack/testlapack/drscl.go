package testlapack

import (
	"math"
	"testing"

	"github.com/gonum/floats"
)

type Drscler interface {
	Drscl(n int, a float64, x []float64, incX int)
}

func DrsclTest(t *testing.T, impl Drscler) {
	for _, test := range []struct {
		x []float64
		a float64
	}{
		{
			x: []float64{1, 2, 3, 4, 5},
			a: 4,
		},
		{
			x: []float64{1, 2, 3, 4, 5},
			a: math.MaxFloat64,
		},
		{
			x: []float64{1, 2, 3, 4, 5},
			a: 1e-307,
		},
	} {
		xcopy := make([]float64, len(test.x))
		copy(xcopy, test.x)

		// Cannot test the scaling directly because of floating point scaling issues
		// (the purpose of Drscl). Instead, check that scaling and scaling back
		// yeilds approximately x. If overflow or underflow occurs then the scaling
		// won't match.
		impl.Drscl(len(test.x), test.a, xcopy, 1)
		if floats.Equal(xcopy, test.x) {
			t.Errorf("x unchanged during call to drscl. a = %v, x = %v.", test.a, test.x)
		}
		impl.Drscl(len(test.x), 1/test.a, xcopy, 1)
		if !floats.EqualApprox(xcopy, test.x, 1e-14) {
			t.Errorf("x not equal after scaling and unscaling. a = %v, x = %v.", test.a, test.x)
		}
	}
}
