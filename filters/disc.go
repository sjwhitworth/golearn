package filters

import (
	"fmt"

	"github.com/amclay/golearn/base"
)

type AbstractDiscretizeFilter struct {
	attrs   map[base.Attribute]bool
	trained bool
	train   base.FixedDataGrid
}

// AddAttribute adds the AttributeSpec of the given attribute `a'
// to the AbstractFloatFilter for discretisation.
func (d *AbstractDiscretizeFilter) AddAttribute(a base.Attribute) error {
	if _, ok := a.(*base.FloatAttribute); !ok {
		return fmt.Errorf("%s is not a FloatAttribute", a)
	}
	_, err := d.train.GetAttribute(a)
	if err != nil {
		return fmt.Errorf("invalid attribute")
	}
	d.attrs[a] = true
	return nil
}

// GetAttributesAfterFiltering gets a list of before/after
// Attributes as base.FilteredAttributes
func (d *AbstractDiscretizeFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	oldAttrs := d.train.AllAttributes()
	ret := make([]base.FilteredAttribute, len(oldAttrs))
	for i, a := range oldAttrs {
		if d.attrs[a] {
			retAttr := new(base.CategoricalAttribute)
			retAttr.SetName(a.GetName())
			ret[i] = base.FilteredAttribute{a, retAttr}
		} else {
			ret[i] = base.FilteredAttribute{a, a}
		}
	}
	return ret
}

func (d *AbstractDiscretizeFilter) getAttributeSpecs() []base.AttributeSpec {
	as := make([]base.AttributeSpec, 0)
	// Set up the AttributeSpecs, and values
	for attr := range d.attrs {
		// If for some reason we've un-added it...
		if !d.attrs[attr] {
			continue
		}
		// Get the AttributeSpec for the training set
		a, err := d.train.GetAttribute(attr)
		if err != nil {
			panic(fmt.Errorf("Attribute resolution error: %s", err))
		}
		// Append to return set
		as = append(as, a)
	}
	return as
}
