// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dlahr2er interface {
	Dlahr2(n, k, nb int, a []float64, lda int, tau, t []float64, ldt int, y []float64, ldy int)
}

type Dlahr2test struct {
	N, K, NB int
	A        []float64

	AWant   []float64
	TWant   []float64
	YWant   []float64
	TauWant []float64
}

func Dlahr2Test(t *testing.T, impl Dlahr2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		n, k, nb int
	}{
		{3, 0, 3},
		{3, 1, 2},
		{3, 1, 1},

		{5, 0, 5},
		{5, 1, 4},
		{5, 1, 3},
		{5, 1, 2},
		{5, 1, 1},
		{5, 2, 3},
		{5, 2, 2},
		{5, 2, 1},
		{5, 3, 2},
		{5, 3, 1},

		{7, 3, 4},
		{7, 3, 3},
		{7, 3, 2},
		{7, 3, 1},

		{10, 0, 10},
		{10, 1, 9},
		{10, 1, 5},
		{10, 1, 1},
		{10, 5, 5},
		{10, 5, 3},
		{10, 5, 1},
	} {
		for cas := 0; cas < 100; cas++ {
			for _, extraStride := range []int{0, 1, 10} {
				n := test.n
				k := test.k
				nb := test.nb

				a := randomGeneral(n, n-k+1, n-k+1+extraStride, rnd)
				aCopy := a
				aCopy.Data = make([]float64, len(a.Data))
				copy(aCopy.Data, a.Data)
				tmat := nanTriangular(blas.Upper, nb, nb+extraStride)
				y := nanGeneral(n, nb, nb+extraStride)
				tau := nanSlice(nb)

				impl.Dlahr2(n, k, nb, a.Data, a.Stride, tau, tmat.Data, tmat.Stride, y.Data, y.Stride)

				prefix := fmt.Sprintf("Case n=%v, k=%v, nb=%v, ldex=%v", n, k, nb, extraStride)

				if !generalOutsideAllNaN(a) {
					t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
				}
				if !triangularOutsideAllNaN(tmat) {
					t.Errorf("%v: out-of-range write to T\n%v", prefix, tmat.Data)
				}
				if !generalOutsideAllNaN(y) {
					t.Errorf("%v: out-of-range write to Y\n%v", prefix, y.Data)
				}

				// Check that A[:k,:] and A[:,nb:] blocks were not modified.
				for i := 0; i < n; i++ {
					for j := 0; j < n-k+1; j++ {
						if i >= k && j < nb {
							continue
						}
						if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
							t.Errorf("%v: unexpected write to A[%v,%v]", prefix, i, j)
						}
					}
				}

				// Check that all elements of tau were assigned.
				for i, v := range tau {
					if math.IsNaN(v) {
						t.Errorf("%v: tau[%v] not assigned", prefix, i)
					}
				}

				// Extract V from a.
				v := blas64.General{
					Rows:   n - k + 1,
					Cols:   nb,
					Stride: nb,
					Data:   make([]float64, (n-k+1)*nb),
				}
				for j := 0; j < v.Cols; j++ {
					v.Data[(j+1)*v.Stride+j] = 1
					for i := j + 2; i < v.Rows; i++ {
						v.Data[i*v.Stride+j] = a.Data[(i+k-1)*a.Stride+j]
					}
				}

				// VT = V.
				vt := v
				vt.Data = make([]float64, len(v.Data))
				copy(vt.Data, v.Data)
				// VT = V * T.
				blas64.Trmm(blas.Right, blas.NoTrans, 1, tmat, vt)
				// YWant = A * V * T.
				ywant := blas64.General{
					Rows:   n,
					Cols:   nb,
					Stride: nb,
					Data:   make([]float64, n*nb),
				}
				blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aCopy, vt, 0, ywant)

				// Compare Y and YWant.
				for i := 0; i < n; i++ {
					for j := 0; j < nb; j++ {
						diff := math.Abs(ywant.Data[i*ywant.Stride+j] - y.Data[i*y.Stride+j])
						if diff > 1e-14 {
							t.Errorf("%v: unexpected Y[%v,%v], diff=%v", prefix, i, j, diff)
						}
					}
				}

				// Construct Q directly from the first nb columns of a.
				q := constructQ("QR", n-k, nb, a.Data[k*a.Stride:], a.Stride, tau)
				if !isOrthonormal(q) {
					t.Errorf("%v: Q is not orthogonal", prefix)
				}
				// Construct Q as the product Q = I - V*T*V^T.
				qwant := blas64.General{
					Rows:   n - k + 1,
					Cols:   n - k + 1,
					Stride: n - k + 1,
					Data:   make([]float64, (n-k+1)*(n-k+1)),
				}
				for i := 0; i < qwant.Rows; i++ {
					qwant.Data[i*qwant.Stride+i] = 1
				}
				blas64.Gemm(blas.NoTrans, blas.Trans, -1, vt, v, 1, qwant)
				if !isOrthonormal(qwant) {
					t.Errorf("%v: Q = I - V*T*V^T is not orthogonal", prefix)
				}

				// Compare Q and QWant. Note that since Q is
				// (n-k)×(n-k) and QWant is (n-k+1)×(n-k+1), we
				// ignore the first row and column of QWant.
				for i := 0; i < n-k; i++ {
					for j := 0; j < n-k; j++ {
						diff := math.Abs(q.Data[i*q.Stride+j] - qwant.Data[(i+1)*qwant.Stride+j+1])
						if diff > 1e-14 {
							t.Errorf("%v: unexpected Q[%v,%v], diff=%v", prefix, i, j, diff)
						}
					}
				}
			}
		}
	}

	// Go runs tests from the source directory, so unfortunately we need to
	// include the "../testlapack" part.
	file, err := os.Open(filepath.FromSlash("../testlapack/testdata/dlahr2data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tests []Dlahr2test
	json.NewDecoder(r).Decode(&tests)
	for _, test := range tests {
		tau := make([]float64, len(test.TauWant))
		for _, ldex := range []int{0, 1, 20} {
			n := test.N
			k := test.K
			nb := test.NB

			lda := n - k + 1 + ldex
			a := make([]float64, (n-1)*lda+n-k+1)
			copyMatrix(n, n-k+1, a, lda, test.A)

			ldt := nb + ldex
			tmat := make([]float64, (nb-1)*ldt+nb)

			ldy := nb + ldex
			y := make([]float64, (n-1)*ldy+nb)

			impl.Dlahr2(n, k, nb, a, lda, tau, tmat, ldt, y, ldy)

			prefix := fmt.Sprintf("Case n=%v, k=%v, nb=%v, ldex=%v", n, k, nb, ldex)
			if !equalApprox(n, n-k+1, a, lda, test.AWant, 1e-14) {
				t.Errorf("%v: unexpected matrix A\n got=%v\nwant=%v", prefix, a, test.AWant)
			}
			if !equalApproxTriangular(true, nb, tmat, ldt, test.TWant, 1e-14) {
				t.Errorf("%v: unexpected matrix T\n got=%v\nwant=%v", prefix, tmat, test.TWant)
			}
			if !equalApprox(n, nb, y, ldy, test.YWant, 1e-14) {
				t.Errorf("%v: unexpected matrix Y\n got=%v\nwant=%v", prefix, y, test.YWant)
			}
			if !floats.EqualApprox(tau, test.TauWant, 1e-14) {
				t.Errorf("%v: unexpected slice tau\n got=%v\nwant=%v", prefix, tau, test.TauWant)
			}
		}
	}
}
