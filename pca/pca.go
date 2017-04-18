//Implementation of Principal Component Analysis(PCA) with SVD
package pca

import (
	"github.com/gonum/matrix"
	"github.com/gonum/matrix/mat64"
)

type PCA struct {
	Num_components int
}

// Number of components. 0 - by default, use number of features as number of components
func NewPCA(num_components int) *PCA {
	return &PCA{Num_components: num_components}
}

//Need return is base.FixedDataGrid
func (pca *PCA) Transform(X *mat64.Dense) *mat64.Dense {
	//Prepare before PCA

	num_samples, num_features := X.Dims()
	//Mean to input data
	M := mean(X)
	X = matrixSubVector(X, M)

	//Get SVD decomposition from data
	var svd mat64.SVD
	ok := svd.Factorize(X, matrix.SVDThin)
	if !ok {
		panic("Unable to factorize")
	}
	if pca.Num_components < 0 {
		panic("Number of components can't be less than zero")
	}

	vTemp := new(mat64.Dense)
	vTemp.VFromSVD(&svd)
	//Compute to full data
	if pca.Num_components == 0 || pca.Num_components > num_features {
		return compute(X, vTemp)
	}

	X = compute(X, vTemp)
	result := mat64.NewDense(num_samples, pca.Num_components, nil)
	result.Copy(X.View(0, 0, num_samples, pca.Num_components))
	return result
}

//Helpful private functions

//Compute mean of the columns of input matrix
func mean(matrix *mat64.Dense) *mat64.Dense {
	rows, cols := matrix.Dims()
	meanVector := make([]float64, cols)
	for i := 0; i < cols; i++ {
		sum := mat64.Sum(matrix.ColView(i))
		meanVector[i] = sum / float64(rows)
	}
	return mat64.NewDense(1, cols, meanVector)
}

// After computing of mean, compute: X(input matrix)  - X(mean vector)
func matrixSubVector(mat, vec *mat64.Dense) *mat64.Dense {
	rowsm, colsm := mat.Dims()
	_, colsv := vec.Dims()
	if colsv != colsm {
		panic("Error in dimension")
	}
	for i := 0; i < rowsm; i++ {
		for j := 0; j < colsm; j++ {
			mat.Set(i, j, (mat.At(i, j) - vec.At(0, j)))
		}
	}
	return mat
}

//Multiplication of X(input data) and V(from SVD)
func compute(X, Y mat64.Matrix) *mat64.Dense {
	var ret mat64.Dense
	ret.Mul(X, Y)
	return &ret
}
