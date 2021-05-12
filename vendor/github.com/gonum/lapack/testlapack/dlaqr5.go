// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
)

type Dlaqr5er interface {
	Dlaqr5(wantt, wantz bool, kacc22 int, n, ktop, kbot, nshfts int, sr, si []float64, h []float64, ldh int, iloz, ihiz int, z []float64, ldz int, v []float64, ldv int, u []float64, ldu int, nh int, wh []float64, ldwh int, nv int, wv []float64, ldwv int)
}

type Dlaqr5test struct {
	WantT          bool
	N              int
	NShifts        int
	KTop, KBot     int
	ShiftR, ShiftI []float64
	H              []float64

	HWant []float64
	ZWant []float64
}

func Dlaqr5Test(t *testing.T, impl Dlaqr5er) {
	// Test without using reference data.
	rnd := rand.New(rand.NewSource(1))
	for _, n := range []int{1, 2, 3, 4, 5, 6, 10, 30} {
		for _, extra := range []int{0, 1, 20} {
			for _, kacc22 := range []int{0, 1, 2} {
				for cas := 0; cas < 100; cas++ {
					testDlaqr5(t, impl, n, extra, kacc22, rnd)
				}
			}
		}
	}

	// Test using reference data computed by the reference netlib
	// implementation.
	file, err := os.Open(filepath.FromSlash("../testlapack/testdata/dlaqr5data.json.gz"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	r, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tests []Dlaqr5test
	json.NewDecoder(r).Decode(&tests)
	for _, test := range tests {
		wantt := test.WantT
		n := test.N
		nshfts := test.NShifts
		ktop := test.KTop
		kbot := test.KBot
		sr := test.ShiftR
		si := test.ShiftI

		for _, extra := range []int{0, 1, 10} {
			v := randomGeneral(nshfts/2, 3, 3+extra, rnd)
			u := randomGeneral(3*nshfts-3, 3*nshfts-3, 3*nshfts-3+extra, rnd)
			nh := n
			wh := randomGeneral(3*nshfts-3, n, n+extra, rnd)
			nv := n
			wv := randomGeneral(n, 3*nshfts-3, 3*nshfts-3+extra, rnd)

			h := nanGeneral(n, n, n+extra)

			for _, kacc22 := range []int{0, 1, 2} {
				copyMatrix(n, n, h.Data, h.Stride, test.H)
				z := eye(n, n+extra)

				impl.Dlaqr5(wantt, true, kacc22,
					n, ktop, kbot,
					nshfts, sr, si,
					h.Data, h.Stride,
					0, n-1, z.Data, z.Stride,
					v.Data, v.Stride,
					u.Data, u.Stride,
					nh, wh.Data, wh.Stride,
					nv, wv.Data, wv.Stride)

				prefix := fmt.Sprintf("wantt=%v, n=%v, nshfts=%v, ktop=%v, kbot=%v, extra=%v, kacc22=%v",
					wantt, n, nshfts, ktop, kbot, extra, kacc22)
				if !equalApprox(n, n, h.Data, h.Stride, test.HWant, 1e-13) {
					t.Errorf("Case %v: unexpected matrix H\nh    =%v\nhwant=%v", prefix, h.Data, test.HWant)
				}
				if !equalApprox(n, n, z.Data, z.Stride, test.ZWant, 1e-13) {
					t.Errorf("Case %v: unexpected matrix Z\nz    =%v\nzwant=%v", prefix, z.Data, test.ZWant)
				}
			}
		}
	}
}

func testDlaqr5(t *testing.T, impl Dlaqr5er, n, extra, kacc22 int, rnd *rand.Rand) {
	wantt := true
	wantz := true
	nshfts := 2 * n
	sr := make([]float64, nshfts)
	si := make([]float64, nshfts)
	for i := 0; i < n; i++ {
		re := rnd.NormFloat64()
		im := rnd.NormFloat64()
		sr[2*i], sr[2*i+1] = re, re
		si[2*i], si[2*i+1] = im, -im
	}
	ktop := rnd.Intn(n)
	kbot := rnd.Intn(n)
	if kbot < ktop {
		ktop, kbot = kbot, ktop
	}

	v := randomGeneral(nshfts/2, 3, 3+extra, rnd)
	u := randomGeneral(3*nshfts-3, 3*nshfts-3, 3*nshfts-3+extra, rnd)
	nh := n
	wh := randomGeneral(3*nshfts-3, n, n+extra, rnd)
	nv := n
	wv := randomGeneral(n, 3*nshfts-3, 3*nshfts-3+extra, rnd)

	h := randomHessenberg(n, n+extra, rnd)
	if ktop > 0 {
		h.Data[ktop*h.Stride+ktop-1] = 0
	}
	if kbot < n-1 {
		h.Data[(kbot+1)*h.Stride+kbot] = 0
	}
	hCopy := h
	hCopy.Data = make([]float64, len(h.Data))
	copy(hCopy.Data, h.Data)

	z := eye(n, n+extra)

	impl.Dlaqr5(wantt, wantz, kacc22,
		n, ktop, kbot,
		nshfts, sr, si,
		h.Data, h.Stride,
		0, n-1, z.Data, z.Stride,
		v.Data, v.Stride,
		u.Data, u.Stride,
		nh, wh.Data, wh.Stride,
		nv, wv.Data, wv.Stride)

	prefix := fmt.Sprintf("Case n=%v, extra=%v, kacc22=%v", n, extra, kacc22)

	if !generalOutsideAllNaN(h) {
		t.Errorf("%v: out-of-range write to H\n%v", prefix, h.Data)
	}
	if !generalOutsideAllNaN(z) {
		t.Errorf("%v: out-of-range write to Z\n%v", prefix, z.Data)
	}
	if !generalOutsideAllNaN(u) {
		t.Errorf("%v: out-of-range write to U\n%v", prefix, u.Data)
	}
	if !generalOutsideAllNaN(v) {
		t.Errorf("%v: out-of-range write to V\n%v", prefix, v.Data)
	}
	if !generalOutsideAllNaN(wh) {
		t.Errorf("%v: out-of-range write to WH\n%v", prefix, wh.Data)
	}
	if !generalOutsideAllNaN(wv) {
		t.Errorf("%v: out-of-range write to WV\n%v", prefix, wv.Data)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < i-1; j++ {
			if h.Data[i*h.Stride+j] != 0 {
				t.Errorf("%v: H is not Hessenberg, H[%v,%v]!=0", prefix, i, j)
			}
		}
	}
	if !isOrthonormal(z) {
		t.Errorf("%v: Z is not orthogonal", prefix)
	}
	// Construct Z^T * HOrig * Z and check that it is equal to H from Dlaqr5.
	hz := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, hCopy, z, 0, hz)
	zhz := blas64.General{
		Rows:   n,
		Cols:   n,
		Stride: n,
		Data:   make([]float64, n*n),
	}
	blas64.Gemm(blas.Trans, blas.NoTrans, 1, z, hz, 0, zhz)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			diff := zhz.Data[i*zhz.Stride+j] - h.Data[i*h.Stride+j]
			if math.Abs(diff) > 1e-13 {
				t.Errorf("%v: Z^T*HOrig*Z and H are not equal, diff at [%v,%v]=%v", prefix, i, j, diff)
			}
		}
	}
}
