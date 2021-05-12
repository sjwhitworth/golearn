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

type Dlarfber interface {
	Dlarfter
	Dlarfb(side blas.Side, trans blas.Transpose, direct lapack.Direct,
		store lapack.StoreV, m, n, k int, v []float64, ldv int, t []float64, ldt int,
		c []float64, ldc int, work []float64, ldwork int)
}

func DlarfbTest(t *testing.T, impl Dlarfber) {
	rnd := rand.New(rand.NewSource(1))
	for _, store := range []lapack.StoreV{lapack.ColumnWise, lapack.RowWise} {
		for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
			for _, side := range []blas.Side{blas.Left, blas.Right} {
				for _, trans := range []blas.Transpose{blas.Trans, blas.NoTrans} {
					for cas, test := range []struct {
						ma, na, cdim, lda, ldt, ldc int
					}{
						{6, 6, 6, 0, 0, 0},
						{6, 8, 10, 0, 0, 0},
						{6, 10, 8, 0, 0, 0},
						{8, 6, 10, 0, 0, 0},
						{8, 10, 6, 0, 0, 0},
						{10, 6, 8, 0, 0, 0},
						{10, 8, 6, 0, 0, 0},
						{6, 6, 6, 12, 15, 30},
						{6, 8, 10, 12, 15, 30},
						{6, 10, 8, 12, 15, 30},
						{8, 6, 10, 12, 15, 30},
						{8, 10, 6, 12, 15, 30},
						{10, 6, 8, 12, 15, 30},
						{10, 8, 6, 12, 15, 30},
						{6, 6, 6, 15, 12, 30},
						{6, 8, 10, 15, 12, 30},
						{6, 10, 8, 15, 12, 30},
						{8, 6, 10, 15, 12, 30},
						{8, 10, 6, 15, 12, 30},
						{10, 6, 8, 15, 12, 30},
						{10, 8, 6, 15, 12, 30},
					} {
						// Generate a matrix for QR
						ma := test.ma
						na := test.na
						lda := test.lda
						if lda == 0 {
							lda = na
						}
						a := make([]float64, ma*lda)
						for i := 0; i < ma; i++ {
							for j := 0; j < lda; j++ {
								a[i*lda+j] = rnd.Float64()
							}
						}
						k := min(ma, na)

						// H is always ma x ma
						var m, n, rowsWork int
						switch {
						default:
							panic("not implemented")
						case side == blas.Left:
							m = test.ma
							n = test.cdim
							rowsWork = n
						case side == blas.Right:
							m = test.cdim
							n = test.ma
							rowsWork = m
						}

						// Use dgeqr2 to find the v vectors
						tau := make([]float64, na)
						work := make([]float64, na)
						impl.Dgeqr2(ma, k, a, lda, tau, work)

						// Correct the v vectors based on the direct and store
						vMatTmp := extractVMat(ma, na, a, lda, lapack.Forward, lapack.ColumnWise)
						vMat := constructVMat(vMatTmp, store, direct)
						v := vMat.Data
						ldv := vMat.Stride

						// Use dlarft to find the t vector
						ldt := test.ldt
						if ldt == 0 {
							ldt = k
						}
						tm := make([]float64, k*ldt)

						impl.Dlarft(direct, store, ma, k, v, ldv, tau, tm, ldt)

						// Generate c matrix
						ldc := test.ldc
						if ldc == 0 {
							ldc = n
						}
						c := make([]float64, m*ldc)
						for i := 0; i < m; i++ {
							for j := 0; j < ldc; j++ {
								c[i*ldc+j] = rnd.Float64()
							}
						}
						cCopy := make([]float64, len(c))
						copy(cCopy, c)

						ldwork := k
						work = make([]float64, rowsWork*k)

						// Call Dlarfb with this information
						impl.Dlarfb(side, trans, direct, store, m, n, k, v, ldv, tm, ldt, c, ldc, work, ldwork)

						h := constructH(tau, vMat, store, direct)

						cMat := blas64.General{
							Rows:   m,
							Cols:   n,
							Stride: ldc,
							Data:   make([]float64, m*ldc),
						}
						copy(cMat.Data, cCopy)
						ans := blas64.General{
							Rows:   m,
							Cols:   n,
							Stride: ldc,
							Data:   make([]float64, m*ldc),
						}
						copy(ans.Data, cMat.Data)
						switch {
						default:
							panic("not implemented")
						case side == blas.Left && trans == blas.NoTrans:
							blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, h, cMat, 0, ans)
						case side == blas.Left && trans == blas.Trans:
							blas64.Gemm(blas.Trans, blas.NoTrans, 1, h, cMat, 0, ans)
						case side == blas.Right && trans == blas.NoTrans:
							blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, cMat, h, 0, ans)
						case side == blas.Right && trans == blas.Trans:
							blas64.Gemm(blas.NoTrans, blas.Trans, 1, cMat, h, 0, ans)
						}
						if !floats.EqualApprox(ans.Data, c, 1e-14) {
							t.Errorf("Cas %v mismatch. Want %v, got %v.", cas, ans.Data, c)
						}
					}
				}
			}
		}
	}
}
