// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dgelser interface {
	Dgels(trans blas.Transpose, m, n, nrhs int, a []float64, lda int, b []float64, ldb int, work []float64, lwork int) bool
}

func DgelsTest(t *testing.T, impl Dgelser) {
	rnd := rand.New(rand.NewSource(1))
	for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
		for _, test := range []struct {
			m, n, nrhs, lda, ldb int
		}{
			{3, 4, 5, 0, 0},
			{3, 5, 4, 0, 0},
			{4, 3, 5, 0, 0},
			{4, 5, 3, 0, 0},
			{5, 3, 4, 0, 0},
			{5, 4, 3, 0, 0},
			{3, 4, 5, 10, 20},
			{3, 5, 4, 10, 20},
			{4, 3, 5, 10, 20},
			{4, 5, 3, 10, 20},
			{5, 3, 4, 10, 20},
			{5, 4, 3, 10, 20},
			{3, 4, 5, 20, 10},
			{3, 5, 4, 20, 10},
			{4, 3, 5, 20, 10},
			{4, 5, 3, 20, 10},
			{5, 3, 4, 20, 10},
			{5, 4, 3, 20, 10},
			{200, 300, 400, 0, 0},
			{200, 400, 300, 0, 0},
			{300, 200, 400, 0, 0},
			{300, 400, 200, 0, 0},
			{400, 200, 300, 0, 0},
			{400, 300, 200, 0, 0},
			{200, 300, 400, 500, 600},
			{200, 400, 300, 500, 600},
			{300, 200, 400, 500, 600},
			{300, 400, 200, 500, 600},
			{400, 200, 300, 500, 600},
			{400, 300, 200, 500, 600},
			{200, 300, 400, 600, 500},
			{200, 400, 300, 600, 500},
			{300, 200, 400, 600, 500},
			{300, 400, 200, 600, 500},
			{400, 200, 300, 600, 500},
			{400, 300, 200, 600, 500},
		} {
			m := test.m
			n := test.n
			nrhs := test.nrhs

			lda := test.lda
			if lda == 0 {
				lda = n
			}
			a := make([]float64, m*lda)
			for i := range a {
				a[i] = rnd.Float64()
			}
			aCopy := make([]float64, len(a))
			copy(aCopy, a)

			// Size of b is the same trans or no trans, because the number of rows
			// has to be the max of (m,n).
			mb := max(m, n)
			nb := nrhs
			ldb := test.ldb
			if ldb == 0 {
				ldb = nb
			}
			b := make([]float64, mb*ldb)
			for i := range b {
				b[i] = rnd.Float64()
			}
			bCopy := make([]float64, len(b))
			copy(bCopy, b)

			// Find optimal work length.
			work := make([]float64, 1)
			impl.Dgels(trans, m, n, nrhs, a, lda, b, ldb, work, -1)

			// Perform linear solve
			work = make([]float64, int(work[0]))
			lwork := len(work)
			for i := range work {
				work[i] = rnd.Float64()
			}
			impl.Dgels(trans, m, n, nrhs, a, lda, b, ldb, work, lwork)

			// Check that the answer is correct by comparing to the normal equations.
			aMat := blas64.General{
				Rows:   m,
				Cols:   n,
				Stride: lda,
				Data:   make([]float64, len(aCopy)),
			}
			copy(aMat.Data, aCopy)
			szAta := n
			if trans == blas.Trans {
				szAta = m
			}
			aTA := blas64.General{
				Rows:   szAta,
				Cols:   szAta,
				Stride: szAta,
				Data:   make([]float64, szAta*szAta),
			}

			// Compute A^T * A if notrans and A * A^T otherwise.
			if trans == blas.NoTrans {
				blas64.Gemm(blas.Trans, blas.NoTrans, 1, aMat, aMat, 0, aTA)
			} else {
				blas64.Gemm(blas.NoTrans, blas.Trans, 1, aMat, aMat, 0, aTA)
			}

			// Multiply by X.
			X := blas64.General{
				Rows:   szAta,
				Cols:   nrhs,
				Stride: ldb,
				Data:   b,
			}
			ans := blas64.General{
				Rows:   aTA.Rows,
				Cols:   X.Cols,
				Stride: X.Cols,
				Data:   make([]float64, aTA.Rows*X.Cols),
			}
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aTA, X, 0, ans)

			B := blas64.General{
				Rows:   szAta,
				Cols:   nrhs,
				Stride: ldb,
				Data:   make([]float64, len(bCopy)),
			}

			copy(B.Data, bCopy)
			var ans2 blas64.General
			if trans == blas.NoTrans {
				ans2 = blas64.General{
					Rows:   aMat.Cols,
					Cols:   B.Cols,
					Stride: B.Cols,
					Data:   make([]float64, aMat.Cols*B.Cols),
				}
			} else {
				ans2 = blas64.General{
					Rows:   aMat.Rows,
					Cols:   B.Cols,
					Stride: B.Cols,
					Data:   make([]float64, aMat.Rows*B.Cols),
				}
			}

			// Compute A^T B if Trans or A * B otherwise
			if trans == blas.NoTrans {
				blas64.Gemm(blas.Trans, blas.NoTrans, 1, aMat, B, 0, ans2)
			} else {
				blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aMat, B, 0, ans2)
			}
			if !floats.EqualApprox(ans.Data, ans2.Data, 1e-12) {
				t.Errorf("Normal equations not satisfied")
			}
		}
	}
}
