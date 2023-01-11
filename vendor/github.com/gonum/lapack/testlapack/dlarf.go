// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dlarfer interface {
	Dlarf(side blas.Side, m, n int, v []float64, incv int, tau float64, c []float64, ldc int, work []float64)
}

func DlarfTest(t *testing.T, impl Dlarfer) {
	rnd := rand.New(rand.NewSource(1))
	for i, test := range []struct {
		m, n, ldc    int
		incv, lastv  int
		lastr, lastc int
		tau          float64
	}{
		{
			m:   3,
			n:   2,
			ldc: 2,

			incv:  4,
			lastv: 1,

			lastr: 2,
			lastc: 1,

			tau: 2,
		},
		{
			m:   2,
			n:   3,
			ldc: 3,

			incv:  4,
			lastv: 1,

			lastr: 1,
			lastc: 2,

			tau: 2,
		},
		{
			m:   2,
			n:   3,
			ldc: 3,

			incv:  4,
			lastv: 1,

			lastr: 0,
			lastc: 1,

			tau: 2,
		},
		{
			m:   2,
			n:   3,
			ldc: 3,

			incv:  4,
			lastv: 0,

			lastr: 0,
			lastc: 1,

			tau: 2,
		},
		{
			m:   10,
			n:   10,
			ldc: 10,

			incv:  4,
			lastv: 6,

			lastr: 9,
			lastc: 8,

			tau: 2,
		},
	} {
		// Construct a random matrix.
		c := make([]float64, test.ldc*test.m)
		for i := 0; i <= test.lastr; i++ {
			for j := 0; j <= test.lastc; j++ {
				c[i*test.ldc+j] = rnd.Float64()
			}
		}
		cCopy := make([]float64, len(c))
		copy(cCopy, c)
		cCopy2 := make([]float64, len(c))
		copy(cCopy2, c)

		// Test with side right.
		sz := max(test.m, test.n) // so v works for both right and left side.
		v := make([]float64, test.incv*sz+1)
		// Fill with nonzero entries up until lastv.
		for i := 0; i <= test.lastv; i++ {
			v[i*test.incv] = rnd.Float64()
		}
		// Construct h explicitly to compare.
		h := make([]float64, test.n*test.n)
		for i := 0; i < test.n; i++ {
			h[i*test.n+i] = 1
		}
		hMat := blas64.General{
			Rows:   test.n,
			Cols:   test.n,
			Stride: test.n,
			Data:   h,
		}
		vVec := blas64.Vector{
			Inc:  test.incv,
			Data: v,
		}
		blas64.Ger(-test.tau, vVec, vVec, hMat)

		// Apply multiplication (2nd copy is to avoid aliasing).
		cMat := blas64.General{
			Rows:   test.m,
			Cols:   test.n,
			Stride: test.ldc,
			Data:   cCopy,
		}
		cMat2 := blas64.General{
			Rows:   test.m,
			Cols:   test.n,
			Stride: test.ldc,
			Data:   cCopy2,
		}
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, cMat2, hMat, 0, cMat)

		// cMat now stores the true answer. Compare with the function call.
		work := make([]float64, sz)
		impl.Dlarf(blas.Right, test.m, test.n, v, test.incv, test.tau, c, test.ldc, work)
		if !floats.EqualApprox(c, cMat.Data, 1e-14) {
			t.Errorf("Dlarf mismatch right, case %v. Want %v, got %v", i, cMat.Data, c)
		}

		// Test on the left side.
		copy(c, cCopy2)
		copy(cCopy, c)
		// Construct h.
		h = make([]float64, test.m*test.m)
		for i := 0; i < test.m; i++ {
			h[i*test.m+i] = 1
		}
		hMat = blas64.General{
			Rows:   test.m,
			Cols:   test.m,
			Stride: test.m,
			Data:   h,
		}
		blas64.Ger(-test.tau, vVec, vVec, hMat)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hMat, cMat2, 0, cMat)
		impl.Dlarf(blas.Left, test.m, test.n, v, test.incv, test.tau, c, test.ldc, work)
		if !floats.EqualApprox(c, cMat.Data, 1e-14) {
			t.Errorf("Dlarf mismatch left, case %v. Want %v, got %v", i, cMat.Data, c)
		}
	}
}
