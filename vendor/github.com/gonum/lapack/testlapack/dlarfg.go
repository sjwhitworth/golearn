// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dlarfger interface {
	Dlarfg(n int, alpha float64, x []float64, incX int) (beta, tau float64)
}

func DlarfgTest(t *testing.T, impl Dlarfger) {
	rnd := rand.New(rand.NewSource(1))
	for i, test := range []struct {
		alpha float64
		n     int
		x     []float64
	}{
		{
			alpha: 4,
			n:     3,
		},
		{
			alpha: -2,
			n:     3,
		},
		{
			alpha: 0,
			n:     3,
		},
		{
			alpha: 1,
			n:     1,
		},
		{
			alpha: 1,
			n:     2,
			x:     []float64{4, 5, 6},
		},
	} {
		n := test.n
		incX := 1
		var x []float64
		if test.x == nil {
			x = make([]float64, n-1)
			for i := range x {
				x[i] = rnd.Float64()
			}
		} else {
			x = make([]float64, n-1)
			copy(x, test.x)
		}
		xcopy := make([]float64, n-1)
		copy(xcopy, x)
		alpha := test.alpha
		beta, tau := impl.Dlarfg(n, alpha, x, incX)

		// Verify the returns and the values in v. Construct h and perform
		// the explicit multiplication.
		h := make([]float64, n*n)
		for i := 0; i < n; i++ {
			h[i*n+i] = 1
		}
		hmat := blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   h,
		}
		v := make([]float64, n)
		copy(v[1:], x)
		v[0] = 1
		vVec := blas64.Vector{
			Inc:  1,
			Data: v,
		}
		blas64.Ger(-tau, vVec, vVec, hmat)
		eye := blas64.General{
			Rows:   n,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, n*n),
		}
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, hmat, hmat, 0, eye)
		iseye := true
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j {
					if math.Abs(eye.Data[i*n+j]-1) > 1e-14 {
						iseye = false
					}
				} else {
					if math.Abs(eye.Data[i*n+j]) > 1e-14 {
						iseye = false
					}
				}
			}
		}
		if !iseye {
			t.Errorf("H^T * H is not I %v", eye)
		}

		xVec := blas64.Vector{
			Inc:  1,
			Data: make([]float64, n),
		}
		xVec.Data[0] = test.alpha
		copy(xVec.Data[1:], xcopy)

		ans := make([]float64, n)
		ansVec := blas64.Vector{
			Inc:  1,
			Data: ans,
		}
		blas64.Gemv(blas.NoTrans, 1, hmat, xVec, 0, ansVec)
		if math.Abs(ans[0]-beta) > 1e-14 {
			t.Errorf("Case %v, beta mismatch. Want %v, got %v", i, ans[0], beta)
		}
		for i := 1; i < n; i++ {
			if math.Abs(ans[i]) > 1e-14 {
				t.Errorf("Case %v, nonzero answer %v", i, ans)
				break
			}
		}
	}
}
