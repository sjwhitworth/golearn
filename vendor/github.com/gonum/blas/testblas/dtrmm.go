package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dtrmmer interface {
	Dtrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
}

func DtrmmTest(t *testing.T, blasser Dtrmmer) {
	for i, test := range []struct {
		s     blas.Side
		ul    blas.Uplo
		tA    blas.Transpose
		d     blas.Diag
		m     int
		n     int
		alpha float64
		a     [][]float64
		b     [][]float64
		ans   [][]float64
	}{
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3, 4},
				{0, 5, 6, 7},
				{0, 0, 8, 9},
				{0, 0, 0, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{320, 340, 360},
				{588, 624, 660},
				{598, 632, 666},
				{380, 400, 420},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2},
				{0, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{72, 78, 84},
				{130, 140, 150},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3, 4},
				{0, 5, 6, 7},
				{0, 0, 8, 9},
				{0, 0, 0, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{320, 340, 360},
				{484, 512, 540},
				{374, 394, 414},
				{38, 40, 42},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2},
				{0, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{72, 78, 84},
				{26, 28, 30},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0, 0},
				{2, 5, 0, 0},
				{3, 6, 8, 0},
				{4, 7, 9, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 22, 24},
				{170, 184, 198},
				{472, 506, 540},
				{930, 990, 1050},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0},
				{2, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 22, 24},
				{170, 184, 198},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0, 0},
				{2, 5, 0, 0},
				{3, 6, 8, 0},
				{4, 7, 9, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 22, 24},
				{66, 72, 78},
				{248, 268, 288},
				{588, 630, 672},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0},
				{2, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 22, 24},
				{66, 72, 78},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3, 4},
				{0, 5, 6, 7},
				{0, 0, 8, 9},
				{0, 0, 0, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 22, 24},
				{170, 184, 198},
				{472, 506, 540},
				{930, 990, 1050},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2},
				{0, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 22, 24},
				{170, 184, 198},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3, 4},
				{0, 5, 6, 7},
				{0, 0, 8, 9},
				{0, 0, 0, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 22, 24},
				{66, 72, 78},
				{248, 268, 288},
				{588, 630, 672},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2},
				{0, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 22, 24},
				{66, 72, 78},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0, 0},
				{2, 5, 0, 0},
				{3, 6, 8, 0},
				{4, 7, 9, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{320, 340, 360},
				{588, 624, 660},
				{598, 632, 666},
				{380, 400, 420},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0},
				{2, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{72, 78, 84},
				{130, 140, 150},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0, 0},
				{2, 5, 0, 0},
				{3, 6, 8, 0},
				{4, 7, 9, 10},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{320, 340, 360},
				{484, 512, 540},
				{374, 394, 414},
				{38, 40, 42},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0},
				{2, 5},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{72, 78, 84},
				{26, 28, 30},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 128, 314},
				{26, 164, 398},
				{32, 200, 482},
				{38, 236, 566},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 128, 314},
				{26, 164, 398},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 62, 194},
				{26, 80, 248},
				{32, 98, 302},
				{38, 116, 356},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 62, 194},
				{26, 80, 248},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{136, 208, 144},
				{172, 262, 180},
				{208, 316, 216},
				{244, 370, 252},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{136, 208, 144},
				{172, 262, 180},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{136, 142, 24},
				{172, 178, 30},
				{208, 214, 36},
				{244, 250, 42},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{136, 142, 24},
				{172, 178, 30},
			},
		},

		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{136, 208, 144},
				{172, 262, 180},
				{208, 316, 216},
				{244, 370, 252},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{136, 208, 144},
				{172, 262, 180},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{136, 142, 24},
				{172, 178, 30},
				{208, 214, 36},
				{244, 250, 42},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{136, 142, 24},
				{172, 178, 30},
			},
		},

		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 128, 314},
				{26, 164, 398},
				{32, 200, 482},
				{38, 236, 566},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 128, 314},
				{26, 164, 398},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{20, 62, 194},
				{26, 80, 248},
				{32, 98, 302},
				{38, 116, 356},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 2,
			a: [][]float64{
				{1, 0, 0},
				{2, 4, 0},
				{3, 5, 6},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{20, 62, 194},
				{26, 80, 248},
			},
		},
	} {
		aFlat := flatten(test.a)
		bFlat := flatten(test.b)
		ansFlat := flatten(test.ans)
		blasser.Dtrmm(test.s, test.ul, test.tA, test.d, test.m, test.n, test.alpha, aFlat, len(test.a[0]), bFlat, len(test.b[0]))
		if !floats.EqualApprox(ansFlat, bFlat, 1e-14) {
			t.Errorf("Case %v. Want %v, got %v.", i, ansFlat, bFlat)
		}
	}
}
