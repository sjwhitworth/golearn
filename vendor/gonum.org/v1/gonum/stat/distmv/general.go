// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

var (
	badQuantile      = "distmv: quantile not between 0 and 1"
	badReceiver      = "distmv: input slice is not nil or the correct length"
	badSizeMismatch  = "distmv: size mismatch"
	badZeroDimension = "distmv: zero dimensional input"
	nonPosDimension  = "distmv: non-positive dimension input"
)

const logTwoPi = 1.8378770664093454835606594728112352797227949472755668

// useAs gets a slice of size n. If len(x) == n, x is returned, if len(x) == 0
// then a slice is returned of length n.
func reuseAs(x []float64, n int) []float64 {
	if len(x) == n {
		return x
	}
	if len(x) == 0 {
		if cap(x) >= n {
			return x[:n]
		}
		return make([]float64, n)
	}
	panic(badReceiver)
}
