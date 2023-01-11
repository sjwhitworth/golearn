// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

// Dorglq generates an m×n matrix Q with orthonormal columns defined by the
// product of elementary reflectors as computed by Dgelqf.
//  Q = H_0 * H_1 * ... * H_{k-1}
// Dorglq is the blocked version of Dorgl2 that makes greater use of level-3 BLAS
// routines.
//
// len(tau) >= k, 0 <= k <= n, and 0 <= n <= m.
//
// work is temporary storage, and lwork specifies the usable memory length. At minimum,
// lwork >= m, and the amount of blocking is limited by the usable length.
// If lwork == -1, instead of computing Dorglq the optimal work length is stored
// into work[0].
//
// Dorglq will panic if the conditions on input values are not met.
//
// Dorglq is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorglq(m, n, k int, a []float64, lda int, tau, work []float64, lwork int) {
	nb := impl.Ilaenv(1, "DORGLQ", " ", m, n, k, -1)
	// work is treated as an n×nb matrix
	if lwork == -1 {
		work[0] = float64(max(1, m) * nb)
		return
	}
	checkMatrix(m, n, a, lda)
	if k < 0 {
		panic(kLT0)
	}
	if k > m {
		panic(kGTM)
	}
	if m > n {
		panic(nLTM)
	}
	if len(tau) < k {
		panic(badTau)
	}
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < m {
		panic(badWork)
	}
	if m == 0 {
		return
	}
	nbmin := 2 // Minimum number of blocks
	var nx int // Minimum number of rows
	iws := m   // Length of work needed
	var ldwork int
	if nb > 1 && nb < k {
		nx = max(0, impl.Ilaenv(3, "DORGLQ", " ", m, n, k, -1))
		if nx < k {
			ldwork = nb
			iws = m * ldwork
			if lwork < iws {
				nb = lwork / m
				ldwork = nb
				nbmin = max(2, impl.Ilaenv(2, "DORGLQ", " ", m, n, k, -1))
			}
		}
	}
	var ki, kk int
	if nb >= nbmin && nb < k && nx < k {
		// The first kk rows are handled by the blocked method.
		// Note: lapack has nx here, but this means the last nx rows are handled
		// serially which could be quite different than nb.
		ki = ((k - nb - 1) / nb) * nb
		kk = min(k, ki+nb)
		for i := kk; i < m; i++ {
			for j := 0; j < kk; j++ {
				a[i*lda+j] = 0
			}
		}
	}
	if kk < m {
		// Perform the operation on colums kk to the end.
		impl.Dorgl2(m-kk, n-kk, k-kk, a[kk*lda+kk:], lda, tau[kk:], work)
	}
	if kk == 0 {
		return
	}
	// Perform the operation on column-blocks
	for i := ki; i >= 0; i -= nb {
		ib := min(nb, k-i)
		if i+ib < m {
			impl.Dlarft(lapack.Forward, lapack.RowWise,
				n-i, ib,
				a[i*lda+i:], lda,
				tau[i:],
				work, ldwork)

			impl.Dlarfb(blas.Right, blas.Trans, lapack.Forward, lapack.RowWise,
				m-i-ib, n-i, ib,
				a[i*lda+i:], lda,
				work, ldwork,
				a[(i+ib)*lda+i:], lda,
				work[ib*ldwork:], ldwork)
		}
		impl.Dorgl2(ib, n-i, ib, a[i*lda+i:], lda, tau[i:], work)
		for l := i; l < i+ib; l++ {
			for j := 0; j < i; j++ {
				a[l*lda+j] = 0
			}
		}
	}
}
