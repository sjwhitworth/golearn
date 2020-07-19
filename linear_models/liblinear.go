package linear_models

/*
#include "integration.h"
#cgo CFLAGS: -O3
#cgo CXXFLAGS: -std=c++11 -O3
*/
import "C"
import (
	"fmt"
	"runtime"
)

// Problem wraps a libsvm problem struct which describes a classification/
// regression problem. No externally-accessible fields.
type Problem struct {
	c_prob       *C.struct_problem
	featureNodes []C.struct_feature_node
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
	prob := &Problem{C.CreateCProblem(), nil}
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
	c_x := convert_vector(x, 0)
	c_y := C.predict((*C.struct_model)(model.c_model), c_x)
	y := float64(c_y)
	return y
}

// convert_vector is an internal function used for converting
// dense float64 vectors into the sparse input that libsvm accepts.
func convert_vector(x []float64, bias float64) *C.struct_feature_node {
	n_ele := 0
	for i := 0; i < len(x); i++ {
		if x[i] > 0 {
			n_ele++
		}
	}
	n_ele += 2

	c_x := make([]C.struct_feature_node, n_ele)
	j := 0
	for i := 0; i < len(x); i++ {
		if x[i] > 0 {
			c_x[j].index = C.int(i + 1)
			c_x[j].value = C.double(x[i])
			j++
		}
	}
	if bias > 0 {
		c_x[j].index = C.int(0)
		c_x[j].value = C.double(0)
		j++
	}
	c_x[j].index = C.int(-1)
	return &c_x[0]
}

// convert_features is an internal function used for converting
// dense 2D arrays of float values into the sparse format libsvm accepts.
func convert_features(prob *Problem, X [][]float64, bias float64) {
	rowCount := len(X)

	// This structure remembers the start and end elements for each row.
	// We push them back into a global list of C.struct_feature_nodes, then
	// riffle it in C using their indices to form the **C.struct_feature_nodes
	// input. Go retains ownership of struct_feature_nodes, C has ownership of
	// the enclosing **C.struct_feature_nodes feature.
	rowOffsets := make([]C.int, 0)
	featureNodes := make([]C.struct_feature_node, 0)

	// First pass, just counting through each row and counting the number of elements we find.
	for i := 0; i < rowCount; i++ {
		rowOffsets = append(rowOffsets, C.int(len(featureNodes))) // Push back the starting element of this row
		if bias != 0.0 {                                          // Allocate space for a bias node
			featureNodes = append(featureNodes, C.struct_feature_node{
				C.int(0), C.double(bias),
			})
		}
		for j := 0; j < len(X[i]); j++ {
			if X[i][j] != 0.0 {
				// For every non-zero thing in the data grid, allocate a feature node.
				featureNodes = append(featureNodes, C.struct_feature_node{
					C.int(j + 1), C.double(X[i][j]),
				})
			}
		}
		// Finally, add a terminating element which tells libsvm that there's nothing
		// left on this row
		featureNodes = append(featureNodes, C.struct_feature_node{
			C.int(-1), C.double(0),
		})
	}

	// Transform [feature_node, feature_node, feature_node, ...] list into
	// [*feature_node(1), *feature_node(m), ...] through the C integration bridge.
	// C owns that particular memory.
	// int RiffleFeatures(struct problem *p, int num_offsets, int* row_offsets, struct feature_node *features) {
	if C.RiffleFeatures(prob.c_prob, C.int(len(featureNodes)), &rowOffsets[0], &featureNodes[0]) != 0 {
		panic("RiffledFeatures could not allocate memory")
	}
}
