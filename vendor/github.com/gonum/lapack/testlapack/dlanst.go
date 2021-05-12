// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/lapack"
)

type Dlanster interface {
	Dlanst(norm lapack.MatrixNorm, n int, d, e []float64) float64
	Dlanger
}

func DlanstTest(t *testing.T, impl Dlanster) {
	rnd := rand.New(rand.NewSource(1))
	for _, norm := range []lapack.MatrixNorm{lapack.MaxAbs, lapack.MaxColumnSum, lapack.MaxRowSum, lapack.NormFrob} {
		for _, n := range []int{1, 3, 10, 100} {
			for cas := 0; cas < 100; cas++ {
				d := make([]float64, n)
				for i := range d {
					d[i] = rnd.NormFloat64()
				}
				e := make([]float64, n-1)
				for i := range e {
					e[i] = rnd.NormFloat64()
				}

				m := n
				lda := n
				a := make([]float64, m*lda)
				for i := 0; i < n; i++ {
					a[i*lda+i] = d[i]
				}
				for i := 0; i < n-1; i++ {
					a[i*lda+i+1] = e[i]
					a[(i+1)*lda+i] = e[i]
				}
				work := make([]float64, n)
				syNorm := impl.Dlanst(norm, n, d, e)
				geNorm := impl.Dlange(norm, m, n, a, lda, work)
				if math.Abs(syNorm-geNorm) > 1e-12 {
					t.Errorf("Norm mismatch: norm = %v, cas = %v, n = %v. Want %v, got %v.", string(norm), cas, n, geNorm, syNorm)
				}
			}
		}
	}
}
