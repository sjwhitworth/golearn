package testblas

import "testing"

type Dgerer interface {
	Dger(m, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int)
}

func DgerTest(t *testing.T, blasser Dgerer) {
	for _, test := range []struct {
		name        string
		a           [][]float64
		m           int
		n           int
		x           []float64
		y           []float64
		incX        int
		incY        int
		ansAlphaEq1 []float64

		trueAns [][]float64
	}{
		{
			name: "M gt N inc 1",
			m:    5,
			n:    3,
			a: [][]float64{
				{1.3, 2.4, 3.5},
				{2.6, 2.8, 3.3},
				{-1.3, -4.3, -9.7},
				{8, 9, -10},
				{-12, -14, -6},
			},
			x:       []float64{-2, -3, 0, 1, 2},
			y:       []float64{-1.1, 5, 0},
			incX:    1,
			incY:    1,
			trueAns: [][]float64{{3.5, -7.6, 3.5}, {5.9, -12.2, 3.3}, {-1.3, -4.3, -9.7}, {6.9, 14, -10}, {-14.2, -4, -6}},
		},
		{
			name: "M eq N inc 1",
			m:    3,
			n:    3,
			a: [][]float64{
				{1.3, 2.4, 3.5},
				{2.6, 2.8, 3.3},
				{-1.3, -4.3, -9.7},
			},
			x:       []float64{-2, -3, 0},
			y:       []float64{-1.1, 5, 0},
			incX:    1,
			incY:    1,
			trueAns: [][]float64{{3.5, -7.6, 3.5}, {5.9, -12.2, 3.3}, {-1.3, -4.3, -9.7}},
		},

		{
			name: "M lt N inc 1",
			m:    3,
			n:    6,
			a: [][]float64{
				{1.3, 2.4, 3.5, 4.8, 1.11, -9},
				{2.6, 2.8, 3.3, -3.4, 6.2, -8.7},
				{-1.3, -4.3, -9.7, -3.1, 8.9, 8.9},
			},
			x:       []float64{-2, -3, 0},
			y:       []float64{-1.1, 5, 0, 9, 19, 22},
			incX:    1,
			incY:    1,
			trueAns: [][]float64{{3.5, -7.6, 3.5, -13.2, -36.89, -53}, {5.9, -12.2, 3.3, -30.4, -50.8, -74.7}, {-1.3, -4.3, -9.7, -3.1, 8.9, 8.9}},
		},
		{
			name: "M gt N inc not 1",
			m:    5,
			n:    3,
			a: [][]float64{
				{1.3, 2.4, 3.5},
				{2.6, 2.8, 3.3},
				{-1.3, -4.3, -9.7},
				{8, 9, -10},
				{-12, -14, -6},
			},
			x:       []float64{-2, -3, 0, 1, 2, 6, 0, 9, 7},
			y:       []float64{-1.1, 5, 0, 8, 7, -5, 7},
			incX:    2,
			incY:    3,
			trueAns: [][]float64{{3.5, -13.6, -10.5}, {2.6, 2.8, 3.3}, {-3.5, 11.7, 4.3}, {8, 9, -10}, {-19.700000000000003, 42, 43}},
		},
		{
			name: "M eq N inc not 1",
			m:    3,
			n:    3,
			a: [][]float64{
				{1.3, 2.4, 3.5},
				{2.6, 2.8, 3.3},
				{-1.3, -4.3, -9.7},
			},
			x:       []float64{-2, -3, 0, 8, 7, -9, 7, -6, 12, 6, 6, 6, -11},
			y:       []float64{-1.1, 5, 0, 0, 9, 8, 6},
			incX:    4,
			incY:    3,
			trueAns: [][]float64{{3.5, 2.4, -8.5}, {-5.1, 2.8, 45.3}, {-14.5, -4.3, 62.3}},
		},
		{
			name: "M lt N inc not 1",
			m:    3,
			n:    6,
			a: [][]float64{
				{1.3, 2.4, 3.5, 4.8, 1.11, -9},
				{2.6, 2.8, 3.3, -3.4, 6.2, -8.7},
				{-1.3, -4.3, -9.7, -3.1, 8.9, 8.9},
			},
			x:       []float64{-2, -3, 0, 0, 8, 0, 9, -3},
			y:       []float64{-1.1, 5, 0, 9, 19, 22, 11, -8.11, -9.22, 9.87, 7},
			incX:    3,
			incY:    2,
			trueAns: [][]float64{{3.5, 2.4, -34.5, -17.2, 19.55, -23}, {2.6, 2.8, 3.3, -3.4, 6.2, -8.7}, {-11.2, -4.3, 161.3, 95.9, -74.08, 71.9}},
		},
	} {
		// TODO: Add tests where a is longer
		// TODO: Add panic tests
		// TODO: Add negative increment tests

		x := sliceCopy(test.x)
		y := sliceCopy(test.y)

		a := sliceOfSliceCopy(test.a)

		// Test with row major
		alpha := 1.0
		aFlat := flatten(a)
		blasser.Dger(test.m, test.n, alpha, x, test.incX, y, test.incY, aFlat, test.n)
		ans := unflatten(aFlat, test.m, test.n)
		dgercomp(t, x, test.x, y, test.y, ans, test.trueAns, test.name+" row maj")

		// Test with different alpha
		alpha = 4.0
		aFlat = flatten(a)
		blasser.Dger(test.m, test.n, alpha, x, test.incX, y, test.incY, aFlat, test.n)
		ans = unflatten(aFlat, test.m, test.n)
		trueCopy := sliceOfSliceCopy(test.trueAns)
		for i := range trueCopy {
			for j := range trueCopy[i] {
				trueCopy[i][j] = alpha*(trueCopy[i][j]-a[i][j]) + a[i][j]
			}
		}
		dgercomp(t, x, test.x, y, test.y, ans, trueCopy, test.name+" row maj alpha")
	}
}

func dgercomp(t *testing.T, x, xCopy, y, yCopy []float64, ans [][]float64, trueAns [][]float64, name string) {
	if !dSliceEqual(x, xCopy) {
		t.Errorf("case %v: x modified during call to dger", name)
	}
	if !dSliceEqual(y, yCopy) {
		t.Errorf("case %v: x modified during call to dger", name)
	}

	for i := range ans {
		if !dSliceTolEqual(ans[i], trueAns[i]) {
			t.Errorf("case %v: answer mismatch. Expected %v, Found %v", name, trueAns, ans)
			break
		}
	}
}
