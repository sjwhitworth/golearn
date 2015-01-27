package perceptron

import (
	base "github.com/sjwhitworth/golearn/base"
	"math"
)

const MaxEpochs = 10

type AveragePerceptron struct {
	TrainingData base.FixedDataGrid
	weights      []float64
	edges        []float64
	bias         float64
	threshold    float64
	learningRate float64
	trainError   float64
	trained      bool
	count        float64
}

type instance struct {
	class    string
	features []float64
}

type instances []instance

func (p *AveragePerceptron) updateWeights(features []float64, correction float64) {

	for i, _ := range p.weights {
		fv := &features[i]
		if fv != nil {
			update := p.learningRate * correction * *fv
			p.weights[i] = update
			p.edges[i]++
		}
	}

	p.average()
}

func (p *AveragePerceptron) average() {

	for i, fcount := range p.edges {
		wv := &p.weights[i]
		if wv != nil {
			p.weights[i] = (p.count**wv + fcount) / (fcount + 1)
		}
	}
	p.count++
}

func (p *AveragePerceptron) score(datum instance) float64 {
	score := 0.0

	for i, wv := range p.weights {
		score += datum.features[i] * wv
	}

	if score >= p.threshold {
		return 1.0
	}
	return -1.0

}

func (p *AveragePerceptron) Fit(trainingData base.FixedDataGrid) {

	epochs := 0
	p.trainError = 0.1
	learning := true

	data := processData(trainingData)
	for learning {
		for _, datum := range data {
			response := p.score(datum)
			expected := 0.0
			correction := expected - response

			if expected != response {
				p.updateWeights(datum.features, correction)
				p.trainError += math.Abs(correction)
			}
		}

		epochs++

		if epochs >= MaxEpochs {
			learning = false
		}
	}

	p.average()
	p.trained = true
	p.TrainingData = trainingData
}

// param base.IFixedDataGrid
// return base.IFixedDataGrid
func (p *AveragePerceptron) Predict(what base.FixedDataGrid) base.FixedDataGrid {

	if !p.trained {
		panic("Cannot call Predict on an untrained AveragePerceptron")
	}

	data := processData(what)

	allAttrs := base.CheckCompatible(what, p.TrainingData)
	if allAttrs == nil {
		// Don't have the same Attributes
		return nil
	}

	// Remove the Attributes which aren't numeric
	allNumericAttrs := make([]base.Attribute, 0)
	for _, a := range allAttrs {
		if fAttr, ok := a.(*base.FloatAttribute); ok {
			allNumericAttrs = append(allNumericAttrs, fAttr)
		}
	}

	ret := base.GeneratePredictionVector(what)
	classAttr := ret.AllClassAttributes()[0]
	classSpec, err := ret.GetAttribute(classAttr)
	if err != nil {
		panic(err)
	}

	for i, datum := range data {
		result := p.score(datum)
		if result > 0.0 {
			ret.Set(classSpec, i, base.PackU64ToBytes(1))
		} else {
			ret.Set(classSpec, 1, []byte{0, 0, 0, 0, 0, 0, 0, 0})
		}
	}

	return ret
}

func processData(x base.FixedDataGrid) instances {
	_, rows := x.Size()

	result := make(instances, rows)

	// Retrieve numeric non-class Attributes
	numericAttrs := base.NonClassFloatAttributes(x)
	numericAttrSpecs := base.ResolveAttributes(x, numericAttrs)

	// Retrieve class Attributes
	classAttrs := x.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic("Only one classAttribute supported!")
	}

	// Check that the class Attribute is categorical
	// (with two values) or binary
	classAttr := classAttrs[0]
	if attr, ok := classAttr.(*base.CategoricalAttribute); ok {
		if len(attr.GetValues()) != 2 {
			panic("To many values for Attribute!")
		}
	} else if _, ok := classAttr.(*base.BinaryAttribute); ok {
	} else {
		panic("Wrong class Attribute type!")
	}

	// Convert each row
	x.MapOverRows(numericAttrSpecs, func(row [][]byte, rowNo int) (bool, error) {
		// Allocate a new row
		probRow := make([]float64, len(numericAttrSpecs))

		// Read out the row
		for i, _ := range numericAttrSpecs {
			probRow[i] = base.UnpackBytesToFloat(row[i])
		}

		// Get the class for the values
		class := base.GetClass(x, rowNo)
		instance := instance{class, probRow}
		result[rowNo] = instance
		return true, nil
	})
	return result
}

func NewAveragePerceptron(features int, learningRate float64, startingThreshold float64, trainError float64) *AveragePerceptron {

	weights := make([]float64, features)
	edges := make([]float64, features)

	p := AveragePerceptron{nil, weights, edges, startingThreshold, learningRate, trainError, 0.0, false, 0}

	return &p
}
