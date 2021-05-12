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

type Dormhrer interface {
	Dormhr(side blas.Side, trans blas.Transpose, m, n, ilo, ihi int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int)

	Dgehrder
}

func DormhrTest(t *testing.T, impl Dormhrer) {
	rnd := rand.New(rand.NewSource(1))

	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
			for _, m := range []int{1, 2, 3, 4, 5, 8, 9, 10, 23} {
				for _, n := range []int{1, 2, 3, 4, 5, 8, 9, 10, 23} {
					for _, extra := range []int{0, 1, 13} {
						for cas := 0; cas < 10; cas++ {
							nq := m
							if side == blas.Right {
								nq = n
							}
							ilo := rnd.Intn(nq)
							ihi := rnd.Intn(nq)
							if ilo > ihi {
								ilo, ihi = ihi, ilo
							}
							testDormhr(t, impl, side, trans, m, n, ilo, ihi, extra, true, rnd)
							testDormhr(t, impl, side, trans, m, n, ilo, ihi, extra, false, rnd)
						}
					}
				}
			}
		}
	}
	for _, side := range []blas.Side{blas.Left, blas.Right} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans} {
			testDormhr(t, impl, side, trans, 0, 0, 0, -1, 0, true, rnd)
			testDormhr(t, impl, side, trans, 0, 0, 0, -1, 0, false, rnd)
		}
	}
}

func testDormhr(t *testing.T, impl Dormhrer, side blas.Side, trans blas.Transpose, m, n, ilo, ihi, extra int, optwork bool, rnd *rand.Rand) {
	const tol = 1e-14

	var nq, nw int
	switch side {
	case blas.Left:
		nq = m
		nw = n
	case blas.Right:
		nq = n
		nw = m
	}

	// Compute the elementary reflectors and tau.
	a := randomGeneral(nq, nq, nq+extra, rnd)
	var tau []float64
	if nq > 1 {
		tau = nanSlice(nq - 1)
	}
	work := nanSlice(max(1, nq)) // Minimum work for Dgehrd.
	impl.Dgehrd(nq, ilo, ihi, a.Data, a.Stride, tau, work, len(work))

	// Construct Q from the elementary reflectors in a and from tau.
	q := eye(nq, nq)
	qCopy := eye(nq, nq)
	for j := ilo; j < ihi; j++ {
		h := eye(nq, nq)
		v := blas64.Vector{
			Inc:  1,
			Data: make([]float64, nq),
		}
		v.Data[j+1] = 1
		for i := j + 2; i < ihi+1; i++ {
			v.Data[i] = a.Data[i*a.Stride+j]
		}
		blas64.Ger(-tau[j], v, v, h)
		copy(qCopy.Data, q.Data)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, h, 0, q)
	}

	c := randomGeneral(m, n, n+extra, rnd)

	// Compute the product of Q and C explicitly.
	qc := randomGeneral(m, n, n+extra, rnd)
	if side == blas.Left {
		blas64.Gemm(trans, blas.NoTrans, 1, q, c, 0, qc)
	} else {
		blas64.Gemm(blas.NoTrans, trans, 1, c, q, 0, qc)
	}

	// Compute the product of Q and C using Dormhr.
	if optwork {
		work = nanSlice(1)
		impl.Dormhr(side, trans, m, n, ilo, ihi, a.Data, a.Stride, tau, c.Data, c.Stride, work, -1)
		work = nanSlice(int(work[0]))
	} else {
		work = nanSlice(max(1, nw))
	}
	impl.Dormhr(side, trans, m, n, ilo, ihi, a.Data, a.Stride, tau, c.Data, c.Stride, work, len(work))

	// Compare the two answers.
	prefix := fmt.Sprintf("Case side=%v, trans=%v, m=%v, n=%v, ilo=%v, ihi=%v, extra=%v, optwork=%v",
		side, trans, m, n, ilo, ihi, extra, optwork)
	if !generalOutsideAllNaN(c) {
		t.Errorf("%v: out-of-range write to C\n%v", prefix, c.Data)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			cij := c.Data[i*c.Stride+j]
			qcij := qc.Data[i*qc.Stride+j]
			if math.Abs(cij-qcij) > tol {
				t.Errorf("%v: unexpected value of the QC product at [%v,%v]: want %v, got %v", prefix, i, j, qcij, cij)
			}
		}
	}
}
