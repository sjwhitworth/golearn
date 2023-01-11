// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

func conj(c complex64) complex64 { return complex(real(c), -imag(c)) }
