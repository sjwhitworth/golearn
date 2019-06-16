package clustering

import (
	"errors"
	"github.com/sjwhitworth/golearn/base"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distmv"
	"math"
	"math/rand"
)

var (
	NoTrainingDataError         = errors.New("You need to Fit() before you can Predict()")
	InsufficientComponentsError = errors.New("Estimation requires at least one component")
	InsufficientDataError       = errors.New("Estimation requires n_obs >= n_comps")
)

type ExpectationMaximization struct {
	n_comps int
	eps     float64
	Params  Params
	fitted  bool
	attrs   []base.Attribute
}

type Params struct {
	Means *mat.Dense
	Covs  []*mat.SymDense
}

// Number of Gaussians to fit in the mixture
func NewExpectationMaximization(n_comps int) (*ExpectationMaximization, error) {
	if n_comps < 1 {
		return nil, InsufficientComponentsError
	}

	return &ExpectationMaximization{n_comps: n_comps, eps: 0.001}, nil
}

// Fit method - generates the component parameters (means and covariance matrices)
func (em *ExpectationMaximization) Fit(inst base.FixedDataGrid) error {
	// Numeric Attrs
	attrs := base.NonClassAttributes(inst)
	attrSpecs := base.ResolveAttributes(inst, attrs)

	_, n_obs := inst.Size()
	n_feats := len(attrs)

	if n_obs < em.n_comps {
		return InsufficientDataError
	}

	// Build the input matrix
	X := mat.NewDense(n_obs, n_feats, nil)
	inst.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		for j, r := range row {
			X.Set(i, j, base.UnpackBytesToFloat(r))
		}
		return true, nil
	})

	// Initialize the parameter distance
	dist := math.Inf(1)

	// Initialize the parameters
	var p Params
	p.Means = initMeans(X, em.n_comps, n_feats)
	p.Covs = initCovariance(em.n_comps, n_feats)

	// Iterate until convergence
	for {
		if dist < em.eps {
			break
		}
		y_new := expectation(X, p, em.n_comps)
		p_new := maximization(X, y_new, p, em.n_comps)
		dist = distance(p, p_new)
		p = p_new
	}

	em.fitted = true
	em.attrs = attrs
	em.Params = p
	return nil
}

// Predict method - returns a ClusterMap of components and row ids
func (em *ExpectationMaximization) Predict(inst base.FixedDataGrid) (ClusterMap, error) {
	if !em.fitted {
		return nil, NoTrainingDataError
	}

	_, n_obs := inst.Size()
	n_feats := len(em.attrs)

	// Numeric attrs
	attrSpecs := base.ResolveAttributes(inst, em.attrs)

	// Build the input matrix
	X := mat.NewDense(n_obs, n_feats, nil)
	inst.MapOverRows(attrSpecs, func(row [][]byte, i int) (bool, error) {
		for j, r := range row {
			X.Set(i, j, base.UnpackBytesToFloat(r))
		}
		return true, nil
	})

	// Vector of predictions
	preds := estimateLogProb(X, em.Params, em.n_comps)

	clusterMap := make(map[int][]int)
	for ix, pred := range vecToInts(preds) {
		clusterMap[pred] = append(clusterMap[pred], ix)
	}

	return ClusterMap(clusterMap), nil
}

// EM-specific functions
// Expectation step
func expectation(X *mat.Dense, p Params, n_comps int) mat.Vector {
	y_new := estimateLogProb(X, p, n_comps)
	return y_new
}

// Maximization step
func maximization(X *mat.Dense, y mat.Vector, p Params, n_comps int) Params {
	_, n_feats := X.Dims()

	// Initialize the new parameters
	var p_new Params
	p_new.Means = mat.NewDense(n_comps, n_feats, nil)
	p_new.Covs = make([]*mat.SymDense, n_comps)

	// Update the parameters
	for k := 0; k < n_comps; k++ {
		X_yk := where(X, y, k)
		n_obs, _ := X_yk.Dims()
		covs_k_reg := mat.NewSymDense(n_feats, nil)

		if n_obs <= 1 {
			p_new.Means.SetRow(k, p.Means.RawRowView(k))
			covs_k_reg = p.Covs[k]
		} else if n_obs < n_feats {
			p_new.Means.SetRow(k, means(X_yk))
			covs_k_reg = shrunkCovariance(X_yk)
		} else {
			p_new.Means.SetRow(k, means(X_yk))
			stat.CovarianceMatrix(covs_k_reg, X_yk, nil)
		}

		p_new.Covs[k] = covs_k_reg
	}

	return p_new
}

// Creates mat.Vector of most likely component for each observation
func estimateLogProb(X *mat.Dense, p Params, n_comps int) mat.Vector {
	n_obs, n_feats := X.Dims()

	// Cache the component Gaussians
	var N = make([]*distmv.Normal, n_comps)
	for k := 0; k < n_comps; k++ {
		dst := make([]float64, n_feats)
		means := mat.Row(dst, k, p.Means)
		dist, ok := distmv.NewNormal(means, p.Covs[k], nil)
		if !ok {
			panic("Cannot create Normal!")
		}
		N[k] = dist
	}

	// Compute the component probabilities
	y_new := mat.NewVecDense(n_obs, nil)
	for i := 0; i < n_obs; i++ {
		max_ix := 0
		max_pr := math.Inf(-1)
		x := X.RawRowView(i)
		for k := 0; k < n_comps; k++ {
			pr := N[k].LogProb(x)
			if pr > max_pr {
				max_ix = k
				max_pr = pr
			}
		}
		y_new.SetVec(i, float64(max_ix))
	}

	return y_new
}

// Returns a symmetric matrix with variance on the diagonal
func shrunkCovariance(X *mat.Dense) *mat.SymDense {
	n_obs, n_feats := X.Dims()
	size := int(math.Pow(float64(n_feats), 2))
	covs := mat.NewSymDense(n_feats, make([]float64, size, size))
	for j := 0; j < n_feats; j++ {
		// compute the variance for the jth feature
		var points []float64
		for i := 0; i < n_obs; i++ {
			points = append(points, X.At(i, j))
		}
		variance := stat.Variance(points, nil)
		// set the jth diagonal entry to the variance
		covs.SetSym(j, j, variance)
	}
	return covs
}

// Creates an n_comps x n_feats array of means
func initMeans(X *mat.Dense, n_comps, n_feats int) *mat.Dense {
	var results []float64
	for k := 0; k < n_comps; k++ {
		for j := 0; j < n_feats; j++ {
			v := X.ColView(j)
			min := vectorMin(v)
			max := vectorMax(v)
			r := min + rand.Float64()*(max-min)
			results = append(results, r)
		}
	}
	means := mat.NewDense(n_comps, n_feats, results)
	return means
}

// Creates a n_comps array of n_feats x n_feats mat.Symmetrics
func initCovariance(n_comps, n_feats int) []*mat.SymDense {
	var result []*mat.SymDense
	floats := identity(n_feats)
	for k := 0; k < n_comps; k++ {
		matrix := mat.NewSymDense(n_feats, floats)
		result = append(result, matrix)
	}
	return result
}

// Compues the euclidian distance between two parameters
func distance(p Params, p_new Params) float64 {
	dist := 0.0
	n_obs, n_feats := p.Means.Dims()
	for i := 0; i < n_obs; i++ {
		means_i := p.Means.RawRowView(i)
		means_new_i := p_new.Means.RawRowView(i)
		for j := 0; j < n_feats; j++ {
			dist += math.Pow((means_i[j] - means_new_i[j]), 2)
		}
	}
	return math.Sqrt(dist)
}

// Helper functions
// Finds the min value of a mat.Vector
func vectorMin(v mat.Vector) float64 {
	n_obs, _ := v.Dims()
	min := v.At(0, 0)
	for i := 0; i < n_obs; i++ {
		if v.At(i, 0) < min {
			min = v.At(i, 0)
		}
	}
	return min
}

// Find the max value of a mat.Vector
func vectorMax(v mat.Vector) float64 {
	n_obs, _ := v.Dims()
	max := v.At(0, 0)
	for i := 0; i < n_obs; i++ {
		if v.At(i, 0) > max {
			max = v.At(i, 0)
		}
	}
	return max
}

// Converts a mat.Vector to an array of ints
func vecToInts(v mat.Vector) []int {
	n_obs, _ := v.Dims()
	var ints = make([]int, n_obs)
	for i := 0; i < n_obs; i++ {
		ints[i] = int(v.At(i, 0))
	}
	return ints
}

// Computes column Means of a mat.Dense
func means(X *mat.Dense) []float64 {
	n_obs, n_feats := X.Dims()
	var result []float64
	for j := 0; j < n_feats; j++ {
		sum_j := 0.0
		for i := 0; i < n_obs; i++ {
			sum_j = sum_j + X.At(i, j)
		}
		mean := (sum_j / float64(n_obs))
		result = append(result, mean)
	}
	return result
}

// Subest a mat.Dense with rows matching a target value
func where(X *mat.Dense, y mat.Vector, target int) *mat.Dense {
	n_obs, n_feats := X.Dims()
	var result []float64
	rows := 0
	for i := 0; i < n_obs; i++ {
		if int(y.At(i, 0)) == target {
			for j := 0; j < n_feats; j++ {
				result = append(result, X.At(i, j))
			}
			rows++
		}
	}
	X_i := mat.NewDense(rows, n_feats, result)
	return X_i
}

// Returns values for a square array with ones on the main diagonal
func identity(N int) []float64 {
	var results []float64
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if j == i {
				results = append(results, 1)
			} else {
				results = append(results, 0)
			}
		}
	}
	return results
}

// Generates an array of values for symmetric matrix
func symVals(M int, v float64) []float64 {
	var results []float64
	for i := 0; i < M*M; i++ {
		results = append(results, v)
	}
	return results
}
