// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dorgl2er interface {
	Dgelqfer
	Dorgl2(m, n, k int, a []float64, lda int, tau []float64, work []float64)
}

func Dorgl2Test(t *testing.T, impl Dorgl2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{3, 3, 0},
		{3, 4, 0},

		{5, 5, 20},
		{5, 10, 20},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = test.n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}
		k := min(m, n)
		tau := make([]float64, k)
		work := make([]float64, 1)
		impl.Dgelqf(m, n, a, lda, tau, work, -1)
		work = make([]float64, int(work[0]))
		impl.Dgelqf(m, n, a, lda, tau, work, len(work))

		q := constructQ("LQ", m, n, a, lda, tau)

		impl.Dorgl2(m, n, k, a, lda, tau, work)

		// Check that the first m rows match.
		same := true
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if !floats.EqualWithinAbsOrRel(q.Data[i*q.Stride+j], a[i*lda+j], 1e-12, 1e-12) {
					same = false
					break
				}
			}
		}
		if !same {
			t.Errorf("Q mismatch")
		}
	}
}
