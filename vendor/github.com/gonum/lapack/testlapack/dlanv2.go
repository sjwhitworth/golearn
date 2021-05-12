// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

type Dlanv2er interface {
	Dlanv2(a, b, c, d float64) (aa, bb, cc, dd float64, rt1r, rt1i, rt2r, rt2i float64, cs, sn float64)
}

func Dlanv2Test(t *testing.T, impl Dlanv2er) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 1000; i++ {
		a := rnd.NormFloat64()
		b := rnd.NormFloat64()
		c := rnd.NormFloat64()
		d := rnd.NormFloat64()
		aa, bb, cc, dd, rt1r, rt1i, rt2r, rt2i, cs, sn := impl.Dlanv2(a, b, c, d)

		mat := fmt.Sprintf("[%v %v; %v %v]", a, b, c, d)
		if cc == 0 {
			if rt1i != 0 || rt2i != 0 {
				t.Errorf("Unexpected complex eigenvalues for %v", mat)
			}
		} else {
			if aa != dd {
				t.Errorf("Diagonal elements not equal for %v: got [%v %v]", mat, aa, dd)
			}
			if bb*cc >= 0 {
				t.Errorf("Non-diagonal elements have the same sign for %v: got [%v %v]", mat, bb, cc)
			} else {
				im := math.Sqrt(-bb * cc)
				if math.Abs(rt1i-im) > 1e-14 && math.Abs(rt1i+im) > 1e-14 {
					t.Errorf("Unexpected imaginary part of eigenvalue for %v: got %v, want %v or %v", mat, rt1i, im, -im)
				}
				if math.Abs(rt2i-im) > 1e-14 && math.Abs(rt2i+im) > 1e-14 {
					t.Errorf("Unexpected imaginary part of eigenvalue for %v: got %v, want %v or %v", mat, rt2i, im, -im)
				}
			}
		}
		if rt1r != aa && rt1r != dd {
			t.Errorf("Unexpected real part of eigenvalue for %v: got %v, want %v or %v", mat, rt1r, aa, dd)
		}
		if rt2r != aa && rt2r != dd {
			t.Errorf("Unexpected real part of eigenvalue for %v: got %v, want %v or %v", mat, rt2r, aa, dd)
		}
		if math.Abs(math.Hypot(cs, sn)-1) > 1e-14 {
			t.Errorf("Unexpected unitary matrix for %v: got cs %v, sn %v", mat, cs, sn)
		}

		gota := cs*(aa*cs-bb*sn) - sn*(cc*cs-dd*sn)
		gotb := cs*(aa*sn+bb*cs) - sn*(cc*sn+dd*cs)
		gotc := sn*(aa*cs-bb*sn) + cs*(cc*cs-dd*sn)
		gotd := sn*(aa*sn+bb*cs) + cs*(cc*sn+dd*cs)
		if math.Abs(gota-a) > 1e-14 ||
			math.Abs(gotb-b) > 1e-14 ||
			math.Abs(gotc-c) > 1e-14 ||
			math.Abs(gotd-d) > 1e-14 {
			t.Errorf("Unexpected factorization: got [%v %v; %v %v], want [%v %v; %v %v]", gota, gotb, gotc, gotd, a, b, c, d)
		}
	}
}
