// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package blas32 provides a simple interface to the float32 BLAS API.
package blas32

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/native"
)

var blas32 blas.Float32 = native.Implementation{}

// Use sets the BLAS float32 implementation to be used by subsequent BLAS calls.
// The default implementation is native.Implementation.
func Use(b blas.Float32) {
	blas32 = b
}

// Implementation returns the current BLAS float32 implementation.
//
// Implementation allows direct calls to the current the BLAS float32 implementation
// giving finer control of parameters.
func Implementation() blas.Float32 {
	return blas32
}

// Vector represents a vector with an associated element increment.
type Vector struct {
	Inc  int
	Data []float32
}

// General represents a matrix using the conventional storage scheme.
type General struct {
	Rows, Cols int
	Stride     int
	Data       []float32
}

// Band represents a band matrix using the band storage scheme.
type Band struct {
	Rows, Cols int
	KL, KU     int
	Stride     int
	Data       []float32
}

// Triangular represents a triangular matrix using the conventional storage scheme.
type Triangular struct {
	N      int
	Stride int
	Data   []float32
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularBand represents a triangular matrix using the band storage scheme.
type TriangularBand struct {
	N, K   int
	Stride int
	Data   []float32
	Uplo   blas.Uplo
	Diag   blas.Diag
}

// TriangularPacked represents a triangular matrix using the packed storage scheme.
type TriangularPacked struct {
	N    int
	Data []float32
	Uplo blas.Uplo
	Diag blas.Diag
}

// Symmetric represents a symmetric matrix using the conventional storage scheme.
type Symmetric struct {
	N      int
	Stride int
	Data   []float32
	Uplo   blas.Uplo
}

// SymmetricBand represents a symmetric matrix using the band storage scheme.
type SymmetricBand struct {
	N, K   int
	Stride int
	Data   []float32
	Uplo   blas.Uplo
}

// SymmetricPacked represents a symmetric matrix using the packed storage scheme.
type SymmetricPacked struct {
	N    int
	Data []float32
	Uplo blas.Uplo
}

// Level 1

const negInc = "blas32: negative vector increment"

// Dot computes the dot product of the two vectors:
//  \sum_i x[i]*y[i].
func Dot(n int, x, y Vector) float32 {
	return blas32.Sdot(n, x.Data, x.Inc, y.Data, y.Inc)
}

// DDot computes the dot product of the two vectors:
//  \sum_i x[i]*y[i].
func DDot(n int, x, y Vector) float64 {
	return blas32.Dsdot(n, x.Data, x.Inc, y.Data, y.Inc)
}

// SDDot computes the dot product of the two vectors adding a constant:
//  alpha + \sum_i x[i]*y[i].
func SDDot(n int, alpha float32, x, y Vector) float32 {
	return blas32.Sdsdot(n, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

// Nrm2 computes the Euclidean norm of the vector x:
//  sqrt(\sum_i x[i]*x[i]).
//
// Nrm2 will panic if the vector increment is negative.
func Nrm2(n int, x Vector) float32 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas32.Snrm2(n, x.Data, x.Inc)
}

// Asum computes the sum of the absolute values of the elements of x:
//  \sum_i |x[i]|.
//
// Asum will panic if the vector increment is negative.
func Asum(n int, x Vector) float32 {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas32.Sasum(n, x.Data, x.Inc)
}

// Iamax returns the index of an element of x with the largest absolute value.
// If there are multiple such indices the earliest is returned.
// Iamax returns -1 if n == 0.
//
// Iamax will panic if the vector increment is negative.
func Iamax(n int, x Vector) int {
	if x.Inc < 0 {
		panic(negInc)
	}
	return blas32.Isamax(n, x.Data, x.Inc)
}

// Swap exchanges the elements of the two vectors:
//  x[i], y[i] = y[i], x[i] for all i.
func Swap(n int, x, y Vector) {
	blas32.Sswap(n, x.Data, x.Inc, y.Data, y.Inc)
}

// Copy copies the elements of x into the elements of y:
//  y[i] = x[i] for all i.
func Copy(n int, x, y Vector) {
	blas32.Scopy(n, x.Data, x.Inc, y.Data, y.Inc)
}

// Axpy adds x scaled by alpha to y:
//  y[i] += alpha*x[i] for all i.
func Axpy(n int, alpha float32, x, y Vector) {
	blas32.Saxpy(n, alpha, x.Data, x.Inc, y.Data, y.Inc)
}

// Rotg computes the parameters of a Givens plane rotation so that
//  ⎡ c s⎤   ⎡a⎤   ⎡r⎤
//  ⎣-s c⎦ * ⎣b⎦ = ⎣0⎦
// where a and b are the Cartesian coordinates of a given point.
// c, s, and r are defined as
//  r = ±Sqrt(a^2 + b^2),
//  c = a/r, the cosine of the rotation angle,
//  s = a/r, the sine of the rotation angle,
// and z is defined such that
//  if |a| > |b|,        z = s,
//  otherwise if c != 0, z = 1/c,
//  otherwise            z = 1.
func Rotg(a, b float32) (c, s, r, z float32) {
	return blas32.Srotg(a, b)
}

// Rotmg computes the modified Givens rotation. See
// http://www.netlib.org/lapack/explore-html/df/deb/drotmg_8f.html
// for more details.
func Rotmg(d1, d2, b1, b2 float32) (p blas.SrotmParams, rd1, rd2, rb1 float32) {
	return blas32.Srotmg(d1, d2, b1, b2)
}

// Rot applies a plane transformation to n points represented by the vectors x
// and y:
//  x[i] =  c*x[i] + s*y[i],
//  y[i] = -s*x[i] + c*y[i], for all i.
func Rot(n int, x, y Vector, c, s float32) {
	blas32.Srot(n, x.Data, x.Inc, y.Data, y.Inc, c, s)
}

// Rotm applies the modified Givens rotation to n points represented by the
// vectors x and y.
func Rotm(n int, x, y Vector, p blas.SrotmParams) {
	blas32.Srotm(n, x.Data, x.Inc, y.Data, y.Inc, p)
}

// Scal scales the vector x by alpha:
//  x[i] *= alpha for all i.
//
// Scal will panic if the vector increment is negative.
func Scal(n int, alpha float32, x Vector) {
	if x.Inc < 0 {
		panic(negInc)
	}
	blas32.Sscal(n, alpha, x.Data, x.Inc)
}

// Level 2

// Gemv computes
//  y = alpha * A * x + beta * y,   if t == blas.NoTrans,
//  y = alpha * A^T * x + beta * y, if t == blas.Trans or blas.ConjTrans,
// where A is an m×n dense matrix, x and y are vectors, and alpha and beta are scalars.
func Gemv(t blas.Transpose, alpha float32, a General, x Vector, beta float32, y Vector) {
	blas32.Sgemv(t, a.Rows, a.Cols, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

// Gbmv computes
//  y = alpha * A * x + beta * y,   if t == blas.NoTrans,
//  y = alpha * A^T * x + beta * y, if t == blas.Trans or blas.ConjTrans,
// where A is an m×n band matrix, x and y are vectors, and alpha and beta are scalars.
func Gbmv(t blas.Transpose, alpha float32, a Band, x Vector, beta float32, y Vector) {
	blas32.Sgbmv(t, a.Rows, a.Cols, a.KL, a.KU, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

// Trmv computes
//  x = A * x,   if t == blas.NoTrans,
//  x = A^T * x, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular matrix, and x is a vector.
func Trmv(t blas.Transpose, a Triangular, x Vector) {
	blas32.Strmv(a.Uplo, t, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

// Tbmv computes
//  x = A * x,   if t == blas.NoTrans,
//  x = A^T * x, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular band matrix, and x is a vector.
func Tbmv(t blas.Transpose, a TriangularBand, x Vector) {
	blas32.Stbmv(a.Uplo, t, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

// Tpmv computes
//  x = A * x,   if t == blas.NoTrans,
//  x = A^T * x, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular matrix in packed format, and x is a vector.
func Tpmv(t blas.Transpose, a TriangularPacked, x Vector) {
	blas32.Stpmv(a.Uplo, t, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

// Trsv solves
//  A * x = b,   if t == blas.NoTrans,
//  A^T * x = b, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular matrix, and x and b are vectors.
//
// At entry to the function, x contains the values of b, and the result is
// stored in-place into x.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func Trsv(t blas.Transpose, a Triangular, x Vector) {
	blas32.Strsv(a.Uplo, t, a.Diag, a.N, a.Data, a.Stride, x.Data, x.Inc)
}

// Tbsv solves
//  A * x = b,   if t == blas.NoTrans,
//  A^T * x = b, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular band matrix, and x and b are vectors.
//
// At entry to the function, x contains the values of b, and the result is
// stored in place into x.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func Tbsv(t blas.Transpose, a TriangularBand, x Vector) {
	blas32.Stbsv(a.Uplo, t, a.Diag, a.N, a.K, a.Data, a.Stride, x.Data, x.Inc)
}

// Tpsv solves
//  A * x = b,   if t == blas.NoTrans,
//  A^T * x = b, if t == blas.Trans or blas.ConjTrans,
// where A is an n×n triangular matrix in packed format, and x and b are
// vectors.
//
// At entry to the function, x contains the values of b, and the result is
// stored in place into x.
//
// No test for singularity or near-singularity is included in this
// routine. Such tests must be performed before calling this routine.
func Tpsv(t blas.Transpose, a TriangularPacked, x Vector) {
	blas32.Stpsv(a.Uplo, t, a.Diag, a.N, a.Data, x.Data, x.Inc)
}

// Symv computes
//    y = alpha * A * x + beta * y,
// where A is an n×n symmetric matrix, x and y are vectors, and alpha and
// beta are scalars.
func Symv(alpha float32, a Symmetric, x Vector, beta float32, y Vector) {
	blas32.Ssymv(a.Uplo, a.N, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

// Sbmv performs
//  y = alpha * A * x + beta * y,
// where A is an n×n symmetric band matrix, x and y are vectors, and alpha
// and beta are scalars.
func Sbmv(alpha float32, a SymmetricBand, x Vector, beta float32, y Vector) {
	blas32.Ssbmv(a.Uplo, a.N, a.K, alpha, a.Data, a.Stride, x.Data, x.Inc, beta, y.Data, y.Inc)
}

// Spmv performs
//    y = alpha * A * x + beta * y,
// where A is an n×n symmetric matrix in packed format, x and y are vectors,
// and alpha and beta are scalars.
func Spmv(alpha float32, a SymmetricPacked, x Vector, beta float32, y Vector) {
	blas32.Sspmv(a.Uplo, a.N, alpha, a.Data, x.Data, x.Inc, beta, y.Data, y.Inc)
}

// Ger performs a rank-1 update
//  A += alpha * x * y^T,
// where A is an m×n dense matrix, x and y are vectors, and alpha is a scalar.
func Ger(alpha float32, x, y Vector, a General) {
	blas32.Sger(a.Rows, a.Cols, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

// Syr performs a rank-1 update
//  A += alpha * x * x^T,
// where A is an n×n symmetric matrix, x is a vector, and alpha is a scalar.
func Syr(alpha float32, x Vector, a Symmetric) {
	blas32.Ssyr(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data, a.Stride)
}

// Spr performs the rank-1 update
//  A += alpha * x * x^T,
// where A is an n×n symmetric matrix in packed format, x is a vector, and
// alpha is a scalar.
func Spr(alpha float32, x Vector, a SymmetricPacked) {
	blas32.Sspr(a.Uplo, a.N, alpha, x.Data, x.Inc, a.Data)
}

// Syr2 performs a rank-2 update
//  A += alpha * x * y^T + alpha * y * x^T,
// where A is a symmetric n×n matrix, x and y are vectors, and alpha is a scalar.
func Syr2(alpha float32, x, y Vector, a Symmetric) {
	blas32.Ssyr2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data, a.Stride)
}

// Spr2 performs a rank-2 update
//  A += alpha * x * y^T + alpha * y * x^T,
// where A is an n×n symmetric matrix in packed format, x and y are vectors,
// and alpha is a scalar.
func Spr2(alpha float32, x, y Vector, a SymmetricPacked) {
	blas32.Sspr2(a.Uplo, a.N, alpha, x.Data, x.Inc, y.Data, y.Inc, a.Data)
}

// Level 3

// Gemm computes
//  C = alpha * A * B + beta * C,
// where A, B, and C are dense matrices, and alpha and beta are scalars.
// tA and tB specify whether A or B are transposed.
func Gemm(tA, tB blas.Transpose, alpha float32, a, b General, beta float32, c General) {
	var m, n, k int
	if tA == blas.NoTrans {
		m, k = a.Rows, a.Cols
	} else {
		m, k = a.Cols, a.Rows
	}
	if tB == blas.NoTrans {
		n = b.Cols
	} else {
		n = b.Rows
	}
	blas32.Sgemm(tA, tB, m, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

// Symm performs
//  C = alpha * A * B + beta * C, if s == blas.Left,
//  C = alpha * B * A + beta * C, if s == blas.Right,
// where A is an n×n or m×m symmetric matrix, B and C are m×n matrices, and
// alpha is a scalar.
func Symm(s blas.Side, alpha float32, a Symmetric, b General, beta float32, c General) {
	var m, n int
	if s == blas.Left {
		m, n = a.N, b.Cols
	} else {
		m, n = b.Rows, a.N
	}
	blas32.Ssymm(s, a.Uplo, m, n, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

// Syrk performs a symmetric rank-k update
//  C = alpha * A * A^T + beta * C, if t == blas.NoTrans,
//  C = alpha * A^T * A + beta * C, if t == blas.Trans or blas.ConjTrans,
// where C is an n×n symmetric matrix, A is an n×k matrix if t == blas.NoTrans and
// a k×n matrix otherwise, and alpha and beta are scalars.
func Syrk(t blas.Transpose, alpha float32, a General, beta float32, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	blas32.Ssyrk(c.Uplo, t, n, k, alpha, a.Data, a.Stride, beta, c.Data, c.Stride)
}

// Syr2k performs a symmetric rank-2k update
//  C = alpha * A * B^T + alpha * B * A^T + beta * C, if t == blas.NoTrans,
//  C = alpha * A^T * B + alpha * B^T * A + beta * C, if t == blas.Trans or blas.ConjTrans,
// where C is an n×n symmetric matrix, A and B are n×k matrices if t == NoTrans
// and k×n matrices otherwise, and alpha and beta are scalars.
func Syr2k(t blas.Transpose, alpha float32, a, b General, beta float32, c Symmetric) {
	var n, k int
	if t == blas.NoTrans {
		n, k = a.Rows, a.Cols
	} else {
		n, k = a.Cols, a.Rows
	}
	blas32.Ssyr2k(c.Uplo, t, n, k, alpha, a.Data, a.Stride, b.Data, b.Stride, beta, c.Data, c.Stride)
}

// Trmm performs
//  B = alpha * A * B,   if tA == blas.NoTrans and s == blas.Left,
//  B = alpha * A^T * B, if tA == blas.Trans or blas.ConjTrans, and s == blas.Left,
//  B = alpha * B * A,   if tA == blas.NoTrans and s == blas.Right,
//  B = alpha * B * A^T, if tA == blas.Trans or blas.ConjTrans, and s == blas.Right,
// where A is an n×n or m×m triangular matrix, B is an m×n matrix, and alpha is
// a scalar.
func Trmm(s blas.Side, tA blas.Transpose, alpha float32, a Triangular, b General) {
	blas32.Strmm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}

// Trsm solves
//  A * X = alpha * B,   if tA == blas.NoTrans and s == blas.Left,
//  A^T * X = alpha * B, if tA == blas.Trans or blas.ConjTrans, and s == blas.Left,
//  X * A = alpha * B,   if tA == blas.NoTrans and s == blas.Right,
//  X * A^T = alpha * B, if tA == blas.Trans or blas.ConjTrans, and s == blas.Right,
// where A is an n×n or m×m triangular matrix, X and B are m×n matrices, and
// alpha is a scalar.
//
// At entry to the function, X contains the values of B, and the result is
// stored in-place into X.
//
// No check is made that A is invertible.
func Trsm(s blas.Side, tA blas.Transpose, alpha float32, a Triangular, b General) {
	blas32.Strsm(s, a.Uplo, tA, a.Diag, b.Rows, b.Cols, alpha, a.Data, a.Stride, b.Data, b.Stride)
}
