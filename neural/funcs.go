package neural

import (
	"math"
)

// SigmoidForward function does S(t) = \frac{1}{1 + e^{-t}}.
//
// See http://en.wikipedia.org/wiki/Sigmoid_function
var Sigmoid NeuralFunction = NeuralFunction{
	func(v float64) float64 { return 1.0 / (1.0 + math.Exp(-v)) },
	func(v float64) float64 { return v * (1 - v) },
}

// TanhForward function does tanh(t) = \frac{1 - e^{-2t}}{1 + e^{-2t}}
//
// See https://en.wikipedia.org/wiki/Hyperbolic_functions
var Tanh NeuralFunction = NeuralFunction{
	func(v float64) float64 { return math.Tanh(v) },
	func(v float64) float64 { return 1 - (v * v) },
}

// LinearFunction doesn't modify the value
var Linear NeuralFunction = NeuralFunction{
	func(v float64) float64 { return v },
	func(v float64) float64 { return 1.0 },
}

// Rectified Linear function
// https://www.wikiwand.com/en/Rectifier_(neural_networks)
var SoftplusRectifier NeuralFunction = NeuralFunction{
	func(v float64) float64 { return math.Log(1 + math.Exp(v)) },
	func(v float64) float64 { return v * (1 - v) },
}
