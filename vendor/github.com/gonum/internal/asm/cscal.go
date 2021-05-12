// Generated code do not edit. Run `go generate`.

// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

func CscalUnitary(alpha complex64, x []complex64) {
	for i := range x {
		x[i] *= alpha
	}
}

func CscalUnitaryTo(dst []complex64, alpha complex64, x []complex64) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}

// incX must be positive.
func CscalInc(alpha complex64, x []complex64, n, incX uintptr) {
	var ix uintptr
	for i := 0; i < int(n); i++ {
		x[ix] *= alpha
		ix += incX
	}
}

// incDst and incX must be positive.
func CscalIncTo(dst []complex64, incDst uintptr, alpha complex64, x []complex64, n, incX uintptr) {
	var idst, ix uintptr
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha * x[ix]
		ix += incX
		idst += incDst
	}
}
