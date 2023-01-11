// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas/blas64"
)

type Dorg2ler interface {
	Dorg2l(m, n, k int, a []float64, lda int, tau, work []float64)
	Dgeql2er
}

func Dorg2lTest(t *testing.T, impl Dorg2ler) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, k, lda int
	}{
		{5, 4, 3, 0},
		{5, 4, 4, 0},
		{3, 3, 2, 0},
		{5, 5, 5, 0},
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
		tau := nanSlice(max(m, n))
		work := make([]float64, n)
		impl.Dgeql2(m, n, a, lda, tau, work)

		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		impl.Dorg2l(m, n, k, a, lda, tau[n-k:], work)
		if !hasOrthonormalColumns(m, n, a, lda) {
			t.Errorf("Q is not orthonormal. m = %v, n = %v, k = %v", m, n, k)
		}
	}
}

// hasOrthornormalColumns checks that the columns of a are orthonormal.
func hasOrthonormalColumns(m, n int, a []float64, lda int) bool {
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			dot := blas64.Dot(m,
				blas64.Vector{Inc: lda, Data: a[i:]},
				blas64.Vector{Inc: lda, Data: a[j:]},
			)
			if i == j {
				if math.Abs(dot-1) > 1e-10 {
					return false
				}
			} else {
				if math.Abs(dot) > 1e-10 {
					return false
				}
			}
		}
	}
	return true
}
