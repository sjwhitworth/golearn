package distance

import (
	"math"
  )

func eucledianDistance (a []float64, b []float64) float64 {
	sum := 0.0;
	for i := 0; i < cap(a); i++ {
		sum += math.Abs(a[i] - b[i])
	}
	return sum;
}
