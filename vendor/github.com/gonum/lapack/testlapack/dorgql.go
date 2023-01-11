// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dorgqler interface {
	Dorgql(m, n, k int, a []float64, lda int, tau, work []float64, lwork int)
	Dorg2ler
}

func DorgqlTest(t *testing.T, impl Dorgqler) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, k, lda int
	}{
		{5, 4, 3, 0},
		{100, 100, 100, 0},
		{200, 100, 50, 0},
		{200, 200, 50, 0},
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
			a[i] = rnd.NormFloat64()
		}
		tau := nanSlice(min(m, n))
		work := nanSlice(max(m, n))

		impl.Dgeql2(m, n, a, lda, tau, work)

		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		impl.Dorg2l(m, n, k, a, lda, tau, work)
		ans := make([]float64, len(a))
		copy(ans, a)

		impl.Dorgql(m, n, k, a, lda, tau, work, -1)
		work = make([]float64, int(work[0]))
		copy(a, aCopy)
		impl.Dorgql(m, n, k, a, lda, tau, work, len(work))

		if !floats.EqualApprox(a, ans, 1e-8) {
			t.Errorf("Answer mismatch. m = %v, n = %v, k = %v", m, n, k)
		}
	}
}
