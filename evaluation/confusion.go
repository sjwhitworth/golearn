package evaluation

import (
	"bytes"
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

// ConfusionMatrix is a nested map of actual and predicted class counts
type ConfusionMatrix map[string]map[string]int

// GetConfusionMatrix builds a ConfusionMatrix from a set of reference (`ref')
// and generate (`gen') Instances.
func GetConfusionMatrix(ref base.FixedDataGrid, gen base.FixedDataGrid) map[string]map[string]int {

	_, refRows := ref.Size()
	_, genRows := gen.Size()

	if refRows != genRows {
		panic("Row counts should match")
	}

	ret := make(map[string]map[string]int)

	for i := 0; i < int(refRows); i++ {
		referenceClass := base.GetClass(ref, i)
		predictedClass := base.GetClass(gen, i)
		if _, ok := ret[referenceClass]; ok {
			ret[referenceClass][predictedClass] += 1
		} else {
			ret[referenceClass] = make(map[string]int)
			ret[referenceClass][predictedClass] = 1
		}
	}
	return ret
}

// GetTruePositives returns the number of times an entry is
// predicted successfully in a given ConfusionMatrix.
func GetTruePositives(class string, c ConfusionMatrix) float64 {
	return float64(c[class][class])
}

// GetFalsePositives returns the number of times an entry is
// incorrectly predicted as having a given class.
func GetFalsePositives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c {
		if k == class {
			continue
		}
		ret += float64(c[k][class])
	}
	return ret
}

// GetFalseNegatives returns the number of times an entry is
// incorrectly predicted as something other than the given class.
func GetFalseNegatives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c[class] {
		if k == class {
			continue
		}
		ret += float64(c[class][k])
	}
	return ret
}

// GetTrueNegatives returns the number of times an entry is
// correctly predicted as something other than the given class.
func GetTrueNegatives(class string, c ConfusionMatrix) float64 {
	ret := 0.0
	for k := range c {
		if k == class {
			continue
		}
		for l := range c[k] {
			if l == class {
				continue
			}
			ret += float64(c[k][l])
		}
	}
	return ret
}

// GetPrecision returns the fraction of of the total predictions
// for a given class which were correct.
func GetPrecision(class string, c ConfusionMatrix) float64 {
	// Fraction of retrieved instances that are relevant
	truePositives := GetTruePositives(class, c)
	falsePositives := GetFalsePositives(class, c)
	return truePositives / (truePositives + falsePositives)
}

// GetRecall returns the fraction of the total occurrences of a
// given class which were predicted.
func GetRecall(class string, c ConfusionMatrix) float64 {
	// Fraction of relevant instances that are retrieved
	truePositives := GetTruePositives(class, c)
	falseNegatives := GetFalseNegatives(class, c)
	return truePositives / (truePositives + falseNegatives)
}

// GetF1Score computes the harmonic mean of precision and recall
// (equivalently called F-measure)
func GetF1Score(class string, c ConfusionMatrix) float64 {
	precision := GetPrecision(class, c)
	recall := GetRecall(class, c)
	return 2 * (precision * recall) / (precision + recall)
}

// GetAccuracy computes the overall classification accuracy
// That is (number of correctly classified instances) / total instances
func GetAccuracy(c ConfusionMatrix) float64 {
	correct := 0
	total := 0
	for i := range c {
		for j := range c[i] {
			if i == j {
				correct += c[i][j]
			}
			total += c[i][j]
		}
	}
	return float64(correct) / float64(total)
}

// GetMicroPrecision assesses Classifier performance across
// all classes using the total true positives and false positives.
func GetMicroPrecision(c ConfusionMatrix) float64 {
	truePositives := 0.0
	falsePositives := 0.0
	for k := range c {
		truePositives += GetTruePositives(k, c)
		falsePositives += GetFalsePositives(k, c)
	}
	return truePositives / (truePositives + falsePositives)
}

// GetMacroPrecision assesses Classifier performance across all
// classes by averaging the precision measures achieved for each class.
func GetMacroPrecision(c ConfusionMatrix) float64 {
	precisionVals := 0.0
	for k := range c {
		precisionVals += GetPrecision(k, c)
	}
	return precisionVals / float64(len(c))
}

// GetMicroRecall assesses Classifier performance across all
// classes using the total true positives and false negatives.
func GetMicroRecall(c ConfusionMatrix) float64 {
	truePositives := 0.0
	falseNegatives := 0.0
	for k := range c {
		truePositives += GetTruePositives(k, c)
		falseNegatives += GetFalseNegatives(k, c)
	}
	return truePositives / (truePositives + falseNegatives)
}

// GetMacroRecall assesses Classifier performance across all classes
// by averaging the recall measures achieved for each class
func GetMacroRecall(c ConfusionMatrix) float64 {
	recallVals := 0.0
	for k := range c {
		recallVals += GetRecall(k, c)
	}
	return recallVals / float64(len(c))
}

// GetSummary returns a table of precision, recall, true positive,
// false positive, and true negatives for each class for a given
// ConfusionMatrix
func GetSummary(c ConfusionMatrix) string {
	var buffer bytes.Buffer
	for k := range c {
		buffer.WriteString(k)
		buffer.WriteString("\t")
		tp := GetTruePositives(k, c)
		fp := GetFalsePositives(k, c)
		tn := GetTrueNegatives(k, c)
		prec := GetPrecision(k, c)
		rec := GetRecall(k, c)
		f1 := GetF1Score(k, c)
		buffer.WriteString(fmt.Sprintf("%.0f\t%.0f\t%.0f\t%.4f\t%.4f\t%.4f\n", tp, fp, tn, prec, rec, f1))
	}

	buffer.WriteString(fmt.Sprintf("Overall accuracy: %.4f\n", GetAccuracy(c)))

	return buffer.String()
}
