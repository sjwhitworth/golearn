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

type Dpotrfer interface {
	Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool)
}

func DpotrfTest(t *testing.T, impl Dpotrfer) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for i, test := range []struct {
		n int
	}{
		{n: 10},
		{n: 30},
		{n: 63},
		{n: 65},
		{n: 128},
		{n: 1000},
	} {
		n := test.n
		// Construct a positive-definite symmetric matrix
		base := make([]float64, n*n)
		for i := range base {
			base[i] = rnd.Float64()
		}
		a := make([]float64, len(base))
		bi.Dgemm(blas.Trans, blas.NoTrans, n, n, n, 1, base, n, base, n, 0, a, n)

		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		// Test with Upper
		impl.Dpotrf(blas.Upper, n, a, n)

		// zero all the other elements
		for i := 0; i < n; i++ {
			for j := 0; j < i; j++ {
				a[i*n+j] = 0
			}
		}
		// Multiply u^T * u
		ans := make([]float64, len(a))
		bi.Dsyrk(blas.Upper, blas.Trans, n, n, 1, a, n, 0, ans, n)

		match := true
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				if !floats.EqualWithinAbsOrRel(ans[i*n+j], aCopy[i*n+j], 1e-14, 1e-14) {
					match = false
				}
			}
		}
		if !match {
			//fmt.Println(aCopy)
			//fmt.Println(ans)
			t.Errorf("Case %v: Mismatch for upper", i)
		}
	}
}
