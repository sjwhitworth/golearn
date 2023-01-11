package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dtbmver interface {
	Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int)
}

func DtbmvTest(t *testing.T, blasser Dtbmver) {
	for i, test := range []struct {
		ul  blas.Uplo
		tA  blas.Transpose
		d   blas.Diag
		n   int
		k   int
		a   [][]float64
		x   []float64
		ans []float64
	}{
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.Unit,
			n:  3,
			k:  1,
			a: [][]float64{
				{1, 2, 0},
				{0, 1, 4},
				{0, 0, 1},
			},
			x:   []float64{2, 3, 4},
			ans: []float64{8, 19, 4},
		},
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  1,
			a: [][]float64{
				{1, 3, 0, 0, 0},
				{0, 6, 7, 0, 0},
				{0, 0, 2, 1, 0},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 33, 10, 63, -5},
		},
		{
			ul: blas.Lower,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  1,
			a: [][]float64{
				{7, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{0, 7, 2, 0, 0},
				{0, 0, 1, 12, 0},
				{0, 0, 0, 3, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 15, 20, 51, 7},
		},
		{
			ul: blas.Upper,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{7, 3, 9, 0, 0},
				{0, 6, 7, 10, 0},
				{0, 0, 2, 1, 11},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 15, 29, 71, 40},
		},
		{
			ul: blas.Lower,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{7, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{9, 7, 2, 0, 0},
				{0, 10, 1, 12, 0},
				{0, 0, 11, 3, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{40, 73, 65, 63, -5},
		},
	} {
		extra := 0
		var aFlat []float64
		if test.ul == blas.Upper {
			aFlat = flattenBanded(test.a, test.k, 0)
		} else {
			aFlat = flattenBanded(test.a, 0, test.k)
		}
		incTest := func(incX, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			ans := makeIncremented(test.ans, incX, extra)
			lda := test.k + 1
			blasser.Dtbmv(test.ul, test.tA, test.d, test.n, test.k, aFlat, lda, xnew, incX)
			if !dSliceTolEqual(ans, xnew) {
				t.Errorf("Case %v, Inc %v: Want %v, got %v", i, incX, ans, xnew)
			}
		}
		incTest(1, extra)
		incTest(3, extra)
		incTest(-2, extra)
	}
}
