// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
)

type Dsyever interface {
	Dsyev(jobz lapack.EigComp, uplo blas.Uplo, n int, a []float64, lda int, w, work []float64, lwork int) (ok bool)
}

func DsyevTest(t *testing.T, impl Dsyever) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Lower, blas.Upper} {
		for _, test := range []struct {
			n, lda int
		}{
			{1, 0},
			{2, 0},
			{5, 0},
			{10, 0},
			{100, 0},

			{1, 5},
			{2, 5},
			{5, 10},
			{10, 20},
			{100, 110},
		} {
			for cas := 0; cas < 10; cas++ {
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				a := make([]float64, n*lda)
				for i := range a {
					a[i] = rnd.NormFloat64()
				}
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				w := make([]float64, n)
				for i := range w {
					w[i] = rnd.NormFloat64()
				}

				work := make([]float64, 1)
				impl.Dsyev(lapack.EigDecomp, uplo, n, a, lda, w, work, -1)
				work = make([]float64, int(work[0]))
				impl.Dsyev(lapack.EigDecomp, uplo, n, a, lda, w, work, len(work))

				// Check that the decomposition is correct
				orig := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, n*n),
				}
				if uplo == blas.Upper {
					for i := 0; i < n; i++ {
						for j := i; j < n; j++ {
							v := aCopy[i*lda+j]
							orig.Data[i*orig.Stride+j] = v
							orig.Data[j*orig.Stride+i] = v
						}
					}
				} else {
					for i := 0; i < n; i++ {
						for j := 0; j <= i; j++ {
							v := aCopy[i*lda+j]
							orig.Data[i*orig.Stride+j] = v
							orig.Data[j*orig.Stride+i] = v
						}
					}
				}

				V := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: lda,
					Data:   a,
				}

				if !eigenDecompCorrect(w, orig, V) {
					t.Errorf("Decomposition mismatch")
				}

				// Check that the decomposition is correct when the eigenvectors
				// are not computed.
				wAns := make([]float64, len(w))
				copy(wAns, w)
				copy(a, aCopy)
				for i := range w {
					w[i] = rnd.Float64()
				}
				for i := range work {
					work[i] = rnd.Float64()
				}
				impl.Dsyev(lapack.EigValueOnly, uplo, n, a, lda, w, work, len(work))
				if !floats.EqualApprox(w, wAns, 1e-8) {
					t.Errorf("Eigenvalue mismatch when vectors not computed")
				}
			}
		}
	}
}
