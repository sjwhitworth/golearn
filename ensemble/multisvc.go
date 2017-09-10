package ensemble

import (
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/linear_models"
	"github.com/sjwhitworth/golearn/meta"

	"fmt"
)

// MultiLinearSVC implements a multi-class Support Vector Classifier using a one-vs-all
// voting scheme. Only one CategoricalAttribute class is supported.
type MultiLinearSVC struct {
	m          *meta.OneVsAllModel
	parameters *linear_models.LinearSVCParams
	weights    map[string]float64
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

	// Return me...
	ret := &MultiLinearSVC{
		parameters: params,
		weights:    weights,
	}

	ret.initializeOneVsAllModel()
	return ret
}

func (m *MultiLinearSVC) initializeOneVsAllModel() {
	// Classifier creation function
	classifierFunc := func(cls string) base.Classifier {
		var weightVec []float64
		newParams := m.parameters.Copy()
		if m.weights != nil {
			weightVec = make([]float64, 2)
			for i := range m.weights {
				if i != cls {
					weightVec[0] += m.weights[i]
				} else {
					weightVec[1] = m.weights[i]
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
	m.m = meta.NewOneVsAllModel(classifierFunc)
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

func (m *MultiLinearSVC) GetClassifierMetadata() base.ClassifierMetadataV1 {
	return base.ClassifierMetadataV1{
		FormatVersion:      1,
		ClassifierName:     "MultiLinearSVC",
		ClassifierVersion:  "1",
		ClassifierMetadata: nil,
	}
}

func (m *MultiLinearSVC) Save(filePath string) error {
	metadata := m.GetClassifierMetadata()
	serializer, err := base.CreateSerializedClassifierStub(filePath, metadata)
	if err != nil {
		return err
	}
	err = m.SaveWithPrefix(serializer, "")
	if err != nil {
		return fmt.Errorf("Unable to Save(): %v", err)
	}
	serializer.Close()
	return err
}

func (m *MultiLinearSVC) SaveWithPrefix(serializer *base.ClassifierSerializer, prefix string) error {

	p := func(fName string) string {
		return fmt.Sprintf("%s/%s", prefix, fName)
	}

	// Write out the linear parameters
	err := serializer.WriteJSONForKey(p("params"), m.parameters)
	if err != nil {
		return fmt.Errorf("Unable to marshal parameters: %v", err)
	}
	// Write out the weights
	err = serializer.WriteJSONForKey(p("weights"), m.weights)
	if err != nil {
		return fmt.Errorf("Unable to write weights: %v", err)
	}

	// Serialize the model
	err = m.m.SaveWithPrefix(serializer, p("one-vs-all"))
	return err
}

func (m *MultiLinearSVC) GetMetadata() base.ClassifierMetadataV1 {
	return base.ClassifierMetadataV1{
		FormatVersion:      1,
		ClassifierName:     "MultiLinearSVC",
		ClassifierVersion:  "1.0",
		ClassifierMetadata: nil,
	}
}

func (m *MultiLinearSVC) Load(filePath string) error {
	reader, err := base.ReadSerializedClassifierStub(filePath)
	if err != nil {
		return err
	}

	err = m.LoadWithPrefix(reader, "")
	if err != nil {
		return err
	}

	return nil
}

func (m *MultiLinearSVC) LoadWithPrefix(reader *base.ClassifierDeserializer, prefix string) error {
	p := func(fName string) string {
		return fmt.Sprintf("%s/%s", prefix, fName)
	}
	err := reader.GetJSONForKey(p("params"), &m.parameters)
	if err != nil {
		return fmt.Errorf("Can't load parameters: %v", err)
	}

	err = reader.GetJSONForKey(p("weights"), &m.weights)
	if err != nil {
		return fmt.Errorf("Can't load parameters: %v", err)
	}

	m.initializeOneVsAllModel()

	// Load the model
	err = m.m.LoadWithPrefix(reader, p("one-vs-all"))
	if err != nil {
		return err
	}

	return nil
}
