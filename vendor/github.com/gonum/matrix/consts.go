// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matrix

// TriKind represents the triangularity of the matrix.
type TriKind bool

const (
	// Upper specifies an upper triangular matrix.
	Upper TriKind = true
	// Lower specifies a lower triangular matrix.
	Lower TriKind = false
)

// SVDKind specifies the treatment of singular vectors during the
// factorization.
type SVDKind int

const (
	// SVDNone specifies that no singular vectors should be computed during
	// the decomposition.
	SVDNone SVDKind = iota + 1
	// SVDThin computes the thin singular vectors, that is, it computes
	//  A = U~ * Σ * V~^T
	// where U~ is of size m×min(m,n), Σ is a diagonal matrix of size min(m,n)×min(m,n)
	// and V~ is of size n×min(m,n).
	SVDThin
	// SVDFull computes the full singular value decomposition,
	//  A = U * Σ * V^T
	// where U is of size m×m, Σ is an m×n diagonal matrix, and V is an n×n matrix.
	SVDFull
)
