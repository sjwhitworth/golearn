package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

type Dlansyer interface {
	Dlanger
	Dlansy(norm lapack.MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int, work []float64) float64
}

func DlansyTest(t *testing.T, impl Dlansyer) {
	rnd := rand.New(rand.NewSource(1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxColumnSum, lapack.MaxRowSum, lapack.NormFrob} {
		for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
			for _, test := range []struct {
				n, lda int
			}{
				{1, 0},
				{3, 0},

				{1, 10},
				{3, 10},
			} {
				for trial := 0; trial < 100; trial++ {
					n := test.n
					lda := test.lda
					if lda == 0 {
						lda = n
					}
					a := make([]float64, lda*n)
					if trial == 0 {
						for i := range a {
							a[i] = float64(i)
						}
					} else {
						for i := range a {
							a[i] = rnd.NormFloat64()
						}
					}

					aDense := make([]float64, n*n)
					if uplo == blas.Upper {
						for i := 0; i < n; i++ {
							for j := i; j < n; j++ {
								v := a[i*lda+j]
								aDense[i*n+j] = v
								aDense[j*n+i] = v
							}
						}
					} else {
						for i := 0; i < n; i++ {
							for j := 0; j <= i; j++ {
								v := a[i*lda+j]
								aDense[i*n+j] = v
								aDense[j*n+i] = v
							}
						}
					}
					work := make([]float64, n)
					got := impl.Dlansy(norm, uplo, n, a, lda, work)
					want := impl.Dlange(norm, n, n, aDense, n, work)
					if math.Abs(want-got) > 1e-14 {
						t.Errorf("Norm mismatch. norm = %c, upper = %v, n = %v, lda = %v, want %v, got %v.",
							norm, uplo == blas.Upper, n, lda, got, want)
					}
				}
			}
		}
	}
}
