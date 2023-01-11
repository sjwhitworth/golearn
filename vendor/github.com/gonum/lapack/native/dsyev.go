// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"math"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/lapack"
)

// Dsyev computes all eigenvalues and, optionally, the eigenvectors of a real
// symmetric matrix A.
//
// w contains the eigenvalues in ascending order upon return. w must have length
// at least n, and Dsyev will panic otherwise.
//
// On entry, a contains the elements of the symmetric matrix A in the triangular
// portion specified by uplo. If jobz == lapack.EigDecomp a contains the
// orthonormal eigenvectors of A on exit, otherwise on exit the specified
// triangular region is overwritten.
//
// work is temporary storage, and lwork specifies the usable memory length. At minimum,
// lwork >= 3*n-1, and Dsyev will panic otherwise. The amount of blocking is
// limited by the usable length. If lwork == -1, instead of computing Dsyev the
// optimal work length is stored into work[0].
func (impl Implementation) Dsyev(jobz lapack.EigComp, uplo blas.Uplo, n int, a []float64, lda int, w, work []float64, lwork int) (ok bool) {
	checkMatrix(n, n, a, lda)
	upper := uplo == blas.Upper
	wantz := jobz == lapack.EigDecomp
	var opts string
	if upper {
		opts = "U"
	} else {
		opts = "L"
	}
	nb := impl.Ilaenv(1, "DSYTRD", opts, n, -1, -1, -1)
	lworkopt := max(1, (nb+2)*n)
	work[0] = float64(lworkopt)
	if lwork == -1 {
		return
	}
	if len(work) < lwork {
		panic(badWork)
	}
	if lwork < 3*n-1 {
		panic(badWork)
	}
	if n == 0 {
		return true
	}
	if n == 1 {
		w[0] = a[0]
		work[0] = 2
		if wantz {
			a[0] = 1
		}
		return true
	}
	safmin := dlamchS
	eps := dlamchP
	smlnum := safmin / eps
	bignum := 1 / smlnum
	rmin := math.Sqrt(smlnum)
	rmax := math.Sqrt(bignum)

	// Scale matrix to allowable range, if necessary.
	anrm := impl.Dlansy(lapack.MaxAbs, uplo, n, a, lda, work)
	scaled := false
	var sigma float64
	if anrm > 0 && anrm < rmin {
		scaled = true
		sigma = rmin / anrm
	} else if anrm > rmax {
		scaled = true
		sigma = rmax / anrm
	}
	if scaled {
		kind := lapack.LowerTri
		if upper {
			kind = lapack.UpperTri
		}
		impl.Dlascl(kind, 0, 0, 1, sigma, n, n, a, lda)
	}
	var inde int
	indtau := inde + n
	indwork := indtau + n
	llwork := lwork - indwork
	impl.Dsytrd(uplo, n, a, lda, w, work[inde:], work[indtau:], work[indwork:], llwork)

	// For eigenvalues only, call Dsterf. For eigenvectors, first call Dorgtr
	// to generate the orthogonal matrix, then call Dsteqr.
	if !wantz {
		ok = impl.Dsterf(n, w, work[inde:])
	} else {
		impl.Dorgtr(uplo, n, a, lda, work[indtau:], work[indwork:], llwork)
		ok = impl.Dsteqr(jobz, n, w, work[inde:], a, lda, work[indtau:])
	}
	if !ok {
		return false
	}

	// If the matrix was scaled, then rescale eigenvalues appropriately.
	if scaled {
		bi := blas64.Implementation()
		bi.Dscal(n, 1/sigma, w, 1)
	}
	work[0] = float64(lworkopt)
	return true
}
