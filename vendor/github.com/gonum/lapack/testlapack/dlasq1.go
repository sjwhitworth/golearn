// Copyright ©2015 The gonum Authors. All rights reserved.
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

func printDlasq1FortranInput(d, e, work []float64, n int) {
	printFortranArray(d, "d")
	printFortranArray(e, "e")
	printFortranArray(work, "work")
	fmt.Println("n = ", n)
	fmt.Println("info = 0")
}

type Dlasq1er interface {
	Dlasq1(n int, d, e, work []float64) int
	Dgetrfer
}

func Dlasq1Test(t *testing.T, impl Dlasq1er) {
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	// TODO(btracey): Increase the size of this test when we have a more numerically
	// stable way to test the singular values.
	for _, n := range []int{1, 2, 5, 8} {
		work := make([]float64, 4*n)
		d := make([]float64, n)
		e := make([]float64, n-1)
		for cas := 0; cas < 1; cas++ {
			for i := range work {
				work[i] = rnd.Float64()
			}
			for i := range d {
				d[i] = rnd.NormFloat64() + 10
			}
			for i := range e {
				e[i] = rnd.NormFloat64()
			}
			ldm := n
			m := make([]float64, n*ldm)
			// Set up the matrix
			for i := 0; i < n; i++ {
				m[i*ldm+i] = d[i]
				if i != n-1 {
					m[(i+1)*ldm+i] = e[i]
				}
			}

			ldmm := n
			mm := make([]float64, n*ldmm)
			bi.Dgemm(blas.Trans, blas.NoTrans, n, n, n, 1, m, ldm, m, ldm, 0, mm, ldmm)

			impl.Dlasq1(n, d, e, work)

			// Check that they are singular values. The
			// singular values are the square roots of the
			// eigenvalues of X^T * X
			mmCopy := make([]float64, len(mm))
			copy(mmCopy, mm)
			ipiv := make([]int, n)
			for elem, sv := range d[0:n] {
				copy(mm, mmCopy)
				lambda := sv * sv
				for i := 0; i < n; i++ {
					mm[i*ldm+i] -= lambda
				}

				// Compute LU.
				ok := impl.Dgetrf(n, n, mm, ldmm, ipiv)
				if !ok {
					// Definitely singular.
					continue
				}
				// Compute determinant
				var logdet float64
				for i := 0; i < n; i++ {
					v := mm[i*ldm+i]
					logdet += math.Log(math.Abs(v))
				}
				if math.Exp(logdet) > 2 {
					t.Errorf("Incorrect singular value. n = %d, cas = %d, elem = %d, det = %v", n, cas, elem, math.Exp(logdet))
				}
			}
		}
	}
}
