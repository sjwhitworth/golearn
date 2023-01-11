// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package blas provides interfaces for the BLAS linear algebra standard.

All methods must perform appropriate parameter checking and panic if
provided parameters that do not conform to the requirements specified
by the BLAS standard.

Quick Reference Guide to the BLAS from http://www.netlib.org/lapack/lug/node145.html

This version is modified to remove the "order" option. All matrix operations are
on row-order matrices.

Level 1 BLAS

	        dim scalar vector   vector   scalars              5-element prefixes
	                                                          struct

	_rotg (                                      a, b )                S, D
	_rotmg(                              d1, d2, a, b )                S, D
	_rot  ( n,         x, incX, y, incY,               c, s )          S, D
	_rotm ( n,         x, incX, y, incY,                      param )  S, D
	_swap ( n,         x, incX, y, incY )                              S, D, C, Z
	_scal ( n,  alpha, x, incX )                                       S, D, C, Z, Cs, Zd
	_copy ( n,         x, incX, y, incY )                              S, D, C, Z
	_axpy ( n,  alpha, x, incX, y, incY )                              S, D, C, Z
	_dot  ( n,         x, incX, y, incY )                              S, D, Ds
	_dotu ( n,         x, incX, y, incY )                              C, Z
	_dotc ( n,         x, incX, y, incY )                              C, Z
	__dot ( n,  alpha, x, incX, y, incY )                              Sds
	_nrm2 ( n,         x, incX )                                       S, D, Sc, Dz
	_asum ( n,         x, incX )                                       S, D, Sc, Dz
	I_amax( n,         x, incX )                                       s, d, c, z

Level 2 BLAS

	        options                   dim   b-width scalar matrix  vector   scalar vector   prefixes

	_gemv (        trans,      m, n,         alpha, a, lda, x, incX, beta,  y, incY ) S, D, C, Z
	_gbmv (        trans,      m, n, kL, kU, alpha, a, lda, x, incX, beta,  y, incY ) S, D, C, Z
	_hemv ( uplo,                 n,         alpha, a, lda, x, incX, beta,  y, incY ) C, Z
	_hbmv ( uplo,                 n, k,      alpha, a, lda, x, incX, beta,  y, incY ) C, Z
	_hpmv ( uplo,                 n,         alpha, ap,     x, incX, beta,  y, incY ) C, Z
	_symv ( uplo,                 n,         alpha, a, lda, x, incX, beta,  y, incY ) S, D
	_sbmv ( uplo,                 n, k,      alpha, a, lda, x, incX, beta,  y, incY ) S, D
	_spmv ( uplo,                 n,         alpha, ap,     x, incX, beta,  y, incY ) S, D
	_trmv ( uplo, trans, diag,    n,                a, lda, x, incX )                 S, D, C, Z
	_tbmv ( uplo, trans, diag,    n, k,             a, lda, x, incX )                 S, D, C, Z
	_tpmv ( uplo, trans, diag,    n,                ap,     x, incX )                 S, D, C, Z
	_trsv ( uplo, trans, diag,    n,                a, lda, x, incX )                 S, D, C, Z
	_tbsv ( uplo, trans, diag,    n, k,             a, lda, x, incX )                 S, D, C, Z
	_tpsv ( uplo, trans, diag,    n,                ap,     x, incX )                 S, D, C, Z

	        options                   dim   scalar vector   vector   matrix  prefixes

	_ger  (                    m, n, alpha, x, incX, y, incY, a, lda ) S, D
	_geru (                    m, n, alpha, x, incX, y, incY, a, lda ) C, Z
	_gerc (                    m, n, alpha, x, incX, y, incY, a, lda ) C, Z
	_her  ( uplo,                 n, alpha, x, incX,          a, lda ) C, Z
	_hpr  ( uplo,                 n, alpha, x, incX,          ap )     C, Z
	_her2 ( uplo,                 n, alpha, x, incX, y, incY, a, lda ) C, Z
	_hpr2 ( uplo,                 n, alpha, x, incX, y, incY, ap )     C, Z
	_syr  ( uplo,                 n, alpha, x, incX,          a, lda ) S, D
	_spr  ( uplo,                 n, alpha, x, incX,          ap )     S, D
	_syr2 ( uplo,                 n, alpha, x, incX, y, incY, a, lda ) S, D
	_spr2 ( uplo,                 n, alpha, x, incX, y, incY, ap )     S, D

Level 3 BLAS

	        options                                 dim      scalar matrix  matrix  scalar matrix  prefixes

	_gemm (             transA, transB,      m, n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_symm ( side, uplo,                      m, n,    alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_hemm ( side, uplo,                      m, n,    alpha, a, lda, b, ldb, beta,  c, ldc ) C, Z
	_syrk (       uplo, trans,                  n, k, alpha, a, lda,         beta,  c, ldc ) S, D, C, Z
	_herk (       uplo, trans,                  n, k, alpha, a, lda,         beta,  c, ldc ) C, Z
	_syr2k(       uplo, trans,                  n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) S, D, C, Z
	_her2k(       uplo, trans,                  n, k, alpha, a, lda, b, ldb, beta,  c, ldc ) C, Z
	_trmm ( side, uplo, transA,        diag, m, n,    alpha, a, lda, b, ldb )                S, D, C, Z
	_trsm ( side, uplo, transA,        diag, m, n,    alpha, a, lda, b, ldb )                S, D, C, Z

Meaning of prefixes

	S - float32	C - complex64
	D - float64	Z - complex128

Matrix types

	GE - GEneral 		GB - General Band
	SY - SYmmetric 		SB - Symmetric Band 	SP - Symmetric Packed
	HE - HErmitian 		HB - Hermitian Band 	HP - Hermitian Packed
	TR - TRiangular 	TB - Triangular Band 	TP - Triangular Packed

Options

	trans 	= NoTrans, Trans, ConjTrans
	uplo 	= Upper, Lower
	diag 	= Nonunit, Unit
	side 	= Left, Right (A or op(A) on the left, or A or op(A) on the right)

For real matrices, Trans and ConjTrans have the same meaning.
For Hermitian matrices, trans = Trans is not allowed.
For complex symmetric matrices, trans = ConjTrans is not allowed.
*/
package blas

// Flag constants indicate Givens transformation H matrix state.
type Flag int

const (
	Identity    Flag = iota - 2 // H is the identity matrix; no rotation is needed.
	Rescaling                   // H specifies rescaling.
	OffDiagonal                 // Off-diagonal elements of H are units.
	Diagonal                    // Diagonal elements of H are units.
)

// SrotmParams contains Givens transformation parameters returned
// by the Float32 Srotm method.
type SrotmParams struct {
	Flag
	H [4]float32 // Column-major 2 by 2 matrix.
}

// DrotmParams contains Givens transformation parameters returned
// by the Float64 Drotm method.
type DrotmParams struct {
	Flag
	H [4]float64 // Column-major 2 by 2 matrix.
}

// Transpose is used to specify the transposition operation for a
// routine.
type Transpose int

const (
	NoTrans Transpose = 111 + iota
	Trans
	ConjTrans
)

// Uplo is used to specify whether the matrix is an upper or lower
// triangular matrix.
type Uplo int

const (
	All Uplo = 120 + iota
	Upper
	Lower
)

// Diag is used to specify whether the matrix is a unit or non-unit
// triangular matrix.
type Diag int

const (
	NonUnit Diag = 131 + iota
	Unit
)

// Side is used to specify from which side a multiplication operation
// is performed.
type Side int

const (
	Left Side = 141 + iota
	Right
)

// Float32 implements the single precision real BLAS routines.
type Float32 interface {
	Float32Level1
	Float32Level2
	Float32Level3
}

// Float32Level1 implements the single precision real BLAS Level 1 routines.
type Float32Level1 interface {
	Sdsdot(n int, alpha float32, x []float32, incX int, y []float32, incY int) float32
	Dsdot(n int, x []float32, incX int, y []float32, incY int) float64
	Sdot(n int, x []float32, incX int, y []float32, incY int) float32
	Snrm2(n int, x []float32, incX int) float32
	Sasum(n int, x []float32, incX int) float32
	Isamax(n int, x []float32, incX int) int
	Sswap(n int, x []float32, incX int, y []float32, incY int)
	Scopy(n int, x []float32, incX int, y []float32, incY int)
	Saxpy(n int, alpha float32, x []float32, incX int, y []float32, incY int)
	Srotg(a, b float32) (c, s, r, z float32)
	Srotmg(d1, d2, b1, b2 float32) (p SrotmParams, rd1, rd2, rb1 float32)
	Srot(n int, x []float32, incX int, y []float32, incY int, c, s float32)
	Srotm(n int, x []float32, incX int, y []float32, incY int, p SrotmParams)
	Sscal(n int, alpha float32, x []float32, incX int)
}

// Float32Level2 implements the single precision real BLAS Level 2 routines.
type Float32Level2 interface {
	Sgemv(tA Transpose, m, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Sgbmv(tA Transpose, m, n, kL, kU int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Strmv(ul Uplo, tA Transpose, d Diag, n int, a []float32, lda int, x []float32, incX int)
	Stbmv(ul Uplo, tA Transpose, d Diag, n, k int, a []float32, lda int, x []float32, incX int)
	Stpmv(ul Uplo, tA Transpose, d Diag, n int, ap []float32, x []float32, incX int)
	Strsv(ul Uplo, tA Transpose, d Diag, n int, a []float32, lda int, x []float32, incX int)
	Stbsv(ul Uplo, tA Transpose, d Diag, n, k int, a []float32, lda int, x []float32, incX int)
	Stpsv(ul Uplo, tA Transpose, d Diag, n int, ap []float32, x []float32, incX int)
	Ssymv(ul Uplo, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Ssbmv(ul Uplo, n, k int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int)
	Sspmv(ul Uplo, n int, alpha float32, ap []float32, x []float32, incX int, beta float32, y []float32, incY int)
	Sger(m, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int)
	Ssyr(ul Uplo, n int, alpha float32, x []float32, incX int, a []float32, lda int)
	Sspr(ul Uplo, n int, alpha float32, x []float32, incX int, ap []float32)
	Ssyr2(ul Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int)
	Sspr2(ul Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32)
}

// Float32Level3 implements the single precision real BLAS Level 3 routines.
type Float32Level3 interface {
	Sgemm(tA, tB Transpose, m, n, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Ssymm(s Side, ul Uplo, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Ssyrk(ul Uplo, t Transpose, n, k int, alpha float32, a []float32, lda int, beta float32, c []float32, ldc int)
	Ssyr2k(ul Uplo, t Transpose, n, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int)
	Strmm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int)
	Strsm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float32, a []float32, lda int, b []float32, ldb int)
}

// Float64 implements the single precision real BLAS routines.
type Float64 interface {
	Float64Level1
	Float64Level2
	Float64Level3
}

// Float64Level1 implements the double precision real BLAS Level 1 routines.
type Float64Level1 interface {
	Ddot(n int, x []float64, incX int, y []float64, incY int) float64
	Dnrm2(n int, x []float64, incX int) float64
	Dasum(n int, x []float64, incX int) float64
	Idamax(n int, x []float64, incX int) int
	Dswap(n int, x []float64, incX int, y []float64, incY int)
	Dcopy(n int, x []float64, incX int, y []float64, incY int)
	Daxpy(n int, alpha float64, x []float64, incX int, y []float64, incY int)
	Drotg(a, b float64) (c, s, r, z float64)
	Drotmg(d1, d2, b1, b2 float64) (p DrotmParams, rd1, rd2, rb1 float64)
	Drot(n int, x []float64, incX int, y []float64, incY int, c float64, s float64)
	Drotm(n int, x []float64, incX int, y []float64, incY int, p DrotmParams)
	Dscal(n int, alpha float64, x []float64, incX int)
}

// Float64Level2 implements the double precision real BLAS Level 2 routines.
type Float64Level2 interface {
	Dgemv(tA Transpose, m, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dgbmv(tA Transpose, m, n, kL, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dtrmv(ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
	Dtbmv(ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtpmv(ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
	Dtrsv(ul Uplo, tA Transpose, d Diag, n int, a []float64, lda int, x []float64, incX int)
	Dtbsv(ul Uplo, tA Transpose, d Diag, n, k int, a []float64, lda int, x []float64, incX int)
	Dtpsv(ul Uplo, tA Transpose, d Diag, n int, ap []float64, x []float64, incX int)
	Dsymv(ul Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dsbmv(ul Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
	Dspmv(ul Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int)
	Dger(m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
	Dsyr(ul Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int)
	Dspr(ul Uplo, n int, alpha float64, x []float64, incX int, ap []float64)
	Dsyr2(ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
	Dspr2(ul Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64)
}

// Float64Level3 implements the double precision real BLAS Level 3 routines.
type Float64Level3 interface {
	Dgemm(tA, tB Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dsymm(s Side, ul Uplo, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dsyrk(ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int)
	Dsyr2k(ul Uplo, t Transpose, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int)
	Dtrmm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
	Dtrsm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha float64, a []float64, lda int, b []float64, ldb int)
}

// Complex64 implements the single precision complex BLAS routines.
type Complex64 interface {
	Complex64Level1
	Complex64Level2
	Complex64Level3
}

// Complex64Level1 implements the single precision complex BLAS Level 1 routines.
type Complex64Level1 interface {
	Cdotu(n int, x []complex64, incX int, y []complex64, incY int) (dotu complex64)
	Cdotc(n int, x []complex64, incX int, y []complex64, incY int) (dotc complex64)
	Scnrm2(n int, x []complex64, incX int) float32
	Scasum(n int, x []complex64, incX int) float32
	Icamax(n int, x []complex64, incX int) int
	Cswap(n int, x []complex64, incX int, y []complex64, incY int)
	Ccopy(n int, x []complex64, incX int, y []complex64, incY int)
	Caxpy(n int, alpha complex64, x []complex64, incX int, y []complex64, incY int)
	Cscal(n int, alpha complex64, x []complex64, incX int)
	Csscal(n int, alpha float32, x []complex64, incX int)
}

// Complex64Level2 implements the single precision complex BLAS routines Level 2 routines.
type Complex64Level2 interface {
	Cgemv(tA Transpose, m, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Cgbmv(tA Transpose, m, n, kL, kU int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Ctrmv(ul Uplo, tA Transpose, d Diag, n int, a []complex64, lda int, x []complex64, incX int)
	Ctbmv(ul Uplo, tA Transpose, d Diag, n, k int, a []complex64, lda int, x []complex64, incX int)
	Ctpmv(ul Uplo, tA Transpose, d Diag, n int, ap []complex64, x []complex64, incX int)
	Ctrsv(ul Uplo, tA Transpose, d Diag, n int, a []complex64, lda int, x []complex64, incX int)
	Ctbsv(ul Uplo, tA Transpose, d Diag, n, k int, a []complex64, lda int, x []complex64, incX int)
	Ctpsv(ul Uplo, tA Transpose, d Diag, n int, ap []complex64, x []complex64, incX int)
	Chemv(ul Uplo, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Chbmv(ul Uplo, n, k int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Chpmv(ul Uplo, n int, alpha complex64, ap []complex64, x []complex64, incX int, beta complex64, y []complex64, incY int)
	Cgeru(m, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Cgerc(m, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Cher(ul Uplo, n int, alpha float32, x []complex64, incX int, a []complex64, lda int)
	Chpr(ul Uplo, n int, alpha float32, x []complex64, incX int, a []complex64)
	Cher2(ul Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int)
	Chpr2(ul Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, ap []complex64)
}

// Complex64Level3 implements the single precision complex BLAS Level 3 routines.
type Complex64Level3 interface {
	Cgemm(tA, tB Transpose, m, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Csymm(s Side, ul Uplo, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Csyrk(ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, beta complex64, c []complex64, ldc int)
	Csyr2k(ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Ctrmm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int)
	Ctrsm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int)
	Chemm(s Side, ul Uplo, m, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int)
	Cherk(ul Uplo, t Transpose, n, k int, alpha float32, a []complex64, lda int, beta float32, c []complex64, ldc int)
	Cher2k(ul Uplo, t Transpose, n, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta float32, c []complex64, ldc int)
}

// Complex128 implements the double precision complex BLAS routines.
type Complex128 interface {
	Complex128Level1
	Complex128Level2
	Complex128Level3
}

// Complex128Level1 implements the double precision complex BLAS Level 1 routines.
type Complex128Level1 interface {
	Zdotu(n int, x []complex128, incX int, y []complex128, incY int) (dotu complex128)
	Zdotc(n int, x []complex128, incX int, y []complex128, incY int) (dotc complex128)
	Dznrm2(n int, x []complex128, incX int) float64
	Dzasum(n int, x []complex128, incX int) float64
	Izamax(n int, x []complex128, incX int) int
	Zswap(n int, x []complex128, incX int, y []complex128, incY int)
	Zcopy(n int, x []complex128, incX int, y []complex128, incY int)
	Zaxpy(n int, alpha complex128, x []complex128, incX int, y []complex128, incY int)
	Zscal(n int, alpha complex128, x []complex128, incX int)
	Zdscal(n int, alpha float64, x []complex128, incX int)
}

// Complex128Level2 implements the double precision complex BLAS Level 2 routines.
type Complex128Level2 interface {
	Zgemv(tA Transpose, m, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zgbmv(tA Transpose, m, n int, kL int, kU int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Ztrmv(ul Uplo, tA Transpose, d Diag, n int, a []complex128, lda int, x []complex128, incX int)
	Ztbmv(ul Uplo, tA Transpose, d Diag, n, k int, a []complex128, lda int, x []complex128, incX int)
	Ztpmv(ul Uplo, tA Transpose, d Diag, n int, ap []complex128, x []complex128, incX int)
	Ztrsv(ul Uplo, tA Transpose, d Diag, n int, a []complex128, lda int, x []complex128, incX int)
	Ztbsv(ul Uplo, tA Transpose, d Diag, n, k int, a []complex128, lda int, x []complex128, incX int)
	Ztpsv(ul Uplo, tA Transpose, d Diag, n int, ap []complex128, x []complex128, incX int)
	Zhemv(ul Uplo, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zhbmv(ul Uplo, n, k int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zhpmv(ul Uplo, n int, alpha complex128, ap []complex128, x []complex128, incX int, beta complex128, y []complex128, incY int)
	Zgeru(m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zgerc(m, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zher(ul Uplo, n int, alpha float64, x []complex128, incX int, a []complex128, lda int)
	Zhpr(ul Uplo, n int, alpha float64, x []complex128, incX int, a []complex128)
	Zher2(ul Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int)
	Zhpr2(ul Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, ap []complex128)
}

// Complex128Level3 implements the double precision complex BLAS Level 3 routines.
type Complex128Level3 interface {
	Zgemm(tA, tB Transpose, m, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zsymm(s Side, ul Uplo, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zsyrk(ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int)
	Zsyr2k(ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Ztrmm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int)
	Ztrsm(s Side, ul Uplo, tA Transpose, d Diag, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int)
	Zhemm(s Side, ul Uplo, m, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int)
	Zherk(ul Uplo, t Transpose, n, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int)
	Zher2k(ul Uplo, t Transpose, n, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta float64, c []complex128, ldc int)
}
