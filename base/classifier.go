package base

import (
		mat "github.com/skelterjohn/go.matrix"
		)

type BaseClassifier struct {
	Data mat.DenseMatrix
	Name string
	Labels []string
}

type BaseRegressor struct {
	Data mat.DenseMatrix
	Name string
	Labels []float64
}