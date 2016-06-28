package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dtpsver interface {
	Dtpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int)
}

func DtpsvTest(t *testing.T, blasser Dtpsver) {
	for i, test := range []struct {
		n   int
		a   [][]float64
		ul  blas.Uplo
		tA  blas.Transpose
		d   blas.Diag
		x   []float64
		ans []float64
	}{
		{
			n: 3,
			a: [][]float64{
				{1, 2, 3},
				{0, 8, 15},
				{0, 0, 8},
			},
			ul:  blas.Upper,
			tA:  blas.NoTrans,
			d:   blas.NonUnit,
			x:   []float64{5, 6, 7},
			ans: []float64{4.15625, -0.890625, 0.875},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 2, 3},
				{0, 1, 15},
				{0, 0, 1},
			},
			ul:  blas.Upper,
			tA:  blas.NoTrans,
			d:   blas.Unit,
			x:   []float64{5, 6, 7},
			ans: []float64{182, -99, 7},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 0, 0},
				{2, 8, 0},
				{3, 15, 8},
			},
			ul:  blas.Lower,
			tA:  blas.NoTrans,
			d:   blas.NonUnit,
			x:   []float64{5, 6, 7},
			ans: []float64{5, -0.5, -0.0625},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 0, 0},
				{2, 8, 0},
				{3, 15, 8},
			},
			ul:  blas.Lower,
			tA:  blas.NoTrans,
			d:   blas.Unit,
			x:   []float64{5, 6, 7},
			ans: []float64{5, -4, 52},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 2, 3},
				{0, 8, 15},
				{0, 0, 8},
			},
			ul:  blas.Upper,
			tA:  blas.Trans,
			d:   blas.NonUnit,
			x:   []float64{5, 6, 7},
			ans: []float64{5, -0.5, -0.0625},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 2, 3},
				{0, 8, 15},
				{0, 0, 8},
			},
			ul:  blas.Upper,
			tA:  blas.Trans,
			d:   blas.Unit,
			x:   []float64{5, 6, 7},
			ans: []float64{5, -4, 52},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 0, 0},
				{2, 8, 0},
				{3, 15, 8},
			},
			ul:  blas.Lower,
			tA:  blas.Trans,
			d:   blas.NonUnit,
			x:   []float64{5, 6, 7},
			ans: []float64{4.15625, -0.890625, 0.875},
		},
		{
			n: 3,
			a: [][]float64{
				{1, 0, 0},
				{2, 1, 0},
				{3, 15, 1},
			},
			ul:  blas.Lower,
			tA:  blas.Trans,
			d:   blas.Unit,
			x:   []float64{5, 6, 7},
			ans: []float64{182, -99, 7},
		},
	} {
		incTest := func(incX, extra int) {
			aFlat := flattenTriangular(test.a, test.ul)
			x := makeIncremented(test.x, incX, extra)
			blasser.Dtpsv(test.ul, test.tA, test.d, test.n, aFlat, x, incX)
			ans := makeIncremented(test.ans, incX, extra)
			if !floats.EqualApprox(x, ans, 1e-14) {
				t.Errorf("Case %v, incX = %v: Want %v, got %v.", i, incX, ans, x)
			}
		}
		incTest(1, 0)
		incTest(-2, 0)
		incTest(3, 0)
		incTest(-3, 8)
		incTest(4, 2)
	}
}
