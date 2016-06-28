// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dlaqr1er interface {
	Dlaqr1(n int, h []float64, ldh int, sr1, si1, sr2, si2 float64, v []float64)
}

func Dlaqr1Test(t *testing.T, impl Dlaqr1er) {
	rnd := rand.New(rand.NewSource(1))

	for _, n := range []int{2, 3} {
		for _, ldh := range []int{n, n + 1, n + 10} {
			for _, cas := range []int{1, 2} {
				for k := 0; k < 100; k++ {
					v := make([]float64, n)
					for i := range v {
						v[i] = math.NaN()
					}
					h := make([]float64, n*(n-1)*ldh)
					for i := range h {
						h[i] = math.NaN()
					}
					for i := 0; i < n; i++ {
						for j := 0; j < n; j++ {
							h[i*ldh+j] = rnd.NormFloat64()
						}
					}
					var sr1, sr2, si1, si2 float64
					if cas == 1 {
						sr1 = rnd.NormFloat64()
						sr2 = sr1
						si1 = rnd.NormFloat64()
						si2 = -si1
					} else {
						sr1 = rnd.NormFloat64()
						sr2 = rnd.NormFloat64()
						si1 = 0
						si2 = 0
					}
					impl.Dlaqr1(n, h, ldh, sr1, si1, sr2, si2, v)

					// Matrix H - s1*I.
					h1 := make([]complex128, n*n)
					for i := 0; i < n; i++ {
						for j := 0; j < n; j++ {
							h1[i*n+j] = complex(h[i*ldh+j], 0)
							if i == j {
								h1[i*n+j] -= complex(sr1, si1)
							}
						}
					}
					// First column of H - s2*I.
					h2 := make([]complex128, n)
					for i := 0; i < n; i++ {
						h2[i] = complex(h[i*ldh], 0)
					}
					h2[0] -= complex(sr2, si2)

					wantv := make([]float64, n)
					// Multiply (H-s1*I)*(H-s2*I) to get a tentative
					// wantv.
					for i := 0; i < n; i++ {
						for j := 0; j < n; j++ {
							wantv[i] += real(h1[i*n+j] * h2[j])
						}
					}
					// Get the unknown scale.
					scale := v[0] / wantv[0]
					// Compute the actual wantv.
					floats.Scale(scale, wantv)

					// The scale must be the same for all elements.
					if floats.Distance(wantv, v, math.Inf(1)) > 1e-13 {
						t.Errorf("n = %v, ldh = %v, case = %v: Unexpected value of v: got %v, want %v", n, ldh, cas, v, wantv)
					}
				}
			}
		}
	}
}
