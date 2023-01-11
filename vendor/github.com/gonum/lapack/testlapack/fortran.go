// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"fmt"
	"strings"
)

// This file implements types for helping to convert to Fortran testing capabilities.

// fortran64 is a float64 type that prints as a double precision constant in
// Fortran format.
type fortran64 float64

func (f fortran64) String() string {
	// Replace exponent with D
	s := fmt.Sprintf("%0.16E", f)
	s = strings.Replace(s, "E", "D", 1)
	return s
}

// printFortranArray prints a Go slice as an array that can be copied into a
// fortran script.
func printFortranArray(z []float64, name string) {
	fmt.Printf("%s(1:%d) = (/%v, &\n", name, len(z), fortran64(z[0]))
	for i := 1; i < len(z)-1; i++ {
		fmt.Printf("%v, &\n", fortran64(z[i]))
	}
	fmt.Printf("%s/)\n", fortran64(z[len(z)-1]))
}
