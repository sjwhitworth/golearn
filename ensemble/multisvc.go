package ensemble

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/linear_models"
	"github.com/sjwhitworth/golearn/meta"
)

// MultiLinearSVC implements a multi-class Support Vector Classifier using a one-vs-all
// voting scheme. Only one CategoricalAttribute class is supported.
type MultiLinearSVC struct {
	m *meta.OneVsAllModel
}

// NewMultiLinearSVC creates a new MultiLinearSVC using the OneVsAllModel.
// The loss and penalty arguments can be "l1" or "l2". Typical values are
// "l1" for the loss and "l2" for the penalty. The dual parameter controls
// whether the system solves the dual or primal SVM form, true should be used
// in most cases. C is the penalty term, normally 1.0. eps is the convergence
// term, typically 1e-4.
func NewMultiLinearSVC(loss, penalty string, dual bool, C float64, eps float64, weights map[string]float64) *MultiLinearSVC {
	// Set up the training parameters
	params := &linear_models.LinearSVCParams{0, nil, C, eps, false, dual}
	err := params.SetKindFromStrings(loss, penalty)
	if err != nil {
		panic(err)
	}

	// Classifier creation function
	classifierFunc := func(cls string) base.Classifier {
		var weightVec []float64
		newParams := params.Copy()
		if weights != nil {
			weightVec = make([]float64, 2)
			for i := range weights {
				if i != cls {
					weightVec[0] += weights[i]
				} else {
					weightVec[1] = weights[i]
				}
			}
		}
		newParams.ClassWeights = weightVec

		ret, err := linear_models.NewLinearSVCFromParams(newParams)
		if err != nil {
			panic(err)
		}
		return ret
	}

	// Return me...
	return &MultiLinearSVC{
		meta.NewOneVsAllModel(classifierFunc),
	}
}

// Fit builds the MultiLinearSVC by building n (where n is the number of values
// the singular CategoricalAttribute can take) seperate one-vs-rest models.
func (m *MultiLinearSVC) Fit(instances base.FixedDataGrid) error {
	m.m.Fit(instances)
	return nil
}

// Predict issues predictions from the MultiLinearSVC. Each underlying LinearSVC is
// used to predict whether an instance takes on a class or some other class, and the
// model which definitively reports a given class is the one chosen. The result is
// undefined if all underlying models predict that the instance originates from some
// other class.
func (m *MultiLinearSVC) Predict(from base.FixedDataGrid) (base.FixedDataGrid, error) {
	return m.m.Predict(from)
}
