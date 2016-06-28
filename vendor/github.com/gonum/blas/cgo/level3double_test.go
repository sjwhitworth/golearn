package cgo

import (
	"testing"

	"github.com/gonum/blas/testblas"
)

func TestDgemm(t *testing.T) {
	testblas.TestDgemm(t, impl)
}

func TestDsymm(t *testing.T) {
	testblas.DsymmTest(t, impl)
}

func TestDtrsm(t *testing.T) {
	testblas.DtrsmTest(t, impl)
}

func TestDsyrk(t *testing.T) {
	testblas.DsyrkTest(t, impl)
}

func TestDsyr2k(t *testing.T) {
	testblas.Dsyr2kTest(t, impl)
}

func TestDtrmm(t *testing.T) {
	testblas.DtrmmTest(t, impl)
}
