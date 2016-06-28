package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

type Dlantrer interface {
	Dlanger
	Dlantr(norm lapack.MatrixNorm, uplo blas.Uplo, diag blas.Diag, m, n int, a []float64, lda int, work []float64) float64
}

func DlantrTest(t *testing.T, impl Dlantrer) {
	rnd := rand.New(rand.NewSource(1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxColumnSum, lapack.MaxRowSum, lapack.NormFrob} {
		for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
			for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
				for _, test := range []struct {
					m, n, lda int
				}{
					{3, 3, 0},
					{3, 5, 0},
					{10, 5, 0},

					{5, 5, 11},
					{5, 10, 11},
					{10, 5, 11},
				} {
					// Do a couple of random trials since the values change.
					for trial := 0; trial < 100; trial++ {
						m := test.m
						n := test.n
						lda := test.lda
						if lda == 0 {
							lda = n
						}
						a := make([]float64, m*lda)
						if trial == 0 {
							for i := range a {
								a[i] = float64(i)
							}
						} else {
							for i := range a {
								a[i] = rnd.NormFloat64()
							}
						}
						aDense := make([]float64, len(a))
						if uplo == blas.Lower {
							for i := 0; i < m; i++ {
								for j := 0; j <= min(i, n-1); j++ {
									aDense[i*lda+j] = a[i*lda+j]
								}
							}
						} else {
							for i := 0; i < m; i++ {
								for j := i; j < n; j++ {
									aDense[i*lda+j] = a[i*lda+j]
								}
							}
						}
						if diag == blas.Unit {
							for i := 0; i < min(m, n); i++ {
								aDense[i*lda+i] = 1
							}
						}
						work := make([]float64, n+6)
						for i := range work {
							work[i] = rnd.Float64()
						}
						got := impl.Dlantr(norm, uplo, diag, m, n, a, lda, work)
						want := impl.Dlange(norm, m, n, aDense, lda, work)
						if math.Abs(got-want) > 1e-13 {
							t.Errorf("Norm mismatch. norm = %c, unitdiag = %v, upper = %v, m = %v, n = %v, lda = %v, Want %v, got %v.",
								norm, diag == blas.Unit, uplo == blas.Upper, m, n, lda, got, want)
						}
					}
				}
			}
		}
	}
}
