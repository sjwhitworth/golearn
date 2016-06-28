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
	"github.com/gonum/lapack"
)

type Dlarfter interface {
	Dgeqr2er
	Dlarft(direct lapack.Direct, store lapack.StoreV, n, k int, v []float64, ldv int, tau []float64, t []float64, ldt int)
}

func DlarftTest(t *testing.T, impl Dlarfter) {
	rnd := rand.New(rand.NewSource(1))
	for _, store := range []lapack.StoreV{lapack.ColumnWise, lapack.RowWise} {
		for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
			for _, test := range []struct {
				m, n, ldv, ldt int
			}{
				{6, 6, 0, 0},
				{8, 6, 0, 0},
				{6, 8, 0, 0},
				{6, 6, 10, 15},
				{8, 6, 10, 15},
				{6, 8, 10, 15},
				{6, 6, 15, 10},
				{8, 6, 15, 10},
				{6, 8, 15, 10},
			} {
				// Generate a matrix
				m := test.m
				n := test.n
				lda := n
				if lda == 0 {
					lda = n
				}

				a := make([]float64, m*lda)
				for i := 0; i < m; i++ {
					for j := 0; j < lda; j++ {
						a[i*lda+j] = rnd.Float64()
					}
				}
				// Use dgeqr2 to find the v vectors
				tau := make([]float64, n)
				work := make([]float64, n)
				impl.Dgeqr2(m, n, a, lda, tau, work)

				// Construct H using these answers
				vMatTmp := extractVMat(m, n, a, lda, lapack.Forward, lapack.ColumnWise)
				vMat := constructVMat(vMatTmp, store, direct)
				v := vMat.Data
				ldv := vMat.Stride

				h := constructH(tau, vMat, store, direct)

				k := min(m, n)
				ldt := test.ldt
				if ldt == 0 {
					ldt = k
				}
				// Find T from the actual function
				tm := make([]float64, k*ldt)
				for i := range tm {
					tm[i] = 100 + rnd.Float64()
				}
				// The v data has been put into a.
				impl.Dlarft(direct, store, m, k, v, ldv, tau, tm, ldt)

				tData := make([]float64, len(tm))
				copy(tData, tm)
				if direct == lapack.Forward {
					// Zero out the lower traingular portion.
					for i := 0; i < k; i++ {
						for j := 0; j < i; j++ {
							tData[i*ldt+j] = 0
						}
					}
				} else {
					// Zero out the upper traingular portion.
					for i := 0; i < k; i++ {
						for j := i + 1; j < k; j++ {
							tData[i*ldt+j] = 0
						}
					}
				}

				T := blas64.General{
					Rows:   k,
					Cols:   k,
					Stride: ldt,
					Data:   tData,
				}

				vMatT := blas64.General{
					Rows:   vMat.Cols,
					Cols:   vMat.Rows,
					Stride: vMat.Rows,
					Data:   make([]float64, vMat.Cols*vMat.Rows),
				}
				for i := 0; i < vMat.Rows; i++ {
					for j := 0; j < vMat.Cols; j++ {
						vMatT.Data[j*vMatT.Stride+i] = vMat.Data[i*vMat.Stride+j]
					}
				}
				var comp blas64.General
				if store == lapack.ColumnWise {
					// H = I - V * T * V^T
					tmp := blas64.General{
						Rows:   T.Rows,
						Cols:   vMatT.Cols,
						Stride: vMatT.Cols,
						Data:   make([]float64, T.Rows*vMatT.Cols),
					}
					// T * V^T
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, T, vMatT, 0, tmp)
					comp = blas64.General{
						Rows:   vMat.Rows,
						Cols:   tmp.Cols,
						Stride: tmp.Cols,
						Data:   make([]float64, vMat.Rows*tmp.Cols),
					}
					// V * (T * V^T)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, vMat, tmp, 0, comp)
				} else {
					// H = I - V^T * T * V
					tmp := blas64.General{
						Rows:   T.Rows,
						Cols:   vMat.Cols,
						Stride: vMat.Cols,
						Data:   make([]float64, T.Rows*vMat.Cols),
					}
					// T * V
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, T, vMat, 0, tmp)
					comp = blas64.General{
						Rows:   vMatT.Rows,
						Cols:   tmp.Cols,
						Stride: tmp.Cols,
						Data:   make([]float64, vMatT.Rows*tmp.Cols),
					}
					// V^T * (T * V)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, vMatT, tmp, 0, comp)
				}
				// I - V^T * T * V
				for i := 0; i < comp.Rows; i++ {
					for j := 0; j < comp.Cols; j++ {
						comp.Data[i*m+j] *= -1
						if i == j {
							comp.Data[i*m+j] += 1
						}
					}
				}
				if !floats.EqualApprox(comp.Data, h.Data, 1e-14) {
					t.Errorf("T does not construct proper H. Store = %v, Direct = %v.\nWant %v\ngot %v.", string(store), string(direct), h.Data, comp.Data)
				}
			}
		}
	}
}
