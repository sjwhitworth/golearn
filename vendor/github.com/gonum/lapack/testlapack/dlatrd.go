// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dlatrder interface {
	Dlatrd(uplo blas.Uplo, n, nb int, a []float64, lda int, e, tau, w []float64, ldw int)
}

func DlatrdTest(t *testing.T, impl Dlatrder) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, test := range []struct {
			n, nb, lda, ldw int
		}{
			{5, 2, 0, 0},
			{5, 5, 0, 0},

			{5, 3, 10, 11},
			{5, 5, 10, 11},
		} {
			n := test.n
			nb := test.nb
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			ldw := test.ldw
			if ldw == 0 {
				ldw = nb
			}

			a := make([]float64, n*lda)
			for i := range a {
				a[i] = rnd.NormFloat64()
			}

			e := make([]float64, n-1)
			for i := range e {
				e[i] = math.NaN()
			}
			tau := make([]float64, n-1)
			for i := range tau {
				tau[i] = math.NaN()
			}
			w := make([]float64, n*ldw)
			for i := range w {
				w[i] = math.NaN()
			}

			aCopy := make([]float64, len(a))
			copy(aCopy, a)

			impl.Dlatrd(uplo, n, nb, a, lda, e, tau, w, ldw)

			// Construct Q.
			ldq := n
			q := blas64.General{
				Rows:   n,
				Cols:   n,
				Stride: ldq,
				Data:   make([]float64, n*ldq),
			}
			for i := 0; i < n; i++ {
				q.Data[i*ldq+i] = 1
			}
			if uplo == blas.Upper {
				for i := n - 1; i >= n-nb; i-- {
					if i == 0 {
						continue
					}
					h := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					for j := 0; j < n; j++ {
						h.Data[j*n+j] = 1
					}
					v := blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					for j := 0; j < i-1; j++ {
						v.Data[j] = a[j*lda+i]
					}
					v.Data[i-1] = 1

					blas64.Ger(-tau[i-1], v, v, h)

					qTmp := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					copy(qTmp.Data, q.Data)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qTmp, h, 0, q)
				}
			} else {
				for i := 0; i < nb; i++ {
					if i == n-1 {
						continue
					}
					h := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					for j := 0; j < n; j++ {
						h.Data[j*n+j] = 1
					}
					v := blas64.Vector{
						Inc:  1,
						Data: make([]float64, n),
					}
					v.Data[i+1] = 1
					for j := i + 2; j < n; j++ {
						v.Data[j] = a[j*lda+i]
					}
					blas64.Ger(-tau[i], v, v, h)

					qTmp := blas64.General{
						Rows: n, Cols: n, Stride: n, Data: make([]float64, n*n),
					}
					copy(qTmp.Data, q.Data)
					blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qTmp, h, 0, q)
				}
			}
			errStr := fmt.Sprintf("isUpper = %v, n = %v, nb = %v", uplo == blas.Upper, n, nb)
			if !isOrthonormal(q) {
				t.Errorf("Q not orthonormal. %s", errStr)
			}
			aGen := genFromSym(blas64.Symmetric{N: n, Stride: lda, Uplo: uplo, Data: aCopy})
			if !dlatrdCheckDecomposition(t, uplo, n, nb, e, tau, a, lda, aGen, q) {
				t.Errorf("Decomposition mismatch. %s", errStr)
			}
		}
	}
}

// dlatrdCheckDecomposition checks that the first nb rows have been successfully
// reduced.
func dlatrdCheckDecomposition(t *testing.T, uplo blas.Uplo, n, nb int, e, tau, a []float64, lda int, aGen, q blas64.General) bool {
	// Compute Q^T * A * Q.
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

	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, aGen, 0, tmp)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, q, 0, ans)

	// Compare with T.
	if uplo == blas.Upper {
		for i := n - 1; i >= n-nb; i-- {
			for j := 0; j < n; j++ {
				v := ans.Data[i*ans.Stride+j]
				switch {
				case i == j:
					if math.Abs(v-a[i*lda+j]) > 1e-10 {
						return false
					}
				case i == j-1:
					if math.Abs(a[i*lda+j]-1) > 1e-10 {
						return false
					}
					if math.Abs(v-e[i]) > 1e-10 {
						return false
					}
				case i == j+1:
				default:
					if math.Abs(v) > 1e-10 {
						return false
					}
				}
			}
		}
	} else {
		for i := 0; i < nb; i++ {
			for j := 0; j < n; j++ {
				v := ans.Data[i*ans.Stride+j]
				switch {
				case i == j:
					if math.Abs(v-a[i*lda+j]) > 1e-10 {
						return false
					}
				case i == j-1:
				case i == j+1:
					if math.Abs(a[i*lda+j]-1) > 1e-10 {
						return false
					}
					if math.Abs(v-e[i-1]) > 1e-10 {
						return false
					}
				default:
					if math.Abs(v) > 1e-10 {
						return false
					}
				}
			}
		}
	}
	return true
}

// genFromSym constructs a (symmetric) general matrix from the data in the
// symmetric.
// TODO(btracey): Replace other constructions of this with a call to this function.
func genFromSym(a blas64.Symmetric) blas64.General {
	n := a.N
	lda := a.Stride
	uplo := a.Uplo
	b := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			v := a.Data[i*lda+j]
			if uplo == blas.Lower {
				v = a.Data[j*lda+i]
			}
			b.Data[i*n+j] = v
			b.Data[j*n+i] = v
		}
	}
	return b
}
