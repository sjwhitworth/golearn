// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dlarfxer interface {
	Dlarfx(side blas.Side, m, n int, v []float64, tau float64, c []float64, ldc int, work []float64)
}

func DlarfxTest(t *testing.T, impl Dlarfxer) {
	rnd := rand.New(rand.NewSource(1))
	for _, side := range []blas.Side{blas.Right, blas.Left} {
		for m := 1; m < 12; m++ {
			for n := 1; n < 12; n++ {
				for _, extra := range []int{0, 1, 11} {
					for cas := 0; cas < 10; cas++ {
						testDlarfx(t, impl, side, m, n, extra, rnd)
					}
				}
			}
		}
	}
}

func testDlarfx(t *testing.T, impl Dlarfxer, side blas.Side, m, n, extra int, rnd *rand.Rand) {
	const tol = 1e-13

	c := randomGeneral(m, n, n+extra, rnd)
	cWant := randomGeneral(m, n, n+extra, rnd)
	tau := rnd.NormFloat64()

	var (
		v []float64
		h blas64.General
	)
	if side == blas.Left {
		v = randomSlice(m, rnd)
		h = eye(m, m+extra)
	} else {
		v = randomSlice(n, rnd)
		h = eye(n, n+extra)
	}
	blas64.Ger(-tau, blas64.Vector{Inc: 1, Data: v}, blas64.Vector{Inc: 1, Data: v}, h)
	if side == blas.Left {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, h, c, 0, cWant)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, c, h, 0, cWant)
	}

	var work []float64
	if h.Rows > 10 {
		// Allocate work only if H has order > 10.
		if side == blas.Left {
			work = make([]float64, n)
		} else {
			work = make([]float64, m)
		}
	}

	impl.Dlarfx(side, m, n, v, tau, c.Data, c.Stride, work)

	prefix := fmt.Sprintf("Case side=%v, m=%v, n=%v, extra=%v", side, m, n, extra)

	// Check any invalid modifications of c.
	if !generalOutsideAllNaN(c) {
		t.Errorf("%v: out-of-range write to C\n%v", prefix, c.Data)
	}

	if !equalApproxGeneral(c, cWant, tol) {
		t.Errorf("%v: unexpected C\n%v", prefix, c.Data)
	}
}
