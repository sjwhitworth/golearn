//Implementation of Principal Component Analysis(PCA) with SVD
package pca

import
(
	"github.com/gonum/matrix/mat64"
	"math"
)

type PCA struct {
	Num_components int
}


// Number of components. 0 - by default, use number of features as number of components
func NewPCA(num_components int)*PCA {
	return &PCA {Num_components: num_components}
}


//Need return is base.FixedDataGrid
func (pca*PCA) Transform(X *mat64.Dense) (*mat64.Dense){
	//Prepare before PCA

	num_samples, num_features := X.Dims()
	//Mean to input data
	M := mean(X)
	X = matrixSubVector(X, M)

	//Get SVD decomposition from data
	svd := mat64.SVD(mat64.DenseCopyOf(X),math.Pow(2, -52.0),math.Pow(2, -966.0),false, true)
	if pca.Num_components < 0 {
		panic("Number of components can't be less than zero")
	}

	//Compute to full data
	if pca.Num_components == 0 || pca.Num_components > num_features{
		return compute(X,svd.V)
	}

	X = compute(X,svd.V)
	result := mat64.NewDense(num_samples, pca.Num_components, nil)
	result.Copy(X.View(0,0,num_samples,pca.Num_components))
	return result
}


//Helpful private functions

//Compute mean of the columns of input matrix
func mean(matrix *mat64.Dense)*mat64.Dense{
	rows, cols := matrix.Dims()
	meanVector := make([]float64, cols)
	for i := 0; i < cols; i++ {
		sum := 0.0
		for _, it := range matrix.Col(nil,i){
			sum += it
		}
		meanVector[i] = sum/float64(rows)
	}
	return mat64.NewDense(1, cols, meanVector)
}

// After computing of mean, compute: X(input matrix)  - X(mean vector)
func matrixSubVector(mat, vec *mat64.Dense)*mat64.Dense {
	rowsm, colsm := mat.Dims()
	_, colsv := vec.Dims()
	if colsv != colsm {
		panic("Error in dimension")
	}
	for i := 0; i < rowsm; i++ {
		for j := 0; j < colsm; j++ {
			mat.Set(i,j, (mat.At(i,j) - vec.At(0,j)))
		}
	}
	return mat
}

//Multiplication of X(input data) and V(from SVD)
func compute (X, Y *mat64.Dense) *mat64.Dense{
	X.Mul(X,Y)
	return X
}