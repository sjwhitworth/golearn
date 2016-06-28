// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This program generates test data for Dlaqr5. Test cases are stored in
// gzip-compressed JSON file testlapack/testdata/dlaqr5data.json.gz which is
// read during testing by testlapack/dlaqr5.go.
//
// This program uses cgo to call Fortran version of DLAQR5. Therefore, matrices
// passed to the Fortran routine are in column-major format but are written into
// the output file in row-major format.
package main

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/gonum/lapack/internal/testdata/netlib"
)

type Dlaqr5Test struct {
	WantT          bool
	N              int
	NShifts        int
	KTop, KBot     int
	ShiftR, ShiftI []float64
	H              []float64

	HWant []float64
	ZWant []float64
}

func main() {
	file, err := os.Create(filepath.FromSlash("../../../testlapack/testdata/dlaqr5data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := gzip.NewWriter(file)

	rnd := rand.New(rand.NewSource(1))

	var tests []Dlaqr5Test
	for _, wantt := range []bool{true, false} {
		for _, n := range []int{2, 3, 4, 5, 6, 7, 11} {
			for k := 0; k <= min(5, n); k++ {
				npairs := k
				if npairs == 0 {
					npairs = 2 * n
				}
				for ktop := 0; ktop < n-1; ktop++ {
					for kbot := ktop + 1; kbot < n; kbot++ {
						sr, si := shiftpairs(npairs, rnd)
						nshfts := len(sr)

						v := genrand(nshfts/2, 3, rnd)
						u := genrand(3*nshfts-3, 3*nshfts-3, rnd)
						wh := genrand(3*nshfts-3, n, rnd)
						nh := n
						wv := genrand(n, 3*nshfts-3, rnd)
						nv := n

						h := hessrand(n, rnd)
						if ktop > 0 {
							h[ktop+(ktop-1)*n] = 0
						}
						if kbot < n-1 {
							h[kbot+1+kbot*n] = 0
						}
						hin := make([]float64, len(h))
						copy(hin, h)
						z := eye(n)

						netlib.Dlaqr5(wantt, true, 2,
							n, ktop+1, kbot+1,
							nshfts, sr, si,
							h, n,
							1, n, z, n,
							v, 3,
							u, 3*nshfts-3,
							nh, wh, nh,
							nv, wv, 3*nshfts-3)

						tests = append(tests, Dlaqr5Test{
							WantT:   wantt,
							N:       n,
							NShifts: nshfts,
							KTop:    ktop,
							KBot:    kbot,
							ShiftR:  sr,
							ShiftI:  si,
							H:       rowMajor(n, n, hin),
							HWant:   rowMajor(n, n, h),
							ZWant:   rowMajor(n, n, z),
						})
					}
				}
			}
		}
	}
	json.NewEncoder(w).Encode(tests)

	err = w.Close()
	if err != nil {
		log.Fatal(err)
	}
}

// genrand returns a general r×c matrix with random entries.
func genrand(r, c int, rnd *rand.Rand) []float64 {
	m := make([]float64, r*c)
	for i := range m {
		m[i] = rnd.NormFloat64()
	}
	return m
}

// eye returns an identity matrix of order n.
func eye(n int) []float64 {
	m := make([]float64, n*n)
	for i := 0; i < n*n; i += n + 1 {
		m[i] = 1
	}
	return m
}

// hessrand returns a Hessenberg matrix of order n with random non-zero entries
// in column-major format.
func hessrand(n int, rnd *rand.Rand) []float64 {
	h := make([]float64, n*n)
	for j := 0; j < n; j++ {
		for i := 0; i <= min(j+1, n-1); i++ {
			h[i+j*n] = rnd.NormFloat64()
		}
	}
	return h
}

// shiftpairs generates k real and complex conjugate shift pairs. That is, the
// length of sr and si is 2*k.
func shiftpairs(k int, rnd *rand.Rand) (sr, si []float64) {
	sr = make([]float64, 2*k)
	si = make([]float64, 2*k)
	for i := 0; i < len(sr); {
		if rnd.Float64() < 0.5 || i == len(sr)-1 {
			sr[i] = rnd.NormFloat64()
			i++
			continue
		}
		// Generate a complex conjugate pair.
		r := rnd.NormFloat64()
		c := rnd.NormFloat64()
		sr[i] = r
		si[i] = c
		sr[i+1] = r
		si[i+1] = -c
		i += 2
	}
	return sr, si
}

// rowMajor returns the given r×c column-major matrix a in row-major format.
func rowMajor(r, c int, a []float64) []float64 {
	if len(a) != r*c {
		panic("testdata: slice length mismatch")
	}
	m := make([]float64, len(a))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			m[i*c+j] = a[i+j*r]
		}
	}
	return m
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
