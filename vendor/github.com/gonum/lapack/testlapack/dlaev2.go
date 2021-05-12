package testlapack

import (
	"math"
	"math/rand"
	"testing"
)

type Dlaev2er interface {
	Dlaev2(a, b, c float64) (rt1, rt2, cs1, sn1 float64)
}

func Dlaev2Test(t *testing.T, impl Dlaev2er) {
	rnd := rand.New(rand.NewSource(1))
	for trial := 0; trial < 100; trial++ {
		a := rnd.NormFloat64()
		b := rnd.NormFloat64()
		c := rnd.NormFloat64()

		rt1, rt2, cs1, sn1 := impl.Dlaev2(a, b, c)
		tmp := mul2by2([2][2]float64{{cs1, sn1}, {-sn1, cs1}}, [2][2]float64{{a, b}, {b, c}})
		ans := mul2by2(tmp, [2][2]float64{{cs1, -sn1}, {sn1, cs1}})
		if math.Abs(ans[0][0]-rt1) > 1e-14 {
			t.Errorf("Largest eigenvalue mismatch. Returned %v, mul %v", rt1, ans[0][0])
		}
		if math.Abs(ans[1][0]) > 1e-14 || math.Abs(ans[0][1]) > 1e-14 {
			t.Errorf("Non-zero off diagonal. ans[1][0] = %v, ans[0][1] = %v", ans[1][0], ans[0][1])
		}
		if math.Abs(ans[1][1]-rt2) > 1e-14 {
			t.Errorf("Smallest eigenvalue mismatch. Returned %v, mul %v", rt2, ans[1][1])
		}
	}
}

func mul2by2(a, b [2][2]float64) [2][2]float64 {
	var c [2][2]float64
	c[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0]
	c[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1]
	c[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0]
	c[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1]
	return c
}
