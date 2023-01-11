// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !amd64 noasm appengine

package asm

func DscalUnitary(alpha float64, x []float64) {
	for i := range x {
		x[i] *= alpha
	}
}

func DscalUnitaryTo(dst []float64, alpha float64, x []float64) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

// incX must be positive.
func DscalInc(alpha float64, x []float64, n, incX uintptr) {
	var ix uintptr
	for i := 0; i < int(n); i++ {
		x[ix] *= alpha
		ix += incX
	}
}

// incDst and incX must be positive.
func DscalIncTo(dst []float64, incDst uintptr, alpha float64, x []float64, n, incX uintptr) {
	var idst, ix uintptr
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha * x[ix]
		ix += incX
		idst += incDst
	}
}
