package linear_models

/*
#include "integration.h"
#cgo CFLAGS:
#cgo CXXFLAGS: -std=c++11 -g -O0
#cgo LDFLAGS: -g -llinear
*/
import "C"
import (
	"fmt"
	"runtime"
)

// Problem wraps a libsvm problem struct which describes a classification/
// regression problem. No externally-accessible fields.
type Problem struct {
	c_prob *C.struct_problem
}

// Free releases resources associated with a libsvm problem.
func (p *Problem) Free() {
	C.FreeCProblem(p.c_prob)
}

// Parameter encasulates all the possible libsvm training options.
// TODO: make user control of these more extensive.
type Parameter struct {
	c_param *C.struct_parameter
}

// Free releases resources associated with a Parameter.
func (p *Parameter) Free() {
	C.FreeCParameter(p.c_param)
}

// Model encapsulates a trained libsvm model.
type Model struct {
	c_model *C.struct_model
}

// Free releases resources associated with a trained libsvm model.
func (m *Model) Free() {
	C.FreeCModel(m.c_model)
}

const (
	L2R_LR              = C.L2R_LR
	L2R_L2LOSS_SVC_DUAL = C.L2R_L2LOSS_SVC_DUAL
	L2R_L2LOSS_SVC      = C.L2R_L2LOSS_SVC
	L2R_L1LOSS_SVC_DUAL = C.L2R_L1LOSS_SVC_DUAL
	MCSVM_CS            = C.MCSVM_CS
	L1R_L2LOSS_SVC      = C.L1R_L2LOSS_SVC
	L1R_LR              = C.L1R_LR
	L2R_LR_DUAL         = C.L2R_LR_DUAL
)

// NewParameter creates a libsvm parameter structure, which controls
// various aspects of libsvm training.
// For more information on what these parameters do, consult the
// "`train` usage" section of
// https://github.com/cjlin1/liblinear/blob/master/README
func NewParameter(solver_type int, C float64, eps float64) *Parameter {
	param := &Parameter{C.CreateCParameter()}
	runtime.SetFinalizer(param, (*Parameter).Free)
	param.c_param.solver_type = C.int(solver_type)
	param.c_param.eps = C.double(eps)
	param.c_param.C = C.double(C)
	param.c_param.nr_weight = C.int(0)
	param.c_param.weight_label = nil
	param.c_param.weight = nil

	return param
}

// NewProblem creates input to libsvm which describes a particular
// regression/classification problem. It requires an array of float values
// and an array of y values.
func NewProblem(X [][]float64, y []float64, bias float64) *Problem {
	prob := &Problem{C.CreateCProblem()}
	runtime.SetFinalizer(prob, (*Problem).Free)
	prob.c_prob.l = C.int(len(X))
	prob.c_prob.n = C.int(len(X[0]) + 1)

	convert_features(prob, X, bias)
	C.AllocateLabelsForProblem(prob.c_prob, C.int(len(y)))
	for i := 0; i < len(y); i++ {
		C.AssignLabelForProblem(prob.c_prob, C.int(i), C.double(y[i]))
	}
	// Should not go out of scope until the Problem struct
	// is cleaned up.
	prob.c_prob.bias = C.double(-1)

	return prob
}

// Train invokes libsvm and returns a trained model.
func Train(prob *Problem, param *Parameter) *Model {
	libLinearHookPrintFunc() // Sets up logging
	out := C.train(prob.c_prob, param.c_param)
	m := &Model{out}
	runtime.SetFinalizer(m, (*Model).Free)
	return m
}

func Export(model *Model, filePath string) error {
	status := C.save_model(C.CString(filePath), (*C.struct_model)(model.c_model))
	if status != 0 {
		return fmt.Errorf("Problem occured during export to %s (status was %d)", filePath, status)
	}
	return nil
}

func Load(model *Model, filePath string) error {
	model.c_model = C.load_model(C.CString(filePath))
	if model.c_model == nil {
		return fmt.Errorf("Something went wrong")
	}
	return nil

}

// Predict takes a row of float values corresponding to a particular
// input and returns the regression result.
func Predict(model *Model, x []float64) float64 {
	cX := convertVector(x, 0)
	cY := C.predict((*C.struct_model)(model.c_model), &cX[0])
	y := float64(cY)
	return y
}

// convertVector is an internal function used for converting
// dense float64 vectors into the sparse input that libsvm accepts.
func convertVector(x []float64, bias float64) []C.struct_feature_node {
	// Count the number of non-zero elements
	nElements := 0
	for i := 0; i < len(x); i++ {
		if x[i] > 0 {
			nElements++
		}
	}
	// Add one at the end for the -1 terminator
	nElements++
	if bias >= 0 {
		// And one for the bias, if we have it
		nElements++
	}

	cX := make([]C.struct_feature_node, nElements)
	j := 0
	for i := 0; i < len(x); i++ {
		if x[i] > 0 {
			cX[j].index = C.int(i + 1)
			cX[j].value = C.double(x[i])
			j++
		}
	}
	if bias >= 0 {
		cX[j].index = C.int(0)
		cX[j].value = C.double(0)
		j++
	}
	cX[j].index = C.int(-1)
	return cX
}

// convert_features is an internal function used for converting
// dense 2D arrays of float values into the sparse format libsvm accepts.
func convert_features(prob *Problem, X [][]float64, bias float64) {

	nonZeroRowElements := make([]C.int, len(X))
	totalElements := 0

	for i := 0; i < len(X); i++ {
		// For each row of input data, we count how many non-zero things are in the row
		nonZeroElementsInRow := 1 // Initially one, because we need the -1 null terminator
		for j := 0; j < len(X[i]); j++ {
			if X[i][j] != 0.0 {
				nonZeroElementsInRow++
			}
			if bias >= 0 {
				nonZeroElementsInRow++
			}
		}
		nonZeroRowElements[i] = C.int(nonZeroElementsInRow)
		totalElements += nonZeroElementsInRow
	}
	// Allocate one feature vector for each row, total number
	C.AllocateFeatureNodesForProblem(prob.c_prob, C.int(len(X)), C.int(totalElements), &nonZeroRowElements[0])

	for i := 0; i < len(X); i++ {
		nonZeroElementCounter := 0
		for j := 0; j < len(X[i]); j++ {
			if X[i][j] != 0.0 {
				xSpace := C.GetFeatureNodeForIndex(prob.c_prob, C.int(i), C.int(nonZeroElementCounter))
				xSpace.index = C.int(j + 1)
				xSpace.value = C.double(X[i][j])
				nonZeroElementCounter++
			}
		}
		if bias >= 0 {
			xSpace := C.GetFeatureNodeForIndex(prob.c_prob, C.int(i), C.int(nonZeroElementCounter))
			xSpace.index = C.int(len(X[i]) + 1)
			xSpace.value = C.double(bias)
			nonZeroElementCounter++
		}
		xSpace := C.GetFeatureNodeForIndex(prob.c_prob, C.int(i), C.int(nonZeroElementCounter))
		xSpace.index = C.int(-1)
		xSpace.value = C.double(0)
	}
}
