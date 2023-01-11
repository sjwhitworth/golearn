// Generated code do not edit. Run `go generate`.

// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"errors"
	"fmt"
	math "github.com/gonum/blas/native/internal/math32"
)

func newGeneral32(r, c int) general32 {
	return general32{
		data:   make([]float32, r*c),
		rows:   r,
		cols:   c,
		stride: c,
	}
}

type general32 struct {
	data       []float32
	rows, cols int
	stride     int
}

// adds element-wise into receiver. rows and columns must match
func (g general32) add(h general32) {
	if debug {
		if g.rows != h.rows {
			panic("blas: row size mismatch")
		}
		if g.cols != h.cols {
			panic("blas: col size mismatch")
		}
	}
	for i := 0; i < g.rows; i++ {
		gtmp := g.data[i*g.stride : i*g.stride+g.cols]
		for j, v := range h.data[i*h.stride : i*h.stride+h.cols] {
			gtmp[j] += v
		}
	}
}

// at returns the value at the ith row and jth column. For speed reasons, the
// rows and columns are not bounds checked.
func (g general32) at(i, j int) float32 {
	if debug {
		if i < 0 || i >= g.rows {
			panic("blas: row out of bounds")
		}
		if j < 0 || j >= g.cols {
			panic("blas: col out of bounds")
		}
	}
	return g.data[i*g.stride+j]
}

func (g general32) check(c byte) error {
	if g.rows < 0 {
		return errors.New("blas: rows < 0")
	}
	if g.cols < 0 {
		return errors.New("blas: cols < 0")
	}
	if g.stride < 1 {
		return errors.New("blas: stride < 1")
	}
	if g.stride < g.cols {
		return errors.New("blas: illegal stride")
	}
	if (g.rows-1)*g.stride+g.cols > len(g.data) {
		return fmt.Errorf("blas: index of %c out of range", c)
	}
	return nil
}

func (g general32) clone() general32 {
	data := make([]float32, len(g.data))
	copy(data, g.data)
	return general32{
		data:   data,
		rows:   g.rows,
		cols:   g.cols,
		stride: g.stride,
	}
}

// assumes they are the same size
func (g general32) copy(h general32) {
	if debug {
		if g.rows != h.rows {
			panic("blas: row mismatch")
		}
		if g.cols != h.cols {
			panic("blas: col mismatch")
		}
	}
	for k := 0; k < g.rows; k++ {
		copy(g.data[k*g.stride:(k+1)*g.stride], h.data[k*h.stride:(k+1)*h.stride])
	}
}

func (g general32) equal(a general32) bool {
	if g.rows != a.rows || g.cols != a.cols || g.stride != a.stride {
		return false
	}
	for i, v := range g.data {
		if a.data[i] != v {
			return false
		}
	}
	return true
}

/*
// print is to aid debugging. Commented out to avoid fmt import
func (g general32) print() {
	fmt.Println("r = ", g.rows, "c = ", g.cols, "stride: ", g.stride)
	for i := 0; i < g.rows; i++ {
		fmt.Println(g.data[i*g.stride : (i+1)*g.stride])
	}

}
*/

func (g general32) view(i, j, r, c int) general32 {
	if debug {
		if i < 0 || i+r > g.rows {
			panic("blas: row out of bounds")
		}
		if j < 0 || j+c > g.cols {
			panic("blas: col out of bounds")
		}
	}
	return general32{
		data:   g.data[i*g.stride+j : (i+r-1)*g.stride+j+c],
		rows:   r,
		cols:   c,
		stride: g.stride,
	}
}

func (g general32) equalWithinAbs(a general32, tol float32) bool {
	if g.rows != a.rows || g.cols != a.cols || g.stride != a.stride {
		return false
	}
	for i, v := range g.data {
		if math.Abs(a.data[i]-v) > tol {
			return false
		}
	}
	return true
}
