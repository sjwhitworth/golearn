// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

// Quantiler returns the multi-dimensional inverse cumulative distribution function.
// len(x) must equal len(p), and if x is non-nil, len(x) must also equal len(p).
// If x is nil, a new slice will be allocated and returned, otherwise the quantile
// will be stored in-place into x. All of the values of p must be between 0 and 1,
// or Quantile will panic.
type Quantiler interface {
	Quantile(x, p []float64) []float64
}

// LogProber computes the log of the probability of the point x.
type LogProber interface {
	LogProb(x []float64) float64
}

// Rander generates a random number according to the distributon.
// If the input is non-nil, len(x) must equal len(p) and the dimension of the distribution,
// otherwise Quantile will panic.
// If the input is nil, a new slice will be allocated and returned.
type Rander interface {
	Rand(x []float64) []float64
}

// RandLogProber is both a Rander and a LogProber.
type RandLogProber interface {
	Rander
	LogProber
}
