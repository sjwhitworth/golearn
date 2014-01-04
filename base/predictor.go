package golearn

import (
		mat "github.com/skelterjohn/go.matrix"
		)

type BasePredictor struct {
	Data mat.DenseMatrix
	Name string
	Labels []string
}