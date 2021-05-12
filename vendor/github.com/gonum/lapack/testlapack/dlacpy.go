package testlapack

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
)

type Dlacpyer interface {
	Dlacpy(uplo blas.Uplo, m, n int, a []float64, lda int, b []float64, ldb int)
}

func DlacpyTest(t *testing.T, impl Dlacpyer) {
	rnd := rand.New(rand.NewSource(1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower, blas.All} {
		for _, test := range []struct {
			m, n, lda, ldb int
		}{
			{3, 5, 0, 0},
			{5, 5, 0, 0},
			{7, 5, 0, 0},

			{3, 5, 10, 12},
			{5, 5, 10, 12},
			{7, 5, 10, 12},
		} {
			m := test.m
			n := test.n
			lda := test.lda
			if lda == 0 {
				lda = n
			}
			ldb := test.ldb
			if ldb == 0 {
				ldb = n
			}
			a := make([]float64, m*lda)
			for i := range a {
				a[i] = rnd.Float64()
			}
			b := make([]float64, m*ldb)
			for i := range b {
				b[i] = rnd.Float64()
			}
			impl.Dlacpy(uplo, m, n, a, lda, b, ldb)
			equal := true
			switch uplo {
			case blas.Upper:
				for i := 0; i < m; i++ {
					for j := i; j < n; j++ {
						if b[i*ldb+j] != a[i*lda+j] {
							equal = false
							goto DoneCheck
						}
					}
				}
			case blas.Lower:
				for i := 0; i < m; i++ {
					for j := 0; j < min(i, n); j++ {
						if b[i*ldb+j] != a[i*lda+j] {
							equal = false
							goto DoneCheck
						}
					}
				}
			case blas.All:
				for i := 0; i < m; i++ {
					for j := 0; j < n; j++ {
						if b[i*ldb+j] != a[i*lda+j] {
							equal = false
							goto DoneCheck
						}
					}
				}
			}
		DoneCheck:
			if !equal {
				fmt.Println(blas.Lower)
				t.Errorf("Matrices not equal after copy. Uplo = %d, m = %d, n = %d", uplo, m, n)
			}
		}
	}
}
