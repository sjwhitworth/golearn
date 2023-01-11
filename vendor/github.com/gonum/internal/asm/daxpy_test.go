// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"math/rand"
	"testing"
)

var daxpyTests = []struct {
	alpha float64
	xData []float64
	yData []float64

	want    []float64
	wantRev []float64 // Result when x is traversed in reverse direction.
}{
	{
		alpha:   0,
		xData:   []float64{2},
		yData:   []float64{-3},
		want:    []float64{-3},
		wantRev: []float64{-3},
	},
	{
		alpha:   1,
		xData:   []float64{2},
		yData:   []float64{-3},
		want:    []float64{-1},
		wantRev: []float64{-1},
	},
	{
		alpha:   3,
		xData:   []float64{2},
		yData:   []float64{-3},
		want:    []float64{3},
		wantRev: []float64{3},
	},
	{
		alpha:   -3,
		xData:   []float64{2},
		yData:   []float64{-3},
		want:    []float64{-9},
		wantRev: []float64{-9},
	},
	{
		alpha:   0,
		xData:   []float64{0, 0, 1, 1, 2, -3, -4},
		yData:   []float64{0, 1, 0, 3, -4, 5, -6},
		want:    []float64{0, 1, 0, 3, -4, 5, -6},
		wantRev: []float64{0, 1, 0, 3, -4, 5, -6},
	},
	{
		alpha:   1,
		xData:   []float64{0, 0, 1, 1, 2, -3, -4},
		yData:   []float64{0, 1, 0, 3, -4, 5, -6},
		want:    []float64{0, 1, 1, 4, -2, 2, -10},
		wantRev: []float64{-4, -2, 2, 4, -3, 5, -6},
	},
	{
		alpha:   3,
		xData:   []float64{0, 0, 1, 1, 2, -3, -4},
		yData:   []float64{0, 1, 0, 3, -4, 5, -6},
		want:    []float64{0, 1, 3, 6, 2, -4, -18},
		wantRev: []float64{-12, -8, 6, 6, -1, 5, -6},
	},
	{
		alpha:   -3,
		xData:   []float64{0, 0, 1, 1, 2, -3, -4},
		yData:   []float64{0, 1, 0, 3, -4, 5, -6},
		want:    []float64{0, 1, -3, 0, -10, 14, 6},
		wantRev: []float64{12, 10, -6, 0, -7, 5, -6},
	},
	{
		alpha:   -5,
		xData:   []float64{0, 0, 1, 1, 2, -3, -4, 5},
		yData:   []float64{0, 1, 0, 3, -4, 5, -6, 7},
		want:    []float64{0, 1, -5, -2, -14, 20, 14, -18},
		wantRev: []float64{-25, 21, 15, -7, -9, 0, -6, 7},
	},
}

func TestDaxpyUnitary(t *testing.T) {
	for i, test := range daxpyTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		prefix := fmt.Sprintf("test %v (y+=a*x)", i)
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		DaxpyUnitary(test.alpha, x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, y, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, y)
		}
	}
}

func TestDaxpyUnitaryTo(t *testing.T) {
	for i, test := range daxpyTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		// Test dst = alpha * x + y.
		prefix := fmt.Sprintf("test %v (dst=a*x+y)", i)
		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		dst, dstFront, dstBack := newGuardedVector(test.xData, 1)
		DaxpyUnitaryTo(dst, test.alpha, x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !allNaN(dstFront) || !allNaN(dstBack) {
			t.Errorf(msgGuard, prefix, "dst", dstFront, dstBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, dst, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, dst)
		}

		// Test y = alpha * x + y.
		prefix = fmt.Sprintf("test %v (y=a*x+y)", i)
		x, xFront, xBack = newGuardedVector(test.xData, 1)
		y, yFront, yBack = newGuardedVector(test.yData, 1)
		DaxpyUnitaryTo(y, test.alpha, x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, y, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, y)
		}

		// Test x = alpha * x + y.
		prefix = fmt.Sprintf("test %v (x=a*x+y)", i)
		x, xFront, xBack = newGuardedVector(test.xData, 1)
		y, yFront, yBack = newGuardedVector(test.yData, 1)

		DaxpyUnitaryTo(x, test.alpha, x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, prefix, "y", yFront, yBack)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("%v: modified read-only y argument", prefix)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDaxpyInc(t *testing.T) {
	for i, test := range daxpyTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"
		n := len(test.xData)

		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}

				prefix := fmt.Sprintf("test %v, incX = %v, incY = %v", i, incX, incY)
				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)
				DaxpyInc(test.alpha, x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msgGuard, prefix, "y", yFront, yBack)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(y, incY) {
					t.Errorf("%v: modified y argument at non-stride position", prefix)
				}

				want := test.want
				if incX*incY < 0 {
					want = test.wantRev
				}
				if !equalStrided(want, y, incY) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, want, y)
				}
			}
		}
	}
}

func TestDaxpyIncTo(t *testing.T) {
	for i, test := range daxpyTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"
		n := len(test.xData)
		want := make([]float64, n)

		for _, incX := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
			for _, incY := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}

				// Test y = alpha * x + y.
				prefix := fmt.Sprintf("test %v (y=a*x+y), incX = %v, incY = %v", i, incX, incY)
				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)
				DaxpyIncTo(y, uintptr(incY), uintptr(iy),
					test.alpha, x, y,
					uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msgGuard, prefix, "y", yFront, yBack)
				}
				if !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}

				if incX*incY < 0 {
					copy(want, test.wantRev)
				} else {
					copy(want, test.want)
				}
				if !equalStrided(want, y, incY) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, want, y)
				}

				// Test x = alpha * x + y.
				prefix = fmt.Sprintf("test %v (x=a*x+y), incX = %v, incY = %v", i, incX, incY)
				x, xFront, xBack = newGuardedVector(test.xData, incX)
				y, yFront, yBack = newGuardedVector(test.yData, incY)

				DaxpyIncTo(x, uintptr(incX), uintptr(ix),
					test.alpha, x, y,
					uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msgGuard, prefix, "y", yFront, yBack)
				}
				if !equalStrided(test.yData, y, incY) {
					t.Errorf("%v: modified read-only y argument", prefix)
				}

				if incX*incY < 0 {
					copy(want, test.wantRev)
					for i := 0; i < n/2; i++ {
						want[i], want[n-i-1] = want[n-i-1], want[i]
					}
				} else {
					copy(want, test.want)
				}
				if !equalStrided(want, x, incX) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, want, x)
				}

				for _, incDst := range []int{-7, -4, -3, -2, -1, 1, 2, 3, 4, 7} {
					var idst int
					if incDst < 0 {
						idst = (-n + 1) * incDst
					}

					// Test dst = alpha * x + y.
					prefix = fmt.Sprintf("test %v (dst=a*x+y), incX = %v, incY = %v, incDst =%v", i, incX, incY, incDst)
					x, xFront, xBack = newGuardedVector(test.xData, incX)
					y, yFront, yBack = newGuardedVector(test.yData, incY)
					dst, dstFront, dstBack := newGuardedVector(test.xData, incDst)
					DaxpyIncTo(dst, uintptr(incDst), uintptr(idst),
						test.alpha, x, y,
						uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

					if !allNaN(xFront) || !allNaN(xBack) {
						t.Errorf(msgGuard, prefix, "x", xFront, xBack)
					}
					if !allNaN(yFront) || !allNaN(yBack) {
						t.Errorf(msgGuard, prefix, "y", yFront, yBack)
					}
					if !allNaN(dstFront) || !allNaN(dstBack) {
						t.Errorf(msgGuard, prefix, "dst", dstFront, dstBack)
					}
					if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
						t.Errorf("%v: modified read-only x argument", prefix)
					}
					if nonStridedWrite(y, incY) || !equalStrided(test.yData, y, incY) {
						t.Errorf("%v: modified read-only y argument", prefix)
					}
					if nonStridedWrite(dst, incDst) {
						t.Errorf("%v: modified dst argument at non-stride position", prefix)
					}

					if incX*incY < 0 {
						copy(want, test.wantRev)
					} else {
						copy(want, test.want)
					}
					if incY*incDst < 0 {
						for i := 0; i < n/2; i++ {
							want[i], want[n-i-1] = want[n-i-1], want[i]
						}
					}
					if !equalStrided(want, dst, incDst) {
						t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, want, dst)
					}
				}
			}
		}
	}
}

var gs []float64

func BenchmarkDaxpyUnitaryN1(b *testing.B)      { daxpyUnitaryBenchmark(b, 1) }
func BenchmarkDaxpyUnitaryN2(b *testing.B)      { daxpyUnitaryBenchmark(b, 2) }
func BenchmarkDaxpyUnitaryN3(b *testing.B)      { daxpyUnitaryBenchmark(b, 3) }
func BenchmarkDaxpyUnitaryN4(b *testing.B)      { daxpyUnitaryBenchmark(b, 4) }
func BenchmarkDaxpyUnitaryN10(b *testing.B)     { daxpyUnitaryBenchmark(b, 10) }
func BenchmarkDaxpyUnitaryN100(b *testing.B)    { daxpyUnitaryBenchmark(b, 100) }
func BenchmarkDaxpyUnitaryN1000(b *testing.B)   { daxpyUnitaryBenchmark(b, 1000) }
func BenchmarkDaxpyUnitaryN10000(b *testing.B)  { daxpyUnitaryBenchmark(b, 10000) }
func BenchmarkDaxpyUnitaryN100000(b *testing.B) { daxpyUnitaryBenchmark(b, 100000) }

func daxpyUnitaryBenchmark(b *testing.B, n int) {
	x := randomSlice(n, 1)
	y := randomSlice(n, 1)
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DaxpyUnitary(a, x, y)
	}
	gs = y
}

func BenchmarkDaxpyUnitaryToYN1(b *testing.B)      { daxpyUnitaryToYBenchmark(b, 1) }
func BenchmarkDaxpyUnitaryToYN2(b *testing.B)      { daxpyUnitaryToYBenchmark(b, 2) }
func BenchmarkDaxpyUnitaryToYN3(b *testing.B)      { daxpyUnitaryToYBenchmark(b, 3) }
func BenchmarkDaxpyUnitaryToYN4(b *testing.B)      { daxpyUnitaryToYBenchmark(b, 4) }
func BenchmarkDaxpyUnitaryToYN10(b *testing.B)     { daxpyUnitaryToYBenchmark(b, 10) }
func BenchmarkDaxpyUnitaryToYN100(b *testing.B)    { daxpyUnitaryToYBenchmark(b, 100) }
func BenchmarkDaxpyUnitaryToYN1000(b *testing.B)   { daxpyUnitaryToYBenchmark(b, 1000) }
func BenchmarkDaxpyUnitaryToYN10000(b *testing.B)  { daxpyUnitaryToYBenchmark(b, 10000) }
func BenchmarkDaxpyUnitaryToYN100000(b *testing.B) { daxpyUnitaryToYBenchmark(b, 100000) }

func daxpyUnitaryToYBenchmark(b *testing.B, n int) {
	x := randomSlice(n, 1)
	y := randomSlice(n, 1)
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DaxpyUnitaryTo(y, a, x, y)
	}
	gs = y
}

func BenchmarkDaxpyUnitaryToN1(b *testing.B)      { daxpyUnitaryToBenchmark(b, 1) }
func BenchmarkDaxpyUnitaryToN2(b *testing.B)      { daxpyUnitaryToBenchmark(b, 2) }
func BenchmarkDaxpyUnitaryToN3(b *testing.B)      { daxpyUnitaryToBenchmark(b, 3) }
func BenchmarkDaxpyUnitaryToN4(b *testing.B)      { daxpyUnitaryToBenchmark(b, 4) }
func BenchmarkDaxpyUnitaryToN10(b *testing.B)     { daxpyUnitaryToBenchmark(b, 10) }
func BenchmarkDaxpyUnitaryToN100(b *testing.B)    { daxpyUnitaryToBenchmark(b, 100) }
func BenchmarkDaxpyUnitaryToN1000(b *testing.B)   { daxpyUnitaryToBenchmark(b, 1000) }
func BenchmarkDaxpyUnitaryToN10000(b *testing.B)  { daxpyUnitaryToBenchmark(b, 10000) }
func BenchmarkDaxpyUnitaryToN100000(b *testing.B) { daxpyUnitaryToBenchmark(b, 100000) }

func daxpyUnitaryToBenchmark(b *testing.B, n int) {
	x := randomSlice(n, 1)
	y := randomSlice(n, 1)
	dst := randomSlice(n, 1)
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DaxpyUnitaryTo(dst, a, x, y)
	}
	gs = dst
}

func BenchmarkDaxpyIncN1Inc1(b *testing.B) { daxpyIncBenchmark(b, 1, 1) }

func BenchmarkDaxpyIncN2Inc1(b *testing.B)  { daxpyIncBenchmark(b, 2, 1) }
func BenchmarkDaxpyIncN2Inc2(b *testing.B)  { daxpyIncBenchmark(b, 2, 2) }
func BenchmarkDaxpyIncN2Inc4(b *testing.B)  { daxpyIncBenchmark(b, 2, 4) }
func BenchmarkDaxpyIncN2Inc10(b *testing.B) { daxpyIncBenchmark(b, 2, 10) }

func BenchmarkDaxpyIncN3Inc1(b *testing.B)  { daxpyIncBenchmark(b, 3, 1) }
func BenchmarkDaxpyIncN3Inc2(b *testing.B)  { daxpyIncBenchmark(b, 3, 2) }
func BenchmarkDaxpyIncN3Inc4(b *testing.B)  { daxpyIncBenchmark(b, 3, 4) }
func BenchmarkDaxpyIncN3Inc10(b *testing.B) { daxpyIncBenchmark(b, 3, 10) }

func BenchmarkDaxpyIncN4Inc1(b *testing.B)  { daxpyIncBenchmark(b, 4, 1) }
func BenchmarkDaxpyIncN4Inc2(b *testing.B)  { daxpyIncBenchmark(b, 4, 2) }
func BenchmarkDaxpyIncN4Inc4(b *testing.B)  { daxpyIncBenchmark(b, 4, 4) }
func BenchmarkDaxpyIncN4Inc10(b *testing.B) { daxpyIncBenchmark(b, 4, 10) }

func BenchmarkDaxpyIncN10Inc1(b *testing.B)  { daxpyIncBenchmark(b, 10, 1) }
func BenchmarkDaxpyIncN10Inc2(b *testing.B)  { daxpyIncBenchmark(b, 10, 2) }
func BenchmarkDaxpyIncN10Inc4(b *testing.B)  { daxpyIncBenchmark(b, 10, 4) }
func BenchmarkDaxpyIncN10Inc10(b *testing.B) { daxpyIncBenchmark(b, 10, 10) }

func BenchmarkDaxpyIncN1000Inc1(b *testing.B)  { daxpyIncBenchmark(b, 1000, 1) }
func BenchmarkDaxpyIncN1000Inc2(b *testing.B)  { daxpyIncBenchmark(b, 1000, 2) }
func BenchmarkDaxpyIncN1000Inc4(b *testing.B)  { daxpyIncBenchmark(b, 1000, 4) }
func BenchmarkDaxpyIncN1000Inc10(b *testing.B) { daxpyIncBenchmark(b, 1000, 10) }

func BenchmarkDaxpyIncN100000Inc1(b *testing.B)  { daxpyIncBenchmark(b, 100000, 1) }
func BenchmarkDaxpyIncN100000Inc2(b *testing.B)  { daxpyIncBenchmark(b, 100000, 2) }
func BenchmarkDaxpyIncN100000Inc4(b *testing.B)  { daxpyIncBenchmark(b, 100000, 4) }
func BenchmarkDaxpyIncN100000Inc10(b *testing.B) { daxpyIncBenchmark(b, 100000, 10) }

func BenchmarkDaxpyIncN100000IncM1(b *testing.B)  { daxpyIncBenchmark(b, 100000, -1) }
func BenchmarkDaxpyIncN100000IncM2(b *testing.B)  { daxpyIncBenchmark(b, 100000, -2) }
func BenchmarkDaxpyIncN100000IncM4(b *testing.B)  { daxpyIncBenchmark(b, 100000, -4) }
func BenchmarkDaxpyIncN100000IncM10(b *testing.B) { daxpyIncBenchmark(b, 100000, -10) }

func daxpyIncBenchmark(b *testing.B, n, inc int) {
	x := randomSlice(n, inc)
	y := randomSlice(n, inc)
	var ini int
	if inc < 0 {
		ini = (-n + 1) * inc
	}
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DaxpyInc(a, x, y, uintptr(n), uintptr(inc), uintptr(inc), uintptr(ini), uintptr(ini))
	}
	gs = y
}

func BenchmarkDaxpyIncToN1Inc1(b *testing.B) { daxpyIncToBenchmark(b, 1, 1) }

func BenchmarkDaxpyIncToN2Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 2, 1) }
func BenchmarkDaxpyIncToN2Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 2, 2) }
func BenchmarkDaxpyIncToN2Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 2, 4) }
func BenchmarkDaxpyIncToN2Inc10(b *testing.B) { daxpyIncToBenchmark(b, 2, 10) }

func BenchmarkDaxpyIncToN3Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 3, 1) }
func BenchmarkDaxpyIncToN3Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 3, 2) }
func BenchmarkDaxpyIncToN3Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 3, 4) }
func BenchmarkDaxpyIncToN3Inc10(b *testing.B) { daxpyIncToBenchmark(b, 3, 10) }

func BenchmarkDaxpyIncToN4Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 4, 1) }
func BenchmarkDaxpyIncToN4Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 4, 2) }
func BenchmarkDaxpyIncToN4Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 4, 4) }
func BenchmarkDaxpyIncToN4Inc10(b *testing.B) { daxpyIncToBenchmark(b, 4, 10) }

func BenchmarkDaxpyIncToN10Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 10, 1) }
func BenchmarkDaxpyIncToN10Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 10, 2) }
func BenchmarkDaxpyIncToN10Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 10, 4) }
func BenchmarkDaxpyIncToN10Inc10(b *testing.B) { daxpyIncToBenchmark(b, 10, 10) }

func BenchmarkDaxpyIncToN1000Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 1000, 1) }
func BenchmarkDaxpyIncToN1000Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 1000, 2) }
func BenchmarkDaxpyIncToN1000Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 1000, 4) }
func BenchmarkDaxpyIncToN1000Inc10(b *testing.B) { daxpyIncToBenchmark(b, 1000, 10) }

func BenchmarkDaxpyIncToN100000Inc1(b *testing.B)  { daxpyIncToBenchmark(b, 100000, 1) }
func BenchmarkDaxpyIncToN100000Inc2(b *testing.B)  { daxpyIncToBenchmark(b, 100000, 2) }
func BenchmarkDaxpyIncToN100000Inc4(b *testing.B)  { daxpyIncToBenchmark(b, 100000, 4) }
func BenchmarkDaxpyIncToN100000Inc10(b *testing.B) { daxpyIncToBenchmark(b, 100000, 10) }

func BenchmarkDaxpyIncToN100000IncM1(b *testing.B)  { daxpyIncToBenchmark(b, 100000, -1) }
func BenchmarkDaxpyIncToN100000IncM2(b *testing.B)  { daxpyIncToBenchmark(b, 100000, -2) }
func BenchmarkDaxpyIncToN100000IncM4(b *testing.B)  { daxpyIncToBenchmark(b, 100000, -4) }
func BenchmarkDaxpyIncToN100000IncM10(b *testing.B) { daxpyIncToBenchmark(b, 100000, -10) }

func daxpyIncToBenchmark(b *testing.B, n, inc int) {
	x := randomSlice(n, inc)
	y := randomSlice(n, inc)
	dst := randomSlice(n, inc)
	var ini int
	if inc < 0 {
		ini = (-n + 1) * inc
	}
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DaxpyIncTo(dst, uintptr(inc), uintptr(ini), a, x, y,
			uintptr(n), uintptr(inc), uintptr(inc), uintptr(ini), uintptr(ini))
	}
	gs = y
}

func randomSlice(n, inc int) []float64 {
	if inc < 0 {
		inc = -inc
	}
	x := make([]float64, (n-1)*inc+1)
	for i := range x {
		x[i] = rand.Float64()
	}
	return x
}
