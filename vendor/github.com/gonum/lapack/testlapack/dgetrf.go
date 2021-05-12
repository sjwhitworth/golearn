package testlapack

import (
	"math/rand"
	"testing"
)

type Dgetrfer interface {
	Dgetrf(m, n int, a []float64, lda int, ipiv []int) bool
}

func DgetrfTest(t *testing.T, impl Dgetrfer) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{10, 5, 0},
		{5, 10, 0},
		{10, 10, 0},
		{300, 5, 0},
		{3, 500, 0},
		{4, 5, 0},
		{300, 200, 0},
		{204, 300, 0},
		{1, 3000, 0},
		{3000, 1, 0},
		{10, 5, 20},
		{5, 10, 20},
		{10, 10, 20},
		{300, 5, 400},
		{3, 500, 600},
		{200, 200, 300},
		{300, 200, 300},
		{204, 300, 400},
		{1, 3000, 4000},
		{3000, 1, 4000},
	} {
		m := test.m
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		a := make([]float64, m*lda)
		for i := range a {
			a[i] = rnd.Float64()
		}
		mn := min(m, n)
		ipiv := make([]int, mn)
		for i := range ipiv {
			ipiv[i] = rnd.Int()
		}

		// Cannot compare the outputs of Dgetrf and Dgetf2 because the pivoting may
		// happen differently. Instead check that the LPQ factorization is correct.
		aCopy := make([]float64, len(a))
		copy(aCopy, a)
		ok := impl.Dgetrf(m, n, a, lda, ipiv)
		checkPLU(t, ok, m, n, lda, ipiv, a, aCopy, 1e-10, false)
	}
}
