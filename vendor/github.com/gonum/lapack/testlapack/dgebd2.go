// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"
)

type Dgebd2er interface {
	Dgebd2(m, n int, a []float64, lda int, d, e, tauq, taup, work []float64)
}

func Dgebd2Test(t *testing.T, impl Dgebd2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{3, 4, 0},
		{4, 3, 0},
		{3, 4, 10},
		{4, 3, 10},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		nb := min(m, n) // 'nb' name parallel with Dlabrd code.
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}
		d := nanSlice(nb)
		e := nanSlice(nb - 1)
		tauP := nanSlice(nb)
		tauQ := nanSlice(nb)
		work := nanSlice(max(m, n))
		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		impl.Dgebd2(m, n, a, lda, d, e, tauQ, tauP, work)
		if m >= n && nb == n {
			tauP[n-1] = 0
		}
		if m < n && nb == m {
			tauQ[m-1] = 0
		}

		checkBidiagonal(t, m, n, nb, a, lda, d, e, tauP, tauQ, aCopy)
	}
}
