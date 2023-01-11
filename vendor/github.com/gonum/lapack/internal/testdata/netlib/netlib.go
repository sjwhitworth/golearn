// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netlib

// void dlahr2_(int* n, int* k, int* nb, double* a, int* lda, double* tau, double* t, int* ldt, double* y, int* ldy);
//
// void dlaqr5_(int* wantt, int* wantz, int* kacc22, int* n, int* ktop, int* kbot, int* nshfts,
//              double* sr, double* si, double* h, int* ldh, int* iloz, int* ihiz,
//              double* z, int* ldz, double* v, int* ldv, double* u, int* ldu,
//              int* nv, double* wv, int* ldwv, int* nh, double* wh, int* ldwh);
import "C"

func Dlahr2(n, k, nb int, a []float64, lda int, tau, t []float64, ldt int, y []float64, ldy int) {
	func() {
		n := C.int(n)
		k := C.int(k)
		nb := C.int(nb)
		lda := C.int(lda)
		ldt := C.int(ldt)
		ldy := C.int(ldy)
		C.dlahr2_((*C.int)(&n), (*C.int)(&k), (*C.int)(&nb),
			(*C.double)(&a[0]), (*C.int)(&lda),
			(*C.double)(&tau[0]),
			(*C.double)(&t[0]), (*C.int)(&ldt),
			(*C.double)(&y[0]), (*C.int)(&ldy))
	}()
}

func Dlaqr5(wantt, wantz bool, kacc22 int, n, ktop, kbot int, nshfts int, sr, si []float64, h []float64,
	ldh int, iloz, ihiz int, z []float64, ldz int, v []float64, ldv int,
	u []float64, ldu int, nh int, wh []float64, ldwh int, nv int, wv []float64, ldwv int) {
	func() {
		wt := C.int(0)
		if wantt {
			wt = 1
		}
		wz := C.int(0)
		if wantz {
			wz = 1
		}
		kacc22 := C.int(kacc22)
		n := C.int(n)
		ktop := C.int(ktop)
		kbot := C.int(kbot)
		nshfts := C.int(nshfts)
		ldh := C.int(ldh)
		iloz := C.int(iloz)
		ihiz := C.int(ihiz)
		ldz := C.int(ldz)
		ldv := C.int(ldv)
		ldu := C.int(ldu)
		nh := C.int(nh)
		ldwh := C.int(ldwh)
		nv := C.int(nv)
		ldwv := C.int(ldwv)
		C.dlaqr5_((*C.int)(&wt), (*C.int)(&wz), (*C.int)(&kacc22),
			(*C.int)(&n), (*C.int)(&ktop), (*C.int)(&kbot),
			(*C.int)(&nshfts), (*C.double)(&sr[0]), (*C.double)(&si[0]),
			(*C.double)(&h[0]), (*C.int)(&ldh),
			(*C.int)(&iloz), (*C.int)(&ihiz), (*C.double)(&z[0]), (*C.int)(&ldz),
			(*C.double)(&v[0]), (*C.int)(&ldv),
			(*C.double)(&u[0]), (*C.int)(&ldu),
			(*C.int)(&nh), (*C.double)(&wh[0]), (*C.int)(&ldwh),
			(*C.int)(&nv), (*C.double)(&wv[0]), (*C.int)(&ldwv))
	}()
}
