package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dtxmver interface {
	Dtrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int)
	Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, x []float64, incX int)
}

type vec struct {
	data []float64
	inc  int
}

var cases = []struct {
	n, k       int
	ul         blas.Uplo
	d          blas.Diag
	ldab       int
	tr, tb, tp []float64
	ins        []vec
	solNoTrans []float64
	solTrans   []float64
}{
	{
		n:    3,
		k:    1,
		ul:   blas.Upper,
		d:    blas.NonUnit,
		tr:   []float64{1, 2, 0, 0, 3, 4, 0, 0, 5},
		tb:   []float64{1, 2, 3, 4, 5, 0},
		ldab: 2,
		tp:   []float64{1, 2, 0, 3, 4, 5},
		ins: []vec{
			{[]float64{2, 3, 4}, 1},
			{[]float64{2, 1, 3, 1, 4}, 2},
			{[]float64{4, 1, 3, 1, 2}, -2},
		},
		solNoTrans: []float64{8, 25, 20},
		solTrans:   []float64{2, 13, 32},
	},
	{
		n:    3,
		k:    1,
		ul:   blas.Upper,
		d:    blas.Unit,
		tr:   []float64{1, 2, 0, 0, 3, 4, 0, 0, 5},
		tb:   []float64{1, 2, 3, 4, 5, 0},
		ldab: 2,
		tp:   []float64{1, 2, 0, 3, 4, 5},
		ins: []vec{
			{[]float64{2, 3, 4}, 1},
			{[]float64{2, 1, 3, 1, 4}, 2},
			{[]float64{4, 1, 3, 1, 2}, -2},
		},
		solNoTrans: []float64{8, 19, 4},
		solTrans:   []float64{2, 7, 16},
	},
	{
		n:    3,
		k:    1,
		ul:   blas.Lower,
		d:    blas.NonUnit,
		tr:   []float64{1, 0, 0, 2, 3, 0, 0, 4, 5},
		tb:   []float64{0, 1, 2, 3, 4, 5},
		ldab: 2,
		tp:   []float64{1, 2, 3, 0, 4, 5},
		ins: []vec{
			{[]float64{2, 3, 4}, 1},
			{[]float64{2, 1, 3, 1, 4}, 2},
			{[]float64{4, 1, 3, 1, 2}, -2},
		},
		solNoTrans: []float64{2, 13, 32},
		solTrans:   []float64{8, 25, 20},
	},
	{
		n:    3,
		k:    1,
		ul:   blas.Lower,
		d:    blas.Unit,
		tr:   []float64{1, 0, 0, 2, 3, 0, 0, 4, 5},
		tb:   []float64{0, 1, 2, 3, 4, 5},
		ldab: 2,
		tp:   []float64{1, 2, 3, 0, 4, 5},
		ins: []vec{
			{[]float64{2, 3, 4}, 1},
			{[]float64{2, 1, 3, 1, 4}, 2},
			{[]float64{4, 1, 3, 1, 2}, -2},
		},
		solNoTrans: []float64{2, 7, 16},
		solTrans:   []float64{8, 19, 4},
	},
}

func DtxmvTest(t *testing.T, blasser Dtxmver) {

	for nc, c := range cases {
		for nx, x := range c.ins {
			in := make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtrmv(c.ul, blas.NoTrans, c.d, c.n, c.tr, c.n, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solNoTrans, 1) {
				t.Error("Wrong Dtrmv result for: NoTrans  in Case:", nc, "input:", nx)
			}

			in = make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtrmv(c.ul, blas.Trans, c.d, c.n, c.tr, c.n, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solTrans, 1) {
				t.Error("Wrong Dtrmv result for: Trans in Case:", nc, "input:", nx)
			}
			in = make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtbmv(c.ul, blas.NoTrans, c.d, c.n, c.k, c.tb, c.ldab, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solNoTrans, 1) {
				t.Error("Wrong Dtbmv result for: NoTrans  in Case:", nc, "input:", nx)
			}

			in = make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtbmv(c.ul, blas.Trans, c.d, c.n, c.k, c.tb, c.ldab, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solTrans, 1) {
				t.Error("Wrong Dtbmv result for: Trans in Case:", nc, "input:", nx)
			}
			in = make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtpmv(c.ul, blas.NoTrans, c.d, c.n, c.tp, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solNoTrans, 1) {
				t.Error("Wrong Dtpmv result for:  NoTrans  in Case:", nc, "input:", nx)
			}

			in = make([]float64, len(x.data))
			copy(in, x.data)
			blasser.Dtpmv(c.ul, blas.Trans, c.d, c.n, c.tp, in, x.inc)
			if !dStridedSliceTolEqual(c.n, in, x.inc, c.solTrans, 1) {
				t.Error("Wrong Dtpmv result for: Trans in Case:", nc, "input:", nx)
			}
		}
	}
}
