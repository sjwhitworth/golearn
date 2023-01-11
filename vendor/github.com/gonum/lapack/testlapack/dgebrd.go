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

type Dgebrder interface {
	Dgebrd(m, n int, a []float64, lda int, d, e, tauQ, tauP, work []float64, lwork int)
	Dgebd2er
}

func DgebrdTest(t *testing.T, impl Dgebrder) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{100, 100, 0},
		{100, 150, 0},
		{150, 100, 0},
		{100, 100, 200},
		{100, 150, 200},
		{150, 100, 200},

		{300, 300, 0},
		{300, 400, 0},
		{400, 300, 0},
		{300, 300, 500},
		{300, 400, 500},
		{300, 400, 500},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		minmn := min(m, n)
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}

		d := make([]float64, minmn)
		e := make([]float64, minmn-1)
		tauP := make([]float64, minmn)
		tauQ := make([]float64, minmn)
		work := make([]float64, max(m, n))
		for i := range work {
			work[i] = math.NaN()
		}

		// Store a.
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		// Compute the true answer with the unblocked algorithm.
		impl.Dgebd2(m, n, a, lda, d, e, tauQ, tauP, work)
		aAns := make([]float64, len(a))
		copy(aAns, a)
		dAns := make([]float64, len(d))
		copy(dAns, d)
		eAns := make([]float64, len(e))
		copy(eAns, e)
		tauQAns := make([]float64, len(tauQ))
		copy(tauQAns, tauQ)
		tauPAns := make([]float64, len(tauP))
		copy(tauPAns, tauP)

		// Test with optimal work.
		lwork := -1
		copy(a, aCopy)
		impl.Dgebrd(m, n, a, lda, d, e, tauQ, tauP, work, lwork)
		work = make([]float64, int(work[0]))
		lwork = len(work)
		for i := range work {
			work[i] = math.NaN()
		}
		for i := range d {
			d[i] = math.NaN()
		}
		for i := range e {
			e[i] = math.NaN()
		}
		for i := range tauQ {
			tauQ[i] = math.NaN()
		}
		for i := range tauP {
			tauP[i] = math.NaN()
		}
		impl.Dgebrd(m, n, a, lda, d, e, tauQ, tauP, work, lwork)

		// Test answers
		if !floats.EqualApprox(a, aAns, 1e-10) {
			t.Errorf("a mismatch")
		}
		if !floats.EqualApprox(d, dAns, 1e-10) {
			t.Errorf("d mismatch")
		}
		if !floats.EqualApprox(e, eAns, 1e-10) {
			t.Errorf("e mismatch")
		}
		if !floats.EqualApprox(tauQ, tauQAns, 1e-10) {
			t.Errorf("tauQ mismatch")
		}
		if !floats.EqualApprox(tauP, tauPAns, 1e-10) {
			t.Errorf("tauP mismatch")
		}

		// Test with shorter than optimal work.
		lwork--
		copy(a, aCopy)
		for i := range d {
			d[i] = 0
		}
		for i := range e {
			e[i] = 0
		}
		for i := range tauP {
			tauP[i] = 0
		}
		for i := range tauQ {
			tauQ[i] = 0
		}
		impl.Dgebrd(m, n, a, lda, d, e, tauQ, tauP, work, lwork)

		// Test answers
		if !floats.EqualApprox(a, aAns, 1e-10) {
			t.Errorf("a mismatch")
		}
		if !floats.EqualApprox(d, dAns, 1e-10) {
			t.Errorf("d mismatch")
		}
		if !floats.EqualApprox(e, eAns, 1e-10) {
			t.Errorf("e mismatch")
		}
		if !floats.EqualApprox(tauQ, tauQAns, 1e-10) {
			t.Errorf("tauQ mismatch")
		}
		if !floats.EqualApprox(tauP, tauPAns, 1e-10) {
			t.Errorf("tauP mismatch")
		}
	}
}
