package base

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"

	"github.com/gonum/matrix/mat64"
)

// SortDirection specifies sorting direction...
type SortDirection int

const (
	// Descending says that Instances should be sorted high to low...
	Descending SortDirection = 1
	// Ascending states that Instances should be sorted low to high...
	Ascending SortDirection = 2
)

const highBit int64 = -1 << 63

// Instances represents a grid of numbers (typed by Attributes)
// stored internally in mat.DenseMatrix as float64's.
// See docs/instances.md for more information.
type Instances struct {
	storage    *mat64.Dense
	attributes []Attribute
	Rows       int
	Cols       int
	ClassIndex int
}

func xorFloatOp(item float64) float64 {
	var ret float64
	var tmp int64
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, item)
	binary.Read(buf, binary.LittleEndian, &tmp)
	tmp ^= -1 << 63
	binary.Write(buf, binary.LittleEndian, tmp)
	binary.Read(buf, binary.LittleEndian, &ret)
	return ret
}

func printFloatByteArr(arr [][]byte) {
	buf := bytes.NewBuffer(nil)
	var f float64
	for _, b := range arr {
		buf.Write(b)
		binary.Read(buf, binary.LittleEndian, &f)
		f = xorFloatOp(f)
		fmt.Println(f)
	}
}

// Sort does an in-place radix sort of Instances, using SortDirection
// direction (Ascending or Descending) with attrs as a slice of Attribute
// indices that you want to sort by.
//
// IMPORTANT: Radix sort is not stable, so ordering outside
// the attributes used for sorting is arbitrary.
func (inst *Instances) Sort(direction SortDirection, attrs []int) {
	// Create a buffer
	buf := bytes.NewBuffer(nil)
	ds := make([][]byte, inst.Rows)
	rs := make([]int, inst.Rows)
	for i := 0; i < inst.Rows; i++ {
		byteBuf := make([]byte, 8*len(attrs))
		for _, a := range attrs {
			x := inst.storage.At(i, a)
			binary.Write(buf, binary.LittleEndian, xorFloatOp(x))
		}
		buf.Read(byteBuf)
		ds[i] = byteBuf
		rs[i] = i
	}
	// Sort viua
	valueBins := make([][][]byte, 256)
	rowBins := make([][]int, 256)
	for i := 0; i < 8*len(attrs); i++ {
		for j := 0; j < len(ds); j++ {
			// Address each row value by it's ith byte
			b := ds[j]
			valueBins[b[i]] = append(valueBins[b[i]], b)
			rowBins[b[i]] = append(rowBins[b[i]], rs[j])
		}
		j := 0
		for k := 0; k < 256; k++ {
			bs := valueBins[k]
			rc := rowBins[k]
			copy(ds[j:], bs)
			copy(rs[j:], rc)
			j += len(bs)
			valueBins[k] = bs[:0]
			rowBins[k] = rc[:0]
		}
	}

	for _, b := range ds {
		var v float64
		buf.Write(b)
		binary.Read(buf, binary.LittleEndian, &v)
	}

	done := make([]bool, inst.Rows)
	for index := range rs {
		if done[index] {
			continue
		}
		j := index
		for {
			done[j] = true
			if rs[j] != index {
				inst.swapRows(j, rs[j])
				j = rs[j]
			} else {
				break
			}
		}
	}

	if direction == Descending {
		// Reverse the matrix
		for i, j := 0, inst.Rows-1; i < j; i, j = i+1, j-1 {
			inst.swapRows(i, j)
		}
	}
}

// NewInstances returns a preallocated Instances structure
// with some helful values pre-filled.
func NewInstances(attrs []Attribute, rows int) *Instances {
	rawStorage := make([]float64, rows*len(attrs))
	return NewInstancesFromRaw(attrs, rows, rawStorage)
}

// CheckNewInstancesFromRaw checks whether a call to NewInstancesFromRaw
// is likely to produce an error-free result.
func CheckNewInstancesFromRaw(attrs []Attribute, rows int, data []float64) error {
	size := rows * len(attrs)
	if size < len(data) {
		return errors.New("base: data length is larger than the rows * attribute space.")
	} else if size > len(data) {
		return errors.New("base: data is smaller than the rows * attribute space")
	}
	return nil
}

// NewInstancesFromRaw wraps a slice of float64 numbers in a
// mat64.Dense structure, reshaping it with the given number of rows
// and representing it with the given attrs (Attribute slice)
//
// IMPORTANT: if the |attrs| * |rows| value doesn't equal len(data)
// then panic()s may occur. Use CheckNewInstancesFromRaw to confirm.
func NewInstancesFromRaw(attrs []Attribute, rows int, data []float64) *Instances {
	rawStorage := mat64.NewDense(rows, len(attrs), data)
	return NewInstancesFromDense(attrs, rows, rawStorage)
}

// NewInstancesFromDense creates a set of Instances from a mat64.Dense
// matrix
func NewInstancesFromDense(attrs []Attribute, rows int, mat *mat64.Dense) *Instances {
	return &Instances{mat, attrs, rows, len(attrs), len(attrs) - 1}
}

// InstancesTrainTestSplit takes a given Instances (src) and a train-test fraction
// (prop) and returns an array of two new Instances, one containing approximately
// that fraction and the other containing what's left.
//
// IMPORTANT: this function is only meaningful when prop is between 0.0 and 1.0.
// Using any other values may result in odd behaviour.
func InstancesTrainTestSplit(src *Instances, prop float64) (*Instances, *Instances) {
	trainingRows := make([]int, 0)
	testingRows := make([]int, 0)
	numAttrs := len(src.attributes)
	src.Shuffle()
	for i := 0; i < src.Rows; i++ {
		trainOrTest := rand.Intn(101)
		if trainOrTest > int(100*prop) {
			trainingRows = append(trainingRows, i)
		} else {
			testingRows = append(testingRows, i)
		}
	}

	rawTrainMatrix := mat64.NewDense(len(trainingRows), numAttrs, make([]float64, len(trainingRows)*numAttrs))
	rawTestMatrix := mat64.NewDense(len(testingRows), numAttrs, make([]float64, len(testingRows)*numAttrs))

	for i, row := range trainingRows {
		rowDat := src.storage.RowView(row)
		rawTrainMatrix.SetRow(i, rowDat)
	}
	for i, row := range testingRows {
		rowDat := src.storage.RowView(row)
		rawTestMatrix.SetRow(i, rowDat)
	}

	trainingRet := NewInstancesFromDense(src.attributes, len(trainingRows), rawTrainMatrix)
	testRet := NewInstancesFromDense(src.attributes, len(testingRows), rawTestMatrix)
	return trainingRet, testRet
}

// CountAttrValues returns the distribution of values of a given
// Attribute.
// IMPORTANT: calls panic() if the attribute index of a cannot be
// determined. Call GetAttrIndex(a) and check for a -1 return value.
func (inst *Instances) CountAttrValues(a Attribute) map[string]int {
	ret := make(map[string]int)
	attrIndex := inst.GetAttrIndex(a)
	if attrIndex == -1 {
		panic("Invalid attribute")
	}
	for i := 0; i < inst.Rows; i++ {
		sysVal := inst.Get(i, attrIndex)
		stringVal := a.GetStringFromSysVal(sysVal)
		ret[stringVal]++
	}
	return ret
}

// CountClassValues returns the class distribution of this
// Instances set
func (inst *Instances) CountClassValues() map[string]int {
	a := inst.GetAttr(inst.ClassIndex)
	return inst.CountAttrValues(a)
}

// DecomposeOnAttributeValues divides the instance set depending on the
// value of a given Attribute, constructs child instances, and returns
// them in a map keyed on the string value of that Attribute.
// IMPORTANT: calls panic() if the attribute index of at cannot be determined.
// Use GetAttrIndex(at) and check for a non-zero return value.
func (inst *Instances) DecomposeOnAttributeValues(at Attribute) map[string]*Instances {
	// Find the attribute we're decomposing on
	attrIndex := inst.GetAttrIndex(at)
	if attrIndex == -1 {
		panic("Invalid attribute index")
	}
	// Construct the new attribute set
	newAttrs := make([]Attribute, 0)
	for i := range inst.attributes {
		a := inst.attributes[i]
		if a.Equals(at) {
			continue
		}
		newAttrs = append(newAttrs, a)
	}
	// Create the return map, several counting maps
	ret := make(map[string]*Instances)
	counts := inst.CountAttrValues(at) // So we know what to allocate
	rows := make(map[string]int)
	for k := range counts {
		tmp := NewInstances(newAttrs, counts[k])
		ret[k] = tmp
	}
	for i := 0; i < inst.Rows; i++ {
		newAttrCounter := 0
		classVar := at.GetStringFromSysVal(inst.Get(i, attrIndex))
		dest := ret[classVar]
		destRow := rows[classVar]
		for j := 0; j < inst.Cols; j++ {
			a := inst.attributes[j]
			if a.Equals(at) {
				continue
			}
			dest.Set(destRow, newAttrCounter, inst.Get(i, j))
			newAttrCounter++
		}
		rows[classVar]++
	}
	return ret
}

func (inst *Instances) GetClassDistributionAfterSplit(at Attribute) map[string]map[string]int {

	ret := make(map[string]map[string]int)

	// Find the attribute we're decomposing on
	attrIndex := inst.GetAttrIndex(at)
	if attrIndex == -1 {
		panic("Invalid attribute index")
	}

	// Get the class index
	classAttr := inst.GetAttr(inst.ClassIndex)

	for i := 0; i < inst.Rows; i++ {
		splitVar := at.GetStringFromSysVal(inst.Get(i, attrIndex))
		classVar := classAttr.GetStringFromSysVal(inst.Get(i, inst.ClassIndex))
		if _, ok := ret[splitVar]; !ok {
			ret[splitVar] = make(map[string]int)
			i--
			continue
		}
		ret[splitVar][classVar]++
	}

	return ret
}

// Get returns the system representation (float64) of the value
// stored at the given row and col coordinate.
func (inst *Instances) Get(row int, col int) float64 {
	return inst.storage.At(row, col)
}

// Set sets the system representation (float64) to val at the
// given row and column coordinate.
func (inst *Instances) Set(row int, col int, val float64) {
	inst.storage.Set(row, col, val)
}

// GetRowVector returns a row of system representation
// values at the given row index.
func (inst *Instances) GetRowVector(row int) []float64 {
	return inst.storage.RowView(row)
}

// GetRowVector returns a row of system representation
// values at the given row index, excluding the class attribute
func (inst *Instances) GetRowVectorWithoutClass(row int) []float64 {
	rawRow := make([]float64, inst.Cols)
	copy(rawRow, inst.GetRowVector(row))
	return append(rawRow[0:inst.ClassIndex], rawRow[inst.ClassIndex+1:inst.Cols]...)
}

// GetClass returns the string representation of the given
// row's class, as determined by the Attribute at the ClassIndex
// position from GetAttr
func (inst *Instances) GetClass(row int) string {
	attr := inst.GetAttr(inst.ClassIndex)
	val := inst.Get(row, inst.ClassIndex)
	return attr.GetStringFromSysVal(val)
}

// GetClassDist returns a map containing the count of each
// class type (indexed by the class' string representation)
func (inst *Instances) GetClassDistribution() map[string]int {
	ret := make(map[string]int)
	attr := inst.GetAttr(inst.ClassIndex)
	for i := 0; i < inst.Rows; i++ {
		val := inst.Get(i, inst.ClassIndex)
		cls := attr.GetStringFromSysVal(val)
		ret[cls]++
	}

	return ret
}

func (Inst *Instances) GetClassAttrPtr() *Attribute {
	attr := Inst.GetAttr(Inst.ClassIndex)
	return &attr
}

func (Inst *Instances) GetClassAttr() Attribute {
	return Inst.GetAttr(Inst.ClassIndex)
}

//
// Attribute functions
//

// GetAttributeCount returns the number of attributes represented.
func (inst *Instances) GetAttributeCount() int {
	// Return the number of attributes attached to this Instance set
	return len(inst.attributes)
}

// SetAttrStr sets the system-representation value of row in column attr
// to value val, implicitly converting the string to system-representation
// via the appropriate Attribute function.
func (inst *Instances) SetAttrStr(row int, attr int, val string) {
	// Set an attribute on a particular row from a string value
	a := inst.attributes[attr]
	sysVal := a.GetSysValFromString(val)
	inst.storage.Set(row, attr, sysVal)
}

// GetAttrStr returns a human-readable string value stored in column `attr'
// and row `row', as determined by the appropriate Attribute function.
func (inst *Instances) GetAttrStr(row int, attr int) string {
	// Get a human-readable value from a particular row
	a := inst.attributes[attr]
	usrVal := a.GetStringFromSysVal(inst.Get(row, attr))
	return usrVal
}

// GetAttr returns information about an attribute at given index
// in the attributes slice.
func (inst *Instances) GetAttr(attrIndex int) Attribute {
	// Return a copy of an attribute attached to this Instance set
	return inst.attributes[attrIndex]
}

// GetAttrIndex returns the offset of a given Attribute `a' to an
// index in the attributes slice
func (inst *Instances) GetAttrIndex(of Attribute) int {
	// Finds the offset of an Attribute in this instance set
	// Returns -1 if no Attribute matches
	for i, a := range inst.attributes {
		if a.Equals(of) {
			return i
		}
	}
	return -1
}

// ReplaceAttr overwrites the attribute at `index' with `a'
func (inst *Instances) ReplaceAttr(index int, a Attribute) {
	// Replace an Attribute at index with another
	// DOESN'T CONVERT ANY EXISTING VALUES
	inst.attributes[index] = a
}

//
// Printing functions
//

// RowStr returns a human-readable representation of a given row.
func (inst *Instances) RowStr(row int) string {
	// Prints a given row
	var buffer bytes.Buffer
	for j := 0; j < inst.Cols; j++ {
		val := inst.storage.At(row, j)
		a := inst.attributes[j]
		postfix := " "
		if j == inst.Cols-1 {
			postfix = ""
		}
		buffer.WriteString(fmt.Sprintf("%s%s", a.GetStringFromSysVal(val), postfix))
	}
	return buffer.String()
}

func (inst *Instances) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("Instances with ")
	buffer.WriteString(fmt.Sprintf("%d row(s) ", inst.Rows))
	buffer.WriteString(fmt.Sprintf("%d attribute(s)\n", inst.Cols))

	buffer.WriteString(fmt.Sprintf("Attributes: \n"))
	for i, a := range inst.attributes {
		prefix := "\t"
		if i == inst.ClassIndex {
			prefix = "*\t"
		}
		buffer.WriteString(fmt.Sprintf("%s%s\n", prefix, a))
	}

	buffer.WriteString("\nData:\n")
	maxRows := 30
	if inst.Rows < maxRows {
		maxRows = inst.Rows
	}

	for i := 0; i < maxRows; i++ {
		buffer.WriteString("\t")
		for j := 0; j < inst.Cols; j++ {
			val := inst.storage.At(i, j)
			a := inst.attributes[j]
			buffer.WriteString(fmt.Sprintf("%s ", a.GetStringFromSysVal(val)))
		}
		buffer.WriteString("\n")
	}

	missingRows := inst.Rows - maxRows
	if missingRows != 0 {
		buffer.WriteString(fmt.Sprintf("\t...\n%d row(s) undisplayed", missingRows))
	} else {
		buffer.WriteString("All rows displayed")
	}

	return buffer.String()
}

// SelectAttributes returns a new instance set containing
// the values from this one with only the Attributes specified
func (inst *Instances) SelectAttributes(attrs []Attribute) *Instances {
	ret := NewInstances(attrs, inst.Rows)
	attrIndices := make([]int, 0)
	for _, a := range attrs {
		attrIndex := inst.GetAttrIndex(a)
		attrIndices = append(attrIndices, attrIndex)
	}
	for i := 0; i < inst.Rows; i++ {
		for j, a := range attrIndices {
			ret.Set(i, j, inst.Get(i, a))
		}
	}
	return ret
}

// GeneratePredictionVector generates a new set of Instances
// with the same number of rows, but only this Instance set's
// class Attribute.
func (inst *Instances) GeneratePredictionVector() *Instances {
	attrs := make([]Attribute, 1)
	attrs[0] = inst.GetClassAttr()
	ret := NewInstances(attrs, inst.Rows)
	return ret
}

// Shuffle randomizes the row order in place
func (inst *Instances) Shuffle() {
	for i := 0; i < inst.Rows; i++ {
		j := rand.Intn(i + 1)
		inst.swapRows(i, j)
	}
}

// SampleWithReplacement returns a new set of Instances of size `size'
// containing random rows from this set of Instances.
//
// IMPORTANT: There's a high chance of seeing duplicate rows
// whenever size is close to the row count.
func (inst *Instances) SampleWithReplacement(size int) *Instances {
	ret := NewInstances(inst.attributes, size)
	for i := 0; i < size; i++ {
		srcRow := rand.Intn(inst.Rows)
		for j := 0; j < inst.Cols; j++ {
			ret.Set(i, j, inst.Get(srcRow, j))
		}
	}
	return ret
}

// Equal checks whether a given Instance set is exactly the same
// as another: same size and same values (as determined by the Attributes)
//
// IMPORTANT: does not explicitly check if the Attributes are considered equal.
func (inst *Instances) Equal(other *Instances) bool {
	if inst.Rows != other.Rows {
		return false
	}
	if inst.Cols != other.Cols {
		return false
	}
	for i := 0; i < inst.Rows; i++ {
		for j := 0; j < inst.Cols; j++ {
			if inst.GetAttrStr(i, j) != other.GetAttrStr(i, j) {
				return false
			}
		}
	}
	return true
}

func (inst *Instances) swapRows(r1 int, r2 int) {
	row1buf := make([]float64, inst.Cols)
	row2buf := make([]float64, inst.Cols)
	row1 := inst.storage.RowView(r1)
	row2 := inst.storage.RowView(r2)
	copy(row1buf, row1)
	copy(row2buf, row2)
	inst.storage.SetRow(r1, row2buf)
	inst.storage.SetRow(r2, row1buf)
}
