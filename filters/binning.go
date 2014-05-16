package filters

import (
	"fmt"
	base "github.com/sjwhitworth/golearn/base"
	"math"
)

// BinningFilter does equal-width binning for numeric
// Attributes (aka "histogram binning")
type BinningFilter struct {
	Attributes []int
	Instances  *base.Instances
	BinCount   int
	MinVals    map[int]float64
	MaxVals    map[int]float64
	trained    bool
}

// NewBinningFilter creates a BinningFilter structure
// with some helpful default initialisations.
func NewBinningFilter(inst *base.Instances, bins int) BinningFilter {
	return BinningFilter{
		make([]int, 0),
		inst,
		bins,
		make(map[int]float64),
		make(map[int]float64),
		false,
	}
}

// AddAttribute adds the index of the given attribute `a'
// to the BinningFilter for discretisation.
func (b *BinningFilter) AddAttribute(a base.Attribute) {
	attrIndex := b.Instances.GetAttrIndex(a)
	if attrIndex == -1 {
		panic("invalid attribute")
	}
	b.Attributes = append(b.Attributes, attrIndex)
}

// AddAllNumericAttributes adds every suitable attribute
// to the BinningFilter for discretiation
func (b *BinningFilter) AddAllNumericAttributes() {
	for i := 0; i < b.Instances.Cols; i++ {
		if i == b.Instances.ClassIndex {
			continue
		}
		attr := b.Instances.GetAttr(i)
		if attr.GetType() != base.Float64Type {
			continue
		}
		b.Attributes = append(b.Attributes, i)
	}
}

// Build computes and stores the bin values
// for the training instances.
func (b *BinningFilter) Build() {
	for _, attr := range b.Attributes {
		maxVal := math.Inf(-1)
		minVal := math.Inf(1)
		for i := 0; i < b.Instances.Rows; i++ {
			val := b.Instances.Get(i, attr)
			if val > maxVal {
				maxVal = val
			}
			if val < minVal {
				minVal = val
			}
		}
		b.MaxVals[attr] = maxVal
		b.MinVals[attr] = minVal
		b.trained = true
	}
}

// Run applies a trained BinningFilter to a set of Instances,
// discretising any numeric attributes added.
//
// IMPORTANT: Run discretises in-place, so make sure to take
// a copy if the original instances are still needed
//
// IMPORTANT: This function panic()s if the filter has not been
// trained. Call Build() before running this function
//
// IMPORTANT: Call Build() after adding any additional attributes.
// Otherwise, the training structure will be out of date from
// the values expected and could cause a panic.
func (b *BinningFilter) Run(on *base.Instances) {
	if !b.trained {
		panic("Call Build() beforehand")
	}
	for attr := range b.Attributes {
		minVal := b.MinVals[attr]
		maxVal := b.MaxVals[attr]
		disc := 0
		// Casts to float32 to replicate a floating point precision error
		delta := float32(maxVal - minVal)
		delta /= float32(b.BinCount)
		for i := 0; i < on.Rows; i++ {
			val := on.Get(i, attr)
			if val <= minVal {
				disc = 0
			} else {
				disc = int(math.Floor(float64(float32(val-minVal) / delta)))
				if disc >= b.BinCount {
					disc = b.BinCount - 1
				}
			}
			on.Set(i, attr, float64(disc))
		}
		newAttribute := new(base.CategoricalAttribute)
		newAttribute.SetName(on.GetAttr(attr).GetName())
		for i := 0; i < b.BinCount; i++ {
			newAttribute.GetSysValFromString(fmt.Sprintf("%d", i))
		}
		on.ReplaceAttr(attr, newAttribute)
	}
}
