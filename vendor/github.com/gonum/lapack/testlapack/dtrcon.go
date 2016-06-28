package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
)

type Dtrconer interface {
	Dgeconer
	Dtrcon(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int, work []float64, iwork []int) float64
}

func DtrconTest(t *testing.T, impl Dtrconer) {
	rnd := rand.New(rand.NewSource(1))
	// Hand crafted tests.
	for _, test := range []struct {
		a       []float64
		n       int
		uplo    blas.Uplo
		diag    blas.Diag
		condOne float64
		condInf float64
	}{
		{
			a: []float64{
				8, 5, 6,
				0, 7, 8,
				0, 0, 6,
			},
			n:       3,
			uplo:    blas.Upper,
			diag:    blas.Unit,
			condOne: 1.0 / 645,
			condInf: 1.0 / 480,
		},
		{
			a: []float64{
				8, 5, 6,
				0, 7, 8,
				0, 0, 6,
			},
			n:       3,
			uplo:    blas.Upper,
			diag:    blas.NonUnit,
			condOne: 0.137704918032787,
			condInf: 0.157894736842105,
		},
		{
			a: []float64{
				8, 0, 0,
				5, 7, 0,
				6, 8, 6,
			},
			n:       3,
			uplo:    blas.Lower,
			diag:    blas.Unit,
			condOne: 1.0 / 480,
			condInf: 1.0 / 645,
		},
		{
			a: []float64{
				8, 0, 0,
				5, 7, 0,
				6, 8, 6,
			},
			n:       3,
			uplo:    blas.Lower,
			diag:    blas.NonUnit,
			condOne: 0.157894736842105,
			condInf: 0.137704918032787,
		},
	} {
		lda := test.n
		work := make([]float64, 3*test.n)
		for i := range work {
			work[i] = rnd.Float64()
		}
		iwork := make([]int, test.n)
		for i := range iwork {
			iwork[i] = int(rnd.Int31())
		}
		aCopy := make([]float64, len(test.a))
		copy(aCopy, test.a)
		condOne := impl.Dtrcon(lapack.MaxColumnSum, test.uplo, test.diag, test.n, test.a, lda, work, iwork)
		if math.Abs(condOne-test.condOne) > 1e-14 {
			t.Errorf("One norm mismatch. Want %v, got %v.", test.condOne, condOne)
		}
		if !floats.Equal(aCopy, test.a) {
			t.Errorf("a modified during call")
		}
		condInf := impl.Dtrcon(lapack.MaxRowSum, test.uplo, test.diag, test.n, test.a, lda, work, iwork)
		if math.Abs(condInf-test.condInf) > 1e-14 {
			t.Errorf("Inf norm mismatch. Want %v, got %v.", test.condInf, condInf)
		}
		if !floats.Equal(aCopy, test.a) {
			t.Errorf("a modified during call")
		}
	}

	// Dtrcon does not match the Dgecon output in many cases. See
	// https://github.com/xianyi/OpenBLAS/issues/636
	// TODO(btracey): Uncomment this when the mismatch between Dgecon and Dtrcon
	// is understood.
	/*
		// Randomized tests against Dgecon.
		for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
			for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
				for _, test := range []struct {
					n, lda int
				}{
					{3, 0},
					{4, 9},
				} {
					for trial := 0; trial < 1; trial++ {
						n := test.n
						lda := test.lda
						if lda == 0 {
							lda = n
						}
						a := make([]float64, n*lda)
						if trial == 0 {
							for i := range a {
								a[i] = float64(i + 2)
							}
						} else {
							for i := range a {
								a[i] = rnd.NormFloat64()
							}
						}

						aDense := make([]float64, len(a))
						if uplo == blas.Upper {
							for i := 0; i < n; i++ {
								for j := i; j < n; j++ {
									aDense[i*lda+j] = a[i*lda+j]
								}
							}
						} else {
							for i := 0; i < n; i++ {
								for j := 0; j <= i; j++ {
									aDense[i*lda+j] = a[i*lda+j]
								}
							}
						}
						if diag == blas.Unit {
							for i := 0; i < n; i++ {
								aDense[i*lda+i] = 1
							}
						}

						ipiv := make([]int, n)
						work := make([]float64, 4*n)
						denseOne := impl.Dlange(lapack.MaxColumnSum, n, n, aDense, lda, work)
						denseInf := impl.Dlange(lapack.MaxRowSum, n, n, aDense, lda, work)

						aDenseLU := make([]float64, len(aDense))
						copy(aDenseLU, aDense)
						impl.Dgetrf(n, n, aDenseLU, lda, ipiv)
						iwork := make([]int, n)
						want := impl.Dgecon(lapack.MaxColumnSum, n, aDenseLU, lda, denseOne, work, iwork)
						got := impl.Dtrcon(lapack.MaxColumnSum, uplo, diag, n, a, lda, work, iwork)
						if math.Abs(want-got) > 1e-14 {
							t.Errorf("One norm mismatch. Upper = %v, unit = %v, want %v, got %v", uplo == blas.Upper, diag == blas.Unit, want, got)
						}
						want = impl.Dgecon(lapack.MaxRowSum, n, aDenseLU, lda, denseInf, work, iwork)
						got = impl.Dtrcon(lapack.MaxRowSum, uplo, diag, n, a, lda, work, iwork)
						if math.Abs(want-got) > 1e-14 {
							t.Errorf("Inf norm mismatch. Upper = %v, unit = %v, want %v, got %v", uplo == blas.Upper, diag == blas.Unit, want, got)
						}
					}
				}
			}
		}
	*/
}
