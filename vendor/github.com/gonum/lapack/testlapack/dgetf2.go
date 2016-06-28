package testlapack

import (
	"math/rand"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

type Dgetf2er interface {
	Dgetf2(m, n int, a []float64, lda int, ipiv []int) bool
}

func Dgetf2Test(t *testing.T, impl Dgetf2er) {
	rnd := rand.New(rand.NewSource(1))
	for _, test := range []struct {
		m, n, lda int
	}{
		{10, 10, 0},
		{10, 5, 0},
		{10, 5, 0},

		{10, 10, 20},
		{5, 10, 20},
		{10, 5, 20},
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
		aCopy := make([]float64, len(a))
		copy(aCopy, a)

		mn := min(m, n)
		ipiv := make([]int, mn)
		for i := range ipiv {
			ipiv[i] = rnd.Int()
		}
		ok := impl.Dgetf2(m, n, a, lda, ipiv)
		checkPLU(t, ok, m, n, lda, ipiv, a, aCopy, 1e-14, true)
	}

	// Test with singular matrices (random matrices are almost surely non-singular).
	for _, test := range []struct {
		m, n, lda int
		a         []float64
	}{
		{
			m:   2,
			n:   2,
			lda: 2,
			a: []float64{
				1, 0,
				0, 0,
			},
		},
		{
			m:   2,
			n:   2,
			lda: 2,
			a: []float64{
				1, 5,
				2, 10,
			},
		},
		{
			m:   3,
			n:   3,
			lda: 3,
			// row 3 = row1 + 2 * row2
			a: []float64{
				1, 5, 7,
				2, 10, -3,
				5, 25, 1,
			},
		},
		{
			m:   3,
			n:   4,
			lda: 4,
			// row 3 = row1 + 2 * row2
			a: []float64{
				1, 5, 7, 9,
				2, 10, -3, 11,
				5, 25, 1, 31,
			},
		},
	} {
		if impl.Dgetf2(test.m, test.n, test.a, test.lda, make([]int, min(test.m, test.n))) {
			t.Log("Returned ok with singular matrix.")
		}
	}
}

// checkPLU checks that the PLU factorization contained in factorize matches
// the original matrix contained in original.
func checkPLU(t *testing.T, ok bool, m, n, lda int, ipiv []int, factorized, original []float64, tol float64, print bool) {
	var hasZeroDiagonal bool
	for i := 0; i < min(m, n); i++ {
		if factorized[i*lda+i] == 0 {
			hasZeroDiagonal = true
			break
		}
	}
	if hasZeroDiagonal && ok {
		t.Error("Has a zero diagonal but returned ok")
	}
	if !hasZeroDiagonal && !ok {
		t.Error("Non-zero diagonal but returned !ok")
	}

	// Check that the LU decomposition is correct.
	mn := min(m, n)
	l := make([]float64, m*mn)
	ldl := mn
	u := make([]float64, mn*n)
	ldu := n
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			v := factorized[i*lda+j]
			switch {
			case i == j:
				l[i*ldl+i] = 1
				u[i*ldu+i] = v
			case i > j:
				l[i*ldl+j] = v
			case i < j:
				u[i*ldu+j] = v
			}
		}
	}

	LU := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, m*n),
	}
	U := blas64.General{
		Rows:   mn,
		Cols:   n,
		Stride: ldu,
		Data:   u,
	}
	L := blas64.General{
		Rows:   m,
		Cols:   mn,
		Stride: ldl,
		Data:   l,
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, L, U, 0, LU)

	p := make([]float64, m*m)
	ldp := m
	for i := 0; i < m; i++ {
		p[i*ldp+i] = 1
	}
	for i := len(ipiv) - 1; i >= 0; i-- {
		v := ipiv[i]
		blas64.Swap(m, blas64.Vector{Inc: 1, Data: p[i*ldp:]}, blas64.Vector{Inc: 1, Data: p[v*ldp:]})
	}
	P := blas64.General{
		Rows:   m,
		Cols:   m,
		Stride: m,
		Data:   p,
	}
	aComp := blas64.General{
		Rows:   m,
		Cols:   n,
		Stride: lda,
		Data:   make([]float64, m*lda),
	}
	copy(aComp.Data, factorized)
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, P, LU, 0, aComp)
	if !floats.EqualApprox(aComp.Data, original, tol) {
		if print {
			t.Errorf("PLU multiplication does not match original matrix.\nWant: %v\nGot: %v", original, aComp.Data)
			return
		}
		t.Error("PLU multiplication does not match original matrix.")
	}
}
