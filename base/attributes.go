package base

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

const (
	// CategoricalType is for Attributes which represent values distinctly.
	CategoricalType = iota
	// Float64Type should be replaced with a FractionalNumeric type [DEPRECATED].
	Float64Type
	// Int64Type is for IntAttributes
	Int64Type
)

// Attribute Attributes disambiguate columns of the feature matrix and declare their types.
type Attribute interface {
	// Returns the general characterstics of this Attribute .
	// to avoid the overhead of casting
	GetType() int
	// Returns the human-readable name of this Attribute.
	GetName() string
	// Sets the human-readable name of this Attribute.
	SetName(string)
	// Gets a human-readable overview of this Attribute for debugging.
	String() string
	// Converts a value given from a human-readable string into a system
	// representation. For example, a CategoricalAttribute with values
	// ["iris-setosa", "iris-virginica"] would return the float64
	// representation of 0 when given "iris-setosa".
	GetSysValFromString(string) float64
	// Converts a given value from a system representation into a human
	// representation. For example, a CategoricalAttribute with values
	// ["iris-setosa", "iris-viriginica"] might return "iris-setosa"
	// when given 0.0 as the argument.
	GetStringFromSysVal(float64) string
	// Tests for equality with another Attribute. Other Attributes are
	// considered equal if:
	// * They have the same type (i.e. FloatAttribute <> CategoricalAttribute)
	// * They have the same name
	// * If applicable, they have the same categorical values (though not
	//   necessarily in the same order).
	Equals(Attribute) bool
}

// IntAttribute is an implementation which stores 64-bit numbers
// in native endianness
type IntAttribute struct {
	Name string
}

// NewIntAttribute returns a new IntAttribute
func NewIntAttribute() *IntAttribute {
	return &IntAttribute{""}
}

// Equals tests if an IntAttribute is the same as another
//
// Returns false if the other attribute has a different name
// or is not IntAttribute
func (Attr *IntAttribute) Equals(other Attribute) bool {
	_, ok := other.(*IntAttribute)
	if !ok {
		return false
	}
	if Attr.GetName() != other.GetName() {
		return false
	}
	return true
}

// GetName returns this IntAttribute's human-readable name
func (Attr *IntAttribute) GetName() string {
	return Attr.Name
}

// SetName sets this IntAttribute's human-readable name
func (Attr *IntAttribute) SetName(name string) {
	Attr.Name = name
}

// GetType returns Int64Type
func (Attr *IntAttribute) GetType() int {
	return Int64Type
}

// String returns a human-readable summary of this Attribute
func (Attr *IntAttribute) String() string {
	return fmt.Sprintf("IntAttribute(%s)", Attr.Name)
}

// CheckSysValFromString confirms whether a given rawVal can
// be converted into a valid system representation.
func (Attr *IntAttribute) CheckSysValFromString(rawVal string) (float64, error) {
	f, err := strconv.ParseInt(rawVal, 10, 64)
	if err != nil {
		return 0.0, err
	}
	bs := make([]byte, 8)
	binary.PutVarint(bs, f)
	fu, amount := binary.Uvarint(bs)
	if amount < 1 {
		panic("Conversion failed")
	}
	float := math.Float64frombits(fu)
	return float, nil
}

// GetSysValFromString parses the given rawVal string to a float64 and returns it.
//
// IMPORTANT: This function panic()s if rawVal is not a valid integer.
// Use CheckSysValFromString to confirm.
func (Attr *IntAttribute) GetSysValFromString(rawVal string) float64 {
	f, err := Attr.CheckSysValFromString(rawVal)
	if err != nil {
		panic(err)
	}
	return f
}

// GetStringFromSysVal converts a given system value to to a string
func (Attr *IntAttribute) GetStringFromSysVal(rawVal float64) string {
	bytes := math.Float64bits(rawVal)
	buf := make([]byte, 8)
	binary.PutUvarint(buf, bytes)
	val, read := binary.Varint(buf)
	if read < 1 {
		panic("Conversion failed!")
	}
	return fmt.Sprintf("%d", val)
}

// FloatAttribute is an implementation which stores floating point
// representations of numbers.
type FloatAttribute struct {
	Name      string
	Precision int
}

// NewFloatAttribute returns a new FloatAttribute with a default
// precision of 2 decimal places
func NewFloatAttribute() *FloatAttribute {
	return &FloatAttribute{"", 2}
}

// Equals tests a FloatAttribute for equality with another Attribute.
//
// Returns false if the other Attribute has a different name
// or if the other Attribute is not a FloatAttribute.
func (Attr *FloatAttribute) Equals(other Attribute) bool {
	// Check whether this FloatAttribute is equal to another
	_, ok := other.(*FloatAttribute)
	if !ok {
		// Not the same type, so can't be equal
		return false
	}
	if Attr.GetName() != other.GetName() {
		return false
	}
	return true
}

// GetName returns this FloatAttribute's human-readable name.
func (Attr *FloatAttribute) GetName() string {
	return Attr.Name
}

// SetName sets this FloatAttribute's human-readable name.
func (Attr *FloatAttribute) SetName(name string) {
	Attr.Name = name
}

// GetType returns Float64Type.
func (Attr *FloatAttribute) GetType() int {
	return Float64Type
}

// String returns a human-readable summary of this Attribute.
// e.g. "FloatAttribute(Sepal Width)"
func (Attr *FloatAttribute) String() string {
	return fmt.Sprintf("FloatAttribute(%s)", Attr.Name)
}

// CheckSysValFromString confirms whether a given rawVal can
// be converted into a valid system representation.
func (Attr *FloatAttribute) CheckSysValFromString(rawVal string) (float64, error) {
	f, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		return 0.0, err
	}
	return f, nil
}

// GetSysValFromString parses the given rawVal string to a float64 and returns it.
//
// float64 happens to be a 1-to-1 mapping to the system representation.
// IMPORTANT: This function panic()s if rawVal is not a valid float.
// Use CheckSysValFromString to confirm.
func (Attr *FloatAttribute) GetSysValFromString(rawVal string) float64 {
	f, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// GetStringFromSysVal converts a given system value to to a string with two decimal
// places of precision [TODO: revise this and allow more precision].
func (Attr *FloatAttribute) GetStringFromSysVal(rawVal float64) string {
	formatString := fmt.Sprintf("%%.%df", Attr.Precision)
	return fmt.Sprintf(formatString, rawVal)
}

// GetSysVal returns the system representation of userVal.
//
// Because FloatAttribute represents float64 types, this
// just returns its argument.
func (Attr *FloatAttribute) GetSysVal(userVal float64) float64 {
	return userVal
}

// GetUsrVal returns the user representation of sysVal.
//
// Because FloatAttribute represents float64 types, this
// just returns its argument.
func (Attr *FloatAttribute) GetUsrVal(sysVal float64) float64 {
	return sysVal
}

// CategoricalAttribute is an Attribute implementation
// which stores discrete string values
// - useful for representing classes.
type CategoricalAttribute struct {
	Name   string
	values []string
}

func NewCategoricalAttribute() *CategoricalAttribute {
	return &CategoricalAttribute{
		"",
		make([]string, 0),
	}
}

// GetName returns the human-readable name assigned to this attribute.
func (Attr *CategoricalAttribute) GetName() string {
	return Attr.Name
}

// SetName sets the human-readable name on this attribute.
func (Attr *CategoricalAttribute) SetName(name string) {
	Attr.Name = name
}

// GetType returns CategoricalType to avoid casting overhead.
func (Attr *CategoricalAttribute) GetType() int {
	return CategoricalType
}

// GetSysVal returns the system representation of userVal as an index into the Values slice
// If the userVal can't be found, it returns -1.
func (Attr *CategoricalAttribute) GetSysVal(userVal string) float64 {
	for idx, val := range Attr.values {
		if val == userVal {
			return float64(idx)
		}
	}
	return -1
}

// GetUsrVal returns a human-readable representation of the given sysVal.
//
// IMPORTANT: this function doesn't check the boundaries of the array.
func (Attr *CategoricalAttribute) GetUsrVal(sysVal float64) string {
	idx := int(sysVal)
	return Attr.values[idx]
}

// GetSysValFromString returns the system representation of rawVal
// as an index into the Values slice. If rawVal is not inside
// the Values slice, it is appended.
//
// IMPORTANT: If no system representation yet exists, this functions adds it.
// If you need to determine whether rawVal exists: use GetSysVal and check
// for a -1 return value.
//
// Example: if the CategoricalAttribute contains the values ["iris-setosa",
// "iris-virginica"] and "iris-versicolor" is provided as the argument,
// the Values slide becomes ["iris-setosa", "iris-virginica", "iris-versicolor"]
// and 2.00 is returned as the system representation.
func (Attr *CategoricalAttribute) GetSysValFromString(rawVal string) float64 {
	// Match in raw values
	catIndex := -1
	for i, s := range Attr.values {
		if s == rawVal {
			catIndex = i
			break
		}
	}
	if catIndex == -1 {
		Attr.values = append(Attr.values, rawVal)
		catIndex = len(Attr.values) - 1
	}
	return float64(catIndex)
}

// String returns a human-readable summary of this Attribute.
//
// Returns a string containing the list of human-readable values this
// CategoricalAttribute can take.
func (Attr *CategoricalAttribute) String() string {
	return fmt.Sprintf("CategoricalAttribute(\"%s\", %s)", Attr.Name, Attr.values)
}

// GetStringFromSysVal returns a human-readable value from the given system-representation
// value val.
//
// IMPORTANT: This function calls panic() if the value is greater than
// the length of the array.
// TODO: Return a user-configurable default instead.
func (Attr *CategoricalAttribute) GetStringFromSysVal(val float64) string {
	convVal := int(val)
	if convVal >= len(Attr.values) {
		panic(fmt.Sprintf("Out of range: %d in %d", convVal, len(Attr.values)))
	}
	return Attr.values[convVal]
}

// Equals checks equality against another Attribute.
//
// Two CategoricalAttributes are considered equal if they contain
// the same values and have the same name. Otherwise, this function
// returns false.
func (Attr *CategoricalAttribute) Equals(other Attribute) bool {
	attribute, ok := other.(*CategoricalAttribute)
	if !ok {
		// Not the same type, so can't be equal
		return false
	}
	if Attr.GetName() != attribute.GetName() {
		return false
	}

	// Check that this CategoricalAttribute has the same
	// values as the other, in the same order
	if len(attribute.values) != len(Attr.values) {
		return false
	}

	for i, a := range Attr.values {
		if a != attribute.values[i] {
			return false
		}
	}

	return true
}
