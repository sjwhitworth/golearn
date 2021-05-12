package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dtrtrier interface {
	Dtrconer
	Dtrtri(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int) bool
}

func DtrtriTest(t *testing.T, impl Dtrtrier) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for _, uplo := range []blas.Uplo{blas.Upper} {
		for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
			for _, test := range []struct {
				n, lda int
			}{
				{3, 0},
				{70, 0},
				{200, 0},
				{3, 5},
				{70, 92},
				{200, 205},
			} {
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				a := make([]float64, n*lda)
				for i := range a {
					a[i] = rnd.Float64() + 1 // This keeps the matrices well conditioned.
				}
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				impl.Dtrtri(uplo, diag, n, a, lda)
				if uplo == blas.Upper {
					for i := 1; i < n; i++ {
						for j := 0; j < i; j++ {
							aCopy[i*lda+j] = 0
							a[i*lda+j] = 0
						}
					}
				} else {
					for i := 1; i < n; i++ {
						for j := i + 1; j < n; j++ {
							aCopy[i*lda+j] = 0
							a[i*lda+j] = 0
						}
					}
				}
				if diag == blas.Unit {
					for i := 0; i < n; i++ {
						a[i*lda+i] = 1
						aCopy[i*lda+i] = 1
					}
				}
				ans := make([]float64, len(a))
				bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, a, lda, aCopy, lda, 0, ans, lda)
				iseye := true
				for i := 0; i < n; i++ {
					for j := 0; j < n; j++ {
						if i == j {
							if math.Abs(ans[i*lda+i]-1) > 1e-4 {
								iseye = false
								break
							}
						} else {
							if math.Abs(ans[i*lda+j]) > 1e-4 {
								iseye = false
								break
							}
						}
					}
				}
				if !iseye {
					t.Errorf("inv(A) * A != I. Upper = %v, unit = %v, n = %v, lda = %v",
						uplo == blas.Upper, diag == blas.Unit, n, lda)
				}
			}
		}
	}
}
