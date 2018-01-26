package linear_models

import "C"

import (
	"encoding/json"
	"fmt"
	"github.com/amclay/golearn/base"
	"io/ioutil"
	"os"
	"unsafe"
)

// LinearSVCParams represnts all available LinearSVC options.
//
// SolverKind: can be linear_models.L2_L1LOSS_SVC_DUAL,
// L2R_L2LOSS_SVC_DUAL, L2R_L2LOSS_SVC, L1R_L2LOSS_SVC.
// It must be set via SetKindFromStrings.
//
// ClassWeights describes how each class is weighted, and can
// be used in class-imabalanced scenarios. If this is nil, then
// all classes will be weighted the same unless WeightClassesAutomatically
// is True.
//
// C is a float64 represnenting the misclassification penalty.
//
// Eps is a float64 convergence threshold.
//
// Dual indicates whether the solution is primary or dual.
type LinearSVCParams struct {
	SolverType                 int
	ClassWeights               []float64
	C                          float64
	Eps                        float64
	WeightClassesAutomatically bool
	Dual                       bool
}

// Copy return s a copy of these parameters
func (p *LinearSVCParams) Copy() *LinearSVCParams {
	ret := &LinearSVCParams{
		p.SolverType,
		nil,
		p.C,
		p.Eps,
		p.WeightClassesAutomatically,
		p.Dual,
	}
	if p.ClassWeights != nil {
		ret.ClassWeights = make([]float64, len(p.ClassWeights))
		copy(ret.ClassWeights, p.ClassWeights)
	}
	return ret
}

// SetKindFromStrings configures the solver kind from strings.
// Penalty and Loss parameters can either be l1 or l2.
func (p *LinearSVCParams) SetKindFromStrings(loss, penalty string) error {
	var ret error
	p.SolverType = 0
	// Loss validation
	if loss == "l1" {
	} else if loss == "l2" {
	} else {
		return fmt.Errorf("loss must be \"l1\" or \"l2\"")
	}

	// Penalty validation
	if penalty == "l2" {
		if loss == "l1" {
			if !p.Dual {
				ret = fmt.Errorf("Important: changed to dual form")
			}
			p.SolverType = L2R_L1LOSS_SVC_DUAL
			p.Dual = true
		} else {
			if p.Dual {
				p.SolverType = L2R_L2LOSS_SVC_DUAL
			} else {
				p.SolverType = L2R_L2LOSS_SVC
			}
		}
	} else if penalty == "l1" {
		if loss == "l2" {
			if p.Dual {
				ret = fmt.Errorf("Important: changed to primary form")
			}
			p.Dual = false
			p.SolverType = L1R_L2LOSS_SVC
		} else {
			return fmt.Errorf("Must have L2 loss with L1 penalty")
		}
	} else {
		return fmt.Errorf("Penalty must be \"l1\" or \"l2\"")
	}

	// Final validation
	if p.SolverType == 0 {
		return fmt.Errorf("Invalid parameter combination")
	}
	return ret
}

// convertToNativeFormat converts the LinearSVCParams given into a format
// for liblinear.
func (p *LinearSVCParams) convertToNativeFormat() *Parameter {
	return NewParameter(p.SolverType, p.C, p.Eps)
}

// LinearSVC represents a linear support-vector classifier.
type LinearSVC struct {
	param *Parameter
	model *Model
	Param *LinearSVCParams
}

// NewLinearSVC creates a new support classifier.
//
// loss and penalty: see LinearSVCParams#SetKindFromString
//
// dual: see LinearSVCParams
//
// eps: see LinearSVCParams
//
// C: see LinearSVCParams
func NewLinearSVC(loss, penalty string, dual bool, C float64, eps float64) (*LinearSVC, error) {

	// Convert and check parameters
	params := &LinearSVCParams{0, nil, C, eps, false, dual}
	err := params.SetKindFromStrings(loss, penalty)
	if err != nil {
		return nil, err
	}

	return NewLinearSVCFromParams(params)
}

// NewLinearSVCFromParams constructs a LinearSVC from the given LinearSVCParams structure.
func NewLinearSVCFromParams(params *LinearSVCParams) (*LinearSVC, error) {
	// Construct model
	lr := LinearSVC{}
	lr.param = params.convertToNativeFormat()
	lr.Param = params
	lr.model = nil
	return &lr, nil
}

// Fit automatically weights the class vector (if configured to do so)
// converts the FixedDataGrid into the right format and trains the model.
func (lr *LinearSVC) Fit(X base.FixedDataGrid) error {

	var weightVec []float64
	var weightClasses []C.int

	// Creates the class weighting
	if lr.Param.ClassWeights == nil {
		if lr.Param.WeightClassesAutomatically {
			weightVec = generateClassWeightVectorFromDist(X)
		} else {
			weightVec = generateClassWeightVectorFromFixed(X)
		}
	} else {
		weightVec = lr.Param.ClassWeights
	}

	weightClasses = make([]C.int, len(weightVec))
	for i := range weightVec {
		weightClasses[i] = C.int(i)
	}

	// Convert the problem
	problemVec := convertInstancesToProblemVec(X)
	labelVec := convertInstancesToLabelVec(X)

	// Train
	prob := NewProblem(problemVec, labelVec, 0)
	lr.param.c_param.nr_weight = C.int(len(weightVec))
	lr.param.c_param.weight_label = &(weightClasses[0])
	lr.param.c_param.weight = (*C.double)(unsafe.Pointer(&weightVec[0]))

	//	lr.param.weights = (*C.double)unsafe.Pointer(&(weightVec[0]));
	lr.model = Train(prob, lr.param)
	return nil
}

// Predict issues predictions from a trained LinearSVC.
func (lr *LinearSVC) Predict(X base.FixedDataGrid) (base.FixedDataGrid, error) {

	// Only support 1 class Attribute
	classAttrs := X.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic(fmt.Sprintf("%d Wrong number of classes", len(classAttrs)))
	}
	// Generate return structure
	ret := base.GeneratePredictionVector(X)
	classAttrSpecs := base.ResolveAttributes(ret, classAttrs)

	// Retrieve numeric non-class Attributes
	numericAttrs := base.NonClassFloatAttributes(X)
	numericAttrSpecs := base.ResolveAttributes(X, numericAttrs)

	// Allocate row storage
	row := make([]float64, len(numericAttrSpecs))
	X.MapOverRows(numericAttrSpecs, func(rowBytes [][]byte, rowNo int) (bool, error) {
		for i, r := range rowBytes {
			row[i] = base.UnpackBytesToFloat(r)
		}
		val := Predict(lr.model, row)
		vals := base.PackFloatToBytes(val)
		ret.Set(classAttrSpecs[0], rowNo, vals)
		return true, nil
	})

	return ret, nil
}

func (lr *LinearSVC) GetMetadata() base.ClassifierMetadataV1 {

	params, err := json.Marshal(lr.Param)
	if err != nil {
		panic(err)
	}

	classifierParams := make(map[string]interface{})
	err = json.Unmarshal(params, &classifierParams)
	if err != nil {
		panic(err)
	}

	return base.ClassifierMetadataV1{
		FormatVersion:      1,
		ClassifierName:     "LinearSVC",
		ClassifierVersion:  "1.0",
		ClassifierMetadata: classifierParams,
	}
}

// Save outputs this classifier
func (lr *LinearSVC) Save(filePath string) error {
	writer, err := base.CreateSerializedClassifierStub(filePath, lr.GetMetadata())
	if err != nil {
		return err
	}
	defer func() {
		writer.Close()
	}()
	fmt.Printf("writer: %v", writer)
	return lr.SaveWithPrefix(writer, "")
}

func (lr *LinearSVC) SaveWithPrefix(writer *base.ClassifierSerializer, prefix string) error {
	params, err := json.Marshal(lr.Param)
	if err != nil {
		return base.DescribeError("Error marshalling parameters", err)
	}

	classifierParams := make(map[string]interface{})
	err = json.Unmarshal(params, &classifierParams)
	if err != nil {
		return base.DescribeError("Error marshalling parameters", err)
	}

	f, err := ioutil.TempFile(os.TempDir(), "liblinear")
	defer func() {
		f.Close()
	}()

	err = Export(lr.model, f.Name())
	if err != nil {
		return base.DescribeError("Error exporting model", err)
	}

	f.Seek(0, os.SEEK_SET)
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return base.DescribeError("Error reading model in again", err)
	}

	err = writer.WriteBytesForKey(writer.Prefix(prefix, "PARAMETERS"), params)
	if err != nil {
		return base.DescribeError("Error writing model parameters", err)
	}
	writer.WriteBytesForKey(writer.Prefix(prefix, "MODEL"), bytes)
	if err != nil {
		return base.DescribeError("Error writing model", err)
	}

	return nil
}

func (lr *LinearSVC) Load(filePath string) error {
	reader, err := base.ReadSerializedClassifierStub(filePath)
	if err != nil {
		return err
	}
	defer func() {
		reader.Close()
	}()

	return lr.LoadWithPrefix(reader, "")
}

func (lr *LinearSVC) LoadWithPrefix(reader *base.ClassifierDeserializer, prefix string) error {
	err := reader.GetJSONForKey(reader.Prefix(prefix, "PARAMETERS"), &lr.Param)
	if err != nil {
		return base.DescribeError("Error reading PARAMETERS", err)
	}

	modelBytes, err := reader.GetBytesForKey(reader.Prefix(prefix, "MODEL"))
	if err != nil {
		return base.DescribeError("Error reading MODEL", err)
	}

	f, err := ioutil.TempFile(os.TempDir(), "linear")
	defer func() {
		f.Close()
	}()

	f.WriteAt(modelBytes, 0)

	lr.model = &Model{}

	err = Load(lr.model, f.Name())
	if err != nil {
		return base.DescribeError("Unable to reload the model", err)
	}
	return nil
}

// String return a humaan-readable version.
func (lr *LinearSVC) String() string {
	return "LogisticSVC"
}
