// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import "math"

// newGuardedVector allocates a new slice and returns it as three subslices.
// v is a strided vector that contains elements of data at indices i*inc and
// NaN elsewhere. frontGuard and backGuard are filled with NaN values, and
// their backing arrays are directly adjacent to v in memory. The three slices
// can be used to detect invalid memory reads and writes.
func newGuardedVector(data []float64, inc int) (v, frontGuard, backGuard []float64) {
	if inc < 0 {
		inc = -inc
	}
	guard := 2 * inc
	size := (len(data)-1)*inc + 1
	whole := make([]float64, size+2*guard)
	v = whole[guard : len(whole)-guard]
	for i := range whole {
		whole[i] = math.NaN()
	}
	for i, d := range data {
		v[i*inc] = d
	}
	return v, whole[:guard], whole[len(whole)-guard:]
}

// allNaN returns true if x contains only NaN values, and false otherwise.
func allNaN(x []float64) bool {
	for _, v := range x {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// equalStrided returns true if the strided vector x contains elements of the
// dense vector ref at indices i*inc, false otherwise.
func equalStrided(ref, x []float64, inc int) bool {
	if inc < 0 {
		inc = -inc
	}
	for i, v := range ref {
		if x[i*inc] != v {
			return false
		}
	}
	return true
}

// nonStridedWrite returns false if all elements of x at non-stride indices are
// equal to NaN, true otherwise.
func nonStridedWrite(x []float64, inc int) bool {
	if inc < 0 {
		inc = -inc
	}
	for i, v := range x {
		if i%inc != 0 && !math.IsNaN(v) {
			return true
		}
	}
	return false
}
