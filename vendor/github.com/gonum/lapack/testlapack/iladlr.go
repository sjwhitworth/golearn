// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import "testing"

type Iladlrer interface {
	Iladlr(m, n int, a []float64, lda int) int
}

func IladlrTest(t *testing.T, impl Iladlrer) {
	for i, test := range []struct {
		a         []float64
		m, n, lda int
		ans       int
	}{
		{
			a:   []float64{0, 0, 0, 0},
			m:   1,
			n:   1,
			lda: 2,
			ans: -1,
		},
		{
			a:   []float64{0, 0, 0, 0},
			m:   2,
			n:   2,
			lda: 2,
			ans: -1,
		},
		{
			a:   []float64{0, 0, 0, 0},
			m:   4,
			n:   1,
			lda: 1,
			ans: -1,
		},
		{
			a:   []float64{0, 0, 0, 0},
			m:   1,
			n:   4,
			lda: 4,
			ans: -1,
		},
		{
			a: []float64{
				1, 2, 3, 4,
				5, 6, 7, 8,
			},
			m:   2,
			n:   4,
			lda: 4,
			ans: 1,
		},
		{
			a: []float64{
				1, 2, 3, 0,
				0, 0, 0, 0,
			},
			m:   2,
			n:   4,
			lda: 4,
			ans: 0,
		},
		{
			a: []float64{
				0, 0, 3, 4,
				0, 0, 0, 0,
			},
			m:   2,
			n:   2,
			lda: 4,
			ans: -1,
		},
	} {
		ans := impl.Iladlr(test.m, test.n, test.a, test.lda)
		if ans != test.ans {
			t.Errorf("Column mismatch case %v. Want: %v, got: %v", i, test.ans, ans)
		}
	}
}
