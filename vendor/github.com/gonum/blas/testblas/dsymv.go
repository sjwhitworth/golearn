package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dsymver interface {
	Dsymv(ul blas.Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DsymvTest(t *testing.T, blasser Dsymver) {
	for i, test := range []struct {
		ul    blas.Uplo
		n     int
		a     [][]float64
		x     []float64
		y     []float64
		alpha float64
		beta  float64
		ans   []float64
	}{
		{
			ul: blas.Upper,
			n:  3,
			a: [][]float64{
				{5, 6, 7},
				{0, 8, 10},
				{0, 0, 13},
			},
			x:     []float64{3, 4, 5},
			y:     []float64{6, 7, 8},
			alpha: 2.1,
			beta:  -3,
			ans:   []float64{137.4, 189, 240.6},
		},
		{
			ul: blas.Lower,
			n:  3,
			a: [][]float64{
				{5, 0, 0},
				{6, 8, 0},
				{7, 10, 13},
			},
			x:     []float64{3, 4, 5},
			y:     []float64{6, 7, 8},
			alpha: 2.1,
			beta:  -3,
			ans:   []float64{137.4, 189, 240.6},
		},
	} {
		incTest := func(incX, incY, extra int) {
			x := makeIncremented(test.x, incX, extra)
			y := makeIncremented(test.y, incY, extra)
			aFlat := flatten(test.a)
			ans := makeIncremented(test.ans, incY, extra)

			blasser.Dsymv(test.ul, test.n, test.alpha, aFlat, test.n, x, incX, test.beta, y, incY)
			if !floats.EqualApprox(ans, y, 1e-14) {
				t.Errorf("Case %v, incX=%v, incY=%v: Want %v, got %v.", i, incX, incY, ans, y)
			}
		}
		incTest(1, 1, 0)
		incTest(2, 3, 0)
		incTest(3, 2, 0)
		incTest(-3, 2, 0)
		incTest(-2, 4, 0)
		incTest(2, -1, 0)
		incTest(-3, -4, 3)
	}
}
