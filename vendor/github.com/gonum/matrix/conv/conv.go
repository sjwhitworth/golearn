// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package conv provides matrix type interconversion utilities.
package conv

import (
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/cmat128"
	"github.com/gonum/matrix/mat64"
)

// Complex is a complex matrix constructed from two real matrices.
type Complex struct {
	// r and i are not exposed to ensure that
	// their dimensions can not be altered by
	// clients behind our back.
	r, i mat64.Matrix

	// imagSign holds the sign of the imaginary
	// part of the Complex. Valid values are 1 and -1.
	imagSign float64
}

var (
	_ Realer = Complex{}
	_ Imager = Complex{}
)

// NewComplex returns a complex matrix constructed from r and i. At least one of
// r or i must be non-nil otherwise NewComplex will panic. If one of the inputs
// is nil, that part of the complex number will be zero when returned by At.
// If both are non-nil but differ in their sizes, NewComplex will panic.
func NewComplex(r, i mat64.Matrix) Complex {
	if r == nil && i == nil {
		panic("conv: no matrix")
	} else if r != nil && i != nil {
		rr, rc := r.Dims()
		ir, ic := i.Dims()
		if rr != ir || rc != ic {
			panic(matrix.ErrShape)
		}
	}
	return Complex{r: r, i: i, imagSign: 1}
}

// Dims returns the number of rows and columns in the matrix.
func (m Complex) Dims() (r, c int) {
	if m.r == nil {
		return m.i.Dims()
	}
	return m.r.Dims()
}

// At returns the element at row i, column j.
func (m Complex) At(i, j int) complex128 {
	if m.i == nil {
		return complex(m.r.At(i, j), 0)
	}
	if m.r == nil {
		return complex(0, m.imagSign*m.i.At(i, j))
	}
	return complex(m.r.At(i, j), m.imagSign*m.i.At(i, j))
}

// H performs an implicit transpose.
func (m Complex) H() cmat128.Matrix {
	if m.i == nil {
		return Complex{r: m.r.T()}
	}
	if m.r == nil {
		return Complex{i: m.i.T(), imagSign: -m.imagSign}
	}
	return Complex{r: m.r.T(), i: m.i.T(), imagSign: -m.imagSign}
}

// Real returns the real part of the receiver.
func (m Complex) Real() mat64.Matrix { return m.r }

// Imag returns the imaginary part of the receiver.
func (m Complex) Imag() mat64.Matrix { return m.i }

// Realer is a complex matrix that can return its real part.
type Realer interface {
	Real() mat64.Matrix
}

// Imager is a complex matrix that can return its imaginary part.
type Imager interface {
	Imag() mat64.Matrix
}

// Real is the real part of a complex matrix.
type Real struct {
	matrix cmat128.Matrix
}

// NewReal returns a mat64.Matrix representing the real part of m. If m is a Realer,
// the real part is returned.
func NewReal(m cmat128.Matrix) mat64.Matrix {
	if m, ok := m.(Realer); ok {
		return m.Real()
	}
	return Real{m}
}

// Dims returns the number of rows and columns in the matrix.
func (m Real) Dims() (r, c int) { return m.matrix.Dims() }

// At returns the element at row i, column j.
func (m Real) At(i, j int) float64 { return real(m.matrix.At(i, j)) }

// T performs an implicit transpose.
func (m Real) T() mat64.Matrix { return Real{m.matrix.H()} }

// Imag is the imaginary part of a complex matrix.
type Imag struct {
	matrix cmat128.Matrix

	// conjSign holds the sign of the matrix.
	// Valid values are 1 and -1.
	conjSign float64
}

// NewImag returns a mat64.Matrix representing the imaginary part of m. If m is an Imager,
// the imaginary part is returned.
func NewImag(m cmat128.Matrix) mat64.Matrix {
	if m, ok := m.(Imager); ok {
		return m.Imag()
	}
	return Imag{matrix: m, conjSign: 1}
}

// Dims returns the number of rows and columns in the matrix.
func (m Imag) Dims() (r, c int) { return m.matrix.Dims() }

// At returns the element at row i, column j.
func (m Imag) At(i, j int) float64 { return m.conjSign * imag(m.matrix.At(i, j)) }

// T performs an implicit transpose.
func (m Imag) T() mat64.Matrix { return Imag{matrix: m.matrix.H(), conjSign: -m.conjSign} }
