// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dsytd2er interface {
	Dsytd2(uplo blas.Uplo, n int, a []float64, lda int, d, e, tau []float64)
}

func Dsytd2Test(t *testing.T, impl Dsytd2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, test := range []struct {
			n, lda int
		}{
			{3, 0},
			{4, 0},
			{5, 0},

			{3, 10},
			{4, 10},
			{5, 10},
		} {
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

			d := make([]float64, n)
			for i := range d {
				d[i] = math.NaN()
			}
			e := make([]float64, n-1)
			for i := range e {
				e[i] = math.NaN()
			}
			tau := make([]float64, n-1)
			for i := range tau {
				tau[i] = math.NaN()
			}

			impl.Dsytd2(uplo, n, a, lda, d, e, tau)

			// Construct Q
			qMat := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, n*n),
			}
			qCopy := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, len(qMat.Data)),
			}
			// Set Q to I.
			for i := 0; i < n; i++ {
				qMat.Data[i*qMat.Stride+i] = 1
			}
			for i := 0; i < n-1; i++ {
				hMat := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, n*n),
				}
				// Set H to I.
				for i := 0; i < n; i++ {
					hMat.Data[i*hMat.Stride+i] = 1
				}
				var vi blas64.Vector
				if uplo == blas.Upper {
					vi = blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					for j := 0; j < i; j++ {
						vi.Data[j] = a[j*lda+i+1]
					}
					vi.Data[i] = 1
				} else {
					vi = blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					vi.Data[i+1] = 1
					for j := i + 2; j < n; j++ {
						vi.Data[j] = a[j*lda+i]
					}
				}
				blas64.Ger(-tau[i], vi, vi, hMat)
				copy(qCopy.Data, qMat.Data)

				// Multiply q by the new h.
				if uplo == blas.Upper {
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hMat, qCopy, 0, qMat)
				} else {
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, hMat, 0, qMat)
				}
			}
			// Check that Q is orthonormal
			othonormal := true
			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					dot := blas64.Dot(n,
						blas64.Vector{Inc: 1, Data: qMat.Data[i*qMat.Stride:]},
						blas64.Vector{Inc: 1, Data: qMat.Data[j*qMat.Stride:]},
					)
					if i == j {
						if math.Abs(dot-1) > 1e-10 {
							othonormal = false
						}
					} else {
						if math.Abs(dot) > 1e-10 {
							othonormal = false
						}
					}
				}
			}
			if !othonormal {
				t.Errorf("Q not orthonormal")
			}

			// Compute Q^T * A * Q.
			aMat := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, len(a)),
			}

			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					v := aCopy[i*lda+j]
					if uplo == blas.Lower {
						v = aCopy[j*lda+i]
					}
					aMat.Data[i*aMat.Stride+j] = v
					aMat.Data[j*aMat.Stride+i] = v
				}
			}

			tmp := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, n*n),
			}

			ans := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, n*n),
			}

			blas64.Gemm(blas.Trans, blas.NoTrans, 1, qMat, aMat, 0, tmp)
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, qMat, 0, ans)

			// Compare with T.
			tMat := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: n,
				Data:   make([]float64, n*n),
			}
			for i := 0; i < n-1; i++ {
				tMat.Data[i*tMat.Stride+i] = d[i]
				tMat.Data[i*tMat.Stride+i+1] = e[i]
				tMat.Data[(i+1)*tMat.Stride+i] = e[i]
			}
			tMat.Data[(n-1)*tMat.Stride+n-1] = d[n-1]

			same := true
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					if math.Abs(ans.Data[i*ans.Stride+j]-tMat.Data[i*tMat.Stride+j]) > 1e-10 {
						same = false
					}
				}
			}
			if !same {
				t.Errorf("Matrix answer mismatch")
			}
		}
	}
}
