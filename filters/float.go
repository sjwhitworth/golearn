package filters

import (
	"fmt"

	"github.com/amclay/golearn/base"
)

// FloatConvertFilters convert a given DataGrid into one which
// only contains BinaryAttributes.
//
// FloatAttributes are discretised into either 0 (if the value is 0)
// or 1 (if the value is not 0).
//
// CategoricalAttributes are discretised into one or more new
// BinaryAttributes.
type FloatConvertFilter struct {
	attrs                          []base.Attribute
	converted                      []base.FilteredAttribute
	twoValuedCategoricalAttributes map[base.Attribute]bool // Two-valued categorical Attributes
	nValuedCategoricalAttributeMap map[base.Attribute]map[uint64]base.Attribute
}

// NewFloatConvertFilter creates a blank FloatConvertFilter
func NewFloatConvertFilter() *FloatConvertFilter {
	ret := &FloatConvertFilter{
		make([]base.Attribute, 0),
		make([]base.FilteredAttribute, 0),
		make(map[base.Attribute]bool),
		make(map[base.Attribute]map[uint64]base.Attribute),
	}
	return ret
}

// AddAttribute adds a new Attribute to this Filter
func (f *FloatConvertFilter) AddAttribute(a base.Attribute) error {
	f.attrs = append(f.attrs, a)
	return nil
}

// GetAttributesAfterFiltering returns the Attributes previously computed via Train()
func (f *FloatConvertFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	return f.converted
}

// String gets a human-readable string
func (f *FloatConvertFilter) String() string {
	return fmt.Sprintf("FloatConvertFilter(%d Attribute(s))", len(f.attrs))
}

// Transform converts the given byte sequence using the old Attribute into the new
// byte sequence.

func (f *FloatConvertFilter) Transform(a base.Attribute, n base.Attribute, attrBytes []byte) []byte {
	ret := make([]byte, 8)
	// Check for CategoricalAttribute
	if _, ok := a.(*base.CategoricalAttribute); ok {
		// Unpack byte value
		val := base.UnpackBytesToU64(attrBytes)
		// If it's a two-valued one, check for non-zero
		if f.twoValuedCategoricalAttributes[a] {
			if val > 0 {
				ret = base.PackFloatToBytes(1.0)
			} else {
				ret = base.PackFloatToBytes(0.0)
			}
		} else if an, ok := f.nValuedCategoricalAttributeMap[a]; ok {
			// If it's an n-valued one, check the new Attribute maps onto
			// the unpacked value
			if af, ok := an[val]; ok {
				if af.Equals(n) {
					ret = base.PackFloatToBytes(1.0)
				} else {
					ret = base.PackFloatToBytes(0.0)
				}
			} else {
				panic("Categorical value not defined!")
			}
		} else {
			panic(fmt.Sprintf("Not a recognised Attribute %v", a))
		}
	} else if _, ok := a.(*base.FloatAttribute); ok {
		// Binary: just return the original value
		ret = attrBytes
	} else if _, ok := a.(*base.BinaryAttribute); ok {
		// Float: check for non-zero
		if attrBytes[0] > 0 {
			ret = base.PackFloatToBytes(1.0)
		} else {
			ret = base.PackFloatToBytes(0.0)
		}
	} else {
		panic(fmt.Sprintf("Unrecognised Attribute: %v", a))
	}
	return ret
}

// Train converts the Attributes into equivalently named FloatAttributes,
// leaves FloatAttributes unmodified and processes
// CategoricalAttributes as follows.
//
// If the CategoricalAttribute has two values, one of them is
// designated 0.0 and the other 1.0, and a single identically-named
// FloatAttribute is returned.
//
// If the CategoricalAttribute has more than two (n) values, the Filter
// generates n FloatAttributes and sets each of them if the value's observed.
func (f *FloatConvertFilter) Train() error {
	for _, a := range f.attrs {
		if ac, ok := a.(*base.CategoricalAttribute); ok {
			vals := ac.GetValues()
			if len(vals) <= 2 {
				nAttr := base.NewFloatAttribute(ac.GetName())
				fAttr := base.FilteredAttribute{ac, nAttr}
				f.converted = append(f.converted, fAttr)
				f.twoValuedCategoricalAttributes[a] = true
			} else {
				if _, ok := f.nValuedCategoricalAttributeMap[a]; !ok {
					f.nValuedCategoricalAttributeMap[a] = make(map[uint64]base.Attribute)
				}
				for i := uint64(0); i < uint64(len(vals)); i++ {
					v := vals[i]
					newName := fmt.Sprintf("%s_%s", ac.GetName(), v)
					newAttr := base.NewFloatAttribute(newName)
					fAttr := base.FilteredAttribute{ac, newAttr}
					f.converted = append(f.converted, fAttr)
					f.nValuedCategoricalAttributeMap[a][i] = newAttr
				}
			}
		} else if ab, ok := a.(*base.FloatAttribute); ok {
			fAttr := base.FilteredAttribute{ab, ab}
			f.converted = append(f.converted, fAttr)
		} else if af, ok := a.(*base.BinaryAttribute); ok {
			newAttr := base.NewFloatAttribute(af.GetName())
			fAttr := base.FilteredAttribute{af, newAttr}
			f.converted = append(f.converted, fAttr)
		} else {
			return fmt.Errorf("Unsupported Attribute type: %v", a)
		}
	}
	return nil
}
