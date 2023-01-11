// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"testing"
)

type Dlabrder interface {
	Dlabrd(m, n, nb int, a []float64, lda int, d, e, tauq, taup, x []float64, ldx int, y []float64, ldy int)
}

func DlabrdTest(t *testing.T, impl Dlabrder) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, nb, lda, ldx, ldy int
	}{
		{4, 5, 2, 0, 0, 0},
		{4, 5, 4, 0, 0, 0},
		{5, 5, 2, 0, 0, 0},
		{5, 5, 5, 0, 0, 0},
		{5, 4, 4, 0, 0, 0},
		{5, 4, 4, 0, 0, 0},

		{4, 5, 2, 10, 11, 12},
		{4, 5, 4, 10, 11, 12},
		{5, 5, 2, 10, 11, 12},
		{5, 5, 5, 10, 11, 12},
		{5, 4, 2, 10, 11, 12},
		{5, 4, 4, 10, 11, 12},

		{4, 5, 2, 11, 12, 10},
		{4, 5, 4, 11, 12, 10},
		{5, 5, 2, 11, 12, 10},
		{5, 5, 5, 11, 12, 10},
		{5, 4, 2, 11, 12, 10},
		{5, 4, 4, 11, 12, 10},

		{4, 5, 2, 12, 11, 10},
		{4, 5, 4, 12, 11, 10},
		{5, 5, 2, 12, 11, 10},
		{5, 5, 5, 12, 11, 10},
		{5, 4, 2, 12, 11, 10},
		{5, 4, 4, 12, 11, 10},
	} {
		m := test.m
		n := test.n
		nb := test.nb
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		ldy := test.ldy
		if ldy == 0 {
			ldy = nb
		}
		ldx := test.ldx
		if ldx == 0 {
			ldx = nb
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}
		d := make([]float64, nb)
		for i := range d {
			d[i] = math.NaN()
		}
		e := make([]float64, nb)
		for i := range e {
			e[i] = math.NaN()
		}
		tauP := make([]float64, nb)
		for i := range tauP {
			tauP[i] = math.NaN()
		}
		tauQ := make([]float64, nb)
		for i := range tauP {
			tauQ[i] = math.NaN()
		}
		x := make([]float64, m*ldx)
		for i := range x {
			x[i] = rnd.NormFloat64()
		}
		y := make([]float64, n*ldy)
		for i := range y {
			y[i] = rnd.NormFloat64()
		}
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		// Compute the reduction.
		impl.Dlabrd(m, n, nb, a, lda, d, e, tauQ, tauP, x, ldx, y, ldy)

		if m >= n && nb == n {
			tauP[n-1] = 0
		}
		if m < n && nb == m {
			tauQ[m-1] = 0
		}
		checkBidiagonal(t, m, n, nb, a, lda, d, e, tauP, tauQ, aCopy)
	}
}
