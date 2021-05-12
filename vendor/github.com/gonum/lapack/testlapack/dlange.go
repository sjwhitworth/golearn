// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

type Dlanger interface {
	Dlange(norm lapack.MatrixNorm, m, n int, a []float64, lda int, work []float64) float64
}

func DlangeTest(t *testing.T, impl Dlanger) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{4, 3, 0},
		{3, 4, 0},
		{4, 3, 100},
		{3, 4, 100},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.Float64() - 0.5
		}
		work := make([]float64, n)
		for i := range work {
			work[i] = rnd.Float64()
		}
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		// Test MaxAbs norm.
		norm := impl.Dlange(lapack.MaxAbs, m, n, a, lda, work)
		var ans float64
		for i := 0; i < m; i++ {
			idx := blas64.Iamax(n, blas64.Vector{Inc: 1, Data: aCopy[i*lda:]})
			ans = math.Max(ans, math.Abs(a[i*lda+idx]))
		}
		// Should be strictly equal because there is no floating point summation error.
		if ans != norm {
			t.Errorf("MaxAbs mismatch. Want %v, got %v.", ans, norm)
		}

		// Test MaxColumnSum norm.
		norm = impl.Dlange(lapack.MaxColumnSum, m, n, a, lda, work)
		ans = 0
		for i := 0; i < n; i++ {
			sum := blas64.Asum(m, blas64.Vector{Inc: lda, Data: aCopy[i:]})
			ans = math.Max(ans, sum)
		}
		if math.Abs(norm-ans) > 1e-14 {
			t.Errorf("MaxColumnSum mismatch. Want %v, got %v.", ans, norm)
		}

		// Test MaxRowSum norm.
		norm = impl.Dlange(lapack.MaxRowSum, m, n, a, lda, work)
		ans = 0
		for i := 0; i < m; i++ {
			sum := blas64.Asum(n, blas64.Vector{Inc: 1, Data: aCopy[i*lda:]})
			ans = math.Max(ans, sum)
		}
		if math.Abs(norm-ans) > 1e-14 {
			t.Errorf("MaxRowSum mismatch. Want %v, got %v.", ans, norm)
		}

		// Test Frobenius norm
		norm = impl.Dlange(lapack.NormFrob, m, n, a, lda, work)
		ans = 0
		for i := 0; i < m; i++ {
			sum := blas64.Nrm2(n, blas64.Vector{Inc: 1, Data: aCopy[i*lda:]})
			ans += sum * sum
		}
		ans = math.Sqrt(ans)
		if math.Abs(norm-ans) > 1e-14 {
			t.Errorf("NormFrob mismatch. Want %v, got %v.", ans, norm)
		}
	}
}
