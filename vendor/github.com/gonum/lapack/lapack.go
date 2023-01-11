// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lapack

import "github.com/gonum/blas"

const None = 'N'

type Job byte

// CompSV determines if the singular values are to be computed in compact form.
type CompSV byte

const (
	Compact  CompSV = 'P'
	Explicit CompSV = 'I'
)

// Complex128 defines the public complex128 LAPACK API supported by gonum/lapack.
type Complex128 interface{}

// Float64 defines the public float64 LAPACK API supported by gonum/lapack.
type Float64 interface {
	Dgecon(norm MatrixNorm, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64
	Dgels(trans blas.Transpose, m, n, nrhs int, a []float64, lda int, b []float64, ldb int, work []float64, lwork int) bool
	Dgelqf(m, n int, a []float64, lda int, tau, work []float64, lwork int)
	Dgeqrf(m, n int, a []float64, lda int, tau, work []float64, lwork int)
	Dgesvd(jobU, jobVT SVDJob, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, lwork int) (ok bool)
	Dgetrf(m, n int, a []float64, lda int, ipiv []int) (ok bool)
	Dgetri(n int, a []float64, lda int, ipiv []int, work []float64, lwork int) (ok bool)
	Dgetrs(trans blas.Transpose, n, nrhs int, a []float64, lda int, ipiv []int, b []float64, ldb int)
	Dlantr(norm MatrixNorm, uplo blas.Uplo, diag blas.Diag, m, n int, a []float64, lda int, work []float64) float64
	Dlange(norm MatrixNorm, m, n int, a []float64, lda int, work []float64) float64
	Dlansy(norm MatrixNorm, uplo blas.Uplo, n int, a []float64, lda int, work []float64) float64
	Dormqr(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int)
	Dormlq(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64, lwork int)
	Dpocon(uplo blas.Uplo, n int, a []float64, lda int, anorm float64, work []float64, iwork []int) float64
	Dpotrf(ul blas.Uplo, n int, a []float64, lda int) (ok bool)
	Dsyev(jobz EigComp, uplo blas.Uplo, n int, a []float64, lda int, w, work []float64, lwork int) (ok bool)
	Dtrcon(norm MatrixNorm, uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int, work []float64, iwork []int) float64
	Dtrtri(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int) (ok bool)
	Dtrtrs(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, nrhs int, a []float64, lda int, b []float64, ldb int) (ok bool)
}

// Direct specifies the direction of the multiplication for the Householder matrix.
type Direct byte

const (
	Forward  Direct = 'F' // Reflectors are right-multiplied, H_0 * H_1 * ... * H_{k-1}.
	Backward Direct = 'B' // Reflectors are left-multiplied, H_{k-1} * ... * H_1 * H_0.
)

// Sort is the sorting order.
type Sort byte

const (
	SortIncreasing Sort = 'I'
	SortDecreasing Sort = 'D'
)

// StoreV indicates the storage direction of elementary reflectors.
type StoreV byte

const (
	ColumnWise StoreV = 'C' // Reflector stored in a column of the matrix.
	RowWise    StoreV = 'R' // Reflector stored in a row of the matrix.
)

// MatrixNorm represents the kind of matrix norm to compute.
type MatrixNorm byte

const (
	MaxAbs       MatrixNorm = 'M' // max(abs(A(i,j)))  ('M')
	MaxColumnSum MatrixNorm = 'O' // Maximum column sum (one norm) ('1', 'O')
	MaxRowSum    MatrixNorm = 'I' // Maximum row sum (infinity norm) ('I', 'i')
	NormFrob     MatrixNorm = 'F' // Frobenius norm (sqrt of sum of squares) ('F', 'f', E, 'e')
)

// MatrixType represents the kind of matrix represented in the data.
type MatrixType byte

const (
	General  MatrixType = 'G' // A dense matrix (like blas64.General).
	UpperTri MatrixType = 'U' // An upper triangular matrix.
	LowerTri MatrixType = 'L' // A lower triangular matrix.
)

// Pivot specifies the pivot type for plane rotations
type Pivot byte

const (
	Variable Pivot = 'V'
	Top      Pivot = 'T'
	Bottom   Pivot = 'B'
)

type DecompUpdate byte

const (
	ApplyP DecompUpdate = 'P'
	ApplyQ DecompUpdate = 'Q'
)

// SVDJob specifies the singular vector computation type for SVD.
type SVDJob byte

const (
	SVDAll       SVDJob = 'A' // Compute all singular vectors
	SVDInPlace          = 'S' // Compute the first singular vectors and store them in provided storage.
	SVDOverwrite        = 'O' // Compute the singular vectors and store them in input matrix
	SVDNone             = 'N' // Do not compute singular vectors
)

// EigComp specifies the type of eigenvalue decomposition.
type EigComp byte

const (
	// EigValueOnly specifies to compute only the eigenvalues of the input matrix.
	EigValueOnly EigComp = 'N'
	// EigDecomp specifies to compute the eigenvalues and eigenvectors of the
	// full symmetric matrix.
	EigDecomp = 'V'
	// EigBoth specifies to compute both the eigenvalues and eigenvectors of the
	// input tridiagonal matrix.
	EigBoth = 'I'
)
