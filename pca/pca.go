//Implementation of Principal Component Analysis(PCA) with SVD
package pca

import (
	"gonum.org/v1/gonum/mat"
)

type PCA struct {
	Num_components int
	svd            *mat.SVD
}

// Number of components. 0 - by default, use number of features as number of components
func NewPCA(num_components int) *PCA {
	return &PCA{Num_components: num_components}
}

// Fit PCA model and transform data
// Need return is base.FixedDataGrid
func (pca *PCA) FitTransform(X *mat.Dense) *mat.Dense {
	return pca.Fit(X).Transform(X)
}

// Fit PCA model
func (pca *PCA) Fit(X *mat.Dense) *PCA {
	// Mean to input data
	M := mean(X)
	X = matrixSubVector(X, M)

	// Get SVD decomposition from data
	pca.svd = &mat.SVD{}
	ok := pca.svd.Factorize(X, mat.SVDThin)
	if !ok {
		panic("Unable to factorize")
	}
	if pca.Num_components < 0 {
		panic("Number of components can't be less than zero")
	}

	return pca
}

// Need return is base.FixedDataGrid
func (pca *PCA) Transform(X *mat.Dense) *mat.Dense {
	if pca.svd == nil {
		panic("You should to fit PCA model first")
	}

	num_samples, num_features := X.Dims()

	vTemp := new(mat.Dense)
	pca.svd.VTo(vTemp)
	//Compute to full data
	if pca.Num_components == 0 || pca.Num_components > num_features {
		return compute(X, vTemp)
	}

	X = compute(X, vTemp)
	result := mat.NewDense(num_samples, pca.Num_components, nil)
	result.Copy(X)
	return result
}

//Helpful private functions

//Compute mean of the columns of input matrix
func mean(matrix *mat.Dense) *mat.Dense {
	rows, cols := matrix.Dims()
	meanVector := make([]float64, cols)
	for i := 0; i < cols; i++ {
		sum := mat.Sum(matrix.ColView(i))
		meanVector[i] = sum / float64(rows)
	}
	return mat.NewDense(1, cols, meanVector)
}

// After computing of mean, compute: X(input matrix)  - X(mean vector)
func matrixSubVector(mat, vec *mat.Dense) *mat.Dense {
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
func compute(X, Y mat.Matrix) *mat.Dense {
	var ret mat.Dense
	ret.Mul(X, Y)
	return &ret
}
