// Copyright ©2016 The gonum Authors. All rights reserved.
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
)

type Dlasy2er interface {
	Dlasy2(tranl, tranr bool, isgn, n1, n2 int, tl []float64, ldtl int, tr []float64, ldtr int, b []float64, ldb int, x []float64, ldx int) (scale, xnorm float64, ok bool)
}

func Dlasy2Test(t *testing.T, impl Dlasy2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, tranl := range []bool{true, false} {
		for _, tranr := range []bool{true, false} {
			for _, isgn := range []int{1, -1} {
				for _, n1 := range []int{0, 1, 2} {
					for _, n2 := range []int{0, 1, 2} {
						for _, extra := range []int{0, 1, 2, 13} {
							for cas := 0; cas < 1000; cas++ {
								testDlasy2(t, impl, tranl, tranr, isgn, n1, n2, extra, rnd)
							}
						}
					}
				}
			}
		}
	}
}

func testDlasy2(t *testing.T, impl Dlasy2er, tranl, tranr bool, isgn, n1, n2, extra int, rnd *rand.Rand) {
	const tol = 1e-11

	tl := randomGeneral(n1, n1, n1+extra, rnd)
	tr := randomGeneral(n2, n2, n2+extra, rnd)
	b := randomGeneral(n1, n2, n2+extra, rnd)
	x := randomGeneral(n1, n2, n2+extra, rnd)

	scale, xnorm, ok := impl.Dlasy2(tranl, tranr, isgn, n1, n2, tl.Data, tl.Stride, tr.Data, tr.Stride, b.Data, b.Stride, x.Data, x.Stride)
	if scale > 1 {
		t.Errorf("invalid value of scale, want <= 1, got %v", scale)
	}
	if n1 == 0 || n2 == 0 {
		return
	}

	prefix := fmt.Sprintf("Case n1=%v, n2=%v, isgn=%v", n1, n2, isgn)

	// Check any invalid modifications of x.
	if !generalOutsideAllNaN(x) {
		t.Errorf("%v: out-of-range write to x\n%v", prefix, x.Data)
	}

	var xnormWant float64
	for i := 0; i < n1; i++ {
		var rowsum float64
		for j := 0; j < n2; j++ {
			rowsum += math.Abs(x.Data[i*x.Stride+j])
		}
		if rowsum > xnormWant {
			xnormWant = rowsum
		}
	}
	if xnormWant != xnorm {
		t.Errorf("%v: unexpected xnorm: want %v, got %v", prefix, xnormWant, xnorm)
	}

	// Multiply b by scale to get the wanted right-hand side.
	for i := 0; i < n1; i++ {
		for j := 0; j < n2; j++ {
			b.Data[i*b.Stride+j] *= scale
		}
	}
	// Compute the wanted left-hand side.
	lhsWant := randomGeneral(n1, n2, n2, rnd)
	if tranl {
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, tl, x, 0, lhsWant)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tl, x, 0, lhsWant)
	}
	if tranr {
		blas64.Gemm(blas.NoTrans, blas.Trans, float64(isgn), x, tr, 1, lhsWant)
	} else {
		blas64.Gemm(blas.NoTrans, blas.NoTrans, float64(isgn), x, tr, 1, lhsWant)
	}
	// Compare them.
	for i := 0; i < n1; i++ {
		for j := 0; j < n2; j++ {
			diff := lhsWant.Data[i*lhsWant.Stride+j] - b.Data[i*b.Stride+j]
			if math.Abs(diff) > tol && ok {
				t.Errorf("%v: unexpected result, diff[%v,%v]=%v", prefix, i, j, diff)
			}
		}
	}
}
