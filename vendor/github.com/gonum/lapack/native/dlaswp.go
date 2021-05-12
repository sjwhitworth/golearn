// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import "github.com/gonum/blas/blas64"

// Dlaswp swaps the rows k1 to k2 of a according to the indices in ipiv.
// a is a matrix with n columns and stride lda. incX is the increment for ipiv.
// k1 and k2 are zero-indexed. If incX is negative, then loops from k2 to k1
//
// Dlaswp is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlaswp(n int, a []float64, lda, k1, k2 int, ipiv []int, incX int) {
	if incX != 1 && incX != -1 {
		panic(absIncNotOne)
	}
	bi := blas64.Implementation()
	if incX == 1 {
		for k := k1; k <= k2; k++ {
			bi.Dswap(n, a[k*lda:], 1, a[ipiv[k]*lda:], 1)
		}
		return
	}
	for k := k2; k >= k1; k-- {
		bi.Dswap(n, a[k*lda:], 1, a[ipiv[k]*lda:], 1)
	}
}
