package linear_models

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

type LinearSVC struct {
	param *Parameter
	model *Model
}

func NewLinearSVC(loss, penalty string, dual bool, C float64, eps float64) (*LinearSVC, error) {
	solver_type := 0
	if penalty == "l2" {
		if loss == "l1" {
			if dual {
				solver_type = L2R_L1LOSS_SVC_DUAL
			}
		} else {
			if dual {
				solver_type = L2R_L2LOSS_SVC_DUAL
			} else {
				solver_type = L2R_L2LOSS_SVC
			}
		}
	} else if penalty == "l1" {
		if loss == "l2" {
			if !dual {
				solver_type = L1R_L2LOSS_SVC
			}
		}
	}
	if solver_type == 0 {
		panic("Parameter combination")
	}

	lr := LinearSVC{}
	lr.param = NewParameter(solver_type, C, eps)
	lr.model = nil
	return &lr, nil
}

func (lr *LinearSVC) Fit(X base.FixedDataGrid) error {
	problemVec := convertInstancesToProblemVec(X)
	labelVec := convertInstancesToLabelVec(X)
	prob := NewProblem(problemVec, labelVec, 0)
	lr.model = Train(prob, lr.param)
	return nil
}

func (lr *LinearSVC) Predict(X base.FixedDataGrid) (base.FixedDataGrid, error) {

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

func (lr *LinearSVC) String() string {
	return "LogisticSVC"
}
