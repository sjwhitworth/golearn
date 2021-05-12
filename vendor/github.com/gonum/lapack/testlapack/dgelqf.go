// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dgelqfer interface {
	Dgelq2er
	Dgelqf(m, n int, a []float64, lda int, tau, work []float64, lwork int)
}

func DgelqfTest(t *testing.T, impl Dgelqfer) {
	rnd := rand.New(rand.NewSource(1))
	for c, test := range []struct {
		m, n, lda int
	}{
		{10, 5, 0},
		{5, 10, 0},
		{10, 10, 0},
		{300, 5, 0},
		{3, 500, 0},
		{200, 200, 0},
		{300, 200, 0},
		{204, 300, 0},
		{1, 3000, 0},
		{3000, 1, 0},
		{10, 5, 30},
		{5, 10, 30},
		{10, 10, 30},
		{300, 5, 500},
		{3, 500, 600},
		{200, 200, 300},
		{300, 200, 300},
		{204, 300, 400},
		{1, 3000, 4000},
		{3000, 1, 4000},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]float64, m*lda)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				a[i*lda+j] = rnd.Float64()
			}
		}
		tau := make([]float64, n)
		for i := 0; i < n; i++ {
			tau[i] = rnd.Float64()
		}
		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		ans := make([]float64, len(a))
		copy(ans, a)
		work := make([]float64, m)
		for i := range work {
			work[i] = rnd.Float64()
		}
		// Compute unblocked QR.
		impl.Dgelq2(m, n, ans, lda, tau, work)
		// Compute blocked QR with small work.
		impl.Dgelqf(m, n, a, lda, tau, work, len(work))
		if !floats.EqualApprox(ans, a, 1e-12) {
			t.Errorf("Case %v, mismatch small work.", c)
		}
		// Try the full length of work.
		impl.Dgelqf(m, n, a, lda, tau, work, -1)
		lwork := int(work[0])
		work = make([]float64, lwork)
		copy(a, aCopy)
		impl.Dgelqf(m, n, a, lda, tau, work, lwork)
		if !floats.EqualApprox(ans, a, 1e-12) {
			t.Errorf("Case %v, mismatch large work.", c)
		}

		// Try a slightly smaller version of work to test blocking code.
		if len(work) <= m {
			continue
		}
		work = work[1:]
		lwork--
		copy(a, aCopy)
		impl.Dgelqf(m, n, a, lda, tau, work, lwork)
		if !floats.EqualApprox(ans, a, 1e-12) {
			t.Errorf("Case %v, mismatch large work.", c)
		}
	}
}
