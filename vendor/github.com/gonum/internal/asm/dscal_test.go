// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"math/rand"
	"testing"
)

var dscalTests = []struct {
	alpha float64
	x     []float64
	want  []float64
}{
	{
		alpha: 0,
		x:     []float64{1},
		want:  []float64{0},
	},
	{
		alpha: 1,
		x:     []float64{1},
		want:  []float64{1},
	},
	{
		alpha: 2,
		x:     []float64{1, -2},
		want:  []float64{2, -4},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3},
		want:  []float64{2, -4, 6},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3, 4},
		want:  []float64{2, -4, 6, 8},
	},
	{
		alpha: 2,
		x:     []float64{1, -2, 3, 4, -5},
		want:  []float64{2, -4, 6, 8, -10},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7, 8},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14, 16},
	},
	{
		alpha: 2,
		x:     []float64{0, 1, -2, 3, 4, -5, 6, -7, 8, 9},
		want:  []float64{0, 2, -4, 6, 8, -10, 12, -14, 16, 18},
	},
}

func TestDscalUnitary(t *testing.T) {
	for i, test := range dscalTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		prefix := fmt.Sprintf("test %v (x*=a)", i)
		x, xFront, xBack := newGuardedVector(test.x, 1)
		DscalUnitary(test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDscalUnitaryTo(t *testing.T) {
	for i, test := range dscalTests {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		// Test dst = alpha * x.
		prefix := fmt.Sprintf("test %v (dst=a*x)", i)
		x, xFront, xBack := newGuardedVector(test.x, 1)
		dst, dstFront, dstBack := newGuardedVector(test.x, 1)
		DscalUnitaryTo(dst, test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}
		if !allNaN(dstFront) || !allNaN(dstBack) {
			t.Errorf(msgGuard, prefix, "dst", dstFront, dstBack)
		}
		if !equalStrided(test.x, x, 1) {
			t.Errorf("%v: modified read-only x argument", prefix)
		}

		if !equalStrided(test.want, dst, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, dst)
		}

		// Test x = alpha * x.
		prefix = fmt.Sprintf("test %v (x=a*x)", i)
		x, xFront, xBack = newGuardedVector(test.x, 1)
		DscalUnitaryTo(x, test.alpha, x)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, prefix, "x", xFront, xBack)
		}

		if !equalStrided(test.want, x, 1) {
			t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
		}
	}
}

func TestDscalInc(t *testing.T) {
	const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

	for i, test := range dscalTests {
		n := len(test.x)
		for _, incX := range []int{1, 2, 3, 4, 7, 10} {
			prefix := fmt.Sprintf("test %v (x*=a), incX = %v", i, incX)
			x, xFront, xBack := newGuardedVector(test.x, incX)
			DscalInc(test.alpha, x, uintptr(n), uintptr(incX))

			if !allNaN(xFront) || !allNaN(xBack) {
				t.Errorf(msgGuard, prefix, "x", xFront, xBack)
			}
			if nonStridedWrite(x, incX) {
				t.Errorf("%v: modified x argument at non-stride position", prefix)
			}

			if !equalStrided(test.want, x, incX) {
				t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
			}
		}
	}
}

func TestDscalIncTo(t *testing.T) {
	const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

	for i, test := range dscalTests {
		n := len(test.x)

		for _, incX := range []int{1, 2, 3, 4, 7, 10} {
			// Test x = alpha * x.
			prefix := fmt.Sprintf("test %v (x=a*x), incX = %v", i, incX)
			x, xFront, xBack := newGuardedVector(test.x, incX)
			DscalIncTo(x, uintptr(incX), test.alpha, x, uintptr(n), uintptr(incX))

			if !allNaN(xFront) || !allNaN(xBack) {
				t.Errorf(msgGuard, prefix, "x", xFront, xBack)
			}
			if nonStridedWrite(x, incX) {
				t.Errorf("%v: modified x argument at non-stride position", prefix)
			}
			if !equalStrided(test.want, x, incX) {
				t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, x)
			}

			for _, incDst := range []int{1, 2, 3, 4, 7, 10} {
				// Test dst = alpha * x.
				prefix = fmt.Sprintf("test %v (dst=a*x), incX = %v, incDst = %v", i, incX, incDst)
				x, xFront, xBack = newGuardedVector(test.x, incX)
				dst, dstFront, dstBack := newGuardedVector(test.x, incDst)
				DscalIncTo(dst, uintptr(incDst), test.alpha, x, uintptr(n), uintptr(incX))

				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(dstFront) || !allNaN(dstBack) {
					t.Errorf(msgGuard, prefix, "dst", dstFront, dstBack)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.x, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(dst, incDst) {
					t.Errorf("%v: modified dst argument at non-stride position", prefix)
				}

				if !equalStrided(test.want, dst, incDst) {
					t.Errorf("%v: unexpected result:\nwant: %v\ngot: %v", prefix, test.want, dst)
				}
			}
		}
	}
}

func BenchmarkDscalUnitaryN1(b *testing.B)      { benchmarkDscalUnitary(b, 1) }
func BenchmarkDscalUnitaryN2(b *testing.B)      { benchmarkDscalUnitary(b, 2) }
func BenchmarkDscalUnitaryN3(b *testing.B)      { benchmarkDscalUnitary(b, 3) }
func BenchmarkDscalUnitaryN4(b *testing.B)      { benchmarkDscalUnitary(b, 4) }
func BenchmarkDscalUnitaryN10(b *testing.B)     { benchmarkDscalUnitary(b, 10) }
func BenchmarkDscalUnitaryN100(b *testing.B)    { benchmarkDscalUnitary(b, 100) }
func BenchmarkDscalUnitaryN1000(b *testing.B)   { benchmarkDscalUnitary(b, 1000) }
func BenchmarkDscalUnitaryN10000(b *testing.B)  { benchmarkDscalUnitary(b, 10000) }
func BenchmarkDscalUnitaryN100000(b *testing.B) { benchmarkDscalUnitary(b, 100000) }

func benchmarkDscalUnitary(b *testing.B, n int) {
	x := randomSlice(n, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		DscalUnitary(2, x)
		DscalUnitary(0.5, x)
	}
	gs = x
}

func BenchmarkDscalUnitaryToN1(b *testing.B)      { benchmarkDscalUnitaryTo(b, 1) }
func BenchmarkDscalUnitaryToN2(b *testing.B)      { benchmarkDscalUnitaryTo(b, 2) }
func BenchmarkDscalUnitaryToN3(b *testing.B)      { benchmarkDscalUnitaryTo(b, 3) }
func BenchmarkDscalUnitaryToN4(b *testing.B)      { benchmarkDscalUnitaryTo(b, 4) }
func BenchmarkDscalUnitaryToN10(b *testing.B)     { benchmarkDscalUnitaryTo(b, 10) }
func BenchmarkDscalUnitaryToN100(b *testing.B)    { benchmarkDscalUnitaryTo(b, 100) }
func BenchmarkDscalUnitaryToN1000(b *testing.B)   { benchmarkDscalUnitaryTo(b, 1000) }
func BenchmarkDscalUnitaryToN10000(b *testing.B)  { benchmarkDscalUnitaryTo(b, 10000) }
func BenchmarkDscalUnitaryToN100000(b *testing.B) { benchmarkDscalUnitaryTo(b, 100000) }

func benchmarkDscalUnitaryTo(b *testing.B, n int) {
	x := randomSlice(n, 1)
	dst := randomSlice(n, 1)
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DscalUnitaryTo(dst, a, x)
	}
	gs = dst
}

func BenchmarkDscalUnitaryToXN1(b *testing.B)      { benchmarkDscalUnitaryToX(b, 1) }
func BenchmarkDscalUnitaryToXN2(b *testing.B)      { benchmarkDscalUnitaryToX(b, 2) }
func BenchmarkDscalUnitaryToXN3(b *testing.B)      { benchmarkDscalUnitaryToX(b, 3) }
func BenchmarkDscalUnitaryToXN4(b *testing.B)      { benchmarkDscalUnitaryToX(b, 4) }
func BenchmarkDscalUnitaryToXN10(b *testing.B)     { benchmarkDscalUnitaryToX(b, 10) }
func BenchmarkDscalUnitaryToXN100(b *testing.B)    { benchmarkDscalUnitaryToX(b, 100) }
func BenchmarkDscalUnitaryToXN1000(b *testing.B)   { benchmarkDscalUnitaryToX(b, 1000) }
func BenchmarkDscalUnitaryToXN10000(b *testing.B)  { benchmarkDscalUnitaryToX(b, 10000) }
func BenchmarkDscalUnitaryToXN100000(b *testing.B) { benchmarkDscalUnitaryToX(b, 100000) }

func benchmarkDscalUnitaryToX(b *testing.B, n int) {
	x := randomSlice(n, 1)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		DscalUnitaryTo(x, 2, x)
		DscalUnitaryTo(x, 0.5, x)
	}
	gs = x
}

func BenchmarkDscalIncN1Inc1(b *testing.B) { benchmarkDscalInc(b, 1, 1) }

func BenchmarkDscalIncN2Inc1(b *testing.B)  { benchmarkDscalInc(b, 2, 1) }
func BenchmarkDscalIncN2Inc2(b *testing.B)  { benchmarkDscalInc(b, 2, 2) }
func BenchmarkDscalIncN2Inc4(b *testing.B)  { benchmarkDscalInc(b, 2, 4) }
func BenchmarkDscalIncN2Inc10(b *testing.B) { benchmarkDscalInc(b, 2, 10) }

func BenchmarkDscalIncN3Inc1(b *testing.B)  { benchmarkDscalInc(b, 3, 1) }
func BenchmarkDscalIncN3Inc2(b *testing.B)  { benchmarkDscalInc(b, 3, 2) }
func BenchmarkDscalIncN3Inc4(b *testing.B)  { benchmarkDscalInc(b, 3, 4) }
func BenchmarkDscalIncN3Inc10(b *testing.B) { benchmarkDscalInc(b, 3, 10) }

func BenchmarkDscalIncN4Inc1(b *testing.B)  { benchmarkDscalInc(b, 4, 1) }
func BenchmarkDscalIncN4Inc2(b *testing.B)  { benchmarkDscalInc(b, 4, 2) }
func BenchmarkDscalIncN4Inc4(b *testing.B)  { benchmarkDscalInc(b, 4, 4) }
func BenchmarkDscalIncN4Inc10(b *testing.B) { benchmarkDscalInc(b, 4, 10) }

func BenchmarkDscalIncN10Inc1(b *testing.B)  { benchmarkDscalInc(b, 10, 1) }
func BenchmarkDscalIncN10Inc2(b *testing.B)  { benchmarkDscalInc(b, 10, 2) }
func BenchmarkDscalIncN10Inc4(b *testing.B)  { benchmarkDscalInc(b, 10, 4) }
func BenchmarkDscalIncN10Inc10(b *testing.B) { benchmarkDscalInc(b, 10, 10) }

func BenchmarkDscalIncN1000Inc1(b *testing.B)  { benchmarkDscalInc(b, 1000, 1) }
func BenchmarkDscalIncN1000Inc2(b *testing.B)  { benchmarkDscalInc(b, 1000, 2) }
func BenchmarkDscalIncN1000Inc4(b *testing.B)  { benchmarkDscalInc(b, 1000, 4) }
func BenchmarkDscalIncN1000Inc10(b *testing.B) { benchmarkDscalInc(b, 1000, 10) }

func BenchmarkDscalIncN100000Inc1(b *testing.B)  { benchmarkDscalInc(b, 100000, 1) }
func BenchmarkDscalIncN100000Inc2(b *testing.B)  { benchmarkDscalInc(b, 100000, 2) }
func BenchmarkDscalIncN100000Inc4(b *testing.B)  { benchmarkDscalInc(b, 100000, 4) }
func BenchmarkDscalIncN100000Inc10(b *testing.B) { benchmarkDscalInc(b, 100000, 10) }

func benchmarkDscalInc(b *testing.B, n, inc int) {
	x := randomSlice(n, inc)
	b.ResetTimer()
	for i := 0; i < b.N; i += 2 {
		DscalInc(2, x, uintptr(n), uintptr(inc))
		DscalInc(0.5, x, uintptr(n), uintptr(inc))
	}
	gs = x
}

func BenchmarkDscalIncToN1Inc1(b *testing.B) { benchmarkDscalIncTo(b, 1, 1) }

func BenchmarkDscalIncToN2Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 2, 1) }
func BenchmarkDscalIncToN2Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 2, 2) }
func BenchmarkDscalIncToN2Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 2, 4) }
func BenchmarkDscalIncToN2Inc10(b *testing.B) { benchmarkDscalIncTo(b, 2, 10) }

func BenchmarkDscalIncToN3Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 3, 1) }
func BenchmarkDscalIncToN3Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 3, 2) }
func BenchmarkDscalIncToN3Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 3, 4) }
func BenchmarkDscalIncToN3Inc10(b *testing.B) { benchmarkDscalIncTo(b, 3, 10) }

func BenchmarkDscalIncToN4Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 4, 1) }
func BenchmarkDscalIncToN4Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 4, 2) }
func BenchmarkDscalIncToN4Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 4, 4) }
func BenchmarkDscalIncToN4Inc10(b *testing.B) { benchmarkDscalIncTo(b, 4, 10) }

func BenchmarkDscalIncToN10Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 10, 1) }
func BenchmarkDscalIncToN10Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 10, 2) }
func BenchmarkDscalIncToN10Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 10, 4) }
func BenchmarkDscalIncToN10Inc10(b *testing.B) { benchmarkDscalIncTo(b, 10, 10) }

func BenchmarkDscalIncToN1000Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 1000, 1) }
func BenchmarkDscalIncToN1000Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 1000, 2) }
func BenchmarkDscalIncToN1000Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 1000, 4) }
func BenchmarkDscalIncToN1000Inc10(b *testing.B) { benchmarkDscalIncTo(b, 1000, 10) }

func BenchmarkDscalIncToN100000Inc1(b *testing.B)  { benchmarkDscalIncTo(b, 100000, 1) }
func BenchmarkDscalIncToN100000Inc2(b *testing.B)  { benchmarkDscalIncTo(b, 100000, 2) }
func BenchmarkDscalIncToN100000Inc4(b *testing.B)  { benchmarkDscalIncTo(b, 100000, 4) }
func BenchmarkDscalIncToN100000Inc10(b *testing.B) { benchmarkDscalIncTo(b, 100000, 10) }

func benchmarkDscalIncTo(b *testing.B, n, inc int) {
	x := randomSlice(n, inc)
	dst := randomSlice(n, inc)
	a := rand.Float64()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DscalIncTo(dst, uintptr(inc), a, x, uintptr(n), uintptr(inc))
	}
	gs = dst
}
