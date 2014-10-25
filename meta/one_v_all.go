package meta

import (
	"fmt"
	"github.com/sjwhitworth/golearn/base"
)

// OneVsAllModel replaces class Attributes with numeric versions
// and trains n wrapped classifiers. The actual class is chosen
// by whichever is most confident. Only one CategoricalAttribute
// class variable is supported.
type OneVsAllModel struct {
	NewClassifierFunction func(string) base.Classifier
	filters               []*oneVsAllFilter
	classifiers           []base.Classifier
	maxClassVal           uint64
}

// NewOneVsAllModel creates a new OneVsAllModel. The argument
// must be a function which returns a base.Classifier ready for training.
func NewOneVsAllModel(f func(string) base.Classifier) *OneVsAllModel {
	return &OneVsAllModel{
		f,
		nil,
		nil,
		0,
	}
}

func (m *OneVsAllModel) generateAttributes(from base.FixedDataGrid) map[base.Attribute]base.Attribute {
	attrs := from.AllAttributes()
	classAttrs := from.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic("Only 1 class Attribute is supported!")
	}
	ret := make(map[base.Attribute]base.Attribute)
	for _, a := range attrs {
		ret[a] = a
		for _, b := range classAttrs {
			if a.Equals(b) {
				cur := base.NewFloatAttribute(b.GetName())
				ret[a] = cur
			}
		}
	}
	return ret
}

// Fit creates n filtered datasets (where n is the number of values
// a CategoricalAttribute can take) and uses them to train the
// underlying classifiers.
func (m *OneVsAllModel) Fit(using base.FixedDataGrid) {
	var classAttr *base.CategoricalAttribute
	// Do some validation
	classAttrs := using.AllClassAttributes()
	for _, a := range classAttrs {
		if c, ok := a.(*base.CategoricalAttribute); !ok {
			panic("Unsupported ClassAttribute type")
		} else {
			classAttr = c
		}
	}
	attrs := m.generateAttributes(using)

	// Find the highest stored value
	val := uint64(0)
	classVals := classAttr.GetValues()
	for _, s := range classVals {
		cur := base.UnpackBytesToU64(classAttr.GetSysValFromString(s))
		if cur > val {
			val = cur
		}
	}
	if val == 0 {
		panic("Must have more than one class!")
	}
	m.maxClassVal = val

	// Create individual filtered instances for training
	filters := make([]*oneVsAllFilter, val+1)
	classifiers := make([]base.Classifier, val+1)
	for i := uint64(0); i <= val; i++ {
		f := &oneVsAllFilter{
			attrs,
			classAttr,
			i,
		}
		filters[i] = f
		classifiers[i] = m.NewClassifierFunction(classVals[int(i)])
		classifiers[i].Fit(base.NewLazilyFilteredInstances(using, f))
	}

	m.filters = filters
	m.classifiers = classifiers
}

// Predict issues predictions. Each class-specific classifier is expected
// to output a value between 0 (indicating that a given instance is not
// a given class) and 1 (indicating that the given instance is definitely
// that class). For each instance, the class with the highest value is chosen.
// The result is undefined if several underlying models output the same value.
func (m *OneVsAllModel) Predict(what base.FixedDataGrid) (base.FixedDataGrid, error) {
	ret := base.GeneratePredictionVector(what)
	vecs := make([]base.FixedDataGrid, m.maxClassVal+1)
	specs := make([]base.AttributeSpec, m.maxClassVal+1)
	for i := uint64(0); i <= m.maxClassVal; i++ {
		f := m.filters[i]
		c := base.NewLazilyFilteredInstances(what, f)
		p, err := m.classifiers[i].Predict(c)
		if err != nil {
			return nil, err
		}
		vecs[i] = p
		specs[i] = base.ResolveAttributes(p, p.AllClassAttributes())[0]
	}
	_, rows := ret.Size()
	spec := base.ResolveAttributes(ret, ret.AllClassAttributes())[0]
	for i := 0; i < rows; i++ {
		class := uint64(0)
		best := 0.0
		for j := uint64(0); j <= m.maxClassVal; j++ {
			val := base.UnpackBytesToFloat(vecs[j].Get(specs[j], i))
			if val > best {
				class = j
				best = val
			}
		}
		ret.Set(spec, i, base.PackU64ToBytes(class))
	}
	return ret, nil
}

//
// Filter implementation
//
type oneVsAllFilter struct {
	attrs        map[base.Attribute]base.Attribute
	classAttr    base.Attribute
	classAttrVal uint64
}

func (f *oneVsAllFilter) AddAttribute(a base.Attribute) error {
	return fmt.Errorf("Not supported")
}

func (f *oneVsAllFilter) GetAttributesAfterFiltering() []base.FilteredAttribute {
	ret := make([]base.FilteredAttribute, len(f.attrs))
	cnt := 0
	for i := range f.attrs {
		ret[cnt] = base.FilteredAttribute{i, f.attrs[i]}
		cnt++
	}
	return ret
}

func (f *oneVsAllFilter) String() string {
	return "oneVsAllFilter"
}

func (f *oneVsAllFilter) Transform(old, to base.Attribute, seq []byte) []byte {
	if !old.Equals(f.classAttr) {
		return seq
	}
	val := base.UnpackBytesToU64(seq)
	if val == f.classAttrVal {
		return base.PackFloatToBytes(1.0)
	}
	return base.PackFloatToBytes(0.0)
}

func (f *oneVsAllFilter) Train() error {
	return fmt.Errorf("Unsupported")
}
