// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
	"github.com/gonum/lapack/lapack64"
	"github.com/gonum/matrix"
)

// Matrix is the basic matrix interface type.
type Matrix interface {
	// Dims returns the dimensions of a Matrix.
	Dims() (r, c int)

	// At returns the value of a matrix element at row i, column j.
	// It will panic if i or j are out of bounds for the matrix.
	At(i, j int) float64

	// T returns the transpose of the Matrix. Whether T returns a copy of the
	// underlying data is implementation dependent.
	// This method may be implemented using the Transpose type, which
	// provides an implicit matrix transpose.
	T() Matrix
}

var (
	_ Matrix       = Transpose{}
	_ Untransposer = Transpose{}
)

// Transpose is a type for performing an implicit matrix transpose. It implements
// the Matrix interface, returning values from the transpose of the matrix within.
type Transpose struct {
	Matrix Matrix
}

// At returns the value of the element at row i and column j of the transposed
// matrix, that is, row j and column i of the Matrix field.
func (t Transpose) At(i, j int) float64 {
	return t.Matrix.At(j, i)
}

// Dims returns the dimensions of the transposed matrix. The number of rows returned
// is the number of columns in the Matrix field, and the number of columns is
// the number of rows in the Matrix field.
func (t Transpose) Dims() (r, c int) {
	c, r = t.Matrix.Dims()
	return r, c
}

// T performs an implicit transpose by returning the Matrix field.
func (t Transpose) T() Matrix {
	return t.Matrix
}

// Untranspose returns the Matrix field.
func (t Transpose) Untranspose() Matrix {
	return t.Matrix
}

// Untransposer is a type that can undo an implicit transpose.
type Untransposer interface {
	// Note: This interface is needed to unify all of the Transpose types. In
	// the mat64 methods, we need to test if the Matrix has been implicitly
	// transposed. If this is checked by testing for the specific Transpose type
	// then the behavior will be different if the user uses T() or TTri() for a
	// triangular matrix.

	// Untranspose returns the underlying Matrix stored for the implicit transpose.
	Untranspose() Matrix
}

// Mutable is a matrix interface type that allows elements to be altered.
type Mutable interface {
	// Set alters the matrix element at row i, column j to v.
	// It will panic if i or j are out of bounds for the matrix.
	Set(i, j int, v float64)

	Matrix
}

// A RowViewer can return a Vector reflecting a row that is backed by the matrix
// data. The Vector returned will have length equal to the number of columns.
type RowViewer interface {
	RowView(i int) *Vector
}

// A RawRowViewer can return a slice of float64 reflecting a row that is backed by the matrix
// data.
type RawRowViewer interface {
	RawRowView(i int) []float64
}

// A ColViewer can return a Vector reflecting a column that is backed by the matrix
// data. The Vector returned will have length equal to the number of rows.
type ColViewer interface {
	ColView(j int) *Vector
}

// A RawColViewer can return a slice of float64 reflecting a column that is backed by the matrix
// data.
type RawColViewer interface {
	RawColView(j int) []float64
}

// A Cloner can make a copy of a into the receiver, overwriting the previous value of the
// receiver. The clone operation does not make any restriction on shape and will not cause
// shadowing.
type Cloner interface {
	Clone(a Matrix)
}

// A Reseter can reset the matrix so that it can be reused as the receiver of a dimensionally
// restricted operation. This is commonly used when the matrix is being used as a workspace
// or temporary matrix.
//
// If the matrix is a view, using the reset matrix may result in data corruption in elements
// outside the view.
type Reseter interface {
	Reset()
}

// A Copier can make a copy of elements of a into the receiver. The submatrix copied
// starts at row and column 0 and has dimensions equal to the minimum dimensions of
// the two matrices. The number of row and columns copied is returned.
// Copy will copy from a source that aliases the receiver unless the source is transposed;
// an aliasing transpose copy will panic with the exception for a special case when
// the source data has a unitary increment or stride.
type Copier interface {
	Copy(a Matrix) (r, c int)
}

// A Viewer returns a submatrix view of the Matrix parameter, starting at row i, column j
// and extending r rows and c columns. If i or j are out of range, or r or c are zero or
// extend beyond the bounds of the matrix View will panic with ErrIndexOutOfRange. The
// returned matrix must retain the receiver's reference to the original matrix such that
// changes in the elements of the submatrix are reflected in the original and vice versa.
type Viewer interface {
	View(i, j, r, c int) Matrix
}

// A Grower can grow the size of the represented matrix by the given number of rows and columns.
// Growing beyond the size given by the Caps method will result in the allocation of a new
// matrix and copying of the elements. If Grow is called with negative increments it will
// panic with ErrIndexOutOfRange.
type Grower interface {
	Caps() (r, c int)
	Grow(r, c int) Matrix
}

// A BandWidther represents a banded matrix and can return the left and right half-bandwidths, k1 and
// k2.
type BandWidther interface {
	BandWidth() (k1, k2 int)
}

// A RawMatrixSetter can set the underlying blas64.General used by the receiver. There is no restriction
// on the shape of the receiver. Changes to the receiver's elements will be reflected in the blas64.General.Data.
type RawMatrixSetter interface {
	SetRawMatrix(a blas64.General)
}

// A RawMatrixer can return a blas64.General representation of the receiver. Changes to the blas64.General.Data
// slice will be reflected in the original matrix, changes to the Rows, Cols and Stride fields will not.
type RawMatrixer interface {
	RawMatrix() blas64.General
}

// A RawVectorer can return a blas64.Vector representation of the receiver. Changes to the blas64.Vector.Data
// slice will be reflected in the original matrix, changes to the Inc field will not.
type RawVectorer interface {
	RawVector() blas64.Vector
}

// TODO(btracey): Consider adding CopyCol/CopyRow if the behavior seems useful.
// TODO(btracey): Add in fast paths to Row/Col for the other concrete types
// (TriDense, etc.) as well as relevant interfaces (RowColer, RawRowViewer, etc.)

// Col copies the elements in the jth column of the matrix into the slice dst.
// The length of the provided slice must equal the number of rows, unless the
// slice is nil in which case a new slice is first allocated.
func Col(dst []float64, j int, a Matrix) []float64 {
	r, c := a.Dims()
	if j < 0 || j >= c {
		panic(matrix.ErrColAccess)
	}
	if dst == nil {
		dst = make([]float64, r)
	} else {
		if len(dst) != r {
			panic(matrix.ErrColLength)
		}
	}
	aU, aTrans := untranspose(a)
	if rm, ok := aU.(RawMatrixer); ok {
		m := rm.RawMatrix()
		if aTrans {
			copy(dst, m.Data[j*m.Stride:j*m.Stride+m.Cols])
			return dst
		}
		blas64.Copy(r,
			blas64.Vector{Inc: m.Stride, Data: m.Data[j:]},
			blas64.Vector{Inc: 1, Data: dst},
		)
		return dst
	}
	for i := 0; i < r; i++ {
		dst[i] = a.At(i, j)
	}
	return dst
}

// Row copies the elements in the jth column of the matrix into the slice dst.
// The length of the provided slice must equal the number of columns, unless the
// slice is nil in which case a new slice is first allocated.
func Row(dst []float64, i int, a Matrix) []float64 {
	r, c := a.Dims()
	if i < 0 || i >= r {
		panic(matrix.ErrColAccess)
	}
	if dst == nil {
		dst = make([]float64, c)
	} else {
		if len(dst) != c {
			panic(matrix.ErrRowLength)
		}
	}
	aU, aTrans := untranspose(a)
	if rm, ok := aU.(RawMatrixer); ok {
		m := rm.RawMatrix()
		if aTrans {
			blas64.Copy(c,
				blas64.Vector{Inc: m.Stride, Data: m.Data[i:]},
				blas64.Vector{Inc: 1, Data: dst},
			)
			return dst
		}
		copy(dst, m.Data[i*m.Stride:i*m.Stride+m.Cols])
		return dst
	}
	for j := 0; j < c; j++ {
		dst[j] = a.At(i, j)
	}
	return dst
}

// Cond returns the condition number of the given matrix under the given norm.
// The condition number must be based on the 1-norm, 2-norm or ∞-norm.
// Cond will panic with matrix.ErrShape if the matrix has zero size.
//
// BUG(btracey): The computation of the 1-norm and ∞-norm for non-square matrices
// is innacurate, although is typically the right order of magnitude. See
// https://github.com/xianyi/OpenBLAS/issues/636. While the value returned will
// change with the resolution of this bug, the result from Cond will match the
// condition number used internally.
func Cond(a Matrix, norm float64) float64 {
	m, n := a.Dims()
	if m == 0 || n == 0 {
		panic(matrix.ErrShape)
	}
	var lnorm lapack.MatrixNorm
	switch norm {
	default:
		panic("mat64: bad norm value")
	case 1:
		lnorm = lapack.MaxColumnSum
	case 2:
		var svd SVD
		ok := svd.Factorize(a, matrix.SVDNone)
		if !ok {
			return math.Inf(1)
		}
		return svd.Cond()
	case math.Inf(1):
		lnorm = lapack.MaxRowSum
	}
	if m == n {
		// Use the LU decomposition to compute the condition number.
		tmp := getWorkspace(m, n, false)
		tmp.Copy(a)
		work := make([]float64, 4*n)
		aNorm := lapack64.Lange(lnorm, tmp.mat, work)
		pivot := make([]int, m)
		lapack64.Getrf(tmp.mat, pivot)
		iwork := make([]int, n)
		v := lapack64.Gecon(lnorm, tmp.mat, aNorm, work, iwork)
		putWorkspace(tmp)
		return 1 / v
	}
	if m > n {
		// Use the QR factorization to compute the condition number.
		tmp := getWorkspace(m, n, false)
		tmp.Copy(a)
		work := make([]float64, 3*n)
		tau := make([]float64, min(m, n))
		lapack64.Geqrf(tmp.mat, tau, work, -1)
		if int(work[0]) > len(work) {
			work = make([]float64, int(work[0]))
		}
		lapack64.Geqrf(tmp.mat, tau, work, len(work))

		iwork := make([]int, n)
		r := tmp.asTriDense(n, blas.NonUnit, blas.Upper)
		v := lapack64.Trcon(lnorm, r.mat, work, iwork)
		putWorkspace(tmp)
		return 1 / v
	}
	// Use the LQ factorization to compute the condition number.
	tmp := getWorkspace(m, n, false)
	tmp.Copy(a)
	work := make([]float64, 3*m)
	tau := make([]float64, min(m, n))
	lapack64.Gelqf(tmp.mat, tau, work, -1)
	if int(work[0]) > len(work) {
		work = make([]float64, int(work[0]))
	}
	lapack64.Gelqf(tmp.mat, tau, work, len(work))

	iwork := make([]int, m)
	l := tmp.asTriDense(m, blas.NonUnit, blas.Lower)
	v := lapack64.Trcon(lnorm, l.mat, work, iwork)
	putWorkspace(tmp)
	return 1 / v
}

// Det returns the determinant of the matrix a. In many expressions using LogDet
// will be more numerically stable.
func Det(a Matrix) float64 {
	det, sign := LogDet(a)
	return math.Exp(det) * sign
}

// Dot returns the sum of the element-wise products of the elements of a and b.
// Dot panics if the matrix sizes are unequal.
func Dot(a, b Matrix) float64 {
	r, c := a.Dims()
	rb, cb := b.Dims()
	if r != rb || c != cb {
		panic(matrix.ErrShape)
	}
	aU, aTrans := untranspose(a)
	bU, bTrans := untranspose(b)
	if rma, ok := aU.(RawVectorer); ok {
		if rmb, ok := bU.(RawVectorer); ok {
			if c > r {
				r = c
			}
			return blas64.Dot(r, rma.RawVector(), rmb.RawVector())
		}
	}
	var sum float64
	if rma, ok := aU.(RawMatrixer); ok {
		if rmb, ok := bU.(RawMatrixer); ok {
			ra := rma.RawMatrix()
			rb := rmb.RawMatrix()
			if aTrans == bTrans {
				for i := 0; i < ra.Rows; i++ {
					sum += blas64.Dot(ra.Cols,
						blas64.Vector{Inc: 1, Data: ra.Data[i*ra.Stride:]},
						blas64.Vector{Inc: 1, Data: rb.Data[i*rb.Stride:]},
					)
				}
				return sum
			}
			for i := 0; i < ra.Rows; i++ {
				sum += blas64.Dot(ra.Cols,
					blas64.Vector{Inc: 1, Data: ra.Data[i*ra.Stride:]},
					blas64.Vector{Inc: rb.Stride, Data: rb.Data[i:]},
				)
			}
			return sum
		}
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			sum += a.At(i, j) * b.At(i, j)
		}
	}
	return sum
}

// Equal returns whether the matrices a and b have the same size
// and are element-wise equal.
func Equal(a, b Matrix) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	aU, aTrans := untranspose(a)
	bU, bTrans := untranspose(b)
	if rma, ok := aU.(RawMatrixer); ok {
		if rmb, ok := bU.(RawMatrixer); ok {
			ra := rma.RawMatrix()
			rb := rmb.RawMatrix()
			if aTrans == bTrans {
				for i := 0; i < ra.Rows; i++ {
					for j := 0; j < ra.Cols; j++ {
						if ra.Data[i*ra.Stride+j] != rb.Data[i*rb.Stride+j] {
							return false
						}
					}
				}
				return true
			}
			for i := 0; i < ra.Rows; i++ {
				for j := 0; j < ra.Cols; j++ {
					if ra.Data[i*ra.Stride+j] != rb.Data[j*rb.Stride+i] {
						return false
					}
				}
			}
			return true
		}
	}
	if rma, ok := aU.(RawSymmetricer); ok {
		if rmb, ok := bU.(RawSymmetricer); ok {
			ra := rma.RawSymmetric()
			rb := rmb.RawSymmetric()
			// Symmetric matrices are always upper and equal to their transpose.
			for i := 0; i < ra.N; i++ {
				for j := i; j < ra.N; j++ {
					if ra.Data[i*ra.Stride+j] != rb.Data[i*rb.Stride+j] {
						return false
					}
				}
			}
			return true
		}
	}
	if ra, ok := aU.(*Vector); ok {
		if rb, ok := bU.(*Vector); ok {
			// If the raw vectors are the same length they must either both be
			// transposed or both not transposed (or have length 1).
			for i := 0; i < ra.n; i++ {
				if ra.mat.Data[i*ra.mat.Inc] != rb.mat.Data[i*rb.mat.Inc] {
					return false
				}
			}
			return true
		}
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.At(i, j) != b.At(i, j) {
				return false
			}
		}
	}
	return true
}

// EqualApprox returns whether the matrices a and b have the same size and contain all equal
// elements with tolerance for element-wise equality specified by epsilon. Matrices
// with non-equal shapes are not equal.
func EqualApprox(a, b Matrix, epsilon float64) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || ac != bc {
		return false
	}
	aU, aTrans := untranspose(a)
	bU, bTrans := untranspose(b)
	if rma, ok := aU.(RawMatrixer); ok {
		if rmb, ok := bU.(RawMatrixer); ok {
			ra := rma.RawMatrix()
			rb := rmb.RawMatrix()
			if aTrans == bTrans {
				for i := 0; i < ra.Rows; i++ {
					for j := 0; j < ra.Cols; j++ {
						if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[i*rb.Stride+j], epsilon, epsilon) {
							return false
						}
					}
				}
				return true
			}
			for i := 0; i < ra.Rows; i++ {
				for j := 0; j < ra.Cols; j++ {
					if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[j*rb.Stride+i], epsilon, epsilon) {
						return false
					}
				}
			}
			return true
		}
	}
	if rma, ok := aU.(RawSymmetricer); ok {
		if rmb, ok := bU.(RawSymmetricer); ok {
			ra := rma.RawSymmetric()
			rb := rmb.RawSymmetric()
			// Symmetric matrices are always upper and equal to their transpose.
			for i := 0; i < ra.N; i++ {
				for j := i; j < ra.N; j++ {
					if !floats.EqualWithinAbsOrRel(ra.Data[i*ra.Stride+j], rb.Data[i*rb.Stride+j], epsilon, epsilon) {
						return false
					}
				}
			}
			return true
		}
	}
	if ra, ok := aU.(*Vector); ok {
		if rb, ok := bU.(*Vector); ok {
			// If the raw vectors are the same length they must either both be
			// transposed or both not transposed (or have length 1).
			for i := 0; i < ra.n; i++ {
				if !floats.EqualWithinAbsOrRel(ra.mat.Data[i*ra.mat.Inc], rb.mat.Data[i*rb.mat.Inc], epsilon, epsilon) {
					return false
				}
			}
			return true
		}
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if !floats.EqualWithinAbsOrRel(a.At(i, j), b.At(i, j), epsilon, epsilon) {
				return false
			}
		}
	}
	return true
}

// LogDet returns the log of the determinant and the sign of the determinant
// for the matrix that has been factorized. Numerical stability in product and
// division expressions is generally improved by working in log space.
func LogDet(a Matrix) (det float64, sign float64) {
	// TODO(btracey): Add specialized routines for TriDense, etc.
	var lu LU
	lu.Factorize(a)
	return lu.LogDet()
}

// Max returns the largest element value of the matrix A.
// Max will panic with matrix.ErrShape if the matrix has zero size.
func Max(a Matrix) float64 {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		panic(matrix.ErrShape)
	}
	// Max(A) = Max(A^T)
	aU, _ := untranspose(a)
	switch m := aU.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		max := math.Inf(-1)
		for i := 0; i < rm.Rows; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+rm.Cols] {
				if v > max {
					max = v
				}
			}
		}
		return max
	case RawTriangular:
		rm := m.RawTriangular()
		// The max of a triangular is at least 0 unless the size is 1.
		if rm.N == 1 {
			return rm.Data[0]
		}
		max := 0.0
		if rm.Uplo == blas.Upper {
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v > max {
						max = v
					}
				}
			}
			return max
		}
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v > max {
					max = v
				}
			}
		}
		return max
	case RawSymmetricer:
		rm := m.RawSymmetric()
		if rm.Uplo != blas.Upper {
			panic(badSymTriangle)
		}
		max := math.Inf(-1)
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
				if v > max {
					max = v
				}
			}
		}
		return max
	default:
		r, c := aU.Dims()
		max := math.Inf(-1)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				v := aU.At(i, j)
				if v > max {
					max = v
				}
			}
		}
		return max
	}
}

// Min returns the smallest element value of the matrix A.
// Min will panic with matrix.ErrShape if the matrix has zero size.
func Min(a Matrix) float64 {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		panic(matrix.ErrShape)
	}
	// Min(A) = Min(A^T)
	aU, _ := untranspose(a)
	switch m := aU.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		min := math.Inf(1)
		for i := 0; i < rm.Rows; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+rm.Cols] {
				if v < min {
					min = v
				}
			}
		}
		return min
	case RawTriangular:
		rm := m.RawTriangular()
		// The min of a triangular is at most 0 unless the size is 1.
		if rm.N == 1 {
			return rm.Data[0]
		}
		min := 0.0
		if rm.Uplo == blas.Upper {
			for i := 0; i < rm.N; i++ {
				for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
					if v < min {
						min = v
					}
				}
			}
			return min
		}
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+i+1] {
				if v < min {
					min = v
				}
			}
		}
		return min
	case RawSymmetricer:
		rm := m.RawSymmetric()
		if rm.Uplo != blas.Upper {
			panic(badSymTriangle)
		}
		min := math.Inf(1)
		for i := 0; i < rm.N; i++ {
			for _, v := range rm.Data[i*rm.Stride+i : i*rm.Stride+rm.N] {
				if v < min {
					min = v
				}
			}
		}
		return min
	default:
		r, c := aU.Dims()
		min := math.Inf(1)
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				v := aU.At(i, j)
				if v < min {
					min = v
				}
			}
		}
		return min
	}
}

// Norm returns the specified (induced) norm of the matrix a. See
// https://en.wikipedia.org/wiki/Matrix_norm for the definition of an induced norm.
//
// Valid norms are:
//    1 - The maximum absolute column sum
//    2 - Frobenius norm, the square root of the sum of the squares of the elements.
//  Inf - The maximum absolute row sum.
// Norm will panic with ErrNormOrder if an illegal norm order is specified and
// with matrix.ErrShape if the matrix has zero size.
func Norm(a Matrix, norm float64) float64 {
	r, c := a.Dims()
	if r == 0 || c == 0 {
		panic(matrix.ErrShape)
	}
	aU, aTrans := untranspose(a)
	var work []float64
	switch rma := aU.(type) {
	case RawMatrixer:
		rm := rma.RawMatrix()
		n := normLapack(norm, aTrans)
		if n == lapack.MaxColumnSum {
			work = make([]float64, rm.Cols)
		}
		return lapack64.Lange(n, rm, work)
	case RawTriangular:
		rm := rma.RawTriangular()
		n := normLapack(norm, aTrans)
		if n == lapack.MaxRowSum || n == lapack.MaxColumnSum {
			work = make([]float64, rm.N)
		}
		return lapack64.Lantr(n, rm, work)
	case RawSymmetricer:
		rm := rma.RawSymmetric()
		n := normLapack(norm, aTrans)
		if n == lapack.MaxRowSum || n == lapack.MaxColumnSum {
			work = make([]float64, rm.N)
		}
		return lapack64.Lansy(n, rm, work)
	case *Vector:
		rv := rma.RawVector()
		switch norm {
		default:
			panic("unreachable")
		case 1:
			if aTrans {
				imax := blas64.Iamax(rma.n, rv)
				return math.Abs(rma.At(imax, 0))
			}
			return blas64.Asum(rma.n, rv)
		case 2:
			return blas64.Nrm2(rma.n, rv)
		case math.Inf(1):
			if aTrans {
				return blas64.Asum(rma.n, rv)
			}
			imax := blas64.Iamax(rma.n, rv)
			return math.Abs(rma.At(imax, 0))
		}
	}
	switch norm {
	default:
		panic("unreachable")
	case 1:
		var max float64
		for j := 0; j < c; j++ {
			var sum float64
			for i := 0; i < r; i++ {
				sum += math.Abs(a.At(i, j))
			}
			if sum > max {
				max = sum
			}
		}
		return max
	case 2:
		var sum float64
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				v := a.At(i, j)
				sum += v * v
			}
		}
		return math.Sqrt(sum)
	case math.Inf(1):
		var max float64
		for i := 0; i < r; i++ {
			var sum float64
			for j := 0; j < c; j++ {
				sum += math.Abs(a.At(i, j))
			}
			if sum > max {
				max = sum
			}
		}
		return max
	}
}

// normLapack converts the float64 norm input in Norm to a lapack.MatrixNorm.
func normLapack(norm float64, aTrans bool) lapack.MatrixNorm {
	switch norm {
	case 1:
		n := lapack.MaxColumnSum
		if aTrans {
			n = lapack.MaxRowSum
		}
		return n
	case 2:
		return lapack.NormFrob
	case math.Inf(1):
		n := lapack.MaxRowSum
		if aTrans {
			n = lapack.MaxColumnSum
		}
		return n
	default:
		panic(matrix.ErrNormOrder)
	}
}

// Sum returns the sum of the elements of the matrix.
func Sum(a Matrix) float64 {
	// TODO(btracey): Add a fast path for the other supported matrix types.

	r, c := a.Dims()
	var sum float64
	aU, _ := untranspose(a)
	if rma, ok := aU.(RawMatrixer); ok {
		rm := rma.RawMatrix()
		for i := 0; i < rm.Rows; i++ {
			for _, v := range rm.Data[i*rm.Stride : i*rm.Stride+rm.Cols] {
				sum += v
			}
		}
		return sum
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			sum += a.At(i, j)
		}
	}
	return sum
}

// Trace returns the trace of the matrix. Trace will panic if the
// matrix is not square.
func Trace(a Matrix) float64 {
	r, c := a.Dims()
	if r != c {
		panic(matrix.ErrSquare)
	}

	aU, _ := untranspose(a)
	switch m := aU.(type) {
	case RawMatrixer:
		rm := m.RawMatrix()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	case RawTriangular:
		rm := m.RawTriangular()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	case RawSymmetricer:
		rm := m.RawSymmetric()
		var t float64
		for i := 0; i < r; i++ {
			t += rm.Data[i*rm.Stride+i]
		}
		return t
	default:
		var t float64
		for i := 0; i < r; i++ {
			t += a.At(i, i)
		}
		return t
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// use returns a float64 slice with l elements, using f if it
// has the necessary capacity, otherwise creating a new slice.
func use(f []float64, l int) []float64 {
	if l <= cap(f) {
		return f[:l]
	}
	return make([]float64, l)
}

// useZeroed returns a float64 slice with l elements, using f if it
// has the necessary capacity, otherwise creating a new slice. The
// elements of the returned slice are guaranteed to be zero.
func useZeroed(f []float64, l int) []float64 {
	if l <= cap(f) {
		f = f[:l]
		zero(f)
		return f
	}
	return make([]float64, l)
}

// zero zeros the given slice's elements.
func zero(f []float64) {
	for i := range f {
		f[i] = 0
	}
}
