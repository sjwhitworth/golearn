// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"github.com/gonum/lapack"
)

type Dgesvder interface {
	Dgesvd(jobU, jobVT lapack.SVDJob, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, lwork int) (ok bool)
}

func DgesvdTest(t *testing.T, impl Dgesvder) {
	rnd := rand.New(rand.NewSource(1))
	// TODO(btracey): Add tests for all of the cases when the SVD implementation
	// is finished.
	// TODO(btracey): Add tests for m > mnthr and n > mnthr when other SVD
	// conditions are implemented. Right now mnthr is 5,000,000 which is too
	// large to create a square matrix of that size.
	for _, test := range []struct {
		m, n, lda, ldu, ldvt int
	}{
		{5, 5, 0, 0, 0},
		{5, 6, 0, 0, 0},
		{6, 5, 0, 0, 0},
		{5, 9, 0, 0, 0},
		{9, 5, 0, 0, 0},

		{5, 5, 10, 11, 12},
		{5, 6, 10, 11, 12},
		{6, 5, 10, 11, 12},
		{5, 5, 10, 11, 12},
		{5, 9, 10, 11, 12},
		{9, 5, 10, 11, 12},

		{300, 300, 0, 0, 0},
		{300, 400, 0, 0, 0},
		{400, 300, 0, 0, 0},
		{300, 600, 0, 0, 0},
		{600, 300, 0, 0, 0},

		{300, 300, 400, 450, 460},
		{300, 400, 500, 550, 560},
		{400, 300, 550, 550, 560},
		{300, 600, 700, 750, 760},
		{600, 300, 700, 750, 760},
	} {
		jobU := lapack.SVDAll
		jobVT := lapack.SVDAll

		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		ldu := test.ldu
		if ldu == 0 {
			ldu = m
		}
		ldvt := test.ldvt
		if ldvt == 0 {
			ldvt = n
		}

		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.NormFloat64()
		}

		u := make([]float64, m*ldu)
		for i := range u {
			u[i] = rnd.NormFloat64()
		}

		vt := make([]float64, n*ldvt)
		for i := range vt {
			vt[i] = rnd.NormFloat64()
		}

		uAllOrig := make([]float64, len(u))
		copy(uAllOrig, u)
		vtAllOrig := make([]float64, len(vt))
		copy(vtAllOrig, vt)
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		s := make([]float64, min(m, n))

		work := make([]float64, 1)
		impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, vt, ldvt, work, -1)

		if !floats.Equal(a, aCopy) {
			t.Errorf("a changed during call to get work length")
		}

		work = make([]float64, int(work[0]))
		impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, vt, ldvt, work, len(work))

		errStr := fmt.Sprintf("m = %v, n = %v, lda = %v, ldu = %v, ldv = %v", m, n, lda, ldu, ldvt)
		svdCheck(t, false, errStr, m, n, s, a, u, ldu, vt, ldvt, aCopy, lda)
		svdCheckPartial(t, impl, lapack.SVDAll, errStr, uAllOrig, vtAllOrig, aCopy, m, n, a, lda, s, u, ldu, vt, ldvt, work, false)

		// Test InPlace
		jobU = lapack.SVDInPlace
		jobVT = lapack.SVDInPlace
		copy(a, aCopy)
		copy(u, uAllOrig)
		copy(vt, vtAllOrig)

		impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, vt, ldvt, work, len(work))
		svdCheck(t, true, errStr, m, n, s, a, u, ldu, vt, ldvt, aCopy, lda)
		svdCheckPartial(t, impl, lapack.SVDInPlace, errStr, uAllOrig, vtAllOrig, aCopy, m, n, a, lda, s, u, ldu, vt, ldvt, work, false)
	}
}

// svdCheckPartial checks that the singular values and vectors are computed when
// not all of them are computed.
func svdCheckPartial(t *testing.T, impl Dgesvder, job lapack.SVDJob, errStr string, uAllOrig, vtAllOrig, aCopy []float64, m, n int, a []float64, lda int, s, u []float64, ldu int, vt []float64, ldvt int, work []float64, shortWork bool) {
	rnd := rand.New(rand.NewSource(1))
	jobU := job
	jobVT := job
	// Compare the singular values when computed with {SVDNone, SVDNone.}
	sCopy := make([]float64, len(s))
	copy(sCopy, s)
	copy(a, aCopy)
	for i := range s {
		s[i] = rnd.Float64()
	}
	tmp1 := make([]float64, 1)
	tmp2 := make([]float64, 1)
	jobU = lapack.SVDNone
	jobVT = lapack.SVDNone

	impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, tmp1, ldu, tmp2, ldvt, work, -1)
	work = make([]float64, int(work[0]))
	lwork := len(work)
	if shortWork {
		lwork--
	}
	ok := impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, tmp1, ldu, tmp2, ldvt, work, lwork)
	if !ok {
		t.Errorf("Dgesvd did not complete successfully")
	}
	if !floats.EqualApprox(s, sCopy, 1e-10) {
		t.Errorf("Singular value mismatch when singular vectors not computed: %s", errStr)
	}
	// Check that the singular vectors are correctly computed when the other
	// is none.
	uAll := make([]float64, len(u))
	copy(uAll, u)
	vtAll := make([]float64, len(vt))
	copy(vtAll, vt)

	// Copy the original vectors so the data outside the matrix bounds is the same.
	copy(u, uAllOrig)
	copy(vt, vtAllOrig)

	jobU = job
	jobVT = lapack.SVDNone
	copy(a, aCopy)
	for i := range s {
		s[i] = rnd.Float64()
	}
	impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, tmp2, ldvt, work, -1)
	work = make([]float64, int(work[0]))
	lwork = len(work)
	if shortWork {
		lwork--
	}
	impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, u, ldu, tmp2, ldvt, work, len(work))
	if !floats.EqualApprox(uAll, u, 1e-10) {
		t.Errorf("U mismatch when VT is not computed: %s", errStr)
	}
	if !floats.EqualApprox(s, sCopy, 1e-10) {
		t.Errorf("Singular value mismatch when U computed VT not")
	}
	jobU = lapack.SVDNone
	jobVT = job
	copy(a, aCopy)
	for i := range s {
		s[i] = rnd.Float64()
	}
	impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, tmp1, ldu, vt, ldvt, work, -1)
	work = make([]float64, int(work[0]))
	lwork = len(work)
	if shortWork {
		lwork--
	}
	impl.Dgesvd(jobU, jobVT, m, n, a, lda, s, tmp1, ldu, vt, ldvt, work, len(work))
	if !floats.EqualApprox(vtAll, vt, 1e-10) {
		t.Errorf("VT mismatch when U is not computed: %s", errStr)
	}
	if !floats.EqualApprox(s, sCopy, 1e-10) {
		t.Errorf("Singular value mismatch when VT computed U not")
	}
}

// svdCheck checks that the singular value decomposition correctly multiplies back
// to the original matrix.
func svdCheck(t *testing.T, thin bool, errStr string, m, n int, s, a, u []float64, ldu int, vt []float64, ldvt int, aCopy []float64, lda int) {
	sigma := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	for i := 0; i < min(m, n); i++ {
		sigma.Data[i*sigma.Stride+i] = s[i]
	}

	uMat := blas64.General{
		Rows:   m,
		Cols:   m,
		Stride: ldu,
		Data:   u,
	}
	vTMat := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: ldvt,
		Data:   vt,
	}
	if thin {
		sigma.Rows = min(m, n)
		sigma.Cols = min(m, n)
		uMat.Cols = min(m, n)
		vTMat.Rows = min(m, n)
	}

	tmp := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	ans := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: lda,
		Data:   make([]float64, m*lda),
	}
	copy(ans.Data, a)

	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, uMat, sigma, 0, tmp)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, tmp, vTMat, 0, ans)

	if !floats.EqualApprox(ans.Data, aCopy, 1e-8) {
		t.Errorf("Decomposition mismatch. Trim = %v, %s", thin, errStr)
	}

	if !thin {
		// Check that U and V are orthogonal.
		for i := 0; i < uMat.Rows; i++ {
			for j := i + 1; j < uMat.Rows; j++ {
				dot := blas64.Dot(uMat.Cols,
					blas64.Vector{Inc: 1, Data: uMat.Data[i*uMat.Stride:]},
					blas64.Vector{Inc: 1, Data: uMat.Data[j*uMat.Stride:]},
				)
				if dot > 1e-8 {
					t.Errorf("U not orthogonal %s", errStr)
				}
			}
		}
		for i := 0; i < vTMat.Rows; i++ {
			for j := i + 1; j < vTMat.Rows; j++ {
				dot := blas64.Dot(vTMat.Cols,
					blas64.Vector{Inc: 1, Data: vTMat.Data[i*vTMat.Stride:]},
					blas64.Vector{Inc: 1, Data: vTMat.Data[j*vTMat.Stride:]},
				)
				if dot > 1e-8 {
					t.Errorf("V not orthogonal %s", errStr)
				}
			}
		}
	}
}
