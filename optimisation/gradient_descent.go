package optimisation

import "github.com/gonum/matrix/mat64"

// BatchGradientDescent finds the local minimum of a function.
// See http://en.wikipedia.org/wiki/Gradient_descent for more details.
func BatchGradientDescent(x, y, theta *mat64.Dense, alpha float64, epoch int) *mat64.Dense {
	m, _ := y.Dims()
	// Helper function for scalar multiplication
	mult := func(r, c int, v float64) float64 { return v * 1.0 / float64(m) * alpha }

	for i := 0; i < epoch; i++ {
		grad := mat64.DenseCopyOf(x)
		grad.TCopy(grad)
		temp := mat64.DenseCopyOf(x)

		// Calculate our best prediction, given theta
		temp.Mul(temp, theta)

		// Calculate our error from the real values
		temp.Sub(temp, y)
		grad.Mul(grad, temp)

		// Multiply by scalar factor
		grad.Apply(mult, grad)

		// Take a step in gradient direction
		theta.Sub(theta, grad)
	}

	return theta
}

// StochasticGradientDescent updates the parameters of theta on a random row selection from a matrix.
// It is faster as it does not compute the cost function over the entire dataset every time.
// It instead calculates the error parameters over only one row of the dataset at a time.
// In return, there is a trade off for accuracy. This is minimised by running multiple SGD processes
// (the number of goroutines spawned is specified by the procs variable) in parallel and taking an average of the result.
func StochasticGradientDescent(x, y, theta *mat64.Dense, alpha float64, epoch, procs int) *mat64.Dense {
	m, _ := y.Dims()
	resultPipe := make(chan *mat64.Dense)
	results := make([]*mat64.Dense, 0)
	// Helper function for scalar multiplication
	mult := func(r, c int, v float64) float64 { return v * 1.0 / float64(m) * alpha }

	for p := 0; p < procs; p++ {
		go func() {
			// Is this just a pointer to theta?
			thetaCopy := mat64.DenseCopyOf(theta)
			for i := 0; i < epoch; i++ {
				for k := 0; k < m; k++ {
					datXtemp := x.RowView(k)
					datYtemp := y.RowView(k)
					datX := mat64.NewDense(1, len(datXtemp), datXtemp)
					datY := mat64.NewDense(1, 1, datYtemp)
					grad := mat64.DenseCopyOf(datX)
					grad.TCopy(grad)
					datX.Mul(datX, thetaCopy)
					datX.Sub(datX, datY)
					grad.Mul(grad, datX)

					// Multiply by scalar factor
					grad.Apply(mult, grad)

					// Take a step in gradient direction
					thetaCopy.Sub(thetaCopy, grad)
				}

			}
			resultPipe <- thetaCopy
		}()
	}

	for {
		select {
		case d := <-resultPipe:
			results = append(results, d)
			if len(results) == procs {
				return averageTheta(results)
			}
		}
	}
}

func averageTheta(matrices []*mat64.Dense) *mat64.Dense {
	if len(matrices) < 2 {
		panic("Must provide at least two matrices to average")
	}
	invLen :=  1.0 / float64(len(matrices))
	// Helper function for scalar multiplication
	mult := func(r, c int, v float64) float64 { return v * invLen}
	// Sum matrices
	average := matrices[0]
	for i := 1; i < len(matrices); i++ {
		average.Add(average, matrices[i])
	}

	// Calculate the average
	average.Apply(mult, average)
	return average
}
