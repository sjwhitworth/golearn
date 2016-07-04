package filters

import (
	"fmt"
	"math"

	"github.com/sjwhitworth/golearn/base"
)

// BinningFilter does equal-width binning for numeric
// Attributes (aka "histogram binning")
type BinningFilter struct {
	AbstractDiscretizeFilter
	bins    int
	minVals map[base.Attribute]float64
	maxVals map[base.Attribute]float64
}

// NewBinningFilter creates a BinningFilter structure
// with some helpful default initialisations.
func NewBinningFilter(d base.FixedDataGrid, bins int) *BinningFilter {
	return &BinningFilter{
		AbstractDiscretizeFilter{
			make(map[base.Attribute]bool),
			false,
			d,
		},
		bins,
		make(map[base.Attribute]float64),
		make(map[base.Attribute]float64),
	}
}

func (b *BinningFilter) String() string {
	return fmt.Sprintf("BinningFilter(%d Attribute(s), %d bin(s)", len(b.attrs), b.bins)
}

// Train computes and stores the bin values
// for the training instances.
func (b *BinningFilter) Train() error {

	as := b.getAttributeSpecs()
	// Set up the AttributeSpecs, and values
	for attr := range b.attrs {
		if !b.attrs[attr] {
			continue
		}
		b.minVals[attr] = float64(math.Inf(1))
		b.maxVals[attr] = float64(math.Inf(-1))
	}

	err := b.train.MapOverRows(as, func(row [][]byte, rowNo int) (bool, error) {
		for i, a := range row {
			attr := as[i].GetAttribute()
			attrf := attr.(*base.FloatAttribute)
			val := float64(attrf.GetFloatFromSysVal(a))
			if val > b.maxVals[attr] {
				b.maxVals[attr] = val
			}
			if val < b.minVals[attr] {
				b.minVals[attr] = val
			}
		}
		return true, nil
	})

	if err != nil {
		return fmt.Errorf("Training error: %s", err)
	}
	b.trained = true
	return nil
}

// Transform takes an Attribute and byte sequence and returns
// the transformed byte sequence.
func (b *BinningFilter) Transform(a base.Attribute, n base.Attribute, field []byte) []byte {

	if !b.attrs[a] {
		return field
	}
	af, ok := a.(*base.FloatAttribute)
	if !ok {
		panic("Attribute is the wrong type")
	}
	minVal := b.minVals[a]
	maxVal := b.maxVals[a]
	disc := 0
	// Casts to float64 to replicate a floating point precision error
	delta := float64(maxVal-minVal) / float64(b.bins)
	val := float64(af.GetFloatFromSysVal(field))
	if val <= minVal {
		disc = 0
	} else {
		disc = int(math.Floor(float64(float64(val-minVal)/delta + 0.0001)))
	}
	return base.PackU64ToBytes(uint64(disc))
}

// GetAttributesAfterFiltering gets a list of before/after
// Attributes as base.FilteredAttributes
func (b *BinningFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	oldAttrs := b.train.AllAttributes()
	ret := make([]base.FilteredAttribute, len(oldAttrs))
	for i, a := range oldAttrs {
		if b.attrs[a] {
			retAttr := new(base.CategoricalAttribute)
			minVal := b.minVals[a]
			maxVal := b.maxVals[a]
			delta := float64(maxVal-minVal) / float64(b.bins)
			retAttr.SetName(a.GetName())
			for i := 0; i <= b.bins; i++ {
				floatVal := float64(i)*delta + minVal
				fmtStr := fmt.Sprintf("%%.%df", a.(*base.FloatAttribute).Precision)
				binVal := fmt.Sprintf(fmtStr, floatVal)
				retAttr.GetSysValFromString(binVal)
			}
			ret[i] = base.FilteredAttribute{a, retAttr}
		} else {
			ret[i] = base.FilteredAttribute{a, a}
		}
	}
	return ret
}
