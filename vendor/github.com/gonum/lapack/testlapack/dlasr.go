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
	"github.com/gonum/lapack"
)

type Dlasrer interface {
	Dlasr(side blas.Side, pivot lapack.Pivot, direct lapack.Direct, m, n int, c, s, a []float64, lda int)
}

func DlasrTest(t *testing.T, impl Dlasrer) {
	rnd := rand.New(rand.NewSource(1))
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, pivot := range []lapack.Pivot{lapack.Variable, lapack.Top, lapack.Bottom} {
			for _, direct := range []lapack.Direct{lapack.Forward, lapack.Backward} {
				for _, test := range []struct {
					m, n, lda int
				}{
					{5, 5, 0},
					{5, 10, 0},
					{10, 5, 0},

					{5, 5, 20},
					{5, 10, 20},
					{10, 5, 20},
				} {
					m := test.m
					n := test.n
					lda := test.lda
					if lda == 0 {
						lda = n
					}
					a := make([]float64, m*lda)
					for i := range a {
						a[i] = rnd.Float64()
					}
					var s, c []float64
					if side == blas.Left {
						s = make([]float64, m-1)
						c = make([]float64, m-1)
					} else {
						s = make([]float64, n-1)
						c = make([]float64, n-1)
					}
					for k := range s {
						theta := rnd.Float64() * 2 * math.Pi
						s[k] = math.Sin(theta)
						c[k] = math.Cos(theta)
					}
					aCopy := make([]float64, len(a))
					copy(a, aCopy)
					impl.Dlasr(side, pivot, direct, m, n, c, s, a, lda)

					pSize := m
					if side == blas.Right {
						pSize = n
					}
					p := blas64.General{
						Rows:   pSize,
						Cols:   pSize,
						Stride: pSize,
						Data:   make([]float64, pSize*pSize),
					}
					pk := blas64.General{
						Rows:   pSize,
						Cols:   pSize,
						Stride: pSize,
						Data:   make([]float64, pSize*pSize),
					}
					ptmp := blas64.General{
						Rows:   pSize,
						Cols:   pSize,
						Stride: pSize,
						Data:   make([]float64, pSize*pSize),
					}
					for i := 0; i < pSize; i++ {
						p.Data[i*p.Stride+i] = 1
						ptmp.Data[i*p.Stride+i] = 1
					}
					// Compare to direct computation.
					for k := range s {
						for i := range p.Data {
							pk.Data[i] = 0
						}
						for i := 0; i < pSize; i++ {
							pk.Data[i*p.Stride+i] = 1
						}
						if pivot == lapack.Variable {
							pk.Data[k*p.Stride+k] = c[k]
							pk.Data[k*p.Stride+k+1] = s[k]
							pk.Data[(k+1)*p.Stride+k] = -s[k]
							pk.Data[(k+1)*p.Stride+k+1] = c[k]
						} else if pivot == lapack.Top {
							pk.Data[0] = c[k]
							pk.Data[k+1] = s[k]
							pk.Data[(k+1)*p.Stride] = -s[k]
							pk.Data[(k+1)*p.Stride+k+1] = c[k]
						} else {
							pk.Data[(pSize-1-k)*p.Stride+pSize-k-1] = c[k]
							pk.Data[(pSize-1-k)*p.Stride+pSize-1] = s[k]
							pk.Data[(pSize-1)*p.Stride+pSize-1-k] = -s[k]
							pk.Data[(pSize-1)*p.Stride+pSize-1] = c[k]
						}
						if direct == lapack.Forward {
							blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, pk, ptmp, 0, p)
						} else {
							blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, ptmp, pk, 0, p)
						}
						copy(ptmp.Data, p.Data)
					}

					aMat := blas64.General{
						Rows:   m,
						Cols:   n,
						Stride: lda,
						Data:   make([]float64, m*lda),
					}
					copy(a, aCopy)
					newA := blas64.General{
						Rows:   m,
						Cols:   n,
						Stride: lda,
						Data:   make([]float64, m*lda),
					}
					if side == blas.Left {
						blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, p, aMat, 0, newA)
					} else {
						blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aMat, p, 0, newA)
					}
					if !floats.EqualApprox(newA.Data, a, 1e-12) {
						t.Errorf("A update mismatch")
					}
				}
			}
		}
	}
}
