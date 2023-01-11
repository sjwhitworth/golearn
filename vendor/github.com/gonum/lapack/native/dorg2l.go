// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

// Dorg2l generates an m×n matrix Q with orthonormal columns which is defined
// as the last n columns of a product of k elementary reflectors of order m.
//  Q = H_{k-1} * ... * H_1 * H_0
// See Dgelqf for more information. It must be that m >= n >= k.
//
// tau contains the scalar reflectors computed by Dgeqlf. tau must have length
// at least k, and Dorg2l will panic otherwise.
//
// work contains temporary memory, and must have length at least n. Dorg2l will
// panic otherwise.
//
// Dorg2l is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorg2l(m, n, k int, a []float64, lda int, tau, work []float64) {
	checkMatrix(m, n, a, lda)
	if len(tau) < k {
		panic(badTau)
	}
	if len(work) < n {
		panic(badWork)
	}
	if m < n {
		panic(mLTN)
	}
	if k > n {
		panic(kGTN)
	}
	if n == 0 {
		return
	}

	// Initialize columns 0:n-k to columns of the unit matrix.
	for j := 0; j < n-k; j++ {
		for l := 0; l < m; l++ {
			a[l*lda+j] = 0
		}
		a[(m-n+j)*lda+j] = 1
	}

	bi := blas64.Implementation()
	for i := 0; i < k; i++ {
		ii := n - k + i

		// Apply H_i to A[0:m-k+i, 0:n-k+i] from the left.
		a[(m-n+ii)*lda+ii] = 1
		impl.Dlarf(blas.Left, m-n+ii+1, ii, a[ii:], lda, tau[i], a, lda, work)
		bi.Dscal(m-n+ii, -tau[i], a[ii:], lda)
		a[(m-n+ii)*lda+ii] = 1 - tau[i]

		// Set A[m-k+i:m, n-k+i+1] to zero.
		for l := m - n + ii + 1; l < m; l++ {
			a[l*lda+ii] = 0
		}
	}
}
