package pairwise

import (
	"github.com/gonum/matrix/mat64"
	"math"
)

// Cosine implements Cosine distance, which measures the cosine
// of the angle between two vectors.
//
// Further reading: http://en.wikipedia.org/wiki/Cosine_similarity
type Cosine struct{}

func NewCosine() *Cosine {
	return &Cosine{}
}

// Distance computes the cosine distance between two vectors.
func (c *Cosine) Distance(vectorX *mat64.Dense, vectorY *mat64.Dense) float64 {

	r1, c1 := vectorX.Dims()
	r2, c2 := vectorY.Dims()
	if r1 != r2 || c1 != c2 {
		panic(mat64.ErrShape)
	}

	// Copute A.B
	top := vectorX.Dot(vectorY)

	// Compute ||A|| ||B||
	tmp := 0.0
	magFunc := func(r, c int, v float64) float64 {
		tmp += v * v
		return v
	}
	vectorX.Apply(magFunc, vectorX)
	bottom := math.Sqrt(tmp)
	tmp = 0
	vectorY.Apply(magFunc, vectorY)
	bottom *= math.Sqrt(tmp)

	return top / bottom
}
