package base

import (
    "fmt"
    "strconv"
)

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

// Compatable checks whether this FloatAttribute can be ponded with another
// Attribute (checks if they're both FloatAttributes)
func (Attr *FloatAttribute) Compatable(other Attribute) bool {
    _, ok := other.(*FloatAttribute)
    return ok
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
// be converted into a valid system representation. If it can't,
// the returned value is nil.
func (Attr *FloatAttribute) CheckSysValFromString(rawVal string) ([]byte, error) {
    f, err := strconv.ParseFloat(rawVal, 64)
    if err != nil {
        return nil, err
    }

    ret := PackFloatToBytes(f)
    return ret, nil
}

// GetSysValFromString parses the given rawVal string to a float64 and returns it.
//
// float64 happens to be a 1-to-1 mapping to the system representation.
// IMPORTANT: This function panic()s if rawVal is not a valid float.
// Use CheckSysValFromString to confirm.
func (Attr *FloatAttribute) GetSysValFromString(rawVal string) []byte {
    f, err := Attr.CheckSysValFromString(rawVal)
    if err != nil {
        panic(err)
    }
    return f
}

// GetFloatFromSysVal converts a given system value to a float
func (Attr *FloatAttribute) GetFloatFromSysVal(rawVal []byte) float64 {
    return UnpackBytesToFloat(rawVal)
}

// GetStringFromSysVal converts a given system value to to a string with two decimal
// places of precision.
func (Attr *FloatAttribute) GetStringFromSysVal(rawVal []byte) string {
    f := UnpackBytesToFloat(rawVal)
    formatString := fmt.Sprintf("%%.%df", Attr.Precision)
    return fmt.Sprintf(formatString, f)
}
