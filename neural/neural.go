//Package neural contains Neural Network functions.
package neural

import (
	"github.com/gonum/matrix/mat64"
)

type ActivationFunction func(float64) float64

// First function is always the forward activation function
// Second function is always the backward activation function
type NeuralFunction struct {
	Forward  ActivationFunction
	Backward ActivationFunction
}

// LayerFuncs are vectorised layer value transformation functions
// (e.g. sigmoid). They must operate in-place.
type LayerFunc func(*mat64.Dense)
