// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestDdotUnitary(t *testing.T) {
	for i, test := range []struct {
		xData []float64
		yData []float64

		want float64
	}{
		{
			xData: []float64{2},
			yData: []float64{-3},
			want:  -6,
		},
		{
			xData: []float64{2, 3},
			yData: []float64{-3, 4},
			want:  6,
		},
		{
			xData: []float64{2, 3, -4},
			yData: []float64{-3, 4, 5},
			want:  -14,
		},
		{
			xData: []float64{2, 3, -4, -5},
			yData: []float64{-3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 2, 3, -4, -5},
			yData: []float64{0, -3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 0, 2, 3, -4, -5},
			yData: []float64{0, 1, -3, 4, 5, -6},
			want:  16,
		},
		{
			xData: []float64{0, 0, 1, 1, 2, -3, -4},
			yData: []float64{0, 1, 0, 3, -4, 5, -6},
			want:  4,
		},
		{
			xData: []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData: []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:  39,
		},
	} {
		const msgGuard = "test %v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		x, xFront, xBack := newGuardedVector(test.xData, 1)
		y, yFront, yBack := newGuardedVector(test.yData, 1)
		got := DdotUnitary(x, y)

		if !allNaN(xFront) || !allNaN(xBack) {
			t.Errorf(msgGuard, i, "x", xFront, xBack)
		}
		if !allNaN(yFront) || !allNaN(yBack) {
			t.Errorf(msgGuard, i, "y", yFront, yBack)
		}
		if !equalStrided(test.xData, x, 1) {
			t.Errorf("test %v: modified read-only x argument", i)
		}
		if !equalStrided(test.yData, y, 1) {
			t.Errorf("test %v: modified read-only y argument", i)
		}
		if math.IsNaN(got) {
			t.Errorf("test %v: invalid memory read", i)
			continue
		}

		if got != test.want {
			t.Errorf("test %v: unexpected result. want %v, got %v", i, test.want, got)
		}
	}
}

func TestDdotInc(t *testing.T) {
	for i, test := range []struct {
		xData []float64
		yData []float64

		want    float64
		wantRev float64 // Result when one of the vectors is reversed.
	}{
		{
			xData:   []float64{2},
			yData:   []float64{-3},
			want:    -6,
			wantRev: -6,
		},
		{
			xData:   []float64{2, 3},
			yData:   []float64{-3, 4},
			want:    6,
			wantRev: -1,
		},
		{
			xData:   []float64{2, 3, -4},
			yData:   []float64{-3, 4, 5},
			want:    -14,
			wantRev: 34,
		},
		{
			xData:   []float64{2, 3, -4, -5},
			yData:   []float64{-3, 4, 5, -6},
			want:    16,
			wantRev: 2,
		},
		{
			xData:   []float64{0, 2, 3, -4, -5},
			yData:   []float64{0, -3, 4, 5, -6},
			want:    16,
			wantRev: 34,
		},
		{
			xData:   []float64{0, 0, 2, 3, -4, -5},
			yData:   []float64{0, 1, -3, 4, 5, -6},
			want:    16,
			wantRev: -5,
		},
		{
			xData:   []float64{0, 0, 1, 1, 2, -3, -4},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6},
			want:    4,
			wantRev: -4,
		},
		{
			xData:   []float64{0, 0, 1, 1, 2, -3, -4, 5},
			yData:   []float64{0, 1, 0, 3, -4, 5, -6, 7},
			want:    39,
			wantRev: 3,
		},
	} {
		const msgGuard = "%v: out-of-bounds write to %v argument\nfront guard: %v\nback guard: %v"

		for _, incX := range []int{-7, -3, -2, -1, 1, 2, 3, 7} {
			for _, incY := range []int{-7, -3, -2, -1, 1, 2, 3, 7} {
				n := len(test.xData)
				x, xFront, xBack := newGuardedVector(test.xData, incX)
				y, yFront, yBack := newGuardedVector(test.yData, incY)

				var ix, iy int
				if incX < 0 {
					ix = (-n + 1) * incX
				}
				if incY < 0 {
					iy = (-n + 1) * incY
				}
				got := DdotInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(ix), uintptr(iy))

				prefix := fmt.Sprintf("test %v, incX = %v, incY = %v", i, incX, incY)
				if !allNaN(xFront) || !allNaN(xBack) {
					t.Errorf(msgGuard, prefix, "x", xFront, xBack)
				}
				if !allNaN(yFront) || !allNaN(yBack) {
					t.Errorf(msgGuard, prefix, "y", yFront, yBack)
				}
				if nonStridedWrite(x, incX) || !equalStrided(test.xData, x, incX) {
					t.Errorf("%v: modified read-only x argument", prefix)
				}
				if nonStridedWrite(y, incY) || !equalStrided(test.yData, y, incY) {
					t.Errorf("%v: modified read-only y argument", prefix)
				}
				if math.IsNaN(got) {
					t.Errorf("%v: invalid memory read", prefix)
					continue
				}

				want := test.want
				if incX*incY < 0 {
					want = test.wantRev
				}
				if got != want {
					t.Errorf("%v: unexpected result. want %v, got %v", prefix, want, got)
				}
			}
		}
	}
}

func BenchmarkDdotUnitaryN1(b *testing.B)      { ddotUnitaryBenchmark(b, 1) }
func BenchmarkDdotUnitaryN2(b *testing.B)      { ddotUnitaryBenchmark(b, 2) }
func BenchmarkDdotUnitaryN3(b *testing.B)      { ddotUnitaryBenchmark(b, 3) }
func BenchmarkDdotUnitaryN4(b *testing.B)      { ddotUnitaryBenchmark(b, 4) }
func BenchmarkDdotUnitaryN10(b *testing.B)     { ddotUnitaryBenchmark(b, 10) }
func BenchmarkDdotUnitaryN100(b *testing.B)    { ddotUnitaryBenchmark(b, 100) }
func BenchmarkDdotUnitaryN1000(b *testing.B)   { ddotUnitaryBenchmark(b, 1000) }
func BenchmarkDdotUnitaryN10000(b *testing.B)  { ddotUnitaryBenchmark(b, 10000) }
func BenchmarkDdotUnitaryN100000(b *testing.B) { ddotUnitaryBenchmark(b, 100000) }

var r float64

func ddotUnitaryBenchmark(b *testing.B, n int) {
	x := make([]float64, n)
	for i := range x {
		x[i] = rand.Float64()
	}
	y := make([]float64, n)
	for i := range y {
		y[i] = rand.Float64()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = DdotUnitary(x, y)
	}
}

func BenchmarkDdotIncN1Inc1(b *testing.B) { ddotIncBenchmark(b, 1, 1) }

func BenchmarkDdotIncN2Inc1(b *testing.B)  { ddotIncBenchmark(b, 2, 1) }
func BenchmarkDdotIncN2Inc2(b *testing.B)  { ddotIncBenchmark(b, 2, 2) }
func BenchmarkDdotIncN2Inc4(b *testing.B)  { ddotIncBenchmark(b, 2, 4) }
func BenchmarkDdotIncN2Inc10(b *testing.B) { ddotIncBenchmark(b, 2, 10) }

func BenchmarkDdotIncN3Inc1(b *testing.B)  { ddotIncBenchmark(b, 3, 1) }
func BenchmarkDdotIncN3Inc2(b *testing.B)  { ddotIncBenchmark(b, 3, 2) }
func BenchmarkDdotIncN3Inc4(b *testing.B)  { ddotIncBenchmark(b, 3, 4) }
func BenchmarkDdotIncN3Inc10(b *testing.B) { ddotIncBenchmark(b, 3, 10) }

func BenchmarkDdotIncN4Inc1(b *testing.B)  { ddotIncBenchmark(b, 4, 1) }
func BenchmarkDdotIncN4Inc2(b *testing.B)  { ddotIncBenchmark(b, 4, 2) }
func BenchmarkDdotIncN4Inc4(b *testing.B)  { ddotIncBenchmark(b, 4, 4) }
func BenchmarkDdotIncN4Inc10(b *testing.B) { ddotIncBenchmark(b, 4, 10) }

func BenchmarkDdotIncN10Inc1(b *testing.B)  { ddotIncBenchmark(b, 10, 1) }
func BenchmarkDdotIncN10Inc2(b *testing.B)  { ddotIncBenchmark(b, 10, 2) }
func BenchmarkDdotIncN10Inc4(b *testing.B)  { ddotIncBenchmark(b, 10, 4) }
func BenchmarkDdotIncN10Inc10(b *testing.B) { ddotIncBenchmark(b, 10, 10) }

func BenchmarkDdotIncN1000Inc1(b *testing.B)  { ddotIncBenchmark(b, 1000, 1) }
func BenchmarkDdotIncN1000Inc2(b *testing.B)  { ddotIncBenchmark(b, 1000, 2) }
func BenchmarkDdotIncN1000Inc4(b *testing.B)  { ddotIncBenchmark(b, 1000, 4) }
func BenchmarkDdotIncN1000Inc10(b *testing.B) { ddotIncBenchmark(b, 1000, 10) }

func BenchmarkDdotIncN100000Inc1(b *testing.B)  { ddotIncBenchmark(b, 100000, 1) }
func BenchmarkDdotIncN100000Inc2(b *testing.B)  { ddotIncBenchmark(b, 100000, 2) }
func BenchmarkDdotIncN100000Inc4(b *testing.B)  { ddotIncBenchmark(b, 100000, 4) }
func BenchmarkDdotIncN100000Inc10(b *testing.B) { ddotIncBenchmark(b, 100000, 10) }

func BenchmarkDdotIncN100000IncM1(b *testing.B)  { ddotIncBenchmark(b, 100000, -1) }
func BenchmarkDdotIncN100000IncM2(b *testing.B)  { ddotIncBenchmark(b, 100000, -2) }
func BenchmarkDdotIncN100000IncM4(b *testing.B)  { ddotIncBenchmark(b, 100000, -4) }
func BenchmarkDdotIncN100000IncM10(b *testing.B) { ddotIncBenchmark(b, 100000, -10) }

func ddotIncBenchmark(b *testing.B, n, inc int) {
	absInc := inc
	if inc < 0 {
		absInc = -inc
	}
	x := make([]float64, (n-1)*absInc+1)
	for i := range x {
		x[i] = rand.Float64()
	}
	y := make([]float64, (n-1)*absInc+1)
	for i := range y {
		y[i] = rand.Float64()
	}
	var ini int
	if inc < 0 {
		ini = (-n + 1) * inc
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r = DdotInc(x, y, uintptr(n), uintptr(inc), uintptr(inc), uintptr(ini), uintptr(ini))
	}
}
