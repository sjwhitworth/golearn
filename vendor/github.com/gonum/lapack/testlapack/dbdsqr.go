// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dbdsqrer interface {
	Dbdsqr(uplo blas.Uplo, n, ncvt, nru, ncc int, d, e, vt []float64, ldvt int, u []float64, ldu int, c []float64, ldc int, work []float64) (ok bool)
}

func DbdsqrTest(t *testing.T, impl Dbdsqrer) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	_ = bi
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, test := range []struct {
			n, ncvt, nru, ncc, ldvt, ldu, ldc int
		}{
			{5, 5, 5, 5, 0, 0, 0},
			{10, 10, 10, 10, 0, 0, 0},
			{10, 11, 12, 13, 0, 0, 0},
			{20, 13, 12, 11, 0, 0, 0},

			{5, 5, 5, 5, 6, 7, 8},
			{10, 10, 10, 10, 30, 40, 50},
			{10, 12, 11, 13, 30, 40, 50},
			{20, 12, 13, 11, 30, 40, 50},

			{130, 130, 130, 500, 900, 900, 500},
		} {
			for cas := 0; cas < 100; cas++ {
				n := test.n
				ncvt := test.ncvt
				nru := test.nru
				ncc := test.ncc
				ldvt := test.ldvt
				ldu := test.ldu
				ldc := test.ldc
				if ldvt == 0 {
					ldvt = ncvt
				}
				if ldu == 0 {
					ldu = n
				}
				if ldc == 0 {
					ldc = ncc
				}

				d := make([]float64, n)
				for i := range d {
					d[i] = rnd.NormFloat64()
				}
				e := make([]float64, n-1)
				for i := range e {
					e[i] = rnd.NormFloat64()
				}
				dCopy := make([]float64, len(d))
				copy(dCopy, d)
				eCopy := make([]float64, len(e))
				copy(eCopy, e)
				work := make([]float64, 4*n)
				for i := range work {
					work[i] = rnd.NormFloat64()
				}

				// First test the decomposition of the bidiagonal matrix. Set
				// pt and u equal to I with the correct size. At the result
				// of Dbdsqr, p and u  will contain the data of P^T and Q, which
				// will be used in the next step to test the multiplication
				// with Q and VT.

				q := make([]float64, n*n)
				ldq := n
				pt := make([]float64, n*n)
				ldpt := n
				for i := 0; i < n; i++ {
					q[i*ldq+i] = 1
				}
				for i := 0; i < n; i++ {
					pt[i*ldpt+i] = 1
				}

				ok := impl.Dbdsqr(uplo, n, n, n, 0, d, e, pt, ldpt, q, ldq, nil, 0, work)

				isUpper := uplo == blas.Upper
				errStr := fmt.Sprintf("isUpper = %v, n = %v, ncvt = %v, nru = %v, ncc = %v", isUpper, n, ncvt, nru, ncc)
				if !ok {
					t.Errorf("Unexpected Dbdsqr failure: %s", errStr)
				}

				bMat := constructBidiagonal(uplo, n, dCopy, eCopy)
				sMat := constructBidiagonal(uplo, n, d, e)

				tmp := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, n*n),
				}
				ansMat := blas64.General{
					Rows:   n,
					Cols:   n,
					Stride: n,
					Data:   make([]float64, n*n),
				}

				bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, q, ldq, sMat.Data, sMat.Stride, 0, tmp.Data, tmp.Stride)
				bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, tmp.Data, tmp.Stride, pt, ldpt, 0, ansMat.Data, ansMat.Stride)

				same := true
				for i := 0; i < n; i++ {
					for j := 0; j < n; j++ {
						if !floats.EqualWithinAbsOrRel(ansMat.Data[i*ansMat.Stride+j], bMat.Data[i*bMat.Stride+j], 1e-8, 1e-8) {
							same = false
						}
					}
				}
				if !same {
					t.Errorf("Bidiagonal mismatch. %s", errStr)
				}
				if !sort.IsSorted(sort.Reverse(sort.Float64Slice(d))) {
					t.Errorf("D is not sorted. %s", errStr)
				}

				// The above computed the real P and Q. Now input data for V^T,
				// U, and C to check that the multiplications happen properly.
				dAns := make([]float64, len(d))
				copy(dAns, d)
				eAns := make([]float64, len(e))
				copy(eAns, e)

				u := make([]float64, nru*ldu)
				for i := range u {
					u[i] = rnd.NormFloat64()
				}
				uCopy := make([]float64, len(u))
				copy(uCopy, u)
				vt := make([]float64, n*ldvt)
				for i := range vt {
					vt[i] = rnd.NormFloat64()
				}
				vtCopy := make([]float64, len(vt))
				copy(vtCopy, vt)
				c := make([]float64, n*ldc)
				for i := range c {
					c[i] = rnd.NormFloat64()
				}
				cCopy := make([]float64, len(c))
				copy(cCopy, c)

				// Reset input data
				copy(d, dCopy)
				copy(e, eCopy)
				impl.Dbdsqr(uplo, n, ncvt, nru, ncc, d, e, vt, ldvt, u, ldu, c, ldc, work)

				// Check result.
				if !floats.EqualApprox(d, dAns, 1e-14) {
					t.Errorf("D mismatch second time. %s", errStr)
				}
				if !floats.EqualApprox(e, eAns, 1e-14) {
					t.Errorf("E mismatch second time. %s", errStr)
				}
				ans := make([]float64, len(vtCopy))
				copy(ans, vtCopy)
				ldans := ldvt
				bi.Dgemm(blas.NoTrans, blas.NoTrans, n, ncvt, n, 1, pt, ldpt, vtCopy, ldvt, 0, ans, ldans)
				if !floats.EqualApprox(ans, vt, 1e-10) {
					t.Errorf("Vt result mismatch. %s", errStr)
				}
				ans = make([]float64, len(uCopy))
				copy(ans, uCopy)
				ldans = ldu
				bi.Dgemm(blas.NoTrans, blas.NoTrans, nru, n, n, 1, uCopy, ldu, q, ldq, 0, ans, ldans)
				if !floats.EqualApprox(ans, u, 1e-10) {
					t.Errorf("U result mismatch. %s", errStr)
				}
				ans = make([]float64, len(cCopy))
				copy(ans, cCopy)
				ldans = ldc
				bi.Dgemm(blas.Trans, blas.NoTrans, n, ncc, n, 1, q, ldq, cCopy, ldc, 0, ans, ldans)
				if !floats.EqualApprox(ans, c, 1e-10) {
					t.Errorf("C result mismatch. %s", errStr)
				}
			}
		}
	}
}
