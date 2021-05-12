package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dgetrier interface {
	Dgetrfer
	Dgetri(n int, a []float64, lda int, ipiv []int, work []float64, lwork int) bool
}

func DgetriTest(t *testing.T, impl Dgetrier) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for _, test := range []struct {
		n, lda int
	}{
		{5, 0},
		{5, 8},
		{45, 0},
		{45, 50},
		{65, 0},
		{65, 70},
		{150, 0},
		{150, 250},
	} {
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		// Generate a random well conditioned matrix
		perm := rnd.Perm(n)
		a := make([]float64, n*lda)
		for i := 0; i < n; i++ {
			a[i*lda+perm[i]] = 1
		}
		for i := range a {
			a[i] += 0.01 * rnd.Float64()
		}
		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		ipiv := make([]int, n)
		// Compute LU decomposition.
		impl.Dgetrf(n, n, a, lda, ipiv)
		// Compute inverse.
		work := make([]float64, 1)
		impl.Dgetri(n, a, lda, ipiv, work, -1)
		work = make([]float64, int(work[0]))
		lwork := len(work)

		ok := impl.Dgetri(n, a, lda, ipiv, work, lwork)
		if !ok {
			t.Errorf("Unexpected singular matrix.")
		}

		// Check that A(inv) * A = I.
		ans := make([]float64, len(a))
		bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, aCopy, lda, a, lda, 0, ans, lda)
		isEye := true
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					// This tolerance is so high because computing matrix inverses
					// is very unstable.
					if math.Abs(ans[i*lda+j]-1) > 5e-2 {
						isEye = false
					}
				} else {
					if math.Abs(ans[i*lda+j]) > 5e-2 {
						isEye = false
					}
				}
			}
		}
		if !isEye {
			t.Errorf("Inv(A) * A != I. n = %v, lda = %v", n, lda)
		}
	}
}
