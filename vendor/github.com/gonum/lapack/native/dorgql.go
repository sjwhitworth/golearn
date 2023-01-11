// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

// Dorgql generates the m×n matrix Q with orthonormal columns defined as the
// last n columns of a product of k elementary reflectors of order m
//  Q = H_{k-1} * ... * H_1 * H_0
// as returned by Dgelqf. See Dgelqf for more information.
//
// tau must have length at least k, and Dorgql will panic otherwise.
//
// work is temporary storage, and lwork specifies the usable memory length. At minimum,
// lwork >= n, and Dorgql will panic otherwise. The amount of blocking is
// limited by the usable length.
// If lwork == -1, instead of computing Dorgql the optimal work length is stored
// into work[0].
//
// Dorgql is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dorgql(m, n, k int, a []float64, lda int, tau, work []float64, lwork int) {
	checkMatrix(m, n, a, lda)
	if len(tau) < k {
		panic(badTau)
	}
	nb := impl.Ilaenv(1, "DORGQL", " ", m, n, k, -1)
	lworkopt := n * nb
	work[0] = float64(lworkopt)
	if lwork == -1 {
		return
	}
	if lwork < n {
		panic(badWork)
	}
	if n == 0 {
		return
	}

	nbmin := 2
	var nx, ldwork int
	iws := n
	if nb > 1 && nb < k {
		// Determine when to cross over from blocked to unblocked code.
		nx = max(0, impl.Ilaenv(3, "DORGQL", " ", m, n, k, -1))
		if nx < k {
			// Determine if workspace is large enough for blocked code.
			ldwork = nb
			iws = n * nb
			if lwork < iws {
				// Not enough workspace to use optimal nb: reduce nb and determine
				// the minimum value of nb.
				nb = lwork / n
				nbmin = max(2, impl.Ilaenv(2, "DORGQL", " ", m, n, k, -1))
			}
		}
	}

	var kk int
	if nb >= nbmin && nb < k && nx < k {
		// Use blocked code after the first block. The last kk columns are handled
		// by the block method.
		kk = min(k, ((k-nx+nb-1)/nb)*nb)

		// Set A(m-kk:m, 0:n-kk) to zero.
		for i := m - kk; i < m; i++ {
			for j := 0; j < n-kk; j++ {
				a[i*lda+j] = 0
			}
		}
	}

	// Use unblocked code for the first or only block.
	impl.Dorg2l(m-kk, n-kk, k-kk, a, lda, tau, work)
	if kk > 0 {
		// Use blocked code.
		for i := k - kk; i < k; i += nb {
			ib := min(nb, k-i)
			if n-k+i > 0 {
				// Form the triangular factor of the block reflector
				// H = H_{i+ib-1} * ... * H_{i+1} * H_i.
				impl.Dlarft(lapack.Backward, lapack.ColumnWise, m-k+i+ib, ib,
					a[n-k+i:], lda, tau[i:], work, ldwork)

				// Apply H to A[0:m-k+i+ib, 0:n-k+i] from the left.
				impl.Dlarfb(blas.Left, blas.NoTrans, lapack.Backward, lapack.ColumnWise,
					m-k+i+ib, n-k+i, ib, a[n-k+i:], lda, work, ldwork,
					a, lda, work[ib*ldwork:], ldwork)
			}

			// Apply H to rows 0:m-k+i+ib of current block.
			impl.Dorg2l(m-k+i+ib, ib, ib, a[n-k+i:], lda, tau[i:], work)

			// Set rows m-k+i+ib:m of current block to zero.
			for j := n - k + i; j < n-k+i+ib; j++ {
				for l := m - k + i + ib; l < m; l++ {
					a[l*lda+j] = 0
				}
			}
		}
	}
}
