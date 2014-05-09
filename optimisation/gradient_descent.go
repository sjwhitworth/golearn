package optimisation

import "github.com/gonum/matrix/mat64"

// Batch gradient descent finds the local minimum of a function.
// See http://en.wikipedia.org/wiki/Gradient_descent for more details.
func BatchGradientDescent(x, y, theta *mat64.Dense, alpha float64, epoch int) *mat64.Dense {
	m, _ := y.Dims()
	for i := 0; i < epoch; i++ {
		xFlat := mat64.DenseCopyOf(x)
		xFlat.TCopy(xFlat)
		temp := mat64.DenseCopyOf(x)

		// Calculate our best prediction, given theta
		temp.Mul(temp, theta)

		// Calculate our error from the real values
		temp.Sub(temp, y)
		xFlat.Mul(xFlat, temp)

		// Temporary hack to get around the fact there is no scalar division in mat64
		xFlatRow, _ := xFlat.Dims()
		gradient := make([]float64, 0)
		for k := 0; k < xFlatRow; k++ {
			row := xFlat.RowView(k)
			for v := range row {
				divd := row[v] / float64(m) * alpha
				gradient = append(gradient, divd)
			}
		}
		grows := len(gradient)
		grad := mat64.NewDense(grows, 1, gradient)
		theta.Sub(theta, grad)
	}
	return theta
}

// Stochastic gradient descent updates the parameters of theta on a random row selection from a matrix.
// It is faster as it does not compute the cost function over the entire dataset every time.
// It instead calculates the error parameters over only one row of the dataset at a time.
// In return, there is a trade off for accuracy. This is minimised by running multiple SGD processes
// (the number of goroutines spawned is specified by the procs variable) in parallel and taking an average of the result.
func StochasticGradientDescent(x, y, theta *mat64.Dense, alpha float64, epoch, procs int) *mat64.Dense {
	m, _ := y.Dims()
	resultPipe := make(chan *mat64.Dense)
	results := make([]*mat64.Dense, 0)

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
					datXFlat := mat64.DenseCopyOf(datX)
					datXFlat.TCopy(datXFlat)
					datX.Mul(datX, thetaCopy)
					datX.Sub(datX, datY)
					datXFlat.Mul(datXFlat, datX)

					// Horrible hack to get around the fact there is no elementwise division in mat64
					xFlatRow, _ := datXFlat.Dims()
					gradient := make([]float64, 0)
					for i := 0; i < xFlatRow; i++ {
						row := datXFlat.RowView(i)
						for i := range row {
							divd := row[i] / float64(m) * alpha
							gradient = append(gradient, divd)
						}
					}
					grows := len(gradient)
					grad := mat64.NewDense(grows, 1, gradient)
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
	base := matrices[0]
	rows, _ := base.Dims()

	for i := 1; i < len(matrices); i++ {
		base.Add(base, matrices[i])
	}

	averaged := make([]float64, 0)
	for i := 0; i < rows; i++ {
		row := base.RowView(i)
		for i := range row {
			divd := row[i] / float64(len(matrices))
			averaged = append(averaged, divd)
		}
	}

	baseDense := mat64.NewDense(rows, 1, averaged)
	return baseDense
}
