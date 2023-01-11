// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

// Dgeqrf computes the QR factorization of the m×n matrix A using a blocked
// algorithm. See the documentation for Dgeqr2 for a description of the
// parameters at entry and exit.
//
// work is temporary storage, and lwork specifies the usable memory length.
// At minimum, lwork >= m and this function will panic otherwise.
// Dgeqrf is a blocked QR factorization, but the block size is limited
// by the temporary space available. If lwork == -1, instead of performing Dgeqrf,
// the optimal work length will be stored into work[0].
//
// tau must have length at least min(m,n), and this function will panic otherwise.
func (impl Implementation) Dgeqrf(m, n int, a []float64, lda int, tau, work []float64, lwork int) {
	// nb is the optimal blocksize, i.e. the number of columns transformed at a time.
	nb := impl.Ilaenv(1, "DGEQRF", " ", m, n, -1, -1)
	lworkopt := n * max(nb, 1)
	lworkopt = max(n, lworkopt)
	if lwork == -1 {
		work[0] = float64(lworkopt)
		return
	}
	checkMatrix(m, n, a, lda)
	if len(work) < lwork {
		panic(shortWork)
	}
	if lwork < n {
		panic(badWork)
	}
	k := min(m, n)
	if len(tau) < k {
		panic(badTau)
	}
	if k == 0 {
		return
	}
	nbmin := 2 // Minimal block size.
	var nx int // Use unblocked (unless changed in the next for loop)
	iws := n
	ldwork := nb
	// Only consider blocked if the suggested block size is > 1 and the
	// number of rows or columns is sufficiently large.
	if 1 < nb && nb < k {
		// nx is the block size at which the code switches from blocked
		// to unblocked.
		nx = max(0, impl.Ilaenv(3, "DGEQRF", " ", m, n, -1, -1))
		if k > nx {
			iws = ldwork * n
			if lwork < iws {
				// Not enough workspace to use the optimal block
				// size. Get the minimum block size instead.
				nb = lwork / n
				nbmin = max(2, impl.Ilaenv(2, "DGEQRF", " ", m, n, -1, -1))
			}
		}
	}
	for i := range work {
		work[i] = 0
	}
	// Compute QR using a blocked algorithm.
	var i int
	if nbmin <= nb && nb < k && nx < k {
		for i = 0; i < k-nx; i += nb {
			ib := min(k-i, nb)
			// Compute the QR factorization of the current block.
			impl.Dgeqr2(m-i, ib, a[i*lda+i:], lda, tau[i:], work)
			if i+ib < n {
				// Form the triangular factor of the block reflector and apply H^T
				// In Dlarft, work becomes the T matrix.
				impl.Dlarft(lapack.Forward, lapack.ColumnWise, m-i, ib,
					a[i*lda+i:], lda,
					tau[i:],
					work, ldwork)
				impl.Dlarfb(blas.Left, blas.Trans, lapack.Forward, lapack.ColumnWise,
					m-i, n-i-ib, ib,
					a[i*lda+i:], lda,
					work, ldwork,
					a[i*lda+i+ib:], lda,
					work[ib*ldwork:], ldwork)
			}
		}
	}
	// Call unblocked code on the remaining columns.
	if i < k {
		impl.Dgeqr2(m-i, n-i, a[i*lda+i:], lda, tau[i:], work)
	}
}
