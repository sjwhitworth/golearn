package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dtpmver interface {
	Dtpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int)
}

func DtpmvTest(t *testing.T, blasser Dtpmver) {
	for i, test := range []struct {
		n   int
		a   [][]float64
		x   []float64
		d   blas.Diag
		ul  blas.Uplo
		tA  blas.Transpose
		ans []float64
	}{
		{
			n: 3,
			a: [][]float64{
				{5, 6, 7},
				{0, 9, 10},
				{0, 0, 13},
			},
			x:   []float64{3, 4, 5},
			d:   blas.NonUnit,
			ul:  blas.Upper,
			tA:  blas.NoTrans,
			ans: []float64{74, 86, 65},
		},
		{
			n: 3,
			a: [][]float64{
				{5, 6, 7},
				{0, 9, 10},
				{0, 0, 13},
			},
			x:   []float64{3, 4, 5},
			d:   blas.Unit,
			ul:  blas.Upper,
			tA:  blas.NoTrans,
			ans: []float64{62, 54, 5},
		},
		{
			n: 3,
			a: [][]float64{
				{5, 0, 0},
				{6, 9, 0},
				{7, 10, 13},
			},
			x:   []float64{3, 4, 5},
			d:   blas.NonUnit,
			ul:  blas.Lower,
			tA:  blas.NoTrans,
			ans: []float64{15, 54, 126},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 0, 0},
				{6, 1, 0},
				{7, 10, 1},
			},
			x:   []float64{3, 4, 5},
			d:   blas.Unit,
			ul:  blas.Lower,
			tA:  blas.NoTrans,
			ans: []float64{3, 22, 66},
		},
		{
			n: 3,
			a: [][]float64{
				{5, 6, 7},
				{0, 9, 10},
				{0, 0, 13},
			},
			x:   []float64{3, 4, 5},
			d:   blas.NonUnit,
			ul:  blas.Upper,
			tA:  blas.Trans,
			ans: []float64{15, 54, 126},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 6, 7},
				{0, 1, 10},
				{0, 0, 1},
			},
			x:   []float64{3, 4, 5},
			d:   blas.Unit,
			ul:  blas.Upper,
			tA:  blas.Trans,
			ans: []float64{3, 22, 66},
		},
		{
			n: 3,
			a: [][]float64{
				{5, 0, 0},
				{6, 9, 0},
				{7, 10, 13},
			},
			x:   []float64{3, 4, 5},
			d:   blas.NonUnit,
			ul:  blas.Lower,
			tA:  blas.Trans,
			ans: []float64{74, 86, 65},
		},
	} {
		incTest := func(incX, extra int) {
			aFlat := flattenTriangular(test.a, test.ul)
			x := makeIncremented(test.x, incX, extra)
			blasser.Dtpmv(test.ul, test.tA, test.d, test.n, aFlat, x, incX)
			ans := makeIncremented(test.ans, incX, extra)
			if !floats.EqualApprox(x, ans, 1e-14) {
				t.Errorf("Case %v, idx %v: Want %v, got %v.", i, incX, ans, x)
			}
		}
		incTest(1, 0)
		incTest(-3, 3)
		incTest(4, 3)
	}
}
