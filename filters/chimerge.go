package filters

import (
	"fmt"
	"math"

	"github.com/amclay/golearn/base"
)

// ChiMergeFilter implements supervised discretisation
// by merging successive numeric intervals if the difference
// in their class distribution is not statistically signficant.
// See Bramer, "Principles of Data Mining", 2nd Edition
//  pp 105--115
type ChiMergeFilter struct {
	AbstractDiscretizeFilter
	tables       map[base.Attribute][]*FrequencyTableEntry
	Significance float64
	MinRows      int
	MaxRows      int
}

// NewChiMergeFilter creates a ChiMergeFilter with some helpful intialisations.
func NewChiMergeFilter(d base.FixedDataGrid, significance float64) *ChiMergeFilter {
	_, rows := d.Size()
	return &ChiMergeFilter{
		AbstractDiscretizeFilter{
			make(map[base.Attribute]bool),
			false,
			d,
		},
		make(map[base.Attribute][]*FrequencyTableEntry),
		significance,
		2,
		rows,
	}
}

// Train computes and stores the
// Produces a value mapping table
//   inst: The base.Instances which need discretising
//   sig:  The significance level (e.g. 0.95)
//   minrows: The minimum number of rows required in the frequency table
//   maxrows: The maximum number of rows allowed in the frequency table
//            If the number of rows is above this, statistically signficant
//            adjacent rows will be merged
//   precision: internal number of decimal places to round E value to
//              (useful for verification)
func chiMerge(inst base.FixedDataGrid, attr base.Attribute, sig float64, minrows int, maxrows int) []*FrequencyTableEntry {

	// Parameter sanity checking
	if !(2 <= minrows) {
		minrows = 2
	}
	if !(minrows < maxrows) {
		maxrows = minrows + 1
	}
	if sig == 0 {
		sig = 10
	}

	// Check that the attribute is numeric
	_, ok := attr.(*base.FloatAttribute)
	if !ok {
		panic("only use Chi-M on numeric stuff")
	}

	// Build a frequency table
	freq := ChiMBuildFrequencyTable(attr, inst)
	// Count the number of classes
	classes := chiCountClasses(freq)
	for {
		if len(freq) <= minrows {
			break
		}
		minChiVal := math.Inf(1)
		// There may be more than one index to merge
		minChiIndexes := make([]int, 0)
		for i := 0; i < len(freq)-1; i++ {
			chiVal := chiComputeStatistic(freq[i], freq[i+1])
			if chiVal < minChiVal {
				minChiVal = chiVal
				minChiIndexes = make([]int, 0)
			}
			if chiVal == minChiVal {
				minChiIndexes = append(minChiIndexes, i)
			}
		}
		// Only merge if:
		//  We're above the maximum number of rows
		//  OR the chiVal is significant
		//   AS LONG AS we're above the minimum row count
		merge := false
		if len(freq) > maxrows {
			merge = true
		}
		// Compute the degress of freedom |classes - 1| * |rows - 1|
		degsOfFree := len(classes) - 1
		sigVal := chiSquaredPercentile(degsOfFree, minChiVal)
		if sigVal < sig {
			merge = true
		}
		// If we don't need to merge, then break
		if !merge {
			break
		}
		// Otherwise merge the rows i, i+1 by taking
		//  The higher of the two things as the value
		//  Combining the class frequencies
		for i, v := range minChiIndexes {
			freq = chiMergeMergeZipAdjacent(freq, v-i)
		}
	}
	return freq
}

func (c *ChiMergeFilter) Train() error {
	as := c.getAttributeSpecs()

	for _, a := range as {

		attr := a.GetAttribute()

		// Skip if not set
		if !c.attrs[attr] {
			continue
		}

		// Build sort order
		sortOrder := []base.AttributeSpec{a}

		// Sort
		sorted, err := base.LazySort(c.train, base.Ascending, sortOrder)
		if err != nil {
			panic(err)
		}

		// Perform ChiMerge
		freq := chiMerge(sorted, attr, c.Significance, c.MinRows, c.MaxRows)
		c.tables[attr] = freq
	}

	return nil
}

// GetAttributesAfterFiltering gets a list of before/after
// Attributes as base.FilteredAttributes
func (c *ChiMergeFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	oldAttrs := c.train.AllAttributes()
	ret := make([]base.FilteredAttribute, len(oldAttrs))
	for i, a := range oldAttrs {
		if c.attrs[a] {
			retAttr := new(base.CategoricalAttribute)
			retAttr.SetName(a.GetName())
			for _, k := range c.tables[a] {
				retAttr.GetSysValFromString(fmt.Sprintf("%f", k.Value))
			}
			ret[i] = base.FilteredAttribute{a, retAttr}
		} else {
			ret[i] = base.FilteredAttribute{a, a}
		}
	}
	return ret
}

// Transform returns the byte sequence after discretisation
func (c *ChiMergeFilter) Transform(a base.Attribute, n base.Attribute, field []byte) []byte {
	// Do we use this Attribute?
	if !c.attrs[a] {
		return field
	}
	// Find the Attribute value in the table
	table := c.tables[a]
	dis := 0
	val := base.UnpackBytesToFloat(field)
	for j, k := range table {
		if k.Value < val {
			dis = j
			continue
		}
		break
	}

	return base.PackU64ToBytes(uint64(dis))
}

func (c *ChiMergeFilter) String() string {
	return fmt.Sprintf("ChiMergeFilter(%d Attributes, %.2f Significance)", len(c.tables), c.Significance)
}
