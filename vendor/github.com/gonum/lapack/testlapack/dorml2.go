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

type Dorml2er interface {
	Dgelqfer
	Dorml2(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64)
}

func Dorml2Test(t *testing.T, impl Dorml2er) {
	rnd := rand.New(rand.NewSource(1))
	// TODO(btracey): This test is not complete, because it
	// doesn't test individual values of m, n, and k, instead only testing
	// a specific subset of possible k values.
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
			for _, test := range []struct {
				common, adim, cdim, lda, ldc int
			}{
				{3, 4, 5, 0, 0},
				{3, 5, 4, 0, 0},
				{4, 3, 5, 0, 0},
				{4, 5, 3, 0, 0},
				{5, 3, 4, 0, 0},
				{5, 4, 3, 0, 0},

				{3, 4, 5, 6, 20},
				{3, 5, 4, 6, 20},
				{4, 3, 5, 6, 20},
				{4, 5, 3, 6, 20},
				{5, 3, 4, 6, 20},
				{5, 4, 3, 6, 20},
				{3, 4, 5, 20, 6},
				{3, 5, 4, 20, 6},
				{4, 3, 5, 20, 6},
				{4, 5, 3, 20, 6},
				{5, 3, 4, 20, 6},
				{5, 4, 3, 20, 6},
			} {
				var ma, na, mc, nc int
				if side == blas.Left {
					ma = test.adim
					na = test.common
					mc = test.common
					nc = test.cdim
				} else {
					ma = test.adim
					na = test.common
					mc = test.cdim
					nc = test.common
				}
				// Generate a random matrix
				lda := test.lda
				if lda == 0 {
					lda = na
				}
				a := make([]float64, ma*lda)
				for i := range a {
					a[i] = rnd.Float64()
				}
				ldc := test.ldc
				if ldc == 0 {
					ldc = nc
				}
				// Compute random C matrix
				c := make([]float64, mc*ldc)
				for i := range c {
					c[i] = rnd.Float64()
				}

				// Compute LQ
				k := min(ma, na)
				tau := make([]float64, k)
				work := make([]float64, 1)
				impl.Dgelqf(ma, na, a, lda, tau, work, -1)
				work = make([]float64, int(work[0]))
				impl.Dgelqf(ma, na, a, lda, tau, work, len(work))

				// Build Q from result
				q := constructQ("LQ", ma, na, a, lda, tau)

				cMat := blas64.General{
					Rows:   mc,
					Cols:   nc,
					Stride: ldc,
					Data:   make([]float64, len(c)),
				}
				copy(cMat.Data, c)
				cMatCopy := blas64.General{
					Rows:   cMat.Rows,
					Cols:   cMat.Cols,
					Stride: cMat.Stride,
					Data:   make([]float64, len(cMat.Data)),
				}
				copy(cMatCopy.Data, cMat.Data)
				switch {
				default:
					panic("bad test")
				case side == blas.Left && trans == blas.NoTrans:
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, q, cMatCopy, 0, cMat)
				case side == blas.Left && trans == blas.Trans:
					blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, cMatCopy, 0, cMat)
				case side == blas.Right && trans == blas.NoTrans:
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, cMatCopy, q, 0, cMat)
				case side == blas.Right && trans == blas.Trans:
					blas64.Gemm(blas.NoTrans, blas.Trans, 1, cMatCopy, q, 0, cMat)
				}
				// Do Dorm2r ard compare
				if side == blas.Left {
					work = make([]float64, nc)
				} else {
					work = make([]float64, mc)
				}
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				tauCopy := make([]float64, len(tau))
				copy(tauCopy, tau)
				impl.Dorml2(side, trans, mc, nc, k, a, lda, tau, c, ldc, work)
				if !floats.Equal(a, aCopy) {
					t.Errorf("a changed in call")
				}
				if !floats.Equal(tau, tauCopy) {
					t.Errorf("tau changed in call")
				}
				if !floats.EqualApprox(cMat.Data, c, 1e-14) {
					isLeft := side == blas.Left
					isTrans := trans == blas.Trans
					t.Errorf("Multiplication mismatch. IsLeft = %v. IsTrans = %v", isLeft, isTrans)
				}
			}
		}
	}
}
