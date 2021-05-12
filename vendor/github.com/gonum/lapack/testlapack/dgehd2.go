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

type Dgehd2er interface {
	Dgehd2(n, ilo, ihi int, a []float64, lda int, tau, work []float64)
}

func Dgehd2Test(t *testing.T, impl Dgehd2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{1, 2, 3, 4, 5, 7, 10, 30} {
		for _, extra := range []int{0, 1, 13} {
			for cas := 0; cas < 100; cas++ {
				testDgehd2(t, impl, n, extra, rnd)
			}
		}
	}
}

func testDgehd2(t *testing.T, impl Dgehd2er, n, extra int, rnd *rand.Rand) {
	ilo := rnd.Intn(n)
	ihi := rnd.Intn(n)
	if ilo > ihi {
		ilo, ihi = ihi, ilo
	}

	tau := nanSlice(n - 1)
	work := nanSlice(n)

	a := randomGeneral(n, n, n+extra, rnd)
	// NaN out elements under the diagonal except
	// for the [ilo:ihi,ilo:ihi] block.
	for i := 1; i <= ihi; i++ {
		for j := 0; j < min(ilo, i); j++ {
			a.Data[i*a.Stride+j] = math.NaN()
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := 0; j < i; j++ {
			a.Data[i*a.Stride+j] = math.NaN()
		}
	}
	aCopy := a
	aCopy.Data = make([]float64, len(a.Data))
	copy(aCopy.Data, a.Data)

	impl.Dgehd2(n, ilo, ihi, a.Data, a.Stride, tau, work)

	prefix := fmt.Sprintf("Case n=%v, ilo=%v, ihi=%v, extra=%v", n, ilo, ihi, extra)

	// Check any invalid modifications of a.
	if !generalOutsideAllNaN(a) {
		t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
	}
	for i := ilo; i <= ihi; i++ {
		for j := 0; j < min(ilo, i); j++ {
			if !math.IsNaN(a.Data[i*a.Stride+j]) {
				t.Errorf("%v: expected NaN at A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if !math.IsNaN(a.Data[i*a.Stride+j]) {
				t.Errorf("%v: expected NaN at A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := 0; i <= ilo; i++ {
		for j := i; j < ilo+1; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
		for j := ihi + 1; j < n; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := i; j < n; j++ {
			if a.Data[i*a.Stride+j] != aCopy.Data[i*aCopy.Stride+j] {
				t.Errorf("%v: unexpected modification at A[%v,%v]", prefix, i, j)
			}
		}
	}

	// Check that tau has been assigned properly.
	for i, v := range tau {
		if i < ilo || i >= ihi {
			if !math.IsNaN(v) {
				t.Errorf("%v: expected NaN at tau[%v]", prefix, i)
			}
		} else {
			if math.IsNaN(v) {
				t.Errorf("%v: unexpected NaN at tau[%v]", prefix, i)
			}
		}
	}

	// Extract Q and check that it is orthogonal.
	q := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	for i := 0; i < q.Rows; i++ {
		q.Data[i*q.Stride+i] = 1
	}
	qCopy := q
	qCopy.Data = make([]float64, len(q.Data))
	for j := ilo; j < ihi; j++ {
		h := blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, n*n),
		}
		for i := 0; i < h.Rows; i++ {
			h.Data[i*h.Stride+i] = 1
		}
		v := blas64.Vector{
			Inc:  1,
			Data: make([]float64, n),
		}
		v.Data[j+1] = 1
		for i := j + 2; i < ihi+1; i++ {
			v.Data[i] = a.Data[i*a.Stride+j]
		}
		blas64.Ger(-tau[j], v, v, h)
		copy(qCopy.Data, q.Data)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, h, 0, q)
	}
	if !isOrthonormal(q) {
		t.Errorf("%v: Q is not orthogonal\nQ=%v", prefix, q)
	}

	// Overwrite NaN elements of aCopy with zeros
	// (we will multiply with it below).
	for i := 1; i <= ihi; i++ {
		for j := 0; j < min(ilo, i); j++ {
			aCopy.Data[i*aCopy.Stride+j] = 0
		}
	}
	for i := ihi + 1; i < n; i++ {
		for j := 0; j < i; j++ {
			aCopy.Data[i*aCopy.Stride+j] = 0
		}
	}

	// Construct Q^T * AOrig * Q and check that it is
	// equal to A from Dgehd2.
	aq := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, aCopy, q, 0, aq)
	qaq := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, q, aq, 0, qaq)
	for i := ilo; i <= ihi; i++ {
		for j := ilo; j <= ihi; j++ {
			qaqij := qaq.Data[i*qaq.Stride+j]
			if j < i-1 {
				if math.Abs(qaqij) > 1e-14 {
					t.Errorf("%v: Q^T*A*Q is not upper Hessenberg, [%v,%v]=%v", prefix, i, j, qaqij)
				}
				continue
			}
			diff := qaqij - a.Data[i*a.Stride+j]
			if math.Abs(diff) > 1e-14 {
				t.Errorf("%v: Q^T*AOrig*Q and A are not equal, diff at [%v,%v]=%v", prefix, i, j, diff)
			}
		}
	}
}
