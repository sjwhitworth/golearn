// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"testing"
)

type Dlae2er interface {
	Dlae2(a, b, c float64) (rt1, rt2 float64)
}

func Dlae2Test(t *testing.T, impl Dlae2er) {
	for _, test := range []struct {
		a, b, c float64
	}{
		{-10, 5, 3},
		{3, 5, -10},
		{0, 3, 0},
		{1, 3, 1},
		{1, -3, 1},
		{5, 0, 3},
		{3, 0, -5},
		{1, 3, 1.02},
		{1.02, 3, 1},
		{1, -3, -9},
	} {
		a := test.a
		b := test.b
		c := test.c
		rt1, rt2 := impl.Dlae2(a, b, c)

		errStr := fmt.Sprintf("a = %v, b = %v, c = %v", a, b, c)
		// Check if rt1 and rt2 are eigenvalues by checking if det(a - λI) = 0
		a1 := a - rt1
		c1 := c - rt1
		det := a1*c1 - b*b
		if math.Abs(det) > 1e-10 {
			t.Errorf("First eigenvalue mismatch. %s. Det = %v", errStr, det)
		}

		a2 := a - rt2
		c2 := c - rt2
		det = a2*c2 - b*b
		if math.Abs(det) > 1e-10 {
			t.Errorf("Second eigenvalue mismatch. %s. Det = %v", errStr, det)
		}
	}
}
