package linear_models

import (
	"errors"

	"github.com/amclay/golearn/base"

	"fmt"

	_ "github.com/gonum/blas"
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
	attrs                  []base.Attribute
	cls                    base.Attribute
}

func NewLinearRegression() *LinearRegression {
	return &LinearRegression{fitted: false}
}

func (lr *LinearRegression) Fit(inst base.FixedDataGrid) error {

	// Retrieve row size
	_, rows := inst.Size()

	// Validate class Attribute count
	classAttrs := inst.AllClassAttributes()
	if len(classAttrs) != 1 {
		return fmt.Errorf("Only 1 class variable is permitted")
	}
	classAttrSpecs := base.ResolveAttributes(inst, classAttrs)

	// Retrieve relevant Attributes
	allAttrs := base.NonClassAttributes(inst)
	attrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if _, ok := a.(*base.FloatAttribute); ok {
			attrs = append(attrs, a)
		}
	}

	cols := len(attrs) + 1

	if rows < cols {
		return NotEnoughDataError
	}

	// Retrieve relevant Attribute specifications
	attrSpecs := base.ResolveAttributes(inst, attrs)

	// Split into two matrices, observed results (dependent variable y)
	// and the explanatory variables (X) - see http://en.wikipedia.org/wiki/Linear_regression
	observed := mat64.NewDense(rows, 1, nil)
	explVariables := mat64.NewDense(rows, cols, nil)

	// Build the observed matrix
	inst.MapOverRows(classAttrSpecs, func(row [][]byte, i int) (bool, error) {
		val := base.UnpackBytesToFloat(row[0])
		observed.Set(i, 0, val)
		return true, nil
	})

	// Build the explainatory variables
	inst.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		// Set intercepts to 1.0
		explVariables.Set(i, 0, 1.0)
		for j, r := range row {
			explVariables.Set(i, j+1, base.UnpackBytesToFloat(r))
		}
		return true, nil
	})

	n := cols
	qr := new(mat64.QR)
	qr.Factorize(explVariables)
	var q, reg mat64.Dense
	q.QFromQR(qr)
	reg.RFromQR(qr)

	var transposed, qty mat64.Dense
	transposed.Clone(q.T())
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
	lr.attrs = attrs
	lr.cls = classAttrs[0]
	return nil
}

func (lr *LinearRegression) Predict(X base.FixedDataGrid) (base.FixedDataGrid, error) {
	if !lr.fitted {
		return nil, NoTrainingDataError
	}

	ret := base.GeneratePredictionVector(X)
	attrSpecs := base.ResolveAttributes(X, lr.attrs)
	clsSpec, err := ret.GetAttribute(lr.cls)
	if err != nil {
		return nil, err
	}

	X.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		var prediction float64 = lr.disturbance
		for j, r := range row {
			prediction += base.UnpackBytesToFloat(r) * lr.regressionCoefficients[j]
		}

		ret.Set(clsSpec, i, base.PackFloatToBytes(prediction))
		return true, nil
	})

	return ret, nil
}
