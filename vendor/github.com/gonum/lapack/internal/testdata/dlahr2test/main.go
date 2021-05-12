// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// dlahr2test generates test data for Dlahr2. Test cases are stored in
// gzip-compressed JSON file testlapack/testdata/dlahr2data.json.gz which is
// read during testing by testlapack/dlahr2.go.
//
// This program uses cgo to call Fortran version of DLAHR2. Therefore, matrices
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

type Dlahr2Test struct {
	N, K, NB int
	A        []float64

	AWant   []float64
	TWant   []float64
	YWant   []float64
	TauWant []float64
}

func main() {
	file, err := os.Create(filepath.FromSlash("../../../testlapack/testdata/dlahr2data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	w := gzip.NewWriter(file)

	rnd := rand.New(rand.NewSource(1))

	var tests []Dlahr2Test
	for _, n := range []int{4, 5, 6, 7, 11} {
		for k := 0; k <= n/2; k++ {
			for nb := 1; nb <= k; nb++ {
				ain := genrand(n, n-k+1, rnd)
				a := make([]float64, len(ain))
				copy(a, ain)

				t := genrand(nb, nb, rnd)
				y := genrand(n, nb, rnd)
				tau := genrand(nb, 1, rnd)

				netlib.Dlahr2(n, k, nb, a, n, tau, t, nb, y, n)

				tests = append(tests, Dlahr2Test{
					N:       n,
					K:       k,
					NB:      nb,
					A:       rowMajor(n, n-k+1, ain),
					AWant:   rowMajor(n, n-k+1, a),
					TWant:   rowMajor(nb, nb, t),
					YWant:   rowMajor(n, nb, y),
					TauWant: tau,
				})
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
