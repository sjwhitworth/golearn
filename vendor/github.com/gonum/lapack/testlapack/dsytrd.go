// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dsytrder interface {
	Dsytrd(uplo blas.Uplo, n int, a []float64, lda int, d, e, tau, work []float64, lwork int)
	Dsytd2er
}

func DsytrdTest(t *testing.T, impl Dsytrder) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, test := range []struct {
			n, lda int
		}{
			{10, 0},
			{50, 0},
			{100, 0},
			{150, 0},
			{300, 0},

			{10, 20},
			{50, 70},
			{100, 120},
			{150, 170},
			{300, 320},
		} {
			n := test.n
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			a := make([]float64, n*lda)
			for i := range a {
				a[i] = rnd.NormFloat64()
			}
			d2 := make([]float64, n)
			e2 := make([]float64, n)
			tau2 := make([]float64, n)

			aCopy := make([]float64, len(a))
			copy(aCopy, a)
			impl.Dsytd2(uplo, n, a, lda, d2, e2, tau2)
			aAns := make([]float64, len(a))
			copy(aAns, a)

			copy(a, aCopy)
			d := make([]float64, n)
			e := make([]float64, n)
			tau := make([]float64, n)
			work := make([]float64, 1)
			impl.Dsytrd(uplo, n, a, lda, d, e, tau, work, -1)
			work = make([]float64, int(work[0]))
			impl.Dsytrd(uplo, n, a, lda, d, e, tau, work, len(work))
			errStr := fmt.Sprintf("upper = %v, n = %v", uplo == blas.Upper, n)
			if !floats.EqualApprox(a, aAns, 1e-8) {
				t.Errorf("A mismatch: %s", errStr)
			}
			if !floats.EqualApprox(d, d2, 1e-8) {
				t.Errorf("D mismatch: %s", errStr)
			}
			if !floats.EqualApprox(e, e2, 1e-8) {
				t.Errorf("E mismatch: %s", errStr)
			}
			if !floats.EqualApprox(tau, tau2, 1e-8) {
				t.Errorf("Tau mismatch: %s", errStr)
			}
		}
	}
}
