// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/floats"
)

type Dpotf2er interface {
	Dpotf2(ul blas.Uplo, n int, a []float64, lda int) (ok bool)
}

func Dpotf2Test(t *testing.T, impl Dpotf2er) {
	for _, test := range []struct {
		a   [][]float64
		ul  blas.Uplo
		pos bool
		U   [][]float64
	}{
		{
			a: [][]float64{
				{23, 37, 34, 32},
				{108, 71, 48, 48},
				{109, 109, 67, 58},
				{106, 107, 106, 63},
			},
			pos: true,
			U: [][]float64{
				{4.795831523312719, 7.715033320111766, 7.089490077940543, 6.672461249826393},
				{0, 3.387958215439679, -1.976308959006481, -1.026654004678691},
				{0, 0, 3.582364210034111, 2.419258947036024},
				{0, 0, 0, 3.401680257083044},
			},
		},
		{
			a: [][]float64{
				{8, 2},
				{2, 4},
			},
			pos: true,
			U: [][]float64{
				{2.82842712474619, 0.707106781186547},
				{0, 1.870828693386971},
			},
		},
	} {
		testDpotf2(t, impl, test.pos, test.a, test.U, len(test.a[0]), blas.Upper)
		testDpotf2(t, impl, test.pos, test.a, test.U, len(test.a[0])+5, blas.Upper)
		aT := transpose(test.a)
		L := transpose(test.U)
		testDpotf2(t, impl, test.pos, aT, L, len(test.a[0]), blas.Lower)
		testDpotf2(t, impl, test.pos, aT, L, len(test.a[0])+5, blas.Lower)
	}
}

func testDpotf2(t *testing.T, impl Dpotf2er, testPos bool, a, ans [][]float64, stride int, ul blas.Uplo) {
	aFlat := flattenTri(a, stride, ul)
	ansFlat := flattenTri(ans, stride, ul)
	pos := impl.Dpotf2(ul, len(a[0]), aFlat, stride)
	if pos != testPos {
		t.Errorf("Positive definite mismatch: Want %v, Got %v", testPos, pos)
		return
	}
	if testPos && !floats.EqualApprox(ansFlat, aFlat, 1e-14) {
		t.Errorf("Result mismatch: Want %v, Got  %v", ansFlat, aFlat)
	}
}

// flattenTri  with a certain stride. stride must be >= dimension. Puts repeatable
// nonce values in non-accessed places
func flattenTri(a [][]float64, stride int, ul blas.Uplo) []float64 {
	m := len(a)
	n := len(a[0])
	if stride < n {
		panic("bad stride")
	}
	upper := ul == blas.Upper
	v := make([]float64, m*stride)
	count := 1000.0
	for i := 0; i < m; i++ {
		for j := 0; j < stride; j++ {
			if j >= n || (upper && j < i) || (!upper && j > i) {
				// not accessed, so give a unique crazy number
				v[i*stride+j] = count
				count++
				continue
			}
			v[i*stride+j] = a[i][j]
		}
	}
	return v
}

func transpose(a [][]float64) [][]float64 {
	m := len(a)
	n := len(a[0])
	if m != n {
		panic("not square")
	}
	aNew := make([][]float64, m)
	for i := 0; i < m; i++ {
		aNew[i] = make([]float64, n)
	}
	for i := 0; i < m; i++ {
		if len(a[i]) != n {
			panic("bad n size")
		}
		for j := 0; j < n; j++ {
			aNew[j][i] = a[i][j]
		}
	}
	return aNew
}
