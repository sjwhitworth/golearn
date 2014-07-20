package linear_models

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
)

type LogisticRegression struct {
	param *Parameter
	model *Model
}

func NewLogisticRegression(penalty string, C float64, eps float64) *LogisticRegression {
	solver_type := 0
	if penalty == "l2" {
		solver_type = L2R_LR
	} else if penalty == "l1" {
		solver_type = L1R_LR
	} else {
		fmt.Println("Invalid penalty")
		return nil
	}

	lr := LogisticRegression{}
	lr.param = NewParameter(solver_type, C, eps)
	lr.model = nil
	return &lr
}

func convertInstancesToProblemVec(X *base.Instances) [][]float64 {
	problemVec := make([][]float64, X.Rows)
	for i := 0; i < X.Rows; i++ {
		problemVecCounter := 0
		problemVec[i] = make([]float64, X.Cols-1)
		for j := 0; j < X.Cols; j++ {
			if j == X.ClassIndex {
				continue
			}
			problemVec[i][problemVecCounter] = X.Get(i, j)
			problemVecCounter++
		}
	}
	base.Logger.Println(problemVec, X)
	return problemVec
}

func convertInstancesToLabelVec(X *base.Instances) []float64 {
	labelVec := make([]float64, X.Rows)
	for i := 0; i < X.Rows; i++ {
		labelVec[i] = X.Get(i, X.ClassIndex)
	}
	return labelVec
}

func (lr *LogisticRegression) Fit(X *base.Instances) {
	problemVec := convertInstancesToProblemVec(X)
	labelVec := convertInstancesToLabelVec(X)
	prob := NewProblem(problemVec, labelVec, 0)
	lr.model = Train(prob, lr.param)
}

func (lr *LogisticRegression) Predict(X *base.Instances) *base.Instances {
	ret := X.GeneratePredictionVector()
	row := make([]float64, X.Cols-1)
	for i := 0; i < X.Rows; i++ {
		rowCounter := 0
		for j := 0; j < X.Cols; j++ {
			if j != X.ClassIndex {
				row[rowCounter] = X.Get(i, j)
				rowCounter++
			}
		}
		base.Logger.Println(Predict(lr.model, row), row)
		ret.Set(i, 0, Predict(lr.model, row))
	}
	return ret
}
