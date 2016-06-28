// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/gonum/floats"
)

type Dsterfer interface {
	Dgetrfer
	Dsterf(n int, d, e []float64) (ok bool)
}

func DsterfTest(t *testing.T, impl Dsterfer) {
	// Hand coded tests.
	for cas, test := range []struct {
		d []float64
		e []float64
		n int

		ans []float64
	}{
		// Computed from Fortran code.
		{
			d:   []float64{1, 3, 4, 6},
			e:   []float64{2, 4, 5},
			n:   4,
			ans: []float64{11.046227528488854, 4.795922173417400, -2.546379458290125, 0.704229756383872},
		},
	} {
		n := test.n
		d := make([]float64, len(test.d))
		copy(d, test.d)
		e := make([]float64, len(test.e))
		copy(e, test.e)
		ok := impl.Dsterf(n, d, e)
		if !ok {
			t.Errorf("Case %d, Eigenvalue decomposition failed", cas)
			continue
		}
		ans := make([]float64, len(test.ans))
		copy(ans, test.ans)
		sort.Float64s(ans)
		if !floats.EqualApprox(ans, d, 1e-10) {
			t.Errorf("eigenvalue mismatch")
		}
	}

	rnd := rand.New(rand.NewSource(1))
	// Probabilistic tests.
	for _, n := range []int{4, 6, 10} {
		for cas := 0; cas < 10; cas++ {
			d := make([]float64, n)
			for i := range d {
				d[i] = rnd.NormFloat64()
			}
			dCopy := make([]float64, len(d))
			copy(dCopy, d)
			e := make([]float64, n-1)
			for i := range e {
				e[i] = rnd.NormFloat64()
			}
			eCopy := make([]float64, len(e))
			copy(eCopy, e)

			ok := impl.Dsterf(n, d, e)
			if !ok {
				t.Errorf("Eigenvalue decomposition failed")
				continue
			}

			// Test that the eigenvalues are sorted.
			if !sort.Float64sAreSorted(d) {
				t.Errorf("Values are not sorted")
			}

			// Construct original tridagional matrix.
			lda := n
			a := make([]float64, n*lda)
			for i := 0; i < n; i++ {
				a[i*lda+i] = dCopy[i]
				if i != n-1 {
					a[i*lda+i+1] = eCopy[i]
					a[(i+1)*lda+i] = eCopy[i]
				}
			}

			asub := make([]float64, len(a))
			ipiv := make([]int, n)

			// Test that they are actually eigenvalues by computing the
			// determinant of A - λI.
			// TODO(btracey): Replace this test with a more numerically stable
			// test.
			for _, lambda := range d {
				copy(asub, a)
				for i := 0; i < n; i++ {
					asub[i*lda+i] -= lambda
				}

				// Compute LU.
				ok := impl.Dgetrf(n, n, asub, lda, ipiv)
				if !ok {
					// Definitely singular.
					continue
				}
				// Compute determinant.
				var logdet float64
				for i := 0; i < n; i++ {
					v := asub[i*lda+i]
					logdet += math.Log(math.Abs(v))
				}
				if math.Exp(logdet) > 2 {
					t.Errorf("Incorrect singular value. n = %d, cas = %d, det = %v", n, cas, math.Exp(logdet))
				}
			}
		}
	}
}
