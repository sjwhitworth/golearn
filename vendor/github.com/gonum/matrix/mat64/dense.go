// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/matrix"
)

var (
	dense *Dense

	_ Matrix  = dense
	_ Mutable = dense

	_ Cloner       = dense
	_ Viewer       = dense
	_ RowViewer    = dense
	_ ColViewer    = dense
	_ RawRowViewer = dense
	_ Grower       = dense

	_ RawMatrixSetter = dense
	_ RawMatrixer     = dense

	_ Reseter = dense
)

// Dense is a dense matrix representation.
type Dense struct {
	mat blas64.General

	capRows, capCols int
}

// NewDense creates a new matrix of type Dense with dimensions r and c.
// If the mat argument is nil, a new data slice is allocated.
//
// The data must be arranged in row-major order, i.e. the (i*c + j)-th
// element in mat is the {i, j}-th element in the matrix.
func NewDense(r, c int, mat []float64) *Dense {
	if mat != nil && r*c != len(mat) {
		panic(matrix.ErrShape)
	}
	if mat == nil {
		mat = make([]float64, r*c)
	}
	return &Dense{
		mat: blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   mat,
		},
		capRows: r,
		capCols: c,
	}
}

// reuseAs resizes an empty matrix to a r×c matrix,
// or checks that a non-empty matrix is r×c.
//
// reuseAs must be kept in sync with reuseAsZeroed.
func (m *Dense) reuseAs(r, c int) {
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat64.Error.
		panic("mat64: caps not correctly set")
	}
	if m.isZero() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   use(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
		return
	}
	if r != m.mat.Rows || c != m.mat.Cols {
		panic(matrix.ErrShape)
	}
}

// reuseAsZeroed resizes an empty matrix to a r×c matrix,
// or checks that a non-empty matrix is r×c. It zeroes
// all the elements of the matrix.
//
// reuseAsZeroed must be kept in sync with reuseAs.
func (m *Dense) reuseAsZeroed(r, c int) {
	if m.mat.Rows > m.capRows || m.mat.Cols > m.capCols {
		// Panic as a string, not a mat64.Error.
		panic("mat64: caps not correctly set")
	}
	if m.isZero() {
		m.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			Data:   useZeroed(m.mat.Data, r*c),
		}
		m.capRows = r
		m.capCols = c
		return
	}
	if r != m.mat.Rows || c != m.mat.Cols {
		panic(matrix.ErrShape)
	}
	for i := 0; i < r; i++ {
		zero(m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c])
	}
}

// untranspose untransposes a matrix if applicable. If a is an Untransposer, then
// untranspose returns the underlying matrix and true. If it is not, then it returns
// the input matrix and false.
func untranspose(a Matrix) (Matrix, bool) {
	if ut, ok := a.(Untransposer); ok {
		return ut.Untranspose(), true
	}
	return a, false
}

// isolatedWorkspace returns a new dense matrix w with the size of a and
// returns a callback to defer which performs cleanup at the return of the call.
// This should be used when a method receiver is the same pointer as an input argument.
func (m *Dense) isolatedWorkspace(a Matrix) (w *Dense, restore func()) {
	r, c := a.Dims()
	w = getWorkspace(r, c, false)
	return w, func() {
		m.Copy(w)
		putWorkspace(w)
	}
}

func (m *Dense) isZero() bool {
	// It must be the case that m.Dims() returns
	// zeros in this case. See comment in Reset().
	return m.mat.Stride == 0
}

// asTriDense returns a TriDense with the given size and side. The backing data
// of the TriDense is the same as the receiver.
func (m *Dense) asTriDense(n int, diag blas.Diag, uplo blas.Uplo) *TriDense {
	return &TriDense{
		mat: blas64.Triangular{
			N:      n,
			Stride: m.mat.Stride,
			Data:   m.mat.Data,
			Uplo:   uplo,
			Diag:   diag,
		},
		cap: n,
	}
}

// DenseCopyOf returns a newly allocated copy of the elements of a.
func DenseCopyOf(a Matrix) *Dense {
	d := &Dense{}
	d.Clone(a)
	return d
}

// SetRawMatrix sets the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in b.
func (m *Dense) SetRawMatrix(b blas64.General) {
	m.capRows, m.capCols = b.Rows, b.Cols
	m.mat = b
}

// RawMatrix returns the underlying blas64.General used by the receiver.
// Changes to elements in the receiver following the call will be reflected
// in returned blas64.General.
func (m *Dense) RawMatrix() blas64.General { return m.mat }

// Dims returns the number of rows and columns in the matrix.
func (m *Dense) Dims() (r, c int) { return m.mat.Rows, m.mat.Cols }

// Caps returns the number of rows and columns in the backing matrix.
func (m *Dense) Caps() (r, c int) { return m.capRows, m.capCols }

// T performs an implicit transpose by returning the receiver inside a Transpose.
func (m *Dense) T() Matrix {
	return Transpose{m}
}

// ColView returns a Vector reflecting the column j, backed by the matrix data.
//
// See ColViewer for more information.
func (m *Dense) ColView(j int) *Vector {
	if j >= m.mat.Cols || j < 0 {
		panic(matrix.ErrColAccess)
	}
	return &Vector{
		mat: blas64.Vector{
			Inc:  m.mat.Stride,
			Data: m.mat.Data[j : (m.mat.Rows-1)*m.mat.Stride+j+1],
		},
		n: m.mat.Rows,
	}
}

// SetCol sets the values in the specified column of the matrix to the values
// in src. len(src) must equal the number of rows in the receiver.
func (m *Dense) SetCol(j int, src []float64) {
	if j >= m.mat.Cols || j < 0 {
		panic(matrix.ErrColAccess)
	}
	if len(src) != m.mat.Rows {
		panic(matrix.ErrColLength)
	}

	blas64.Copy(m.mat.Rows,
		blas64.Vector{Inc: 1, Data: src},
		blas64.Vector{Inc: m.mat.Stride, Data: m.mat.Data[j:]},
	)
}

// SetRow sets the values in the specified rows of the matrix to the values
// in src. len(src) must equal the number of columns in the receiver.
func (m *Dense) SetRow(i int, src []float64) {
	if i >= m.mat.Rows || i < 0 {
		panic(matrix.ErrRowAccess)
	}
	if len(src) != m.mat.Cols {
		panic(matrix.ErrRowLength)
	}

	copy(m.rawRowView(i), src)
}

// RowView returns row i of the matrix data represented as a column vector,
// backed by the matrix data.
//
// See RowViewer for more information.
func (m *Dense) RowView(i int) *Vector {
	if i >= m.mat.Rows || i < 0 {
		panic(matrix.ErrRowAccess)
	}
	return &Vector{
		mat: blas64.Vector{
			Inc:  1,
			Data: m.rawRowView(i),
		},
		n: m.mat.Cols,
	}
}

// RawRowView returns a slice backed by the same array as backing the
// receiver.
func (m *Dense) RawRowView(i int) []float64 {
	if i >= m.mat.Rows || i < 0 {
		panic(matrix.ErrRowAccess)
	}
	return m.rawRowView(i)
}

func (m *Dense) rawRowView(i int) []float64 {
	return m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+m.mat.Cols]
}

// View returns a new Matrix that shares backing data with the receiver.
// The new matrix is located from row i, column j extending r rows and c
// columns. View panics if the view is outside the bounds of the receiver.
func (m *Dense) View(i, j, r, c int) Matrix {
	mr, mc := m.Dims()
	if i < 0 || i >= mr || j < 0 || j >= mc || r <= 0 || i+r > mr || c <= 0 || j+c > mc {
		panic(matrix.ErrIndexOutOfRange)
	}
	t := *m
	t.mat.Data = t.mat.Data[i*t.mat.Stride+j : (i+r-1)*t.mat.Stride+(j+c)]
	t.mat.Rows = r
	t.mat.Cols = c
	t.capRows -= i
	t.capCols -= j
	return &t
}

// Grow returns the receiver expanded by r rows and c columns. If the dimensions
// of the expanded matrix are outside the capacities of the receiver a new
// allocation is made, otherwise not. Note the receiver itself is not modified
// during the call to Grow.
func (m *Dense) Grow(r, c int) Matrix {
	if r < 0 || c < 0 {
		panic(matrix.ErrIndexOutOfRange)
	}
	if r == 0 && c == 0 {
		return m
	}

	r += m.mat.Rows
	c += m.mat.Cols

	var t Dense
	switch {
	case m.mat.Rows == 0 || m.mat.Cols == 0:
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: c,
			// We zero because we don't know how the matrix will be used.
			// In other places, the mat is immediately filled with a result;
			// this is not the case here.
			Data: useZeroed(m.mat.Data, r*c),
		}
	case r > m.capRows || c > m.capCols:
		cr := max(r, m.capRows)
		cc := max(c, m.capCols)
		t.mat = blas64.General{
			Rows:   r,
			Cols:   c,
			Stride: cc,
			Data:   make([]float64, cr*cc),
		}
		t.capRows = cr
		t.capCols = cc
		// Copy the complete matrix over to the new matrix.
		// Including elements not currently visible. Use a temporary structure
		// to avoid modifying the receiver.
		var tmp Dense
		tmp.mat = blas64.General{
			Rows:   m.mat.Rows,
			Cols:   m.mat.Cols,
			Stride: m.mat.Stride,
			Data:   m.mat.Data,
		}
		tmp.capRows = m.capRows
		tmp.capCols = m.capCols
		t.Copy(&tmp)
		return &t
	default:
		t.mat = blas64.General{
			Data:   m.mat.Data[:(r-1)*m.mat.Stride+c],
			Rows:   r,
			Cols:   c,
			Stride: m.mat.Stride,
		}
	}
	t.capRows = r
	t.capCols = c
	return &t
}

// Reset zeros the dimensions of the matrix so that it can be reused as the
// receiver of a dimensionally restricted operation.
//
// See the Reseter interface for more information.
func (m *Dense) Reset() {
	// No change of Stride, Rows and Cols to 0
	// may be made unless all are set to 0.
	m.mat.Rows, m.mat.Cols, m.mat.Stride = 0, 0, 0
	m.capRows, m.capCols = 0, 0
	m.mat.Data = m.mat.Data[:0]
}

// Clone makes a copy of a into the receiver, overwriting the previous value of
// the receiver. The clone operation does not make any restriction on shape and
// will not cause shadowing.
//
// See the Cloner interface for more information.
func (m *Dense) Clone(a Matrix) {
	r, c := a.Dims()
	mat := blas64.General{
		Rows:   r,
		Cols:   c,
		Stride: c,
	}
	m.capRows, m.capCols = r, c

	aU, trans := untranspose(a)
	switch aU := aU.(type) {
	case RawMatrixer:
		amat := aU.RawMatrix()
		mat.Data = make([]float64, r*c)
		if trans {
			for i := 0; i < r; i++ {
				blas64.Copy(c,
					blas64.Vector{Inc: amat.Stride, Data: amat.Data[i : i+(c-1)*amat.Stride+1]},
					blas64.Vector{Inc: 1, Data: mat.Data[i*c : (i+1)*c]})
			}
		} else {
			for i := 0; i < r; i++ {
				copy(mat.Data[i*c:(i+1)*c], amat.Data[i*amat.Stride:i*amat.Stride+c])
			}
		}
	case *Vector:
		amat := aU.mat
		mat.Data = make([]float64, aU.n)
		blas64.Copy(aU.n,
			blas64.Vector{Inc: amat.Inc, Data: amat.Data},
			blas64.Vector{Inc: 1, Data: mat.Data})
	default:
		mat.Data = make([]float64, r*c)
		w := *m
		w.mat = mat
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				w.set(i, j, a.At(i, j))
			}
		}
		*m = w
		return
	}
	m.mat = mat
}

// Copy makes a copy of elements of a into the receiver. It is similar to the
// built-in copy; it copies as much as the overlap between the two matrices and
// returns the number of rows and columns it copied. If a aliases the receiver
// and is a transposed Dense or Vector, with a non-unitary increment, Copy will
// panic.
//
// See the Copier interface for more information.
func (m *Dense) Copy(a Matrix) (r, c int) {
	r, c = a.Dims()
	if a == m {
		return r, c
	}
	r = min(r, m.mat.Rows)
	c = min(c, m.mat.Cols)
	if r == 0 || c == 0 {
		return 0, 0
	}

	aU, trans := untranspose(a)
	switch aU := aU.(type) {
	case RawMatrixer:
		amat := aU.RawMatrix()
		if trans {
			if amat.Stride != 1 {
				m.checkOverlap(amat)
			}
			for i := 0; i < r; i++ {
				blas64.Copy(c,
					blas64.Vector{Inc: amat.Stride, Data: amat.Data[i : i+(c-1)*amat.Stride+1]},
					blas64.Vector{Inc: 1, Data: m.mat.Data[i*m.mat.Stride : i*m.mat.Stride+c]})
			}
		} else {
			switch o := offset(m.mat.Data, amat.Data); {
			case o < 0:
				for i := r - 1; i >= 0; i-- {
					copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
				}
			case o > 0:
				for i := 0; i < r; i++ {
					copy(m.mat.Data[i*m.mat.Stride:i*m.mat.Stride+c], amat.Data[i*amat.Stride:i*amat.Stride+c])
				}
			default:
				// Nothing to do.
			}
		}
	case *Vector:
		var n, stride int
		amat := aU.mat
		if trans {
			if amat.Inc != 1 {
				m.checkOverlap(aU.asGeneral())
			}
			n = c
			stride = 1
		} else {
			n = r
			stride = m.mat.Stride
		}
		if amat.Inc == 1 && stride == 1 {
			copy(m.mat.Data, amat.Data[:n])
			break
		}
		switch o := offset(m.mat.Data, amat.Data); {
		case o < 0:
			blas64.Copy(n,
				blas64.Vector{Inc: -amat.Inc, Data: amat.Data},
				blas64.Vector{Inc: -stride, Data: m.mat.Data})
		case o > 0:
			blas64.Copy(n,
				blas64.Vector{Inc: amat.Inc, Data: amat.Data},
				blas64.Vector{Inc: stride, Data: m.mat.Data})
		default:
			// Nothing to do.
		}
	default:
		for i := 0; i < r; i++ {
			for j := 0; j < c; j++ {
				m.set(i, j, a.At(i, j))
			}
		}
	}

	return r, c
}

// Stack appends the rows of b onto the rows of a, placing the result into the
// receiver with b placed in the greater indexed rows. Stack will panic if the
// two input matrices do not have the same number of columns or the constructed
// stacked matrix is not the same shape as the receiver.
func (m *Dense) Stack(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ac != bc || m == a || m == b {
		panic(matrix.ErrShape)
	}

	m.reuseAs(ar+br, ac)

	m.Copy(a)
	w := m.View(ar, 0, br, bc).(*Dense)
	w.Copy(b)
}

// Augment creates the augmented matrix of a and b, where b is placed in the
// greater indexed columns. Augment will panic if the two input matrices do
// not have the same number of rows or the constructed augmented matrix is
// not the same shape as the receiver.
func (m *Dense) Augment(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br || m == a || m == b {
		panic(matrix.ErrShape)
	}

	m.reuseAs(ar, ac+bc)

	m.Copy(a)
	w := m.View(0, ac, br, bc).(*Dense)
	w.Copy(b)
}
