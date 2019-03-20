// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/bound"
)

// Uniform represents a multivariate uniform distribution.
type Uniform struct {
	bounds []bound.Bound
	dim    int
	rnd    *rand.Rand
}

// NewUniform creates a new uniform distribution with the given bounds.
func NewUniform(bnds []bound.Bound, src rand.Source) *Uniform {
	dim := len(bnds)
	if dim == 0 {
		panic(badZeroDimension)
	}
	for _, b := range bnds {
		if b.Max < b.Min {
			panic("uniform: maximum less than minimum")
		}
	}
	u := &Uniform{
		bounds: make([]bound.Bound, dim),
		dim:    dim,
	}
	if src != nil {
		u.rnd = rand.New(src)
	}
	for i, b := range bnds {
		u.bounds[i].Min = b.Min
		u.bounds[i].Max = b.Max
	}
	return u
}

// NewUnitUniform creates a new Uniform distribution over the dim-dimensional
// unit hypercube. That is, a uniform distribution where each dimension has
// Min = 0 and Max = 1.
func NewUnitUniform(dim int, src rand.Source) *Uniform {
	if dim <= 0 {
		panic(nonPosDimension)
	}
	bounds := make([]bound.Bound, dim)
	for i := range bounds {
		bounds[i].Min = 0
		bounds[i].Max = 1
	}
	u := Uniform{
		bounds: bounds,
		dim:    dim,
	}
	if src != nil {
		u.rnd = rand.New(src)
	}
	return &u
}

// Bounds returns the bounds on the variables of the distribution. If the input
// is nil, a new slice is allocated and returned. If the input is non-nil, then
// the bounds are stored in-place into the input argument, and Bounds will panic
// if len(bounds) != u.Dim().
func (u *Uniform) Bounds(bounds []bound.Bound) []bound.Bound {
	if bounds == nil {
		bounds = make([]bound.Bound, u.Dim())
	}
	if len(bounds) != u.Dim() {
		panic(badInputLength)
	}
	copy(bounds, u.bounds)
	return bounds
}

// CDF returns the multidimensional cumulative distribution function of the
// probability distribution at the point x. If p is non-nil, the CDF is stored
// in-place into the first argument, otherwise a new slice is allocated and
// returned.
//
// CDF will panic if len(x) is not equal to the dimension of the distribution,
// or if p is non-nil and len(p) is not equal to the dimension of the distribution.
func (u *Uniform) CDF(p, x []float64) []float64 {
	if len(x) != u.dim {
		panic(badSizeMismatch)
	}
	if p == nil {
		p = make([]float64, u.dim)
	}
	if len(p) != u.dim {
		panic(badSizeMismatch)
	}
	for i, v := range x {
		if v < u.bounds[i].Min {
			p[i] = 0
		} else if v > u.bounds[i].Max {
			p[i] = 1
		} else {
			p[i] = (v - u.bounds[i].Min) / (u.bounds[i].Max - u.bounds[i].Min)
		}
	}
	return p
}

// Dim returns the dimension of the distribution.
func (u *Uniform) Dim() int {
	return u.dim
}

// Entropy returns the differential entropy of the distribution.
func (u *Uniform) Entropy() float64 {
	// Entropy is log of the volume.
	var logVol float64
	for _, b := range u.bounds {
		logVol += math.Log(b.Max - b.Min)
	}
	return logVol
}

// LogProb computes the log of the pdf of the point x.
func (u *Uniform) LogProb(x []float64) float64 {
	dim := u.dim
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	var logprob float64
	for i, b := range u.bounds {
		if x[i] < b.Min || x[i] > b.Max {
			return math.Inf(-1)
		}
		logprob -= math.Log(b.Max - b.Min)
	}
	return logprob
}

// Mean returns the mean of the probability distribution at x. If the
// input argument is nil, a new slice will be allocated, otherwise the result
// will be put in-place into the receiver.
func (u *Uniform) Mean(x []float64) []float64 {
	x = reuseAs(x, u.dim)
	for i, b := range u.bounds {
		x[i] = (b.Max + b.Min) / 2
	}
	return x
}

// Prob computes the value of the probability density function at x.
func (u *Uniform) Prob(x []float64) float64 {
	return math.Exp(u.LogProb(x))
}

// Rand generates a random number according to the distributon.
// If the input slice is nil, new memory is allocated, otherwise the result is stored
// in place.
func (u *Uniform) Rand(x []float64) []float64 {
	x = reuseAs(x, u.dim)
	if u.rnd == nil {
		for i, b := range u.bounds {
			x[i] = rand.Float64()*(b.Max-b.Min) + b.Min
		}
		return x
	}
	for i, b := range u.bounds {
		x[i] = u.rnd.Float64()*(b.Max-b.Min) + b.Min
	}
	return x
}

// Quantile returns the multi-dimensional inverse cumulative distribution function.
// len(x) must equal len(p), and if x is non-nil, len(x) must also equal len(p).
// If x is nil, a new slice will be allocated and returned, otherwise the quantile
// will be stored in-place into x. All of the values of p must be between 0 and 1,
// or Quantile will panic.
func (u *Uniform) Quantile(x, p []float64) []float64 {
	if len(p) != u.dim {
		panic(badSizeMismatch)
	}
	if x == nil {
		x = make([]float64, u.dim)
	}
	if len(x) != u.dim {
		panic(badSizeMismatch)
	}
	for i, v := range p {
		if v < 0 || v > 1 {
			panic(badQuantile)
		}
		x[i] = v*(u.bounds[i].Max-u.bounds[i].Min) + u.bounds[i].Min
	}
	return x
}
