// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/matrix"
)

var (
	symDense *SymDense

	_ Matrix           = symDense
	_ Symmetric        = symDense
	_ RawSymmetricer   = symDense
	_ MutableSymmetric = symDense
)

const (
	badSymTriangle = "mat64: blas64.Symmetric not upper"
	badSymCap      = "mat64: bad capacity for SymDense"
)

// SymDense is a symmetric matrix that uses dense storage. SymDense
// matrices are stored in the upper triangle.
type SymDense struct {
	mat blas64.Symmetric
	cap int
}

// Symmetric represents a symmetric matrix (where the element at {i, j} equals
// the element at {j, i}). Symmetric matrices are always square.
type Symmetric interface {
	Matrix
	// Symmetric returns the number of rows/columns in the matrix.
	Symmetric() int
}

// A RawSymmetricer can return a view of itself as a BLAS Symmetric matrix.
type RawSymmetricer interface {
	RawSymmetric() blas64.Symmetric
}

type MutableSymmetric interface {
	Symmetric
	SetSym(i, j int, v float64)
}

// NewSymDense constructs an n x n symmetric matrix. If len(mat) == n * n,
// mat will be used to hold the underlying data, or if mat == nil, new data will be allocated.
// The underlying data representation is the same as a Dense matrix, except
// the values of the entries in the lower triangular portion are completely ignored.
func NewSymDense(n int, mat []float64) *SymDense {
	if n < 0 {
		panic("mat64: negative dimension")
	}
	if mat != nil && n*n != len(mat) {
		panic(matrix.ErrShape)
	}
	if mat == nil {
		mat = make([]float64, n*n)
	}
	return &SymDense{
		mat: blas64.Symmetric{
			N:      n,
			Stride: n,
			Data:   mat,
			Uplo:   blas.Upper,
		},
		cap: n,
	}
}

func (s *SymDense) Dims() (r, c int) {
	return s.mat.N, s.mat.N
}

// T implements the Matrix interface. Symmetric matrices, by definition, are
// equal to their transpose, and this is a no-op.
func (s *SymDense) T() Matrix {
	return s
}

func (s *SymDense) Symmetric() int {
	return s.mat.N
}

// RawSymmetric returns the matrix as a blas64.Symmetric. The returned
// value must be stored in upper triangular format.
func (s *SymDense) RawSymmetric() blas64.Symmetric {
	return s.mat
}

// SetRawSymmetric sets the underlying blas64.Symmetric used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in b. SetRawSymmetric will panic if b is not an upper-encoded symmetric
// matrix.
func (s *SymDense) SetRawSymmetric(b blas64.Symmetric) {
	if b.Uplo != blas.Upper {
		panic(badSymTriangle)
	}
	s.mat = b
}

func (s *SymDense) isZero() bool {
	return s.mat.N == 0
}

// reuseAs resizes an empty matrix to a n×n matrix,
// or checks that a non-empty matrix is n×n.
func (s *SymDense) reuseAs(n int) {
	if s.mat.N > s.cap {
		panic(badSymCap)
	}
	if s.isZero() {
		s.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Data:   use(s.mat.Data, n*n),
			Uplo:   blas.Upper,
		}
		s.cap = n
		return
	}
	if s.mat.Uplo != blas.Upper {
		panic(badSymTriangle)
	}
	if s.mat.N != n {
		panic(matrix.ErrShape)
	}
}

func (s *SymDense) isolatedWorkspace(a Symmetric) (w *SymDense, restore func()) {
	n := a.Symmetric()
	w = getWorkspaceSym(n, false)
	return w, func() {
		s.CopySym(w)
		putWorkspaceSym(w)
	}
}

func (s *SymDense) AddSym(a, b Symmetric) {
	n := a.Symmetric()
	if n != b.Symmetric() {
		panic(matrix.ErrShape)
	}
	s.reuseAs(n)

	if a, ok := a.(RawSymmetricer); ok {
		if b, ok := b.(RawSymmetricer); ok {
			amat, bmat := a.RawSymmetric(), b.RawSymmetric()
			if s != a {
				s.checkOverlap(amat)
			}
			if s != b {
				s.checkOverlap(bmat)
			}
			for i := 0; i < n; i++ {
				btmp := bmat.Data[i*bmat.Stride+i : i*bmat.Stride+n]
				stmp := s.mat.Data[i*s.mat.Stride+i : i*s.mat.Stride+n]
				for j, v := range amat.Data[i*amat.Stride+i : i*amat.Stride+n] {
					stmp[j] = v + btmp[j]
				}
			}
			return
		}
	}

	for i := 0; i < n; i++ {
		stmp := s.mat.Data[i*s.mat.Stride : i*s.mat.Stride+n]
		for j := i; j < n; j++ {
			stmp[j] = a.At(i, j) + b.At(i, j)
		}
	}
}

func (s *SymDense) CopySym(a Symmetric) int {
	n := a.Symmetric()
	n = min(n, s.mat.N)
	if n == 0 {
		return 0
	}
	switch a := a.(type) {
	case RawSymmetricer:
		amat := a.RawSymmetric()
		if amat.Uplo != blas.Upper {
			panic(badSymTriangle)
		}
		for i := 0; i < n; i++ {
			copy(s.mat.Data[i*s.mat.Stride+i:i*s.mat.Stride+n], amat.Data[i*amat.Stride+i:i*amat.Stride+n])
		}
	default:
		for i := 0; i < n; i++ {
			stmp := s.mat.Data[i*s.mat.Stride : i*s.mat.Stride+n]
			for j := i; j < n; j++ {
				stmp[j] = a.At(i, j)
			}
		}
	}
	return n
}

// SymRankOne performs a symetric rank-one update to the matrix a and stores
// the result in the receiver
//  s = a + alpha * x * x'
func (s *SymDense) SymRankOne(a Symmetric, alpha float64, x *Vector) {
	n := x.Len()
	if a.Symmetric() != n {
		panic(matrix.ErrShape)
	}
	s.reuseAs(n)
	if s != a {
		if rs, ok := a.(RawSymmetricer); ok {
			s.checkOverlap(rs.RawSymmetric())
		}
		s.CopySym(a)
	}
	blas64.Syr(alpha, x.mat, s.mat)
}

// SymRankK performs a symmetric rank-k update to the matrix a and stores the
// result into the receiver. If a is zero, see SymOuterK.
//  s = a + alpha * x * x'
func (s *SymDense) SymRankK(a Symmetric, alpha float64, x Matrix) {
	n := a.Symmetric()
	r, _ := x.Dims()
	if r != n {
		panic(matrix.ErrShape)
	}
	xMat, aTrans := untranspose(x)
	var g blas64.General
	if rm, ok := xMat.(RawMatrixer); ok {
		g = rm.RawMatrix()
	} else {
		g = DenseCopyOf(x).mat
		aTrans = false
	}
	if a != s {
		if rs, ok := a.(RawSymmetricer); ok {
			s.checkOverlap(rs.RawSymmetric())
		}
		s.reuseAs(n)
		s.CopySym(a)
	}
	t := blas.NoTrans
	if aTrans {
		t = blas.Trans
	}
	blas64.Syrk(t, alpha, g, 1, s.mat)
}

// SymOuterK calculates the outer product of a times its transpose and stores
// the result into the receiver. In order to update an existing matrix, see
// SymRankOne
//  s = alpha * x * x'
func (s *SymDense) SymOuterK(alpha float64, x Matrix) {
	n, _ := x.Dims()
	switch {
	case s.isZero():
		s.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Data:   useZeroed(s.mat.Data, n*n),
			Uplo:   blas.Upper,
		}
		s.cap = n
		s.SymRankK(s, alpha, x)
	case s.mat.Uplo != blas.Upper:
		panic(badSymTriangle)
	case s.mat.N == n:
		if s == x {
			w := getWorkspaceSym(n, true)
			w.SymRankK(w, alpha, x)
			s.CopySym(w)
			putWorkspaceSym(w)
		} else {
			if rs, ok := x.(RawSymmetricer); ok {
				s.checkOverlap(rs.RawSymmetric())
			}
			// Only zero the upper triangle.
			for i := 0; i < n; i++ {
				ri := i * s.mat.Stride
				zero(s.mat.Data[ri+i : ri+n])
			}
			s.SymRankK(s, alpha, x)
		}
	default:
		panic(matrix.ErrShape)
	}
}

// RankTwo performs a symmmetric rank-two update to the matrix a and stores
// the result in the receiver
//  m = a + alpha * (x * y' + y * x')
func (s *SymDense) RankTwo(a Symmetric, alpha float64, x, y *Vector) {
	n := s.mat.N
	if x.Len() != n {
		panic(matrix.ErrShape)
	}
	if y.Len() != n {
		panic(matrix.ErrShape)
	}
	var w SymDense
	if s == a {
		w = *s
	}
	w.reuseAs(n)
	if s != a {
		if rs, ok := a.(RawSymmetricer); ok {
			s.checkOverlap(rs.RawSymmetric())
		}
		w.CopySym(a)
	}
	blas64.Syr2(alpha, x.mat, y.mat, w.mat)
	*s = w
	return
}

// ScaleSym multiplies the elements of a by f, placing the result in the receiver.
func (s *SymDense) ScaleSym(f float64, a Symmetric) {
	n := a.Symmetric()
	s.reuseAs(n)
	if a, ok := a.(RawSymmetricer); ok {
		amat := a.RawSymmetric()
		if s != a {
			s.checkOverlap(amat)
		}
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				s.mat.Data[i*s.mat.Stride+j] = f * amat.Data[i*amat.Stride+j]
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			s.mat.Data[i*s.mat.Stride+j] = f * a.At(i, j)
		}
	}
}

// SubsetSym extracts a subset of the rows and columns of the matrix a and stores
// the result in-place into the receiver. The resulting matrix size is
// len(set)×len(set). Specifically, at the conclusion of SubsetSym,
// s.At(i, j) equals a.At(set[i], set[j]). Note that the supplied set does not
// have to be a strict subset, dimension repeats are allowed.
func (s *SymDense) SubsetSym(a Symmetric, set []int) {
	n := len(set)
	na := a.Symmetric()
	s.reuseAs(n)
	var restore func()
	if a == s {
		s, restore = s.isolatedWorkspace(a)
		defer restore()
	}

	if a, ok := a.(RawSymmetricer); ok {
		raw := a.RawSymmetric()
		if s != a {
			s.checkOverlap(raw)
		}
		for i := 0; i < n; i++ {
			ssub := s.mat.Data[i*s.mat.Stride : i*s.mat.Stride+n]
			r := set[i]
			rsub := raw.Data[r*raw.Stride : r*raw.Stride+na]
			for j := i; j < n; j++ {
				c := set[j]
				if r <= c {
					ssub[j] = rsub[c]
				} else {
					ssub[j] = raw.Data[c*raw.Stride+r]
				}
			}
		}
		return
	}
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			s.mat.Data[i*s.mat.Stride+j] = a.At(set[i], set[j])
		}
	}
}

// ViewSquare returns a view of the submatrix starting at {i, i} and extending
// for n rows and columns. ViewSquare panics if the view is outside the bounds
// of the receiver.
func (s *SymDense) ViewSquare(i, n int) Matrix {
	sz := s.Symmetric()
	if i < 0 || i > sz || n < 0 || i+n > sz {
		panic(matrix.ErrIndexOutOfRange)
	}
	v := *s
	v.mat.Data = s.mat.Data[i*s.mat.Stride+i : (i+n-1)*s.mat.Stride+i+n]
	v.mat.N = n
	v.cap = s.cap - i
	return &v
}

// GrowSquare returns the receiver expanded by n rows and n columns. If the
// dimensions of the expanded matrix are outside the capacity of the receiver
// a new allocation is made, otherwise not. Note that the receiver itself is
// not modified during the call to GrowSquare.
func (s *SymDense) GrowSquare(n int) Matrix {
	if n < 0 {
		panic(matrix.ErrIndexOutOfRange)
	}
	if n == 0 {
		return s
	}
	var v SymDense
	n += s.mat.N
	if n > s.cap {
		v.mat = blas64.Symmetric{
			N:      n,
			Stride: n,
			Uplo:   blas.Upper,
			Data:   make([]float64, n*n),
		}
		v.cap = n
		// Copy elements, including those not currently visible. Use a temporary
		// structure to avoid modifying the receiver.
		var tmp SymDense
		tmp.mat = blas64.Symmetric{
			N:      s.cap,
			Stride: s.mat.Stride,
			Data:   s.mat.Data,
			Uplo:   s.mat.Uplo,
		}
		tmp.cap = s.cap
		v.CopySym(&tmp)
		return &v
	}
	v.mat = blas64.Symmetric{
		N:      n,
		Stride: s.mat.Stride,
		Uplo:   blas.Upper,
		Data:   s.mat.Data[:(n-1)*s.mat.Stride+n],
	}
	v.cap = s.cap
	return &v
}
