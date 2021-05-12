// Generated code do not edit. Run `go generate`.

// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

func ZaxpyUnitary(alpha complex128, x, y []complex128) {
	for i, v := range x {
		y[i] += alpha * v
	}
}

func ZaxpyUnitaryTo(dst []complex128, alpha complex128, x, y []complex128) {
	for i, v := range x {
		dst[i] = alpha*v + y[i]
	}
}

func ZaxpyInc(alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr) {
	for i := 0; i < int(n); i++ {
		y[iy] += alpha * x[ix]
		ix += incX
		iy += incY
	}
}

func ZaxpyIncTo(dst []complex128, incDst, idst uintptr, alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr) {
	for i := 0; i < int(n); i++ {
		dst[idst] = alpha*x[ix] + y[iy]
		ix += incX
		iy += incY
		idst += incDst
	}
}
