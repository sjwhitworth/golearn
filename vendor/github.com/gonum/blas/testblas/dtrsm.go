package testblas

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dtrsmer interface {
	Dtrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m, n int,
		alpha float64, a []float64, lda int, b []float64, ldb int)
}

func DtrsmTest(t *testing.T, blasser Dtrsmer) {
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
			m:     3,
			n:     2,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 5},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{1, 3.4},
				{-0.5, -0.5},
				{2, 3.2},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     3,
			n:     2,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 5},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{60, 96},
				{-42, -66},
				{10, 16},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     3,
			n:     4,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 5},
			},
			b: [][]float64{
				{3, 6, 2, 9},
				{4, 7, 1, 3},
				{5, 8, 9, 10},
			},
			ans: [][]float64{
				{1, 3.4, 1.2, 13},
				{-0.5, -0.5, -4, -3.5},
				{2, 3.2, 3.6, 4},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     3,
			n:     4,
			alpha: 2,
			a: [][]float64{
				{1, 2, 3},
				{0, 4, 5},
				{0, 0, 5},
			},
			b: [][]float64{
				{3, 6, 2, 9},
				{4, 7, 1, 3},
				{5, 8, 9, 10},
			},
			ans: [][]float64{
				{60, 96, 126, 146},
				{-42, -66, -88, -94},
				{10, 16, 18, 20},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 7},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{4.5, 9},
				{-0.375, -1.5},
				{-0.75, -12.0 / 7},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 7},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{9, 18},
				{-15, -33},
				{60, 132},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 7},
			},
			b: [][]float64{
				{3, 6, 2, 9},
				{4, 7, 1, 3},
				{5, 8, 9, 10},
			},
			ans: [][]float64{
				{4.5, 9, 3, 13.5},
				{-0.375, -1.5, -1.5, -63.0 / 8},
				{-0.75, -12.0 / 7, 3, 39.0 / 28},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 7},
			},
			b: [][]float64{
				{3, 6, 2, 9},
				{4, 7, 1, 3},
				{5, 8, 9, 10},
			},
			ans: [][]float64{
				{9, 18, 6, 27},
				{-15, -33, -15, -72},
				{60, 132, 87, 327},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{4.5, 9},
				{-0.30, -1.2},
				{-6.0 / 35, -24.0 / 35},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{9, 18},
				{-15, -33},
				{69, 150},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{3, 6, 6, 7},
				{4, 7, 8, 9},
				{5, 8, 10, 11},
			},
			ans: [][]float64{
				{4.5, 9, 9, 10.5},
				{-0.3, -1.2, -0.6, -0.9},
				{-6.0 / 35, -24.0 / 35, -12.0 / 35, -18.0 / 35},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{3, 6, 6, 7},
				{4, 7, 8, 9},
				{5, 8, 10, 11},
			},
			ans: [][]float64{
				{9, 18, 18, 21},
				{-15, -33, -30, -36},
				{69, 150, 138, 165},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 8},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{-0.46875, 0.375},
				{0.1875, 0.75},
				{1.875, 3},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     3,
			n:     2,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 8},
			},
			b: [][]float64{
				{3, 6},
				{4, 7},
				{5, 8},
			},
			ans: [][]float64{
				{168, 267},
				{-78, -123},
				{15, 24},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 8},
			},
			b: [][]float64{
				{3, 6, 2, 3},
				{4, 7, 4, 5},
				{5, 8, 6, 7},
			},
			ans: [][]float64{
				{-0.46875, 0.375, -2.0625, -1.78125},
				{0.1875, 0.75, -0.375, -0.1875},
				{1.875, 3, 2.25, 2.625},
			},
		},
		{
			s:     blas.Left,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     3,
			n:     4,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 4, 0},
				{5, 6, 8},
			},
			b: [][]float64{
				{3, 6, 2, 3},
				{4, 7, 4, 5},
				{5, 8, 6, 7},
			},
			ans: [][]float64{
				{168, 267, 204, 237},
				{-78, -123, -96, -111},
				{15, 24, 18, 21},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{15, -2.4, -48.0 / 35},
				{19.5, -3.3, -66.0 / 35},
				{24, -4.2, -2.4},
				{28.5, -5.1, -102.0 / 35},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{30, -57, 258},
				{39, -75, 339},
				{48, -93, 420},
				{57, -111, 501},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 7},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{15, -2.4, -48.0 / 35},
				{19.5, -3.3, -66.0 / 35},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{30, -57, 258},
				{39, -75, 339},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{4.2, 1.2, 4.5},
				{5.775, 1.65, 5.625},
				{7.35, 2.1, 6.75},
				{8.925, 2.55, 7.875},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{435, -183, 36},
				{543, -228, 45},
				{651, -273, 54},
				{759, -318, 63},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{4.2, 1.2, 4.5},
				{5.775, 1.65, 5.625},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.NoTrans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{435, -183, 36},
				{543, -228, 45},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{4.2, 1.2, 4.5},
				{5.775, 1.65, 5.625},
				{7.35, 2.1, 6.75},
				{8.925, 2.55, 7.875},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{435, -183, 36},
				{543, -228, 45},
				{651, -273, 54},
				{759, -318, 63},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{4.2, 1.2, 4.5},
				{5.775, 1.65, 5.625},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Upper,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 3, 4},
				{0, 5, 6},
				{0, 0, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{435, -183, 36},
				{543, -228, 45},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{15, -2.4, -1.2},
				{19.5, -3.3, -1.65},
				{24, -4.2, -2.1},
				{28.5, -5.1, -2.55},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     4,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
				{16, 17, 18},
				{19, 20, 21},
			},
			ans: [][]float64{
				{30, -57, 258},
				{39, -75, 339},
				{48, -93, 420},
				{57, -111, 501},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.NonUnit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{15, -2.4, -1.2},
				{19.5, -3.3, -1.65},
			},
		},
		{
			s:     blas.Right,
			ul:    blas.Lower,
			tA:    blas.Trans,
			d:     blas.Unit,
			m:     2,
			n:     3,
			alpha: 3,
			a: [][]float64{
				{2, 0, 0},
				{3, 5, 0},
				{4, 6, 8},
			},
			b: [][]float64{
				{10, 11, 12},
				{13, 14, 15},
			},
			ans: [][]float64{
				{30, -57, 258},
				{39, -75, 339},
			},
		},
	} {
		aFlat := flatten(test.a)
		bFlat := flatten(test.b)
		ansFlat := flatten(test.ans)
		var lda int
		if test.s == blas.Left {
			lda = test.m
		} else {
			lda = test.n
		}
		blasser.Dtrsm(test.s, test.ul, test.tA, test.d, test.m, test.n, test.alpha, aFlat, lda, bFlat, test.n)
		if !floats.EqualApprox(ansFlat, bFlat, 1e-13) {
			t.Errorf("Case %v: Want %v, got %v.", i, ansFlat, bFlat)
		}
	}
}
