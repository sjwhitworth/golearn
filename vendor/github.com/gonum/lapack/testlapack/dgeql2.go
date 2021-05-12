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
)

type Dgeql2er interface {
	Dgeql2(m, n int, a []float64, lda int, tau, work []float64)
}

func Dgeql2Test(t *testing.T, impl Dgeql2er) {
	rnd := rand.New(rand.NewSource(1))
	// TODO(btracey): Add tests for m < n.
	for _, test := range []struct {
		m, n, lda int
	}{
		{5, 5, 0},
		{5, 3, 0},
		{5, 4, 0},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}
		tau := nanSlice(min(m, n))
		work := nanSlice(n)

		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		impl.Dgeql2(m, n, a, lda, tau, work)

		k := min(m, n)
		// Construct Q.
		q := blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, m*m),
		}
		for i := 0; i < m; i++ {
			q.Data[i*q.Stride+i] = 1
		}
		for i := 0; i < k; i++ {
			h := blas64.General{Rows: m, Cols: m, Stride: m, Data: make([]float64, m*m)}
			for j := 0; j < m; j++ {
				h.Data[j*h.Stride+j] = 1
			}
			v := blas64.Vector{Inc: 1, Data: make([]float64, m)}
			v.Data[m-k+i] = 1
			for j := 0; j < m-k+i; j++ {
				v.Data[j] = a[j*lda+n-k+i]
			}
			blas64.Ger(-tau[i], v, v, h)
			qTmp := blas64.General{Rows: q.Rows, Cols: q.Cols, Stride: q.Stride, Data: make([]float64, len(q.Data))}
			copy(qTmp.Data, q.Data)
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, h, qTmp, 0, q)
		}
		if !isOrthonormal(q) {
			t.Errorf("Q is not orthonormal")
		}
		l := blas64.General{
			Rows:   m,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, m*n),
		}
		if m >= n {
			for i := m - n; i < m; i++ {
				for j := 0; j <= min(i-(m-n), n-1); j++ {
					l.Data[i*l.Stride+j] = a[i*lda+j]
				}
			}
		} else {
			panic("untested")
		}
		ans := blas64.General{Rows: m, Cols: n, Stride: lda, Data: make([]float64, len(a))}
		copy(ans.Data, a)

		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, l, 0, ans)
		if !floats.EqualApprox(ans.Data, aCopy, 1e-10) {
			t.Errorf("Reconstruction mismatch: m = %v, n = %v", m, n)
		}
	}
}
