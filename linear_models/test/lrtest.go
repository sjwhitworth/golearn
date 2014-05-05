package main

import (
	"fmt"
	lm "golearn/linear_models"
)

type Obj struct {
	a int
}

func main() {
	X := [][]float64{
		{0, 0, 0, 1},
		{0, 0, 1, 0},
		{0, 1, 0, 0},
		{1, 0, 0, 0},
	}
	y := []float64{-1, -1, 1, 1}
	lr := lm.NewLogisticRegression("l2", 1.0, 1e-6)

	fmt.Println("Training")
	lr.Fit(X, y)

	pred_X := [][]float64{
		{1, 1, 0, 0},
		{0, 0, 1, 1},
	}
	fmt.Println("Testing")
	pred_y := lr.Predict(pred_X)

	fmt.Println("Result")
	fmt.Println(pred_y)
}
