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

type Dorghrer interface {
	Dorghr(n, ilo, ihi int, a []float64, lda int, tau, work []float64, lwork int)

	Dgehrder
}

func DorghrTest(t *testing.T, impl Dorghrer) {
	rnd := rand.New(rand.NewSource(1))

	for _, n := range []int{1, 2, 3, 4, 5, 6, 7, 8, 23, 34} {
		for _, extra := range []int{0, 1, 13} {
			for _, optwork := range []bool{true, false} {
				for cas := 0; cas < 100; cas++ {
					ilo := rnd.Intn(n)
					ihi := rnd.Intn(n)
					if ilo > ihi {
						ilo, ihi = ihi, ilo
					}
					testDorghr(t, impl, n, ilo, ihi, extra, optwork, rnd)
				}
			}
		}
	}
	testDorghr(t, impl, 0, 0, -1, 0, false, rnd)
	testDorghr(t, impl, 0, 0, -1, 0, true, rnd)
}

func testDorghr(t *testing.T, impl Dorghrer, n, ilo, ihi, extra int, optwork bool, rnd *rand.Rand) {
	const tol = 1e-14

	// Construct the matrix A with elementary reflectors and scalar factors tau.
	a := randomGeneral(n, n, n+extra, rnd)
	var tau []float64
	if n > 1 {
		tau = nanSlice(n - 1)
	}
	work := nanSlice(max(1, n)) // Minimum work for Dgehrd.
	impl.Dgehrd(n, ilo, ihi, a.Data, a.Stride, tau, work, len(work))

	// Extract Q for later comparison.
	q := eye(n, n)
	qCopy := cloneGeneral(q)
	for j := ilo; j < ihi; j++ {
		h := eye(n, n)
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

	if optwork {
		work = nanSlice(1)
		impl.Dorghr(n, ilo, ihi, a.Data, a.Stride, tau, work, -1)
		work = nanSlice(int(work[0]))
	} else {
		work = nanSlice(max(1, ihi-ilo))
	}
	impl.Dorghr(n, ilo, ihi, a.Data, a.Stride, tau, work, len(work))

	prefix := fmt.Sprintf("Case n=%v, ilo=%v, ihi=%v, extra=%v, optwork=%v", n, ilo, ihi, extra, optwork)
	if !generalOutsideAllNaN(a) {
		t.Errorf("%v: out-of-range write to A\n%v", prefix, a.Data)
	}
	if !isOrthonormal(a) {
		t.Errorf("%v: A is not orthogonal\n%v", prefix, a.Data)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			aij := a.Data[i*a.Stride+j]
			qij := q.Data[i*q.Stride+j]
			if math.Abs(aij-qij) > tol {
				t.Errorf("%v: unexpected value of A[%v,%v]. want %v, got %v", prefix, i, j, qij, aij)
			}
		}
	}
}
