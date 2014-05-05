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
		temp.Mul(temp, theta)
		temp.Sub(temp, y)
		xFlat.Mul(xFlat, temp)

		// Horrible hack to get around the fact there is no scalar division in mat64
		xFlatRow, _ := xFlat.Dims()
		gradient := make([]float64, 0)
		for i := 0; i < xFlatRow; i++ {
			row := xFlat.RowView(i)
			for i := range row {
				divd := row[i] / float64(m) * alpha
				gradient = append(gradient, divd)
			}
		}
		grows := len(gradient)
		grad := mat64.NewDense(grows, 1, gradient)
		theta.Sub(theta, grad)
	}
	return theta
}

// Stochastic Gradient Descent updates the parameters of theta on a random selection from X,Y.
// It is faster as it does not compute the cost function over the entire dataset every time.
// In return, there is a trade off for accuracy.
// @todo: use goroutines to parallelise training.
func StochasticGradientDescent(x, y, theta *mat64.Dense, alpha float64, epoch int) *mat64.Dense {
	m, _ := y.Dims()
	for i := 0; i < epoch; i++ {
		for k := 0; k < m; k++ {
			datXtemp := x.RowView(k)
			datYtemp := y.RowView(k)
			datX := mat64.NewDense(1, len(datXtemp), datXtemp)
			datY := mat64.NewDense(1, 1, datYtemp)
			datXFlat := mat64.DenseCopyOf(datX)
			datXFlat.TCopy(datXFlat)
			datX.Mul(datX, theta)
			datX.Sub(datX, datY)
			datXFlat.Mul(datXFlat, datX)

			// Horrible hack to get around the fact there is no scalar division in mat64
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
			theta.Sub(theta, grad)
		}

	}
	return theta
}
