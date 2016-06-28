package testlapack

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dtrti2er interface {
	Dtrti2(uplo blas.Uplo, diag blas.Diag, n int, a []float64, lda int)
}

func Dtrti2Test(t *testing.T, impl Dtrti2er) {
	for _, test := range []struct {
		a    []float64
		n    int
		uplo blas.Uplo
		diag blas.Diag
		ans  []float64
	}{
		{
			a: []float64{
				2, 3, 4,
				0, 5, 6,
				8, 0, 8},
			n:    3,
			uplo: blas.Upper,
			diag: blas.NonUnit,
			ans: []float64{
				0.5, -0.3, -0.025,
				0, 0.2, -0.15,
				8, 0, 0.125,
			},
		},
		{
			a: []float64{
				5, 3, 4,
				0, 7, 6,
				10, 0, 8},
			n:    3,
			uplo: blas.Upper,
			diag: blas.Unit,
			ans: []float64{
				5, -3, 14,
				0, 7, -6,
				10, 0, 8,
			},
		},
		{
			a: []float64{
				2, 0, 0,
				3, 5, 0,
				4, 6, 8},
			n:    3,
			uplo: blas.Lower,
			diag: blas.NonUnit,
			ans: []float64{
				0.5, 0, 0,
				-0.3, 0.2, 0,
				-0.025, -0.15, 0.125,
			},
		},
		{
			a: []float64{
				1, 0, 0,
				3, 1, 0,
				4, 6, 1},
			n:    3,
			uplo: blas.Lower,
			diag: blas.Unit,
			ans: []float64{
				1, 0, 0,
				-3, 1, 0,
				14, -6, 1,
			},
		},
	} {
		impl.Dtrti2(test.uplo, test.diag, test.n, test.a, test.n)
		if !floats.EqualApprox(test.ans, test.a, 1e-14) {
			t.Errorf("Matrix inverse mismatch. Want %v, got %v.", test.ans, test.a)
		}
	}
	rnd := rand.New(rand.NewSource(1))
	bi := blas64.Implementation()
	for _, uplo := range []blas.Uplo{blas.Upper} {
		for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
			for _, test := range []struct {
				n, lda int
			}{
				{3, 0},
				{3, 5},
			} {
				n := test.n
				lda := test.lda
				if lda == 0 {
					lda = n
				}
				a := make([]float64, n*lda)
				for i := range a {
					a[i] = rnd.Float64()
				}
				aCopy := make([]float64, len(a))
				copy(aCopy, a)
				impl.Dtrti2(uplo, diag, n, a, lda)
				if uplo == blas.Upper {
					for i := 1; i < n; i++ {
						for j := 0; j < i; j++ {
							aCopy[i*lda+j] = 0
							a[i*lda+j] = 0
						}
					}
				} else {
					for i := 1; i < n; i++ {
						for j := i + 1; j < n; j++ {
							aCopy[i*lda+j] = 0
							a[i*lda+j] = 0
						}
					}
				}
				if diag == blas.Unit {
					for i := 0; i < n; i++ {
						a[i*lda+i] = 1
						aCopy[i*lda+i] = 1
					}
				}
				ans := make([]float64, len(a))
				bi.Dgemm(blas.NoTrans, blas.NoTrans, n, n, n, 1, a, lda, aCopy, lda, 0, ans, lda)
				iseye := true
				for i := 0; i < n; i++ {
					for j := 0; j < n; j++ {
						if i == j {
							if math.Abs(ans[i*lda+i]-1) > 1e-14 {
								iseye = false
								break
							}
						} else {
							if math.Abs(ans[i*lda+j]) > 1e-14 {
								iseye = false
								break
							}
						}
					}
				}
				if !iseye {
					t.Errorf("inv(A) * A != I. Upper = %v, unit = %v, ans = %v", uplo == blas.Upper, diag == blas.Unit, ans)
				}
			}
		}
	}
}
