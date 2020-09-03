package neural

import (
	"bytes"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
)

// Network represents the most general neural network possible
// Weights are stored in a dense matrix, each can have its own
// NeuralFunction.
type Network struct {
	origWeights *mat.Dense
	weights     *mat.Dense       // n * n
	biases      []float64        // n for each neuron
	funcs       []NeuralFunction // for each neuron
	size        int
	input       int
}

// NewNetwork creates a new Network containing size neurons,
// with a certain number dedicated to input, and a pre-defined
// neural function applied to the rest.
//
// Input nodes are set to have a Linear NeuralFunction and are
// connected to themselves for propagation.
func NewNetwork(size int, input int, f NeuralFunction) *Network {
	ret := new(Network)
	ret.weights = mat.NewDense(size, size, make([]float64, size*size))
	ret.biases = make([]float64, size)
	ret.funcs = make([]NeuralFunction, size)
	ret.size = size
	ret.input = input
	for i := range ret.funcs {
		ret.funcs[i] = f
	}
	for i := 0; i < input; i++ {
		ret.funcs[i] = Linear
		ret.SetWeight(i+1, i+1, 1.0)
	}
	return ret
}

// String gets a human-readable representation of this network.
func (n *Network) String() string {
	var buf bytes.Buffer
	var biases bytes.Buffer

	for i := 0; i < n.size; i++ {
		for j := 0; j < n.size; j++ {
			v := n.weights.At(j, i)
			if math.Abs(v) > 0 {
				buf.WriteString(fmt.Sprintf("\t(%d %d %.2f)\n", i+1, j+1, v))
			}
		}
	}

	for _, v := range n.biases {
		biases.WriteString(fmt.Sprintf(" %.2f", v))
	}

	return fmt.Sprintf("Network(%d, %s, %s)", n.size, biases.String(), buf.String())
}

// GetWeight returns the weight between a given source and
// target neuron (counted from 1).
func (n *Network) GetWeight(src, target int) float64 {
	src--
	target--
	return n.weights.At(target, src)
}

// SetWeight sets the weight between a given source and
// target neuron (counted from 1).
func (n *Network) SetWeight(src, target int, v float64) {
	src--
	target--
	n.weights.Set(target, src, v)
}

// SetBias sets the bias at a given neuron (counted from 1).
func (n *Network) SetBias(node int, v float64) {
	if node <= n.input {
		return
	}
	node--
	n.biases[node] = v
}

// GetBias returns the bias at a given neuron (counted from 1).
func (n *Network) GetBias(node int) float64 {
	node--
	return n.biases[node]
}

// Activate propagates the given input matrix (with) across the network
// a certain number of times (up to maxIterations).
//
// The with matrix should be size * size elements, with only the values
// of input neurons set (everything else should be zero).
//
// If the network is conceptually organised into layers, maxIterations
// should be set to the number of layers.
//
// This function overwrites whatever's stored in its first argument.
func (n *Network) Activate(with *mat.Dense, maxIterations int) {

	// Add bias and feed to activation
	biasFunc := func(r, c int, v float64) float64 {
		return v + n.biases[r]
	}
	activFunc := func(r, c int, v float64) float64 {
		return n.funcs[r].Forward(v)
	}

	tmp := new(mat.Dense)
	tmp.CloneFrom(with)

	// Main loop
	for i := 0; i < maxIterations; i++ {
		with.Mul(n.weights, with)
		with.Apply(biasFunc, with)
		with.Apply(activFunc, with)
	}
}

// UpdateWeights takes an output size * 1 output vector and a size * 1
// back-propagated error vector, as well as a learnRate and updates
// the internal weights matrix.
func (n *Network) UpdateWeights(out, err *mat.Dense, learnRate float64) {

	if n.origWeights == nil {
		n.origWeights = mat.DenseCopyOf(n.weights)
	}

	// Multiply that by the learning rate
	mulFunc := func(target, source int, v float64) float64 {
		if target == source {
			return v
		}
		if math.Abs(n.origWeights.At(target, source)) > 0.005 {
			return v + learnRate*out.At(source, 0)*err.At(target, 0)
		}
		return 0.00
	}

	// Add that to the weights
	n.weights.Apply(mulFunc, n.weights)
}

// UpdateBias computes B = B + l.E and updates the bias weights
// from a size * 1 back-propagated error vector.
func (n *Network) UpdateBias(err *mat.Dense, learnRate float64) {

	for i, b := range n.biases {
		if i < n.input {
			continue
		}
		n.biases[i] = b + err.At(i, 0)*learnRate
	}

}

// Error computes the back-propagation error from a given size * 1 output
// vector and a size * 1 error vector for a given number of iterations.
//
// outArg should be the response from Activate.
//
// errArg should be the difference between the output neuron's output and
// that expected, and should be zero everywhere else.
//
// If the network is conceptually organised into n layers, maxIterations
// should be set to n.
func (n *Network) Error(outArg, errArg *mat.Dense, maxIterations int) *mat.Dense {

	// Copy the arguments
	out := mat.DenseCopyOf(outArg)
	err := mat.DenseCopyOf(errArg)

	// err should be the difference between observed and expected
	// for observation nodes only (everything else should be zero)

	// Allocate output vector
	outRows, outCols := out.Dims()
	if outCols != 1 {
		panic("Unsupported output size")
	}

	ret := mat.NewDense(outRows, 1, make([]float64, outRows))

	// Do differential calculation
	diffFunc := func(r, c int, v float64) float64 {
		return n.funcs[r].Backward(v)
	}
	out.Apply(diffFunc, out)

	// Transpose weights matrix
	reverseWeights := mat.DenseCopyOf(n.weights)
	reverseWeights.CloneFrom(n.weights.T())

	// We only need a certain number of passes
	for i := 0; i < maxIterations; i++ {

		// Element-wise multiply errors and derivatives
		err.MulElem(err, out)

		// Add the accumulated error
		ret.Add(ret, err)

		if i != maxIterations-1 {
			// Feed the errors backwards through the network
			err.Mul(reverseWeights, err)
		}
	}

	return ret

}
