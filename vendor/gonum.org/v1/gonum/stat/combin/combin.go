// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package combin

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

const (
	badNegInput          = "combin: negative input"
	badSetSize           = "combin: n < k"
	badInput             = "combin: wrong input slice length"
	nonpositiveDimension = "combin: non-positive dimension"
)

// Binomial returns the binomial coefficient of (n,k), also commonly referred to
// as "n choose k".
//
// The binomial coefficient, C(n,k), is the number of unordered combinations of
// k elements in a set that is n elements big, and is defined as
//
//  C(n,k) = n!/((n-k)!k!)
//
// n and k must be non-negative with n >= k, otherwise Binomial will panic.
// No check is made for overflow.
func Binomial(n, k int) int {
	if n < 0 || k < 0 {
		panic(badNegInput)
	}
	if n < k {
		panic(badSetSize)
	}
	// (n,k) = (n, n-k)
	if k > n/2 {
		k = n - k
	}
	b := 1
	for i := 1; i <= k; i++ {
		b = (n - k + i) * b / i
	}
	return b
}

// GeneralizedBinomial returns the generalized binomial coefficient of (n, k),
// defined as
//  Γ(n+1) / (Γ(k+1) Γ(n-k+1))
// where Γ is the Gamma function. GeneralizedBinomial is useful for continuous
// relaxations of the binomial coefficient, or when the binomial coefficient value
// may overflow int. In the latter case, one may use math/big for an exact
// computation.
//
// n and k must be non-negative with n >= k, otherwise GeneralizedBinomial will panic.
func GeneralizedBinomial(n, k float64) float64 {
	return math.Exp(LogGeneralizedBinomial(n, k))
}

// LogGeneralizedBinomial returns the log of the generalized binomial coefficient.
// See GeneralizedBinomial for more information.
func LogGeneralizedBinomial(n, k float64) float64 {
	if n < 0 || k < 0 {
		panic(badNegInput)
	}
	if n < k {
		panic(badSetSize)
	}
	a, _ := math.Lgamma(n + 1)
	b, _ := math.Lgamma(k + 1)
	c, _ := math.Lgamma(n - k + 1)
	return a - b - c
}

// CombinationGenerator generates combinations iteratively. Combinations may be
// called to generate all combinations collectively.
type CombinationGenerator struct {
	n         int
	k         int
	previous  []int
	remaining int
}

// NewCombinationGenerator returns a CombinationGenerator for generating the
// combinations of k elements from a set of size n.
//
// n and k must be non-negative with n >= k, otherwise NewCombinationGenerator
// will panic.
func NewCombinationGenerator(n, k int) *CombinationGenerator {
	return &CombinationGenerator{
		n:         n,
		k:         k,
		remaining: Binomial(n, k),
	}
}

// Next advances the iterator if there are combinations remaining to be generated,
// and returns false if all combinations have been generated. Next must be called
// to initialize the first value before calling Combination or Combination will
// panic. The value returned by Combination is only changed during calls to Next.
func (c *CombinationGenerator) Next() bool {
	if c.remaining <= 0 {
		// Next is called before combination, so c.remaining is set to zero before
		// Combination is called. Thus, Combination cannot panic on zero, and a
		// second sentinel value is needed.
		c.remaining = -1
		return false
	}
	if c.previous == nil {
		c.previous = make([]int, c.k)
		for i := range c.previous {
			c.previous[i] = i
		}
	} else {
		nextCombination(c.previous, c.n, c.k)
	}
	c.remaining--
	return true
}

// Combination generates the next combination. If next is non-nil, it must have
// length k and the result will be stored in-place into combination. If combination
// is nil a new slice will be allocated and returned. If all of the combinations
// have already been constructed (Next() returns false), Combination will panic.
//
// Next must be called to initialize the first value before calling Combination
// or Combination will panic. The value returned by Combination is only changed
// during calls to Next.
func (c *CombinationGenerator) Combination(combination []int) []int {
	if c.remaining == -1 {
		panic("combin: all combinations have been generated")
	}
	if c.previous == nil {
		panic("combin: Combination called before Next")
	}
	if combination == nil {
		combination = make([]int, c.k)
	}
	if len(combination) != c.k {
		panic(badInput)
	}
	copy(combination, c.previous)
	return combination
}

// Combinations generates all of the combinations of k elements from a
// set of size n. The returned slice has length Binomial(n,k) and each inner slice
// has length k.
//
// n and k must be non-negative with n >= k, otherwise Combinations will panic.
//
// CombinationGenerator may alternatively be used to generate the combinations
// iteratively instead of collectively.
func Combinations(n, k int) [][]int {
	combins := Binomial(n, k)
	data := make([][]int, combins)
	if len(data) == 0 {
		return data
	}
	data[0] = make([]int, k)
	for i := range data[0] {
		data[0][i] = i
	}
	for i := 1; i < combins; i++ {
		next := make([]int, k)
		copy(next, data[i-1])
		nextCombination(next, n, k)
		data[i] = next
	}
	return data
}

// nextCombination generates the combination after s, overwriting the input value.
func nextCombination(s []int, n, k int) {
	for j := k - 1; j >= 0; j-- {
		if s[j] == n+j-k {
			continue
		}
		s[j]++
		for l := j + 1; l < k; l++ {
			s[l] = s[j] + l - j
		}
		break
	}
}

// Cartesian returns the cartesian product of the slices in data. The Cartesian
// product of two sets is the set of all combinations of the items. For example,
// given the input
//  [][]float64{{1,2},{3,4},{5,6}}
// the returned matrix will be
//  [ 1 3 5 ]
//  [ 1 3 6 ]
//  [ 1 4 5 ]
//  [ 1 4 6 ]
//  [ 2 3 5 ]
//  [ 2 3 6 ]
//  [ 2 4 5 ]
//  [ 2 4 6 ]
// If dst is nil, a new matrix will be allocated and returned, otherwise the number
// of rows of dst must equal \prod_i len(data[i]), and the number of columns in
// dst must equal len(data). Cartesian also panics if len(data) = 0.
func Cartesian(dst *mat.Dense, data [][]float64) *mat.Dense {
	if len(data) == 0 {
		panic("combin: empty data input")
	}
	cols := len(data)
	rows := 1
	lens := make([]int, cols)
	for i, d := range data {
		v := len(d)
		lens[i] = v
		rows *= v
	}
	if dst == nil {
		dst = mat.NewDense(rows, cols, nil)
	}
	r, c := dst.Dims()
	if r != rows || c != cols {
		panic("combin: destination matrix size mismatch")
	}
	idxs := make([]int, cols)
	for i := 0; i < rows; i++ {
		SubFor(idxs, i, lens)
		for j := 0; j < len(data); j++ {
			dst.Set(i, j, data[j][idxs[j]])
		}
	}
	return dst
}

// IdxFor converts a multi-dimensional index into a linear index for a
// multi-dimensional space. sub specifies the index for each dimension, and dims
// specifies the size of each dimension. IdxFor is the inverse of SubFor.
// IdxFor panics if any of the entries of sub are negative, any of the entries
// of dim are non-positive, or if sub[i] >= dims[i] for any i.
func IdxFor(sub, dims []int) int {
	// The index returned is "row-major", that is the last index of sub is
	// continuous.
	var idx int
	stride := 1
	for i := len(dims) - 1; i >= 0; i-- {
		v := sub[i]
		d := dims[i]
		if d <= 0 {
			panic(nonpositiveDimension)
		}
		if v < 0 || v >= d {
			panic("combin: invalid subscript")
		}
		idx += v * stride
		stride *= d
	}
	return idx
}

// SubFor returns the multi-dimensional subscript for the input linear index to
// the multi-dimensional space. dims specifies the size of each dimension, and
// idx specifies the linear index. SubFor is the inverse of IdxFor.
//
// If sub is non-nil the result is stored in-place into sub, and SubFor will panic
// if len(sub) != len(dims). If sub is nil a new slice of the appropriate length
// is allocated. SubFor panics if idx < 0 or if idx is greater than or equal to
// the product of the dimensions.
func SubFor(sub []int, idx int, dims []int) []int {
	if sub == nil {
		sub = make([]int, len(dims))
	}
	if len(sub) != len(dims) {
		panic(badInput)
	}
	if idx < 0 {
		panic(badNegInput)
	}
	stride := 1
	for i := len(dims) - 1; i >= 1; i-- {
		stride *= dims[i]
	}
	for i := 0; i < len(dims)-1; i++ {
		v := idx / stride
		d := dims[i]
		if d < 0 {
			panic(nonpositiveDimension)
		}
		if v >= dims[i] {
			panic("combin: index too large")
		}
		sub[i] = v
		idx -= v * stride
		stride /= dims[i+1]
	}
	if idx > dims[len(sub)-1] {
		panic("combin: index too large")
	}
	sub[len(sub)-1] = idx
	return sub
}
