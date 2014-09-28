package base

import (
	"fmt"
	"strconv"
)

// BinaryAttributes can only represent 1 or 0.
type BinaryAttribute struct {
	Name string
}

// NewBinaryAttribute creates a BinaryAttribute with the given name
func NewBinaryAttribute(name string) *BinaryAttribute {
	return &BinaryAttribute{
		name,
	}
}

// GetName returns the name of this Attribute.
func (b *BinaryAttribute) GetName() string {
	return b.Name
}

// SetName sets the name of this Attribute.
func (b *BinaryAttribute) SetName(name string) {
	b.Name = name
}

// GetType returns BinaryType.
func (b *BinaryAttribute) GetType() int {
	return BinaryType
}

// GetSysValFromString returns either 1 or 0 in a single byte.
func (b *BinaryAttribute) GetSysValFromString(userVal string) []byte {
	f, err := strconv.ParseFloat(userVal, 64)
	if err != nil {
		panic(err)
	}
	ret := make([]byte, 1)
	if f > 0 {
		ret[0] = 1
	}
	return ret
}

// GetStringFromSysVal returns either 1 or 0.
func (b *BinaryAttribute) GetStringFromSysVal(val []byte) string {
	if val[0] > 0 {
		return "1"
	}
	return "0"
}

// Equals checks for equality with another BinaryAttribute.
func (b *BinaryAttribute) Equals(other Attribute) bool {
	if a, ok := other.(*BinaryAttribute); !ok {
		return false
	} else {
		return a.Name == b.Name
	}
}

// Compatible checks whether this Attribute can be represented
// in the same pond as another.
func (b *BinaryAttribute) Compatible(other Attribute) bool {
	if _, ok := other.(*BinaryAttribute); !ok {
		return false
	} else {
		return true
	}
}

// String returns a human-redable representation.
func (b *BinaryAttribute) String() string {
	return fmt.Sprintf("BinaryAttribute(%s)", b.Name)
}
