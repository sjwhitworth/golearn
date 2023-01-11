package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type DgemvCase struct {
	Name  string
	m     int
	n     int
	A     [][]float64
	tA    blas.Transpose
	x     []float64
	incX  int
	y     []float64
	incY  int
	xCopy []float64
	yCopy []float64

	Subcases []DgemvSubcase
}

type DgemvSubcase struct {
	mulXNeg1 bool
	mulYNeg1 bool
	alpha    float64
	beta     float64
	ans      []float64
}

var DgemvCases = []DgemvCase{
	{
		Name: "M_gt_N_Inc1_NoTrans",
		tA:   blas.NoTrans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3},
		y:    []float64{7, 8, 9, 10, 11},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 8, 9, 10, 11},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{40.8, 43.9, 33, 9, 28},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 303.2, 210, 12, 158},
			},
		},
	},
	{
		Name: "M_gt_N_Inc1_Trans",
		tA:   blas.Trans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3, -4, 5},
		y:    []float64{7, 8, 9},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 8, 9},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{94.3, 40.2, 52.3},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{712.4, 273.6, 364.4},
			},
		},
	},
	{
		Name: "M_eq_N_Inc1_NoTrans",
		tA:   blas.NoTrans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3},
		y:    []float64{7, 2, 2},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 2, 2},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{40.8, 43.9, 33},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{40.8*8 - 6*7, 43.9*8 - 6*2, 33*8 - 6*2},
			},
		},
	},
	{
		Name: "M_eq_N_Inc1_Trans",
		tA:   blas.Trans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3},
		y:    []float64{7, 2, 2},

		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 261.6, 270.4},
			},
		},
	},
	{
		Name: "M_lt_N_Inc1_NoTrans",
		tA:   blas.NoTrans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 7},
			{9.6, 3.5, 9.1, -2, 9},
			{10, 7, 3, 1, -5},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3, -7.6, 8.1},
		y:    []float64{7, 2, 2},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 0, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 2, 2},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{21.5, 132, -15.1},
			},

			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{21.5*8 - 6*7, 132*8 - 6*2, -15.1*8 - 6*2},
			},
		},
	},
	{
		Name: "M_lt_N_Inc1_Trans",
		tA:   blas.Trans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 7},
			{9.6, 3.5, 9.1, -2, 9},
			{10, 7, 3, 1, -5},
		},
		incX: 1,
		incY: 1,
		x:    []float64{1, 2, 3},
		y:    []float64{7, 2, 2, -3, 5},

		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 261.6, 270.4, 90, 50},
			},
		},
	},
	{
		Name: "M_gt_N_Part1_NoTrans",
		tA:   blas.NoTrans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 1,
		incY: 2,
		x:    []float64{1, 2, 3},
		y:    []float64{7, 100, 8, 101, 9, 102, 10, 103, 11},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 100, 0, 101, 0, 102, 0, 103, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 100, 8, 101, 9, 102, 10, 103, 11},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{40.8, 100, 43.9, 101, 33, 102, 9, 103, 28},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 100, 303.2, 101, 210, 102, 12, 103, 158},
			},
		},
	},
	{
		Name: "M_gt_N_Part1_Trans",
		tA:   blas.Trans,
		m:    5,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 1,
		incY: 2,
		x:    []float64{1, 2, 3, -4, 5},
		y:    []float64{7, 100, 8, 101, 9},

		Subcases: []DgemvSubcase{
			{
				alpha: 0,
				beta:  0,
				ans:   []float64{0, 100, 0, 101, 0},
			},
			{
				alpha: 0,
				beta:  1,
				ans:   []float64{7, 100, 8, 101, 9},
			},
			{
				alpha: 1,
				beta:  0,
				ans:   []float64{94.3, 100, 40.2, 101, 52.3},
			},
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{712.4, 100, 273.6, 101, 364.4},
			},
		},
	},
	{
		Name: "M_gt_N_IncNot1_NoTrans",
		tA:   blas.NoTrans,
		m:    5,
		n:    3,

		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3},
		y:    []float64{7, 2, 6, 8, -4, -5, 9, 1, 1, 10, 19, 22, 11},
		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 2, 6, 303.2, -4, -5, 210, 1, 1, 12, 19, 22, 158},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{220.4, 2, 6, 311.2, -4, -5, 322, 1, 1, -4, 19, 22, 222},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{182, 2, 6, 24, -4, -5, 210, 1, 1, 291.2, 19, 22, 260.4},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{246, 2, 6, 8, -4, -5, 322, 1, 1, 299.2, 19, 22, 196.4},
			},
		},
	},
	{
		Name: "M_gt_N_IncNot1_Trans",
		tA:   blas.Trans,
		m:    5,
		n:    3,

		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
			{1, 1, 2},
			{9, 2, 5},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3, 8, -3, 6, 5},
		y:    []float64{7, 2, 6, 8, -4, -5, 9},
		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{720.4, 2, 6, 281.6, -4, -5, 380.4},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{219.6, 2, 6, 316, -4, -5, 195.6},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{392.4, 2, 6, 281.6, -4, -5, 708.4},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{207.6, 2, 6, 316, -4, -5, 207.6},
			},
		},
	},
	{
		Name: "M_eq_N_IncNot1_NoTrans",
		tA:   blas.NoTrans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3},
		y:    []float64{7, 2, 6, 8, -4, -5, 9},
		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{284.4, 2, 6, 303.2, -4, -5, 210},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{220.4, 2, 6, 311.2, -4, -5, 322},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{222, 2, 6, 303.2, -4, -5, 272.4},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{334, 2, 6, 311.2, -4, -5, 208.4},
			},
		},
	},
	{
		Name: "M_eq_N_IncNot1_Trans",
		tA:   blas.Trans,
		m:    3,
		n:    3,
		A: [][]float64{
			{4.1, 6.2, 8.1},
			{9.6, 3.5, 9.1},
			{10, 7, 3},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3},
		y:    []float64{7, 2, 6, 8, -4, -5, 9},

		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 2, 6, 225.6, -4, -5, 228.4},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{290, 2, 6, 212.8, -4, -5, 310},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{240.4, 2, 6, 225.6, -4, -5, 372.4},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{322, 2, 6, 212.8, -4, -5, 278},
			},
		},
	},
	{
		Name: "M_lt_N_IncNot1_NoTrans",
		tA:   blas.NoTrans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 11},
			{9.6, 3.5, 9.1, -3, -2},
			{10, 7, 3, -7, -4},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3, -2, -4, 8, -9},
		y:    []float64{7, 2, 6, 8, -4, -5, 9},

		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{-827.6, 2, 6, 543.2, -4, -5, 722},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{-93.2, 2, 6, -696.8, -4, -5, -1070},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{734, 2, 6, 543.2, -4, -5, -839.6},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{-1058, 2, 6, -696.8, -4, -5, -105.2},
			},
		},
	},
	{
		Name: "M_lt_N_IncNot1_Trans",
		tA:   blas.Trans,
		m:    3,
		n:    5,
		A: [][]float64{
			{4.1, 6.2, 8.1, 10, 11},
			{9.6, 3.5, 9.1, -3, -2},
			{10, 7, 3, -7, -4},
		},
		incX: 2,
		incY: 3,
		x:    []float64{1, 15, 2, 150, 3},
		y:    []float64{7, 2, 6, 8, -4, -5, 9, -4, -1, -9, 1, 1, 2},

		Subcases: []DgemvSubcase{
			{
				alpha: 8,
				beta:  -6,
				ans:   []float64{384.4, 2, 6, 225.6, -4, -5, 228.4, -4, -1, -82, 1, 1, -52},
			},
			{
				mulXNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{290, 2, 6, 212.8, -4, -5, 310, -4, -1, 190, 1, 1, 188},
			},
			{
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{-82, 2, 6, -184, -4, -5, 228.4, -4, -1, 327.6, 1, 1, 414.4},
			},
			{
				mulXNeg1: true,
				mulYNeg1: true,
				alpha:    8,
				beta:     -6,
				ans:      []float64{158, 2, 6, 88, -4, -5, 310, -4, -1, 314.8, 1, 1, 320},
			},
		},
	},

	// TODO: A can be longer than mxn. Add cases where it is longer
	// TODO: x and y can also be longer. Add tests for these
	// TODO: Add tests for dimension mismatch
	// TODO: Add places with a "submatrix view", where lda != m
}

type Dgemver interface {
	Dgemv(tA blas.Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DgemvTest(t *testing.T, blasser Dgemver) {
	for _, test := range DgemvCases {
		for i, cas := range test.Subcases {
			// Test that it passes with row-major
			dgemvcomp(t, test, cas, i, blasser)

			// Test the bad inputs
			dgemvbad(t, test, cas, i, blasser)
		}
	}
}

func dgemvcomp(t *testing.T, test DgemvCase, cas DgemvSubcase, i int, blasser Dgemver) {
	x := sliceCopy(test.x)
	y := sliceCopy(test.y)
	a := sliceOfSliceCopy(test.A)
	aFlat := flatten(a)

	lda := test.n

	incX := test.incX
	if cas.mulXNeg1 {
		incX *= -1
	}
	incY := test.incY
	if cas.mulYNeg1 {
		incY *= -1
	}

	f := func() {
		blasser.Dgemv(test.tA, test.m, test.n, cas.alpha, aFlat, lda, x, incX, cas.beta, y, incY)
	}
	if panics(f) {
		t.Errorf("Test %v case %v: unexpected panic", test.Name, i)
		if throwPanic {
			blasser.Dgemv(test.tA, test.m, test.n, cas.alpha, aFlat, lda, x, incX, cas.beta, y, incY)
		}
		return
	}
	// Check that x and a are unchanged
	if !dSliceEqual(x, test.x) {
		t.Errorf("Test %v, case %v: x modified during call", test.Name, i)
	}
	aFlat2 := flatten(sliceOfSliceCopy(test.A))
	if !dSliceEqual(aFlat2, aFlat) {
		t.Errorf("Test %v, case %v: a modified during call", test.Name, i)
	}

	// Check that the answer matches
	if !dSliceTolEqual(cas.ans, y) {
		t.Errorf("Test %v, case %v: answer mismatch: Expected %v, Found %v", test.Name, i, cas.ans, y)
	}
}

func dgemvbad(t *testing.T, test DgemvCase, cas DgemvSubcase, i int, blasser Dgemver) {
	x := sliceCopy(test.x)
	y := sliceCopy(test.y)
	a := sliceOfSliceCopy(test.A)
	aFlatRow := flatten(a)
	ldaRow := test.n

	f := func() {
		blasser.Dgemv(312, test.m, test.n, cas.alpha, aFlatRow, ldaRow, x, test.incX, cas.beta, y, test.incY)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for bad transpose", test.Name, i)
	}
	f = func() {
		blasser.Dgemv(test.tA, -2, test.n, cas.alpha, aFlatRow, ldaRow, x, test.incX, cas.beta, y, test.incY)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for m negative", test.Name, i)
	}
	f = func() {
		blasser.Dgemv(test.tA, test.m, -4, cas.alpha, aFlatRow, ldaRow, x, test.incX, cas.beta, y, test.incY)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for n negative", test.Name, i)
	}
	f = func() {
		blasser.Dgemv(test.tA, test.m, test.n, cas.alpha, aFlatRow, ldaRow, x, 0, cas.beta, y, test.incY)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for incX zero", test.Name, i)
	}
	f = func() {
		blasser.Dgemv(test.tA, test.m, test.n, cas.alpha, aFlatRow, ldaRow, x, test.incX, cas.beta, y, 0)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for incY zero", test.Name, i)
	}
	f = func() {
		blasser.Dgemv(test.tA, test.m, test.n, cas.alpha, aFlatRow, ldaRow-1, x, test.incX, cas.beta, y, test.incY)
	}
	if !panics(f) {
		t.Errorf("Test %v case %v: no panic for lda too small row major", test.Name, i)
	}
}
