package linear_models

import (
	"errors"
	"fmt"

	"github.com/amclay/golearn/base"
)

type LogisticRegression struct {
	param *Parameter
	model *Model
}

func NewLogisticRegression(penalty string, C float64, eps float64) (*LogisticRegression, error) {
	solver_type := 0
	if penalty == "l2" {
		solver_type = L2R_LR
	} else if penalty == "l1" {
		solver_type = L1R_LR
	} else {
		return nil, errors.New(fmt.Sprintf("Invalid penalty '%s'", penalty))
	}

	lr := LogisticRegression{}
	lr.param = NewParameter(solver_type, C, eps)
	lr.model = nil
	return &lr, nil
}

func (lr *LogisticRegression) Fit(X base.FixedDataGrid) error {
	problemVec := convertInstancesToProblemVec(X)
	labelVec := convertInstancesToLabelVec(X)
	prob := NewProblem(problemVec, labelVec, 0)
	lr.model = Train(prob, lr.param)
	return nil
}

func (lr *LogisticRegression) Predict(X base.FixedDataGrid) (base.FixedDataGrid, error) {

	// Only support 1 class Attribute
	classAttrs := X.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic(fmt.Sprintf("%d Wrong number of classes", len(classAttrs)))
	}
	// Generate return structure
	ret := base.GeneratePredictionVector(X)
	classAttrSpecs := base.ResolveAttributes(ret, classAttrs)
	// Retrieve numeric non-class Attributes
	numericAttrs := base.NonClassFloatAttributes(X)
	numericAttrSpecs := base.ResolveAttributes(X, numericAttrs)

	// Allocate row storage
	row := make([]float64, len(numericAttrSpecs))
	X.MapOverRows(numericAttrSpecs, func(rowBytes [][]byte, rowNo int) (bool, error) {
		for i, r := range rowBytes {
			row[i] = base.UnpackBytesToFloat(r)
		}
		val := Predict(lr.model, row)
		vals := base.PackFloatToBytes(val)
		ret.Set(classAttrSpecs[0], rowNo, vals)
		return true, nil
	})

	return ret, nil
}

func (lr *LogisticRegression) String() string {
	return "LogisticRegression"
}
