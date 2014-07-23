package linear_models

import (
	"errors"

	"github.com/sjwhitworth/golearn/base"

	_ "github.com/gonum/blas"
	"github.com/gonum/blas/cblas"
	"github.com/gonum/matrix/mat64"
)

var (
	NotEnoughDataError  = errors.New("not enough rows to support this many variables.")
	NoTrainingDataError = errors.New("you need to Fit() before you can Predict()")
)

type LinearRegression struct {
	fitted                 bool
	disturbance            float64
	regressionCoefficients []float64
}

func init() {
	mat64.Register(cblas.Blas{})
}

func NewLinearRegression() *LinearRegression {
	return &LinearRegression{fitted: false}
}

func (lr *LinearRegression) Fit(inst *base.Instances) error {
	if inst.Rows < inst.GetAttributeCount() {
		return NotEnoughDataError
	}

	// Split into two matrices, observed results (dependent variable y)
	// and the explanatory variables (X) - see http://en.wikipedia.org/wiki/Linear_regression
	observed := mat64.NewDense(inst.Rows, 1, nil)
	explVariables := mat64.NewDense(inst.Rows, inst.GetAttributeCount(), nil)

	for i := 0; i < inst.Rows; i++ {
		observed.Set(i, 0, inst.Get(i, inst.ClassIndex)) // Set observed data

		for j := 0; j < inst.GetAttributeCount(); j++ {
			if j == 0 {
				// Set intercepts to 1.0
				// Could / should be done better: http://www.theanalysisfactor.com/interpret-the-intercept/
				explVariables.Set(i, 0, 1.0)
			} else {
				explVariables.Set(i, j, inst.Get(i, j-1))
			}
		}
	}

	n := inst.GetAttributeCount()
	qr := mat64.QR(explVariables)
	q := qr.Q()
	reg := qr.R()

	var transposed, qty mat64.Dense
	transposed.TCopy(q)
	qty.Mul(&transposed, observed)

	regressionCoefficients := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		regressionCoefficients[i] = qty.At(i, 0)
		for j := i + 1; j < n; j++ {
			regressionCoefficients[i] -= regressionCoefficients[j] * reg.At(i, j)
		}
		regressionCoefficients[i] /= reg.At(i, i)
	}

	lr.disturbance = regressionCoefficients[0]
	lr.regressionCoefficients = regressionCoefficients[1:]
	lr.fitted = true

	return nil
}

func (lr *LinearRegression) Predict(X *base.Instances) (*base.Instances, error) {
	if !lr.fitted {
		return nil, NoTrainingDataError
	}

	ret := X.GeneratePredictionVector()
	for i := 0; i < X.Rows; i++ {
		var prediction float64 = lr.disturbance
		for j := 0; j < X.Cols; j++ {
			if j != X.ClassIndex {
				prediction += X.Get(i, j) * lr.regressionCoefficients[j]
			}
		}
		ret.Set(i, 0, prediction)
	}

	return ret, nil
}
