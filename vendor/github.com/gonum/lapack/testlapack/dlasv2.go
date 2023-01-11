// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/floats"
)

type Dlasv2er interface {
	Dlasv2(f, g, h float64) (ssmin, ssmax, snr, csr, snl, csl float64)
}

func Dlasv2Test(t *testing.T, impl Dlasv2er) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		f := rnd.NormFloat64()
		g := rnd.NormFloat64()
		h := rnd.NormFloat64()

		ssmin, ssmax, snr, csr, snl, csl := impl.Dlasv2(f, g, h)

		// tmp =
		// [ csl snl] [f g]
		// [-snl csl] [0 h]
		tmp11 := csl * f
		tmp12 := csl*g + snl*h
		tmp21 := -snl * f
		tmp22 := -snl*g + csl*h
		// lhs =
		// [tmp11 tmp12] [csr -snr]
		// [tmp21 tmp22] [snr  csr]
		ans11 := tmp11*csr + tmp12*snr
		ans12 := tmp11*-snr + tmp12*csr
		ans21 := tmp21*csr + tmp22*snr
		ans22 := tmp21*-snr + tmp22*csr

		lhs := []float64{ans11, ans12, ans21, ans22}
		rhs := []float64{ssmax, 0, 0, ssmin}
		if !floats.EqualApprox(rhs, lhs, 1e-12) {
			t.Errorf("SVD mismatch. f = %v, g = %v, h = %v.\nLHS: %v\nRHS: %v", f, g, h, lhs, rhs)
		}
	}
}
