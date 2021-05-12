package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dsyr2er interface {
	Dsyr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
}

func Dsyr2Test(t *testing.T, blasser Dsyr2er) {
	for i, test := range []struct {
		n     int
		a     [][]float64
		ul    blas.Uplo
		x     []float64
		y     []float64
		alpha float64
		ans   [][]float64
	}{
		{
			n: 3,
			a: [][]float64{
				{7, 2, 4},
				{0, 3, 5},
				{0, 0, 6},
			},
			x:     []float64{2, 3, 4},
			y:     []float64{5, 6, 7},
			alpha: 2,
			ul:    blas.Upper,
			ans: [][]float64{
				{47, 56, 72},
				{0, 75, 95},
				{0, 0, 118},
			},
		},
		{
			n: 3,
			a: [][]float64{
				{7, 0, 0},
				{2, 3, 0},
				{4, 5, 6},
			},
			x:     []float64{2, 3, 4},
			y:     []float64{5, 6, 7},
			alpha: 2,
			ul:    blas.Lower,
			ans: [][]float64{
				{47, 0, 0},
				{56, 75, 0},
				{72, 95, 118},
			},
		},
	} {
		incTest := func(incX, incY, extra int) {
			aFlat := flatten(test.a)
			x := makeIncremented(test.x, incX, extra)
			y := makeIncremented(test.y, incY, extra)
			blasser.Dsyr2(test.ul, test.n, test.alpha, x, incX, y, incY, aFlat, test.n)
			ansFlat := flatten(test.ans)
			if !floats.EqualApprox(aFlat, ansFlat, 1e-14) {
				t.Errorf("Case %v, incX = %v, incY = %v. Want %v, got %v.", i, incX, incY, ansFlat, aFlat)
			}
		}
		incTest(1, 1, 0)
		incTest(-2, 1, 0)
		incTest(-2, 3, 0)
		incTest(2, -3, 0)
		incTest(3, -2, 0)
		incTest(-3, -4, 0)
	}
}
