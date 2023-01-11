// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

// Dtrti2 computes the inverse of a triangular matrix, storing the result in place
// into a. This is the BLAS level 2 version of the algorithm.
//
// Dtrti2 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dtrti2(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int) {
	checkMatrix(n, n, a, lda)
	if uplo != blas.Upper && uplo != blas.Lower {
		panic(badUplo)
	}
	if diag != blas.NonUnit && diag != blas.Unit {
		panic(badDiag)
	}
	bi := blas64.Implementation()

	nonUnit := diag == blas.NonUnit
	// TODO(btracey): Replace this with a row-major ordering.
	if uplo == blas.Upper {
		for j := 0; j < n; j++ {
			var ajj float64
			if nonUnit {
				ajj = 1 / a[j*lda+j]
				a[j*lda+j] = ajj
				ajj *= -1
			} else {
				ajj = -1
			}
			bi.Dtrmv(blas.Upper, blas.NoTrans, diag, j, a, lda, a[j:], lda)
			bi.Dscal(j, ajj, a[j:], lda)
		}
		return
	}
	for j := n - 1; j >= 0; j-- {
		var ajj float64
		if nonUnit {
			ajj = 1 / a[j*lda+j]
			a[j*lda+j] = ajj
			ajj *= -1
		} else {
			ajj = -1
		}
		if j < n-1 {
			bi.Dtrmv(blas.Lower, blas.NoTrans, diag, n-j-1, a[(j+1)*lda+j+1:], lda, a[(j+1)*lda+j:], lda)
			bi.Dscal(n-j-1, ajj, a[(j+1)*lda+j:], lda)
		}
	}
}
