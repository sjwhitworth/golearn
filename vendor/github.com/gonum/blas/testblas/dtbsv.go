package testblas

import (
	"testing"

	"github.com/gonum/blas"
)

type Dtbsver interface {
	Dtbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtrsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int)
}

func DtbsvTest(t *testing.T, blasser Dtbsver) {
	for i, test := range []struct {
		ul   blas.Uplo
		tA   blas.Transpose
		d    blas.Diag
		n, k int
		a    [][]float64
		lda  int
		x    []float64
		incX int
		ans  []float64
	}{
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
			x:    []float64{1, 2, 3, 4, 5},
			incX: 1,
			ans:  []float64{2.479166666666667, -0.493055555555556, 0.708333333333333, 1.583333333333333, -5.000000000000000},
		},
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 3, 5, 0, 0},
				{0, 6, 7, 5, 0},
				{0, 0, 2, 1, 5},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:    []float64{1, 2, 3, 4, 5},
			incX: 1,
			ans:  []float64{-15.854166666666664, -16.395833333333336, 13.208333333333334, 1.583333333333333, -5.000000000000000},
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
			x:    []float64{1, -101, 2, -201, 3, -301, 4, -401, 5, -501, -601, -701},
			incX: 2,
			ans:  []float64{2.479166666666667, -101, -0.493055555555556, -201, 0.708333333333333, -301, 1.583333333333333, -401, -5.000000000000000, -501, -601, -701},
		},
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 3, 5, 0, 0},
				{0, 6, 7, 5, 0},
				{0, 0, 2, 1, 5},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:    []float64{1, -101, 2, -201, 3, -301, 4, -401, 5, -501, -601, -701},
			incX: 2,
			ans:  []float64{-15.854166666666664, -101, -16.395833333333336, -201, 13.208333333333334, -301, 1.583333333333333, -401, -5.000000000000000, -501, -601, -701},
		},
		{
			ul: blas.Lower,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{5, 7, 2, 0, 0},
				{0, 5, 1, 12, 0},
				{0, 0, 5, 3, -1},
			},
			x:    []float64{1, 2, 3, 4, 5},
			incX: 1,
			ans:  []float64{1, -0.166666666666667, -0.416666666666667, 0.437500000000000, -5.770833333333334},
		},
		{
			ul: blas.Lower,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{5, 7, 2, 0, 0},
				{0, 5, 1, 12, 0},
				{0, 0, 5, 3, -1},
			},
			x:    []float64{1, -101, 2, -201, 3, -301, 4, -401, 5, -501, -601, -701},
			incX: 2,
			ans:  []float64{1, -101, -0.166666666666667, -201, -0.416666666666667, -301, 0.437500000000000, -401, -5.770833333333334, -501, -601, -701},
		},
		{
			ul: blas.Upper,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 3, 5, 0, 0},
				{0, 6, 7, 5, 0},
				{0, 0, 2, 1, 5},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:    []float64{1, 2, 3, 4, 5},
			incX: 1,
			ans:  []float64{1, -0.166666666666667, -0.416666666666667, 0.437500000000000, -5.770833333333334},
		},
		{
			ul: blas.Upper,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 3, 5, 0, 0},
				{0, 6, 7, 5, 0},
				{0, 0, 2, 1, 5},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:    []float64{1, -101, 2, -201, 3, -301, 4, -401, 5, -501, -601, -701},
			incX: 2,
			ans:  []float64{1, -101, -0.166666666666667, -201, -0.416666666666667, -301, 0.437500000000000, -401, -5.770833333333334, -501, -601, -701},
		},
		{
			ul: blas.Lower,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{5, 7, 2, 0, 0},
				{0, 5, 1, 12, 0},
				{0, 0, 5, 3, -1},
			},
			x:    []float64{1, 2, 3, 4, 5},
			incX: 1,
			ans:  []float64{-15.854166666666664, -16.395833333333336, 13.208333333333334, 1.583333333333333, -5.000000000000000},
		},
		{
			ul: blas.Lower,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{1, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{5, 7, 2, 0, 0},
				{0, 5, 1, 12, 0},
				{0, 0, 5, 3, -1},
			},
			x:    []float64{1, -101, 2, -201, 3, -301, 4, -401, 5, -501, -601, -701},
			incX: 2,
			ans:  []float64{-15.854166666666664, -101, -16.395833333333336, -201, 13.208333333333334, -301, 1.583333333333333, -401, -5.000000000000000, -501, -601, -701},
		},
	} {
		var aFlat []float64
		if test.ul == blas.Upper {
			aFlat = flattenBanded(test.a, test.k, 0)
		} else {
			aFlat = flattenBanded(test.a, 0, test.k)
		}
		xCopy := sliceCopy(test.x)
		// TODO: Have tests where the banded matrix is constructed explicitly
		// to allow testing for lda =! k+1
		blasser.Dtbsv(test.ul, test.tA, test.d, test.n, test.k, aFlat, test.k+1, xCopy, test.incX)
		if !dSliceTolEqual(test.ans, xCopy) {
			t.Errorf("Case %v: Want %v, got %v", i, test.ans, xCopy)
		}
	}

	/*
		// TODO: Uncomment when Dtrsv is fixed
		// Compare with dense for larger matrices
		for _, ul := range [...]blas.Uplo{blas.Upper, blas.Lower} {
			for _, tA := range [...]blas.Transpose{blas.NoTrans, blas.Trans} {
				for _, n := range [...]int{7, 8, 11} {
					for _, d := range [...]blas.Diag{blas.NonUnit, blas.Unit} {
						for _, k := range [...]int{0, 1, 3} {
							for _, incX := range [...]int{1, 3} {
								a := make([][]float64, n)
								for i := range a {
									a[i] = make([]float64, n)
									for j := range a[i] {
										a[i][j] = rand.Float64()
									}
								}
								x := make([]float64, n)
								for i := range x {
									x[i] = rand.Float64()
								}
								extra := 3
								xinc := makeIncremented(x, incX, extra)
								bandX := sliceCopy(xinc)
								var aFlatBand []float64
								if ul == blas.Upper {
									aFlatBand = flattenBanded(a, k, 0)
								} else {
									aFlatBand = flattenBanded(a, 0, k)
								}
								blasser.Dtbsv(ul, tA, d, n, k, aFlatBand, k+1, bandX, incX)

								aFlatDense := flatten(a)
								denseX := sliceCopy(xinc)
								blasser.Dtrsv(ul, tA, d, n, aFlatDense, n, denseX, incX)
								if !dSliceTolEqual(denseX, bandX) {
									t.Errorf("Case %v: dense banded mismatch")
								}
							}
						}
					}
				}
			}
		}
	*/
}
