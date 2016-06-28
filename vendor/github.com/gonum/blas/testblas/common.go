package testblas

import (
	"math"
	"testing"

	"github.com/gonum/blas"
)

// throwPanic will throw unexpected panics if true, or will just report them as errors if false
const throwPanic = true

func dTolEqual(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if a == b {
		return true
	}
	m := math.Max(math.Abs(a), math.Abs(b))
	if m > 1 {
		a /= m
		b /= m
	}
	if math.Abs(a-b) < 1e-14 {
		return true
	}
	return false
}

func dSliceTolEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !dTolEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func dStridedSliceTolEqual(n int, a []float64, inca int, b []float64, incb int) bool {
	ia := 0
	ib := 0
	if inca <= 0 {
		ia = -(n - 1) * inca
	}
	if incb <= 0 {
		ib = -(n - 1) * incb
	}
	for i := 0; i < n; i++ {
		if !dTolEqual(a[ia], b[ib]) {
			return false
		}
		ia += inca
		ib += incb
	}
	return true
}

func dSliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
			return false
		}
	}
	return true
}

func dCopyTwoTmp(x, xTmp, y, yTmp []float64) {
	if len(x) != len(xTmp) {
		panic("x size mismatch")
	}
	if len(y) != len(yTmp) {
		panic("y size mismatch")
	}
	for i, val := range x {
		xTmp[i] = val
	}
	for i, val := range y {
		yTmp[i] = val
	}
}

// returns true if the function panics
func panics(f func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	f()
	return
}

func testpanics(f func(), name string, t *testing.T) {
	b := panics(f)
	if !b {
		t.Errorf("%v should panic and does not", name)
	}
}

func sliceOfSliceCopy(a [][]float64) [][]float64 {
	n := make([][]float64, len(a))
	for i := range a {
		n[i] = make([]float64, len(a[i]))
		copy(n[i], a[i])
	}
	return n
}

func sliceCopy(a []float64) []float64 {
	n := make([]float64, len(a))
	copy(n, a)
	return n
}

func flatten(a [][]float64) []float64 {
	if len(a) == 0 {
		return nil
	}
	m := len(a)
	n := len(a[0])
	s := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			s[i*n+j] = a[i][j]
		}
	}
	return s
}

func unflatten(a []float64, m, n int) [][]float64 {
	s := make([][]float64, m)
	for i := 0; i < m; i++ {
		s[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			s[i][j] = a[i*n+j]
		}
	}
	return s
}

// flattenTriangular turns the upper or lower triangle of a dense slice of slice
// into a single slice with packed storage. a must be a square matrix.
func flattenTriangular(a [][]float64, ul blas.Uplo) []float64 {
	m := len(a)
	aFlat := make([]float64, m*(m+1)/2)
	var k int
	if ul == blas.Upper {
		for i := 0; i < m; i++ {
			k += copy(aFlat[k:], a[i][i:])
		}
		return aFlat
	}
	for i := 0; i < m; i++ {
		k += copy(aFlat[k:], a[i][:i+1])
	}
	return aFlat
}

// flattenBanded turns a dense banded slice of slice into the compact banded matrix format
func flattenBanded(a [][]float64, ku, kl int) []float64 {
	m := len(a)
	n := len(a[0])
	if ku < 0 || kl < 0 {
		panic("testblas: negative band length")
	}
	nRows := m
	nCols := (ku + kl + 1)
	aflat := make([]float64, nRows*nCols)
	for i := range aflat {
		aflat[i] = math.NaN()
	}
	// loop over the rows, and then the bands
	// elements in the ith row stay in the ith row
	// order in bands is kept
	for i := 0; i < nRows; i++ {
		min := -kl
		if i-kl < 0 {
			min = -i
		}
		max := ku
		if i+ku >= n {
			max = n - i - 1
		}
		for j := min; j <= max; j++ {
			col := kl + j
			aflat[i*nCols+col] = a[i][i+j]
		}
	}
	return aflat
}

// makeIncremented takes a slice with inc == 1 and makes an incremented version
// and adds extra values on the end
func makeIncremented(x []float64, inc int, extra int) []float64 {
	if inc == 0 {
		panic("zero inc")
	}
	absinc := inc
	if absinc < 0 {
		absinc = -inc
	}
	xcopy := make([]float64, len(x))
	if inc > 0 {
		copy(xcopy, x)
	} else {
		for i := 0; i < len(x); i++ {
			xcopy[i] = x[len(x)-i-1]
		}
	}

	// don't use NaN because it makes comparison hard
	// Do use a weird unique value for easier debugging
	counter := 100.0
	var xnew []float64
	for i, v := range xcopy {
		xnew = append(xnew, v)
		if i != len(x)-1 {
			for j := 0; j < absinc-1; j++ {
				xnew = append(xnew, counter)
				counter++
			}
		}
	}
	for i := 0; i < extra; i++ {
		xnew = append(xnew, counter)
		counter++
	}
	return xnew
}
