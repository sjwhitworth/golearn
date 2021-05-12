// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dgeqr2er interface {
	Dgeqr2(m, n int, a []float64, lda int, tau []float64, work []float64)
}

func Dgeqr2Test(t *testing.T, impl Dgeqr2er) {
	rnd := rand.New(rand.NewSource(1))
	for c, test := range []struct {
		m, n, lda int
	}{
		{1, 1, 0},
		{2, 2, 0},
		{3, 2, 0},
		{2, 3, 0},
		{1, 12, 0},
		{2, 6, 0},
		{3, 4, 0},
		{4, 3, 0},
		{6, 2, 0},
		{12, 1, 0},
		{1, 1, 20},
		{2, 2, 20},
		{3, 2, 20},
		{2, 3, 20},
		{1, 12, 20},
		{2, 6, 20},
		{3, 4, 20},
		{4, 3, 20},
		{6, 2, 20},
		{12, 1, 20},
	} {
		n := test.n
		m := test.m
		lda := test.lda
		if lda == 0 {
			lda = test.n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.Float64()
		}
		aCopy := make([]float64, len(a))
		k := min(m, n)
		tau := make([]float64, k)
		for i := range tau {
			tau[i] = rnd.Float64()
		}
		work := make([]float64, n)
		for i := range work {
			work[i] = rnd.Float64()
		}
		copy(aCopy, a)
		impl.Dgeqr2(m, n, a, lda, tau, work)

		// Test that the QR factorization has completed successfully. Compute
		// Q based on the vectors.
		q := constructQ("QR", m, n, a, lda, tau)

		// Check that q is orthonormal
		for i := 0; i < m; i++ {
			nrm := blas64.Nrm2(m, blas64.Vector{Inc: 1, Data: q.Data[i*m:]})
			if math.Abs(nrm-1) > 1e-14 {
				t.Errorf("Case %v, q not normal", c)
			}
			for j := 0; j < i; j++ {
				dot := blas64.Dot(m, blas64.Vector{Inc: 1, Data: q.Data[i*m:]}, blas64.Vector{Inc: 1, Data: q.Data[j*m:]})
				if math.Abs(dot) > 1e-14 {
					t.Errorf("Case %v, q not orthogonal", i)
				}
			}
		}
		// Check that A = Q * R
		r := blas64.General{
			Rows:   m,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, m*n),
		}
		for i := 0; i < m; i++ {
			for j := i; j < n; j++ {
				r.Data[i*n+j] = a[i*lda+j]
			}
		}
		atmp := blas64.General{
			Rows:   m,
			Cols:   n,
			Stride: lda,
			Data:   make([]float64, m*lda),
		}
		copy(atmp.Data, a)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, r, 0, atmp)
		if !floats.EqualApprox(atmp.Data, aCopy, 1e-14) {
			t.Errorf("Q*R != a")
		}
	}
}
