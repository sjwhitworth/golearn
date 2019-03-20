// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"math"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

// Dirichlet implements the Dirichlet probability distribution.
//
// The Dirichlet distribution is a continuous probability distribution that
// generates elements over the probability simplex, i.e. ||x||_1 = 1. The Dirichlet
// distribution is the conjugate prior to the categorical distribution and the
// multivariate version of the beta distribution. The probability of a point x is
//  1/Beta(α) \prod_i x_i^(α_i - 1)
// where Beta(α) is the multivariate Beta function (see the mathext package).
//
// For more information see https://en.wikipedia.org/wiki/Dirichlet_distribution
type Dirichlet struct {
	alpha []float64
	dim   int
	src   rand.Source

	lbeta    float64
	sumAlpha float64
}

// NewDirichlet creates a new dirichlet distribution with the given parameters alpha.
// NewDirichlet will panic if len(alpha) == 0, or if any alpha is <= 0.
func NewDirichlet(alpha []float64, src rand.Source) *Dirichlet {
	dim := len(alpha)
	if dim == 0 {
		panic(badZeroDimension)
	}
	for _, v := range alpha {
		if v <= 0 {
			panic("dirichlet: non-positive alpha")
		}
	}
	a := make([]float64, len(alpha))
	copy(a, alpha)
	d := &Dirichlet{
		alpha: a,
		dim:   dim,
		src:   src,
	}
	d.lbeta, d.sumAlpha = d.genLBeta(a)
	return d
}

// CovarianceMatrix returns the covariance matrix of the distribution. Upon
// return, the value at element {i, j} of the covariance matrix is equal to
// the covariance of the i^th and j^th variables.
//  covariance(i, j) = E[(x_i - E[x_i])(x_j - E[x_j])]
// If the input matrix is nil a new matrix is allocated, otherwise the result
// is stored in-place into the input.
func (d *Dirichlet) CovarianceMatrix(cov *mat.SymDense) *mat.SymDense {
	if cov == nil {
		cov = mat.NewSymDense(d.Dim(), nil)
	} else if cov.Symmetric() == 0 {
		*cov = *(cov.GrowSquare(d.dim).(*mat.SymDense))
	} else if cov.Symmetric() != d.dim {
		panic("normal: input matrix size mismatch")
	}
	scale := 1 / (d.sumAlpha * d.sumAlpha * (d.sumAlpha + 1))
	for i := 0; i < d.dim; i++ {
		ai := d.alpha[i]
		v := ai * (d.sumAlpha - ai) * scale
		cov.SetSym(i, i, v)
		for j := i + 1; j < d.dim; j++ {
			aj := d.alpha[j]
			v := -ai * aj * scale
			cov.SetSym(i, j, v)
		}
	}
	return cov
}

// genLBeta computes the generalized LBeta function.
func (d *Dirichlet) genLBeta(alpha []float64) (lbeta, sumAlpha float64) {
	for _, alpha := range d.alpha {
		lg, _ := math.Lgamma(alpha)
		lbeta += lg
		sumAlpha += alpha
	}
	lg, _ := math.Lgamma(sumAlpha)
	return lbeta - lg, sumAlpha
}

// Dim returns the dimension of the distribution.
func (d *Dirichlet) Dim() int {
	return d.dim
}

// LogProb computes the log of the pdf of the point x.
//
// It does not check that ||x||_1 = 1.
func (d *Dirichlet) LogProb(x []float64) float64 {
	dim := d.dim
	if len(x) != dim {
		panic(badSizeMismatch)
	}
	var lprob float64
	for i, x := range x {
		lprob += (d.alpha[i] - 1) * math.Log(x)
	}
	lprob -= d.lbeta
	return lprob
}

// Mean returns the mean of the probability distribution at x. If the
// input argument is nil, a new slice will be allocated, otherwise the result
// will be put in-place into the receiver.
func (d *Dirichlet) Mean(x []float64) []float64 {
	x = reuseAs(x, d.dim)
	copy(x, d.alpha)
	floats.Scale(1/d.sumAlpha, x)
	return x
}

// Prob computes the value of the probability density function at x.
func (d *Dirichlet) Prob(x []float64) float64 {
	return math.Exp(d.LogProb(x))
}

// Rand generates a random number according to the distributon.
// If the input slice is nil, new memory is allocated, otherwise the result is stored
// in place.
func (d *Dirichlet) Rand(x []float64) []float64 {
	x = reuseAs(x, d.dim)
	for i := range x {
		x[i] = distuv.Gamma{Alpha: d.alpha[i], Beta: 1, Src: d.src}.Rand()
	}
	sum := floats.Sum(x)
	floats.Scale(1/sum, x)
	return x
}
