package testblas

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
)

func DgemmBenchmark(b *testing.B, dgemm Dgemmer, m, n, k int, tA, tB blas.Transpose) {
	a := make([]float64, m*k)
	for i := range a {
		a[i] = rand.Float64()
	}
	bv := make([]float64, k*n)
	for i := range bv {
		bv[i] = rand.Float64()
	}
	c := make([]float64, m*n)
	for i := range c {
		c[i] = rand.Float64()
	}
	var lda, ldb int
	if tA == blas.Trans {
		lda = m
	} else {
		lda = k
	}
	if tB == blas.Trans {
		ldb = k
	} else {
		ldb = n
	}
	ldc := n
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dgemm.Dgemm(tA, tB, m, n, k, 3.0, a, lda, bv, ldb, 1.0, c, ldc)
	}
}
