// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dorglqer interface {
	Dorgl2er
	Dorglq(m, n, k int, a []float64, lda int, tau, work []float64, lwork int)
}

func DorglqTest(t *testing.T, impl Dorglqer) {
	rnd := rand.New(rand.NewSource(1))
	// TODO(btracey): Base tests off of nb and nx.
	for _, test := range []struct{ m, n, k, lda int }{
		{10, 10, 10, 0},
		{10, 10, 10, 20},
		{10, 30, 10, 0},
		{20, 30, 10, 0},

		{100, 100, 100, 0},
		{100, 100, 50, 0},
		{100, 130, 100, 0},
		{100, 130, 50, 0},
		{100, 100, 100, 150},
		{100, 100, 50, 150},
		{100, 130, 100, 150},
		{100, 130, 50, 150},

		{200, 200, 200, 0},
		{200, 200, 150, 0},
		{200, 230, 200, 0},
		{200, 230, 150, 0},
		{200, 200, 200, 250},
		{200, 200, 150, 250},
		{200, 230, 200, 250},
		{200, 230, 150, 250},
	} {
		m := test.m
		n := test.n
		k := test.k
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.Float64()
		}
		work := make([]float64, 1)
		tau := make([]float64, m)
		for i := range tau {
			tau[i] = math.NaN()
		}
		// Compute LQ factorization.
		impl.Dgelqf(m, n, a, lda, tau, work, -1)
		work = make([]float64, int(work[0]))
		impl.Dgelqf(m, n, a, lda, tau, work, len(work))

		aUnblocked := make([]float64, len(a))
		copy(aUnblocked, a)
		for i := range work {
			work[i] = math.NaN()
		}
		impl.Dorgl2(m, n, k, aUnblocked, lda, tau, work)
		// make sure work isn't used before initialized
		for i := range work {
			work[i] = math.NaN()
		}
		impl.Dorglq(m, n, k, a, lda, tau, work, len(work))
		if !floats.EqualApprox(a, aUnblocked, 1e-10) {
			t.Errorf("Q Mismatch. m = %d, n = %d, k = %d, lda = %d", m, n, k, lda)
		}
	}
}
