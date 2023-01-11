// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"github.com/gonum/blas"
	"github.com/gonum/lapack"
)

// Dormbr applies a multiplicative update to the matrix C based on a
// decomposition computed by Dgebrd.
//
// Dormbr computes
//  Q * C if vect == lapack.ApplyQ, side == blas.Left, and trans == blas.NoTrans
//  C * Q if vect == lapack.ApplyQ, side == blas.Right, and trans == blas.NoTrans
//  Q^T * C if vect == lapack.ApplyQ, side == blas.Left, and trans == blas.Trans
//  C * Q^T if vect == lapack.ApplyQ, side == blas.Right, and trans == blas.Trans
//
//  P * C if vect == lapack.ApplyP, side == blas.Left, and trans == blas.NoTrans
//  C * P if vect == lapack.ApplyP, side == blas.Left, and trans == blas.NoTrans
//  P^T * C if vect == lapack.ApplyP, side == blas.Right, and trans == blas.Trans
//  C * P^T if vect == lapack.ApplyP, side == blas.Right, and trans == blas.Trans
// where P and Q are the orthogonal matrices determined by Dgebrd, A = Q * B * P^T.
// See Dgebrd for the definitions of Q and P.
//
// If vect == lapack.ApplyQ, A is assumed to have been an nq×k matrix, while if
// vect == lapack.ApplyP, A is assumed to have been a k×nq matrix. nq = m if
// side == blas.Left, while nq = n if side == blas.Right.
//
// C is an m×n matrix. On exit it is updated by the multiplication listed above.
//
// tau must have length min(nq,k), and Dormbr will panic otherwise. tau contains
// the elementary reflectors to construct Q or P depending on the value of
// vect.
//
// Dormbr is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dormbr(vect lapack.DecompUpdate, side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int) {
	if side != blas.Left && side != blas.Right {
		panic(badSide)
	}
	if trans != blas.NoTrans && trans != blas.Trans {
		panic(badTrans)
	}
	if vect != lapack.ApplyP && vect != lapack.ApplyQ {
		panic(badDecompUpdate)
	}
	nq := n
	nw := m
	if side == blas.Left {
		nq = m
		nw = n
	}
	if vect == lapack.ApplyQ {
		checkMatrix(nq, min(nq, k), a, lda)
	} else {
		checkMatrix(min(nq, k), nq, a, lda)
	}

	applyQ := vect == lapack.ApplyQ
	left := side == blas.Left
	var nb int

	// The current implementation does not use opts, but a future change may
	// use these options so construct them.
	var opts string
	if side == blas.Left {
		opts = "L"
	} else {
		opts = "R"
	}
	if trans == blas.Trans {
		opts += "T"
	} else {
		opts += "N"
	}
	if applyQ {
		if left {
			nb = impl.Ilaenv(1, "DORMQR", opts, m-1, n, m-1, -1)
		} else {
			nb = impl.Ilaenv(1, "DORMQR", opts, m, n-1, n-1, -1)
		}
	} else {
		if left {
			nb = impl.Ilaenv(1, "DORMLQ", opts, m-1, n, m-1, -1)
		} else {
			nb = impl.Ilaenv(1, "DORMLQ", opts, m, n-1, n-1, -1)
		}
	}
	lworkopt := max(1, nw) * nb
	if lwork == -1 {
		work[0] = float64(lworkopt)
	}
	if applyQ {
		// Change the operation to get Q depending on the size of the initial
		// matrix to Dgebrd. The size matters due to the storage location of
		// the off-diagonal elements.
		if nq >= k {
			impl.Dormqr(side, trans, m, n, k, a, lda, tau, c, ldc, work, lwork)
		} else if nq > 1 {
			mi := m
			ni := n - 1
			i1 := 0
			i2 := 1
			if left {
				mi = m - 1
				ni = n
				i1 = 1
				i2 = 0
			}
			impl.Dormqr(side, trans, mi, ni, nq-1, a[1*lda:], lda, tau, c[i1*ldc+i2:], ldc, work, lwork)
		}
		return
	}
	transt := blas.Trans
	if trans == blas.Trans {
		transt = blas.NoTrans
	}
	// Change the operation to get P depending on the size of the initial
	// matrix to Dgebrd. The size matters due to the storage location of
	// the off-diagonal elements.
	if nq > k {
		impl.Dormlq(side, transt, m, n, k, a, lda, tau, c, ldc, work, lwork)
	} else if nq > 1 {
		mi := m
		ni := n - 1
		i1 := 0
		i2 := 1
		if left {
			mi = m - 1
			ni = n
			i1 = 1
			i2 = 0
		}
		impl.Dormlq(side, transt, mi, ni, nq-1, a[1:], lda, tau, c[i1*ldc+i2:], ldc, work, lwork)
	}
}
