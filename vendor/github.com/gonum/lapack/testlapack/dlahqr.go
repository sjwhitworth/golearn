// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas/blas64"
)

type Dlahqrer interface {
	Dlahqr(wantt, wantz bool, n, ilo, ihi int, h []float64, ldh int, wr, wi []float64, iloz, ihiz int, z []float64, ldz int) int
}

func DlahqrTest(t *testing.T, impl Dlahqrer) {
	rnd := rand.New(rand.NewSource(1))

	// Tests that choose the [ilo:ihi+1,ilo:ihi+1] and
	// [iloz:ihiz+1,ilo:ihi+1] blocks randomly.
	for _, wantt := range []bool{true, false} {
		for _, wantz := range []bool{true, false} {
			for _, n := range []int{1, 2, 3, 4, 5, 6, 10, 18, 31, 53} {
				for _, extra := range []int{0, 1, 11} {
					for cas := 0; cas < 100; cas++ {
						ilo := rnd.Intn(n)
						ihi := rnd.Intn(n)
						if ilo > ihi {
							ilo, ihi = ihi, ilo
						}
						iloz := rnd.Intn(ilo + 1)
						ihiz := ihi + rnd.Intn(n-ihi)
						testDlahqr(t, impl, wantt, wantz, n, ilo, ihi, iloz, ihiz, extra, rnd)
					}
				}
			}
		}
	}
	// Tests that make sure that some potentially problematic corner cases,
	// like zero-sized matrix, are covered.
	for _, wantt := range []bool{true, false} {
		for _, wantz := range []bool{true, false} {
			for _, extra := range []int{0, 1, 11} {
				for _, test := range []struct {
					n, ilo, ihi, iloz, ihiz int
				}{
					{
						n:    0,
						ilo:  0,
						ihi:  -1,
						iloz: 0,
						ihiz: -1,
					},
					{
						n:    1,
						ilo:  0,
						ihi:  0,
						iloz: 0,
						ihiz: 0,
					},
					{
						n:    2,
						ilo:  1,
						ihi:  1,
						iloz: 1,
						ihiz: 1,
					},
					{
						n:    2,
						ilo:  0,
						ihi:  1,
						iloz: 0,
						ihiz: 1,
					},
					{
						n:    10,
						ilo:  0,
						ihi:  0,
						iloz: 0,
						ihiz: 0,
					},
					{
						n:    10,
						ilo:  0,
						ihi:  9,
						iloz: 0,
						ihiz: 9,
					},
					{
						n:    10,
						ilo:  0,
						ihi:  1,
						iloz: 0,
						ihiz: 1,
					},
					{
						n:    10,
						ilo:  0,
						ihi:  1,
						iloz: 0,
						ihiz: 9,
					},
					{
						n:    10,
						ilo:  9,
						ihi:  9,
						iloz: 0,
						ihiz: 9,
					},
				} {
					testDlahqr(t, impl, wantt, wantz, test.n, test.ilo, test.ihi, test.iloz, test.ihiz, extra, rnd)
				}
			}
		}
	}
}

func testDlahqr(t *testing.T, impl Dlahqrer, wantt, wantz bool, n, ilo, ihi, iloz, ihiz, extra int, rnd *rand.Rand) {
	const tol = 1e-14

	var z, zCopy blas64.General
	if wantz {
		z = eye(n, n+extra)
		zCopy = cloneGeneral(z)
	}
	h := randomHessenberg(n, n+extra, rnd)
	if ilo > 0 {
		h.Data[ilo*h.Stride+ilo-1] = 0
	}
	if ihi < n-1 {
		h.Data[(ihi+1)*h.Stride+ihi] = 0
	}
	wr := nanSlice(n)
	wi := nanSlice(n)

	info := impl.Dlahqr(wantt, wantz, n, ilo, ihi, h.Data, h.Stride, wr, wi, iloz, ihiz, z.Data, z.Stride)

	prefix := fmt.Sprintf("Case n=%v, ilo=%v, ihi=%v, iloz=%v, ihiz=%v, wantt=%v, wantz=%v, extra=%v", n, ilo, ihi, iloz, ihiz, wantt, wantz, extra)

	if !generalOutsideAllNaN(h) {
		t.Errorf("%v: out-of-range write to H\n%v", prefix, h.Data)
	}
	if !generalOutsideAllNaN(z) {
		t.Errorf("%v: out-of-range write to Z\n%v", prefix, z.Data)
	}

	if wantz {
		// Z should contain the orthogonal matrix U.
		if !isOrthonormal(z) {
			t.Errorf("%v: Z is not orthogonal", prefix)
		}
		// Z should have been modified only in the
		// [iloz:ihiz+1:ilo:ihi+1] block.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if iloz <= i && i <= ihiz && ilo <= j && j <= ihi {
					continue
				}
				if z.Data[i*z.Stride+j] != zCopy.Data[i*zCopy.Stride+j] {
					t.Errorf("%v: Z modified outside of [iloz:ihiz+1,ilo:ihi+1] block", prefix)
				}
			}
		}
	}

	start := ilo // Index of the first computed eigenvalue.
	if info > 0 {
		start = info + 1
	}

	// Check that wr and wi have not been modified outside [start:ihi+1].
	for i := 0; i < start; i++ {
		if !math.IsNaN(wr[i]) {
			t.Errorf("%v: wr modified before [ilo:ihiz+1] block", prefix)
		}
		if !math.IsNaN(wi[i]) {
			t.Errorf("%v: wi modified before [ilo:ihiz+1] block", prefix)
		}
	}
	for i := ihi + 1; i < n; i++ {
		if !math.IsNaN(wr[i]) {
			t.Errorf("%v: wr modified after [ilo:ihiz+1] block", prefix)
		}
		if !math.IsNaN(wi[i]) {
			t.Errorf("%v: wi modified after [ilo:ihiz+1] block", prefix)
		}
	}

	var hasReal bool
	for i := start; i <= ihi; {
		if wi[i] == 0 { // Real eigenvalue.
			hasReal = true

			// Check that the eigenvalue corresponds to a 1×1 block
			// on the diagonal of H.
			if wantt {
				if wr[i] != h.Data[i*h.Stride+i] {
					t.Errorf("%v: wr[%v] != H[%v,%v]", prefix, i, i, i)
				}
				for _, index := range []struct{ r, c int }{
					{i, i - 1},     // h   h   h
					{i + 1, i - 1}, // 0 wr[i] h
					{i + 1, i},     // 0   0   h
				} {
					if index.r >= n || index.c < 0 {
						continue
					}
					if h.Data[index.r*h.Stride+index.c] != 0 {
						t.Errorf("%v: H[%v,%v] != 0", prefix, index.r, index.c)
					}
				}
			}

			i++
			continue
		}

		// Complex eigenvalue.

		// In the conjugate pair the real parts must be equal.
		if wr[i] != wr[i+1] {
			t.Errorf("%v: real part of conjugate pair not equal, i=%v", prefix, i)
		}
		// The first imaginary part must be positive.
		if wi[i] < 0 {
			t.Errorf("%v: wi[%v] not positive", prefix, i)
		}
		// The second imaginary part must be negative with the same
		// magnitude.
		if wi[i] != -wi[i+1] {
			t.Errorf("%v: wi[%v] != wi[%v]", prefix, i, i+1)
		}
		if wantt {
			// Check that wi[i] has the correct value.
			if wr[i] != h.Data[i*h.Stride+i] {
				t.Errorf("%v: wr[%v] != H[%v,%v]", prefix, i, i, i)
			}
			if wr[i] != h.Data[(i+1)*h.Stride+i+1] {
				t.Errorf("%v: wr[%v] != H[%v,%v]", prefix, i, i+1, i+1)
			}
			prod := math.Abs(h.Data[(i+1)*h.Stride+i] * h.Data[i*h.Stride+i+1])
			if math.Abs(math.Sqrt(prod)-wi[i]) > tol {
				t.Errorf("%v: unexpected value of wi[%v]: want %v, got %v", prefix, i, math.Sqrt(prod), wi[i])
			}

			// Check that the corresponding diagonal block is 2×2.
			for _, index := range []struct{ r, c int }{
				{i, i - 1},     //     i
				{i + 1, i - 1}, // h   h      h    h
				{i + 2, i - 1}, // 0 wr[i]    b    h   i
				{i + 2, i},     // 0   c   wr[i+1] h
				{i + 2, i + 1}, // 0   0      0    h
			} {
				if index.r >= n || index.c < 0 {
					continue
				}
				if h.Data[index.r*h.Stride+index.c] != 0 {
					t.Errorf("%v: H[%v,%v] != 0", prefix, index.r, index.c)
				}
			}
		}
		i += 2
	}

	// If the number of found eigenvalues is odd, at least one must be real.
	if (ihi+1-start)%2 != 0 && !hasReal {
		t.Errorf("%v: expected at least one real eigenvalue", prefix)
	}
}
