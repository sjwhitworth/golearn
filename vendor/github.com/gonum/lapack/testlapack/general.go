// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// nanSlice allocates a new slice of length n filled with NaN.
func nanSlice(n int) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = math.NaN()
	}
	return s
}

// randomSlice allocates a new slice of length n filled with random values.
func randomSlice(n int, rnd *rand.Rand) []float64 {
	s := make([]float64, n)
	for i := range s {
		s[i] = rnd.NormFloat64()
	}
	return s
}

// nanGeneral allocates a new r×c general matrix filled with NaN values.
func nanGeneral(r, c, stride int) blas64.General {
	if r < 0 || c < 0 {
		panic("bad matrix size")
	}
	if r == 0 || c == 0 {
		return blas64.General{Stride: max(1, stride)}
	}
	if stride < c {
		panic("bad stride")
	}
	return blas64.General{
		Rows:   r,
		Cols:   c,
		Stride: stride,
		Data:   nanSlice((r-1)*stride + c),
	}
}

// randomGeneral allocates a new r×c general matrix filled with random
// numbers. Out-of-range elements are filled with NaN values.
func randomGeneral(r, c, stride int, rnd *rand.Rand) blas64.General {
	ans := nanGeneral(r, c, stride)
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// randomHessenberg allocates a new n×n Hessenberg matrix filled with zeros
// under the first subdiagonal and with random numbers elsewhere. Out-of-range
// elements are filled with NaN values.
func randomHessenberg(n, stride int, rnd *rand.Rand) blas64.General {
	ans := nanGeneral(n, n, stride)
	for i := 0; i < n; i++ {
		for j := 0; j < i-1; j++ {
			ans.Data[i*ans.Stride+j] = 0
		}
		for j := max(0, i-1); j < n; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// nanTriangular allocates a new r×c triangular matrix filled with NaN values.
func nanTriangular(uplo blas.Uplo, n, stride int) blas64.Triangular {
	if n < 0 {
		panic("bad matrix size")
	}
	if n == 0 {
		return blas64.Triangular{
			Stride: max(1, stride),
			Uplo:   uplo,
			Diag:   blas.NonUnit,
		}
	}
	if stride < n {
		panic("bad stride")
	}
	return blas64.Triangular{
		N:      n,
		Stride: stride,
		Data:   nanSlice((n-1)*stride + n),
		Uplo:   uplo,
		Diag:   blas.NonUnit,
	}
}

// randomTriangular allocates a new r×c triangular matrix filled with random
// numbers. Out-of-triangle elements are filled with NaN values.
func randomTriangular(uplo blas.Uplo, n, stride int, rnd *rand.Rand) blas64.Triangular {
	ans := nanTriangular(uplo, n, stride)
	if uplo == blas.Upper {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
			}
		}
		return ans
	}
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			ans.Data[i*ans.Stride+j] = rnd.NormFloat64()
		}
	}
	return ans
}

// generalOutsideAllNaN returns whether all out-of-range elements have NaN
// values.
func generalOutsideAllNaN(a blas64.General) bool {
	// Check after last column.
	for i := 0; i < a.Rows-1; i++ {
		for _, v := range a.Data[i*a.Stride+a.Cols : i*a.Stride+a.Stride] {
			if !math.IsNaN(v) {
				return false
			}
		}
	}
	// Check after last element.
	last := (a.Rows-1)*a.Stride + a.Cols
	if a.Rows == 0 || a.Cols == 0 {
		last = 0
	}
	for _, v := range a.Data[last:] {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// triangularOutsideAllNaN returns whether all out-of-triangle elements have NaN
// values.
func triangularOutsideAllNaN(a blas64.Triangular) bool {
	if a.Uplo == blas.Upper {
		// Check below diagonal.
		for i := 0; i < a.N; i++ {
			for _, v := range a.Data[i*a.Stride : i*a.Stride+i] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
		// Check after last column.
		for i := 0; i < a.N-1; i++ {
			for _, v := range a.Data[i*a.Stride+a.N : i*a.Stride+a.Stride] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
	} else {
		// Check above diagonal.
		for i := 0; i < a.N-1; i++ {
			for _, v := range a.Data[i*a.Stride+i+1 : i*a.Stride+a.Stride] {
				if !math.IsNaN(v) {
					return false
				}
			}
		}
	}
	// Check after last element.
	for _, v := range a.Data[max(0, a.N-1)*a.Stride+a.N:] {
		if !math.IsNaN(v) {
			return false
		}
	}
	return true
}

// transposeGeneral returns a new general matrix that is the transpose of the
// input. Nothing is done with data outside the {rows, cols} limit of the general.
func transposeGeneral(a blas64.General) blas64.General {
	ans := blas64.General{
		Rows:   a.Cols,
		Cols:   a.Rows,
		Stride: a.Rows,
		Data:   make([]float64, a.Cols*a.Rows),
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			ans.Data[j*ans.Stride+i] = a.Data[i*a.Stride+j]
		}
	}
	return ans
}

// extractVMat collects the single reflectors from a into a matrix.
func extractVMat(m, n int, a []float64, lda int, direct lapack.Direct, store lapack.StoreV) blas64.General {
	k := min(m, n)
	switch {
	default:
		panic("not implemented")
	case direct == lapack.Forward && store == lapack.ColumnWise:
		v := blas64.General{
			Rows:   m,
			Cols:   k,
			Stride: k,
			Data:   make([]float64, m*k),
		}
		for i := 0; i < k; i++ {
			for j := 0; j < i; j++ {
				v.Data[j*v.Stride+i] = 0
			}
			v.Data[i*v.Stride+i] = 1
			for j := i + 1; j < m; j++ {
				v.Data[j*v.Stride+i] = a[j*lda+i]
			}
		}
		return v
	case direct == lapack.Forward && store == lapack.RowWise:
		v := blas64.General{
			Rows:   k,
			Cols:   n,
			Stride: n,
			Data:   make([]float64, k*n),
		}
		for i := 0; i < k; i++ {
			for j := 0; j < i; j++ {
				v.Data[i*v.Stride+j] = 0
			}
			v.Data[i*v.Stride+i] = 1
			for j := i + 1; j < n; j++ {
				v.Data[i*v.Stride+j] = a[i*lda+j]
			}
		}
		return v
	}
}

// constructBidiagonal constructs a bidiagonal matrix with the given diagonal
// and off-diagonal elements.
func constructBidiagonal(uplo blas.Uplo, n int, d, e []float64) blas64.General {
	bMat := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}

	for i := 0; i < n-1; i++ {
		bMat.Data[i*bMat.Stride+i] = d[i]
		if uplo == blas.Upper {
			bMat.Data[i*bMat.Stride+i+1] = e[i]
		} else {
			bMat.Data[(i+1)*bMat.Stride+i] = e[i]
		}
	}
	bMat.Data[(n-1)*bMat.Stride+n-1] = d[n-1]
	return bMat
}

// constructVMat transforms the v matrix based on the storage.
func constructVMat(vMat blas64.General, store lapack.StoreV, direct lapack.Direct) blas64.General {
	m := vMat.Rows
	k := vMat.Cols
	switch {
	default:
		panic("not implemented")
	case store == lapack.ColumnWise && direct == lapack.Forward:
		ldv := k
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				if j > i {
					v[i*ldv+j] = 0
				} else if j == i {
					v[i*ldv+i] = 1
				} else {
					v[i*ldv+j] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   m,
			Cols:   k,
			Stride: k,
			Data:   v,
		}
	case store == lapack.RowWise && direct == lapack.Forward:
		ldv := m
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				if j > i {
					v[j*ldv+i] = 0
				} else if j == i {
					v[j*ldv+i] = 1
				} else {
					v[j*ldv+i] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   k,
			Cols:   m,
			Stride: m,
			Data:   v,
		}
	case store == lapack.ColumnWise && direct == lapack.Backward:
		rowsv := m
		ldv := k
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				vrow := rowsv - i - 1
				vcol := k - j - 1
				if j > i {
					v[vrow*ldv+vcol] = 0
				} else if j == i {
					v[vrow*ldv+vcol] = 1
				} else {
					v[vrow*ldv+vcol] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   rowsv,
			Cols:   ldv,
			Stride: ldv,
			Data:   v,
		}
	case store == lapack.RowWise && direct == lapack.Backward:
		rowsv := k
		ldv := m
		v := make([]float64, m*k)
		for i := 0; i < m; i++ {
			for j := 0; j < k; j++ {
				vcol := ldv - i - 1
				vrow := k - j - 1
				if j > i {
					v[vrow*ldv+vcol] = 0
				} else if j == i {
					v[vrow*ldv+vcol] = 1
				} else {
					v[vrow*ldv+vcol] = vMat.Data[i*vMat.Stride+j]
				}
			}
		}
		return blas64.General{
			Rows:   rowsv,
			Cols:   ldv,
			Stride: ldv,
			Data:   v,
		}
	}
}

func constructH(tau []float64, v blas64.General, store lapack.StoreV, direct lapack.Direct) blas64.General {
	m := v.Rows
	k := v.Cols
	if store == lapack.RowWise {
		m, k = k, m
	}
	h := blas64.General{
		Rows:   m,
		Cols:   m,
		Stride: m,
		Data:   make([]float64, m*m),
	}
	for i := 0; i < m; i++ {
		h.Data[i*m+i] = 1
	}
	for i := 0; i < k; i++ {
		vecData := make([]float64, m)
		if store == lapack.ColumnWise {
			for j := 0; j < m; j++ {
				vecData[j] = v.Data[j*v.Cols+i]
			}
		} else {
			for j := 0; j < m; j++ {
				vecData[j] = v.Data[i*v.Cols+j]
			}
		}
		vec := blas64.Vector{
			Inc:  1,
			Data: vecData,
		}

		hi := blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, m*m),
		}
		for i := 0; i < m; i++ {
			hi.Data[i*m+i] = 1
		}
		// hi = I - tau * v * v^T
		blas64.Ger(-tau[i], vec, vec, hi)

		hcopy := blas64.General{
			Rows:   m,
			Cols:   m,
			Stride: m,
			Data:   make([]float64, m*m),
		}
		copy(hcopy.Data, h.Data)
		if direct == lapack.Forward {
			// H = H * H_I in forward mode
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hcopy, hi, 0, h)
		} else {
			// H = H_I * H in backward mode
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hi, hcopy, 0, h)
		}
	}
	return h
}

// constructQ constructs the Q matrix from the result of dgeqrf and dgeqr2.
func constructQ(kind string, m, n int, a []float64, lda int, tau []float64) blas64.General {
	k := min(m, n)
	return constructQK(kind, m, n, k, a, lda, tau)
}

// constructQK constructs the Q matrix from the result of dgeqrf and dgeqr2 using
// the first k reflectors.
func constructQK(kind string, m, n, k int, a []float64, lda int, tau []float64) blas64.General {
	var sz int
	switch kind {
	case "QR":
		sz = m
	case "LQ":
		sz = n
	}

	q := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: sz,
		Data:   make([]float64, sz*sz),
	}
	for i := 0; i < sz; i++ {
		q.Data[i*sz+i] = 1
	}
	qCopy := blas64.General{
		Rows:   q.Rows,
		Cols:   q.Cols,
		Stride: q.Stride,
		Data:   make([]float64, len(q.Data)),
	}
	for i := 0; i < k; i++ {
		h := blas64.General{
			Rows:   sz,
			Cols:   sz,
			Stride: sz,
			Data:   make([]float64, sz*sz),
		}
		for j := 0; j < sz; j++ {
			h.Data[j*sz+j] = 1
		}
		vVec := blas64.Vector{
			Inc:  1,
			Data: make([]float64, sz),
		}
		for j := 0; j < i; j++ {
			vVec.Data[j] = 0
		}
		vVec.Data[i] = 1
		switch kind {
		case "QR":
			for j := i + 1; j < sz; j++ {
				vVec.Data[j] = a[lda*j+i]
			}
		case "LQ":
			for j := i + 1; j < sz; j++ {
				vVec.Data[j] = a[i*lda+j]
			}
		}
		blas64.Ger(-tau[i], vVec, vVec, h)
		copy(qCopy.Data, q.Data)
		// Mulitply q by the new h
		switch kind {
		case "QR":
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, h, 0, q)
		case "LQ":
			blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, h, qCopy, 0, q)
		}
	}
	return q
}

// checkBidiagonal checks the bidiagonal decomposition from dlabrd and dgebd2.
// The input to this function is the answer returned from the routines, stored
// in a, d, e, tauP, and tauQ. The data of original A matrix (before
// decomposition) is input in aCopy.
//
// checkBidiagonal constructs the V and U matrices, and from them constructs Q
// and P. Using these constructions, it checks that Q^T * A * P and checks that
// the result is bidiagonal.
func checkBidiagonal(t *testing.T, m, n, nb int, a []float64, lda int, d, e, tauP, tauQ, aCopy []float64) {
	// Check the answer.
	// Construct V and U.
	qMat := constructQPBidiagonal(lapack.ApplyQ, m, n, nb, a, lda, tauQ)
	pMat := constructQPBidiagonal(lapack.ApplyP, m, n, nb, a, lda, tauP)

	// Compute Q^T * A * P
	aMat := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: lda,
		Data:   make([]float64, len(aCopy)),
	}
	copy(aMat.Data, aCopy)

	tmp1 := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, qMat, aMat, 0, tmp1)
	tmp2 := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp1, pMat, 0, tmp2)

	// Check that the first nb rows and cols of tm2 are upper bidiagonal
	// if m >= n, and lower bidiagonal otherwise.
	correctDiag := true
	matchD := true
	matchE := true
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if i >= nb && j >= nb {
				continue
			}
			v := tmp2.Data[i*tmp2.Stride+j]
			if i == j {
				if math.Abs(d[i]-v) > 1e-12 {
					matchD = false
				}
				continue
			}
			if m >= n && i == j-1 {
				if math.Abs(e[j-1]-v) > 1e-12 {
					matchE = false
				}
				continue
			}
			if m < n && i-1 == j {
				if math.Abs(e[i-1]-v) > 1e-12 {
					matchE = false
				}
				continue
			}
			if math.Abs(v) > 1e-12 {
				correctDiag = false
			}
		}
	}
	if !correctDiag {
		t.Errorf("Updated A not bi-diagonal")
	}
	if !matchD {
		fmt.Println("d = ", d)
		t.Errorf("D Mismatch")
	}
	if !matchE {
		t.Errorf("E mismatch")
	}
}

// constructQPBidiagonal constructs Q or P from the Bidiagonal decomposition
// computed by dlabrd and bgebd2.
func constructQPBidiagonal(vect lapack.DecompUpdate, m, n, nb int, a []float64, lda int, tau []float64) blas64.General {
	sz := n
	if vect == lapack.ApplyQ {
		sz = m
	}

	var ldv int
	var v blas64.General
	if vect == lapack.ApplyQ {
		ldv = nb
		v = blas64.General{
			Rows:   m,
			Cols:   nb,
			Stride: ldv,
			Data:   make([]float64, m*ldv),
		}
	} else {
		ldv = n
		v = blas64.General{
			Rows:   nb,
			Cols:   n,
			Stride: ldv,
			Data:   make([]float64, m*ldv),
		}
	}

	if vect == lapack.ApplyQ {
		if m >= n {
			for i := 0; i < m; i++ {
				for j := 0; j <= min(nb-1, i); j++ {
					if i == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		} else {
			for i := 1; i < m; i++ {
				for j := 0; j <= min(nb-1, i-1); j++ {
					if i-1 == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		}
	} else {
		if m < n {
			for i := 0; i < nb; i++ {
				for j := i; j < n; j++ {
					if i == j {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		} else {
			for i := 0; i < nb; i++ {
				for j := i + 1; j < n; j++ {
					if j-1 == i {
						v.Data[i*ldv+j] = 1
						continue
					}
					v.Data[i*ldv+j] = a[i*lda+j]
				}
			}
		}
	}

	// The variable name is a computation of Q, but the algorithm is mostly the
	// same for computing P (just with different data).
	qMat := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: sz,
		Data:   make([]float64, sz*sz),
	}
	hMat := blas64.General{
		Rows:   sz,
		Cols:   sz,
		Stride: sz,
		Data:   make([]float64, sz*sz),
	}
	// set Q to I
	for i := 0; i < sz; i++ {
		qMat.Data[i*qMat.Stride+i] = 1
	}
	for i := 0; i < nb; i++ {
		qCopy := blas64.General{Rows: qMat.Rows, Cols: qMat.Cols, Stride: qMat.Stride, Data: make([]float64, len(qMat.Data))}
		copy(qCopy.Data, qMat.Data)

		// Set g and h to I
		for i := 0; i < sz; i++ {
			for j := 0; j < sz; j++ {
				if i == j {
					hMat.Data[i*sz+j] = 1
				} else {
					hMat.Data[i*sz+j] = 0
				}
			}
		}
		var vi blas64.Vector
		// H -= tauQ[i] * v[i] * v[i]^t
		if vect == lapack.ApplyQ {
			vi = blas64.Vector{
				Inc:  v.Stride,
				Data: v.Data[i:],
			}
		} else {
			vi = blas64.Vector{
				Inc:  1,
				Data: v.Data[i*v.Stride:],
			}
		}
		blas64.Ger(-tau[i], vi, vi, hMat)
		// Q = Q * G[1]
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, qCopy, hMat, 0, qMat)
	}
	return qMat
}

// printRowise prints the matrix with one row per line. This is useful for debugging.
// If beyond is true, it prints beyond the final column to lda. If false, only
// the columns are printed.
func printRowise(a []float64, m, n, lda int, beyond bool) {
	for i := 0; i < m; i++ {
		end := n
		if beyond {
			end = lda
		}
		fmt.Println(a[i*lda : i*lda+end])
	}
}

// isOrthonormal checks that a general matrix is orthonormal.
func isOrthonormal(q blas64.General) bool {
	n := q.Rows
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			dot := blas64.Dot(n,
				blas64.Vector{Inc: 1, Data: q.Data[i*q.Stride:]},
				blas64.Vector{Inc: 1, Data: q.Data[j*q.Stride:]},
			)
			if i == j {
				if math.Abs(dot-1) > 1e-10 {
					return false
				}
			} else {
				if math.Abs(dot) > 1e-10 {
					return false
				}
			}
		}
	}
	return true
}

// copyMatrix copies an m×n matrix src of stride n into an m×n matrix dst of stride ld.
func copyMatrix(m, n int, dst []float64, ld int, src []float64) {
	for i := 0; i < m; i++ {
		copy(dst[i*ld:i*ld+n], src[i*n:i*n+n])
	}
}

// cloneGeneral allocates and returns an exact copy of the given general matrix.
func cloneGeneral(a blas64.General) blas64.General {
	c := a
	c.Data = make([]float64, len(a.Data))
	copy(c.Data, a.Data)
	return c
}

// equalApprox returns whether the matrices A and B of order n are approximately
// equal within given tolerance.
func equalApprox(m, n int, a []float64, lda int, b []float64, tol float64) bool {
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if math.Abs(a[i*lda+j]-b[i*n+j]) > tol {
				return false
			}
		}
	}
	return true
}

// equalApproxGeneral returns whether the general matrices a and b are
// approximately equal within given tolerance.
func equalApproxGeneral(a, b blas64.General, tol float64) bool {
	if a.Rows != b.Rows || a.Cols != b.Cols {
		panic("bad input")
	}
	for i := 0; i < a.Rows; i++ {
		for j := 0; j < a.Cols; j++ {
			diff := a.Data[i*a.Stride+j] - b.Data[i*b.Stride+j]
			if math.Abs(diff) > tol {
				return false
			}
		}
	}
	return true
}

// equalApproxTriangular returns whether the triangular matrices A and B of
// order n are approximately equal within given tolerance.
func equalApproxTriangular(upper bool, n int, a []float64, lda int, b []float64, tol float64) bool {
	if upper {
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				if math.Abs(a[i*lda+j]-b[i*n+j]) > tol {
					return false
				}
			}
		}
		return true
	}
	for i := 0; i < n; i++ {
		for j := 0; j <= i; j++ {
			if math.Abs(a[i*lda+j]-b[i*n+j]) > tol {
				return false
			}
		}
	}
	return true
}

// eye returns an identity matrix of given order and stride.
func eye(n, stride int) blas64.General {
	ans := nanGeneral(n, n, stride)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ans.Data[i*ans.Stride+j] = 0
		}
		ans.Data[i*ans.Stride+i] = 1
	}
	return ans
}
