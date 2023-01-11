// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package asm

func DaxpyUnitary(alpha float64, x, y []float64)

func DaxpyUnitaryTo(dst []float64, alpha float64, x, y []float64)

func DaxpyInc(alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr)

func DaxpyIncTo(dst []float64, incDst, idst uintptr, alpha float64, x, y []float64, n, incX, incY, ix, iy uintptr)

func DdotUnitary(x, y []float64) (sum float64)

func DdotInc(x, y []float64, n, incX, incY, ix, iy uintptr) (sum float64)

func DscalUnitary(alpha float64, x []float64)

func DscalUnitaryTo(dst []float64, alpha float64, x []float64)

// incX must be positive.
func DscalInc(alpha float64, x []float64, n, incX uintptr)

// incDst and incX must be positive.
func DscalIncTo(dst []float64, incDst uintptr, alpha float64, x []float64, n, incX uintptr)
