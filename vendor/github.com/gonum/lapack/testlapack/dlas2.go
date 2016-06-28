// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"testing"
)

type Dlas2er interface {
	Dlas2(f, g, h float64) (min, max float64)
}

func Dlas2Test(t *testing.T, impl Dlas2er) {
	for i, test := range []struct {
		f, g, h, ssmin, ssmax float64
	}{
		// Singular values computed from Octave.
		{10, 30, 12, 3.567778859365365, 33.634371616111189},
		{10, 30, -12, 3.567778859365365, 33.634371616111189},
		{2, 30, -12, 0.741557056404952, 32.364333658088754},
		{-2, 5, 12, 1.842864429909778, 13.023204317408728},
	} {
		ssmin, ssmax := impl.Dlas2(test.f, test.g, test.h)
		if math.Abs(ssmin-test.ssmin) > 1e-12 {
			t.Errorf("Case %d, minimal singular value mismatch. Want %v, got %v", i, test.ssmin, ssmin)
		}
		if math.Abs(ssmax-test.ssmax) > 1e-12 {
			t.Errorf("Case %d, minimal singular value mismatch. Want %v, got %v", i, test.ssmin, ssmin)
		}
	}
}
