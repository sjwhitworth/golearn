package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dsyker interface {
	Dsyrk(ul blas.Uplo, tA blas.Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int)
}

func DsyrkTest(t *testing.T, blasser Dsyker) {
	for i, test := range []struct {
		ul    blas.Uplo
		tA    blas.Transpose
		n     int
		k     int
		alpha float64
		a     [][]float64
		c     [][]float64
		beta  float64
		ans   [][]float64
	}{
		{
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			n:     3,
			k:     2,
			alpha: 0,
			a: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			c: [][]float64{
				{1, 2, 3},
				{0, 5, 6},
				{0, 0, 9},
			},
			beta: 2,
			ans: [][]float64{
				{2, 4, 6},
				{0, 10, 12},
				{0, 0, 18},
			},
		},
		{
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			n:     3,
			k:     2,
			alpha: 0,
			a: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			c: [][]float64{
				{1, 0, 0},
				{2, 3, 0},
				{4, 5, 6},
			},
			beta: 2,
			ans: [][]float64{
				{2, 0, 0},
				{4, 6, 0},
				{8, 10, 12},
			},
		},
		{
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			n:     3,
			k:     2,
			alpha: 3,
			a: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			c: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			beta: 2,
			ans: [][]float64{
				{17, 37, 57},
				{0, 83, 127},
				{0, 0, 195},
			},
		},
		{
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			n:     3,
			k:     2,
			alpha: 3,
			a: [][]float64{
				{1, 2},
				{3, 4},
				{5, 6},
			},
			c: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			beta: 2,
			ans: [][]float64{
				{17, 0, 0},
				{37, 83, 0},
				{57, 127, 195},
			},
		},
		{
			ul:    blas.Upper,
			tA:    blas.Trans,
			n:     3,
			k:     2,
			alpha: 3,
			a: [][]float64{
				{1, 3, 5},
				{2, 4, 6},
			},
			c: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			beta: 2,
			ans: [][]float64{
				{17, 37, 57},
				{0, 83, 127},
				{0, 0, 195},
			},
		},
		{
			ul:    blas.Lower,
			tA:    blas.Trans,
			n:     3,
			k:     2,
			alpha: 3,
			a: [][]float64{
				{1, 3, 5},
				{2, 4, 6},
			},
			c: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			beta: 2,
			ans: [][]float64{
				{17, 0, 0},
				{37, 83, 0},
				{57, 127, 195},
			},
		},
	} {
		aFlat := flatten(test.a)
		cFlat := flatten(test.c)
		ansFlat := flatten(test.ans)
		blasser.Dsyrk(test.ul, test.tA, test.n, test.k, test.alpha, aFlat, len(test.a[0]), test.beta, cFlat, len(test.c[0]))
		if !floats.EqualApprox(ansFlat, cFlat, 1e-14) {
			t.Errorf("Case %v. Want %v, got %v.", i, ansFlat, cFlat)
		}
	}
}
