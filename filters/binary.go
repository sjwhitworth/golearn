package filters

import (
	"fmt"

	"github.com/amclay/golearn/base"
)

// BinaryConvertFilters convert a given DataGrid into one which
// only contains BinaryAttributes.
//
// FloatAttributes are discretised into either 0 (if the value is 0)
// or 1 (if the value is not 0).
//
// CategoricalAttributes are discretised into one or more new
// BinaryAttributes.
type BinaryConvertFilter struct {
	attrs                          []base.Attribute
	converted                      []base.FilteredAttribute
	twoValuedCategoricalAttributes map[base.Attribute]bool // Two-valued categorical Attributes
	nValuedCategoricalAttributeMap map[base.Attribute]map[uint64]base.Attribute
}

// NewBinaryConvertFilter creates a blank BinaryConvertFilter
func NewBinaryConvertFilter() *BinaryConvertFilter {
	ret := &BinaryConvertFilter{
		make([]base.Attribute, 0),
		make([]base.FilteredAttribute, 0),
		make(map[base.Attribute]bool),
		make(map[base.Attribute]map[uint64]base.Attribute),
	}
	return ret
}

// AddAttribute adds a new Attribute to this Filter
func (b *BinaryConvertFilter) AddAttribute(a base.Attribute) error {
	b.attrs = append(b.attrs, a)
	return nil
}

// GetAttributesAfterFiltering returns the Attributes previously computed via Train()
func (b *BinaryConvertFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	return b.converted
}

// String gets a human-readable string
func (b *BinaryConvertFilter) String() string {
	return fmt.Sprintf("BinaryConvertFilter(%d Attribute(s))", len(b.attrs))
}

// Transform converts the given byte sequence using the old Attribute into the new
// byte sequence.
//
// If the old Attribute has a categorical value of at most two items, then a zero or
// non-zero byte sequence is returned.
//
// If the old Attribute has a categorical value of at most n-items, then a non-zero
// or zero byte sequence is returned based on the value of the new Attribute passed in.
//
// If the old Attribute is a float, it's value's unpacked and we check for non-zeroness
//
// If the old Attribute is a BinaryAttribute, just return the input
func (b *BinaryConvertFilter) Transform(a base.Attribute, n base.Attribute, attrBytes []byte) []byte {
	ret := make([]byte, 1)
	// Check for CategoricalAttribute
	if _, ok := a.(*base.CategoricalAttribute); ok {
		// Unpack byte value
		val := base.UnpackBytesToU64(attrBytes)
		// If it's a two-valued one, check for non-zero
		if b.twoValuedCategoricalAttributes[a] {
			if val > 0 {
				ret[0] = 1
			} else {
				ret[0] = 0
			}
		} else if an, ok := b.nValuedCategoricalAttributeMap[a]; ok {
			// If it's an n-valued one, check the new Attribute maps onto
			// the unpacked value
			if af, ok := an[val]; ok {
				if af.Equals(n) {
					ret[0] = 1
				} else {
					ret[0] = 0
				}
			} else {
				panic("Categorical value not defined!")
			}
		} else {
			panic(fmt.Sprintf("Not a recognised Attribute %v", a))
		}
	} else if _, ok := a.(*base.BinaryAttribute); ok {
		// Binary: just return the original value
		ret = attrBytes
	} else if _, ok := a.(*base.FloatAttribute); ok {
		// Float: check for non-zero
		val := base.UnpackBytesToFloat(attrBytes)
		if val > 0 {
			ret[0] = 1
		} else {
			ret[0] = 0
		}
	} else {
		panic(fmt.Sprintf("Unrecognised Attribute: %v", a))
	}
	return ret
}

// Train converts the FloatAttributesinto equivalently named BinaryAttributes,
// leaves BinaryAttributes unmodified and processes
// CategoricalAttributes as follows.
//
// If the CategoricalAttribute has two values, one of them is
// designated 0 and the other 1, and a single identically-named
// binary Attribute is returned.
//
// If the CategoricalAttribute has more than two (n) values, the Filter
// generates n BinaryAttributes and sets each of them if the value's observed.
func (b *BinaryConvertFilter) Train() error {
	for _, a := range b.attrs {
		if ac, ok := a.(*base.CategoricalAttribute); ok {
			vals := ac.GetValues()
			if len(vals) <= 2 {
				nAttr := base.NewBinaryAttribute(ac.GetName())
				fAttr := base.FilteredAttribute{ac, nAttr}
				b.converted = append(b.converted, fAttr)
				b.twoValuedCategoricalAttributes[a] = true
			} else {
				if _, ok := b.nValuedCategoricalAttributeMap[a]; !ok {
					b.nValuedCategoricalAttributeMap[a] = make(map[uint64]base.Attribute)
				}
				for i := uint64(0); i < uint64(len(vals)); i++ {
					v := vals[i]
					newName := fmt.Sprintf("%s_%s", ac.GetName(), v)
					newAttr := base.NewBinaryAttribute(newName)
					fAttr := base.FilteredAttribute{ac, newAttr}
					b.converted = append(b.converted, fAttr)
					b.nValuedCategoricalAttributeMap[a][i] = newAttr
				}
			}
		} else if ab, ok := a.(*base.BinaryAttribute); ok {
			fAttr := base.FilteredAttribute{ab, ab}
			b.converted = append(b.converted, fAttr)
		} else if af, ok := a.(*base.FloatAttribute); ok {
			newAttr := base.NewBinaryAttribute(af.GetName())
			fAttr := base.FilteredAttribute{af, newAttr}
			b.converted = append(b.converted, fAttr)
		} else {
			return fmt.Errorf("Unsupported Attribute type: %v", a)
		}
	}
	return nil
}
