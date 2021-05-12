package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dgetrser interface {
	Dgetrfer
	Dgetrs(trans blas.Transpose, n, nrhs int, a []float64, lda int, ipiv []int, b []float64, ldb int)
}

func DgetrsTest(t *testing.T, impl Dgetrser) {
	rnd := rand.New(rand.NewSource(1))
	// TODO(btracey): Put more thought into creating more regularized matrices
	// and what correct tolerances should be. Consider also seeding the random
	// number in this test to make it more robust to code changes in other
	// parts of the suite.
	for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
		for _, test := range []struct {
			n, nrhs, lda, ldb int
			tol               float64
		}{
			{3, 3, 0, 0, 1e-12},
			{3, 5, 0, 0, 1e-12},
			{5, 3, 0, 0, 1e-12},

			{3, 3, 8, 10, 1e-12},
			{3, 5, 8, 10, 1e-12},
			{5, 3, 8, 10, 1e-12},

			{300, 300, 0, 0, 1e-8},
			{300, 500, 0, 0, 1e-8},
			{500, 300, 0, 0, 1e-6},

			{300, 300, 700, 600, 1e-8},
			{300, 500, 700, 600, 1e-8},
			{500, 300, 700, 600, 1e-6},
		} {
			n := test.n
			nrhs := test.nrhs
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			ldb := test.ldb
			if ldb == 0 {
				ldb = nrhs
			}
			a := make([]float64, n*lda)
			for i := range a {
				a[i] = rnd.Float64()
			}
			b := make([]float64, n*ldb)
			for i := range b {
				b[i] = rnd.Float64()
			}
			aCopy := make([]float64, len(a))
			copy(aCopy, a)
			bCopy := make([]float64, len(b))
			copy(bCopy, b)

			ipiv := make([]int, n)
			for i := range ipiv {
				ipiv[i] = rnd.Int()
			}

			// Compute the LU factorization.
			impl.Dgetrf(n, n, a, lda, ipiv)
			// Solve the system of equations given the result.
			impl.Dgetrs(trans, n, nrhs, a, lda, ipiv, b, ldb)

			// Check that the system of equations holds.
			A := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: lda,
				Data:   aCopy,
			}
			B := blas64.General{
				Rows:   n,
				Cols:   nrhs,
				Stride: ldb,
				Data:   bCopy,
			}
			X := blas64.General{
				Rows:   n,
				Cols:   nrhs,
				Stride: ldb,
				Data:   b,
			}
			tmp := blas64.General{
				Rows:   n,
				Cols:   nrhs,
				Stride: ldb,
				Data:   make([]float64, n*ldb),
			}
			copy(tmp.Data, bCopy)
			blas64.Gemm(trans, blas.NoTrans, 1, A, X, 0, B)
			if !floats.EqualApprox(tmp.Data, bCopy, test.tol) {
				t.Errorf("Linear solve mismatch. trans = %v, n = %v, nrhs = %v, lda = %v, ldb = %v", trans, n, nrhs, lda, ldb)
			}
		}
	}
}
