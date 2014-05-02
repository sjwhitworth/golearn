package pairwise

import (
	mat "github.com/skelterjohn/go.matrix"
)

type Metric interface {
	InnerProduct(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix)
	Distance(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix)
}

func CheckDimMatch(vectorX *mat.DenseMatrix, vectorY *mat.DenseMatrix) bool {
	if vectorX.Cols() != 1 ||
		vectorY.Cols() != 1 ||
		vectorX.Rows() != vectorY.Rows() {
		return false
	} else {
		return true
	}
}
