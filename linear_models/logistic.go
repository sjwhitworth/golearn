package linear_models

import "fmt"

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

func (lr *LogisticRegression) Fit(X [][]float64, y []float64) {
	prob := NewProblem(X, y, 0)
	lr.model = Train(prob, lr.param)
}

func (lr *LogisticRegression) Predict(X [][]float64) []float64 {
	n_samples := len(X)
	y := make([]float64, n_samples)
	for i, x := range X {
		y[i] = Predict(lr.model, x)
	}
	return y
}
