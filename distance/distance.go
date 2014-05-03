package distance

import (
	"math"
)

func eucledianDistance (p1 []float64, p2 []float64) float64 {
	sum := float64(0)
	for i := 0; i < cap(p1); i++ {
		sum += math.Pow(p1[i] - p2[i], 2)
	}
	return math.Sqrt(sum);
}
