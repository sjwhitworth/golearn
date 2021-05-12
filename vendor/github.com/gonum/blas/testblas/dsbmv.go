package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dsbmver interface {
	Dsbmv(ul blas.Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DsbmvTest(t *testing.T, blasser Dsbmver) {
	for i, test := range []struct {
		ul    blas.Uplo
		n     int
		k     int
		alpha float64
		beta  float64
		a     [][]float64
		x     []float64
		y     []float64

		ans []float64
	}{
		{
			ul:    blas.Upper,
			n:     4,
			k:     2,
			alpha: 2,
			beta:  3,
			a: [][]float64{
				{7, 8, 2, 0},
				{0, 8, 2, -3},
				{0, 0, 3, 6},
				{0, 0, 0, 9},
			},
			x:   []float64{1, 2, 3, 4},
			y:   []float64{-1, -2, -3, -4},
			ans: []float64{55, 30, 69, 84},
		},
		{
			ul:    blas.Lower,
			n:     4,
			k:     2,
			alpha: 2,
			beta:  3,
			a: [][]float64{
				{7, 0, 0, 0},
				{8, 8, 0, 0},
				{2, 2, 3, 0},
				{0, -3, 6, 9},
			},
			x:   []float64{1, 2, 3, 4},
			y:   []float64{-1, -2, -3, -4},
			ans: []float64{55, 30, 69, 84},
		},
	} {
		extra := 0
		var aFlat []float64
		if test.ul == blas.Upper {
			aFlat = flattenBanded(test.a, test.k, 0)
		} else {
			aFlat = flattenBanded(test.a, 0, test.k)
		}
		incTest := func(incX, incY, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			ynew := makeIncremented(test.y, incY, extra)
			ans := makeIncremented(test.ans, incY, extra)
			blasser.Dsbmv(test.ul, test.n, test.k, test.alpha, aFlat, test.k+1, xnew, incX, test.beta, ynew, incY)
			if !dSliceTolEqual(ans, ynew) {
				t.Errorf("Case %v: Want %v, got %v", i, ans, ynew)
			}
		}
		incTest(1, 1, extra)
		incTest(1, 3, extra)
		incTest(1, -3, extra)
		incTest(2, 3, extra)
		incTest(2, -3, extra)
		incTest(3, 2, extra)
		incTest(-3, 2, extra)
	}
}
