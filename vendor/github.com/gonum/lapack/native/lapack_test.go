// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"testing"

	"github.com/gonum/lapack/testlapack"
)

var impl = Implementation{}

func TestDbdsqr(t *testing.T) {
	testlapack.DbdsqrTest(t, impl)
}

func TestDgebd2(t *testing.T) {
	testlapack.Dgebd2Test(t, impl)
}

func TestDgebrd(t *testing.T) {
	testlapack.DgebrdTest(t, impl)
}

func TestDgecon(t *testing.T) {
	testlapack.DgeconTest(t, impl)
}

func TestDgehd2(t *testing.T) {
	testlapack.Dgehd2Test(t, impl)
}

func TestDgehrd(t *testing.T) {
	testlapack.DgehrdTest(t, impl)
}

func TestDgelqf(t *testing.T) {
	testlapack.DgelqfTest(t, impl)
}

func TestDgelq2(t *testing.T) {
	testlapack.Dgelq2Test(t, impl)
}

func TestDgeql2(t *testing.T) {
	testlapack.Dgeql2Test(t, impl)
}

func TestDgels(t *testing.T) {
	testlapack.DgelsTest(t, impl)
}

func TestDgeqr2(t *testing.T) {
	testlapack.Dgeqr2Test(t, impl)
}

func TestDgeqrf(t *testing.T) {
	testlapack.DgeqrfTest(t, impl)
}

func TestDgesvd(t *testing.T) {
	testlapack.DgesvdTest(t, impl)
}

func TestDgetri(t *testing.T) {
	testlapack.DgetriTest(t, impl)
}

func TestDgetf2(t *testing.T) {
	testlapack.Dgetf2Test(t, impl)
}

func TestDgetrf(t *testing.T) {
	testlapack.DgetrfTest(t, impl)
}

func TestDgetrs(t *testing.T) {
	testlapack.DgetrsTest(t, impl)
}

func TestDlabrd(t *testing.T) {
	testlapack.DlabrdTest(t, impl)
}

func TestDlacpy(t *testing.T) {
	testlapack.DlacpyTest(t, impl)
}

func TestDlae2(t *testing.T) {
	testlapack.Dlae2Test(t, impl)
}

func TestDlaev2(t *testing.T) {
	testlapack.Dlaev2Test(t, impl)
}

func TestDlahqr(t *testing.T) {
	testlapack.DlahqrTest(t, impl)
}

func TestDlahr2(t *testing.T) {
	testlapack.Dlahr2Test(t, impl)
}

func TestDlange(t *testing.T) {
	testlapack.DlangeTest(t, impl)
}

func TestDlas2(t *testing.T) {
	testlapack.Dlas2Test(t, impl)
}

func TestDlasy2(t *testing.T) {
	testlapack.Dlasy2Test(t, impl)
}

func TestDlanst(t *testing.T) {
	testlapack.DlanstTest(t, impl)
}

func TestDlansy(t *testing.T) {
	testlapack.DlansyTest(t, impl)
}

func TestDlantr(t *testing.T) {
	testlapack.DlantrTest(t, impl)
}

func TestDlanv2(t *testing.T) {
	testlapack.Dlanv2Test(t, impl)
}

func TestDlaqr1(t *testing.T) {
	testlapack.Dlaqr1Test(t, impl)
}

func TestDlaqr5(t *testing.T) {
	testlapack.Dlaqr5Test(t, impl)
}

func TestDlarf(t *testing.T) {
	testlapack.DlarfTest(t, impl)
}

func TestDlarfb(t *testing.T) {
	testlapack.DlarfbTest(t, impl)
}

func TestDlarfg(t *testing.T) {
	testlapack.DlarfgTest(t, impl)
}

func TestDlarft(t *testing.T) {
	testlapack.DlarftTest(t, impl)
}

func TestDlarfx(t *testing.T) {
	testlapack.DlarfxTest(t, impl)
}

func TestDlartg(t *testing.T) {
	testlapack.DlartgTest(t, impl)
}

func TestDlasq1(t *testing.T) {
	testlapack.Dlasq1Test(t, impl)
}

func TestDlasq2(t *testing.T) {
	testlapack.Dlasq2Test(t, impl)
}

func TestDlasq3(t *testing.T) {
	testlapack.Dlasq3Test(t, impl)
}

func TestDlasq4(t *testing.T) {
	testlapack.Dlasq4Test(t, impl)
}

func TestDlasq5(t *testing.T) {
	testlapack.Dlasq5Test(t, impl)
}

func TestDlasr(t *testing.T) {
	testlapack.DlasrTest(t, impl)
}

func TestDlasv2(t *testing.T) {
	testlapack.Dlasv2Test(t, impl)
}

func TestDlatrd(t *testing.T) {
	testlapack.DlatrdTest(t, impl)
}

func TestDorg2r(t *testing.T) {
	testlapack.Dorg2rTest(t, impl)
}

func TestDorgbr(t *testing.T) {
	testlapack.DorgbrTest(t, impl)
}

func TestDorghr(t *testing.T) {
	testlapack.DorghrTest(t, impl)
}

func TestDorg2l(t *testing.T) {
	testlapack.Dorg2lTest(t, impl)
}

func TestDorgl2(t *testing.T) {
	testlapack.Dorgl2Test(t, impl)
}

func TestDorglq(t *testing.T) {
	testlapack.DorglqTest(t, impl)
}

func TestDorgql(t *testing.T) {
	testlapack.DorgqlTest(t, impl)
}

func TestDorgqr(t *testing.T) {
	testlapack.DorgqrTest(t, impl)
}

func TestDorgtr(t *testing.T) {
	testlapack.DorgtrTest(t, impl)
}

func TestDormbr(t *testing.T) {
	testlapack.DormbrTest(t, impl)
}

func TestDormhr(t *testing.T) {
	testlapack.DormhrTest(t, impl)
}

func TestDorml2(t *testing.T) {
	testlapack.Dorml2Test(t, impl)
}

func TestDormlq(t *testing.T) {
	testlapack.DormlqTest(t, impl)
}

func TestDormqr(t *testing.T) {
	testlapack.DormqrTest(t, impl)
}

func TestDorm2r(t *testing.T) {
	testlapack.Dorm2rTest(t, impl)
}

func TestDpocon(t *testing.T) {
	testlapack.DpoconTest(t, impl)
}

func TestDpotf2(t *testing.T) {
	testlapack.Dpotf2Test(t, impl)
}

func TestDpotrf(t *testing.T) {
	testlapack.DpotrfTest(t, impl)
}

func TestDrscl(t *testing.T) {
	testlapack.DrsclTest(t, impl)
}

func TestDsteqr(t *testing.T) {
	testlapack.DsteqrTest(t, impl)
}

func TestDsterf(t *testing.T) {
	testlapack.DsterfTest(t, impl)
}

func TestDsyev(t *testing.T) {
	testlapack.DsyevTest(t, impl)
}

func TestDsytd2(t *testing.T) {
	testlapack.Dsytd2Test(t, impl)
}

func TestDsytrd(t *testing.T) {
	testlapack.DsytrdTest(t, impl)
}

func TestDtrcon(t *testing.T) {
	testlapack.DtrconTest(t, impl)
}

func TestDtrti2(t *testing.T) {
	testlapack.Dtrti2Test(t, impl)
}

func TestDtrtri(t *testing.T) {
	testlapack.DtrtriTest(t, impl)
}

func TestIladlc(t *testing.T) {
	testlapack.IladlcTest(t, impl)
}

func TestIladlr(t *testing.T) {
	testlapack.IladlrTest(t, impl)
}
