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
	fitOn                 base.FixedDataGrid
	classValues           []string
}

// NewOneVsAllModel creates a new OneVsAllModel. The argument
// must be a function which returns a base.Classifier ready for training.
func NewOneVsAllModel(f func(string) base.Classifier) *OneVsAllModel {
	return &OneVsAllModel{
		f,
		nil,
		nil,
		0,
		nil,
		nil,
	}
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

	// If we're reloading, we may just be fitting to the structure
	_, srcRows := using.Size()
	fittingToStructure := srcRows == 0

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
		if !fittingToStructure {
			classifiers[i].Fit(base.NewLazilyFilteredInstances(using, f))
		}
	}

	m.filters = filters
	m.classifiers = classifiers
	m.fitOn = base.NewStructuralCopy(using)
	m.classValues = classVals
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

	if int(m.maxClassVal) > len(m.filters) || (m.maxClassVal == 0 && len(m.filters) == 0) {
		return nil, base.WrapError(fmt.Errorf("Internal error: m.Filter len = %d, maxClassVal = %d", len(m.filters), m.maxClassVal))
	}

	for i := uint64(0); i <= m.maxClassVal; i++ {
		//log.Printf("i = %d, m.Filter len = %d, maxClassVal = %d", i, len(m.filters), m.maxClassVal)
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

func (m *OneVsAllModel) Load(filePath string) error {
	reader, err := base.ReadSerializedClassifierStub(filePath)
	if err != nil {
		return err
	}

	err = m.LoadWithPrefix(reader, "")
	reader.Close()
	return err
}

func (m *OneVsAllModel) LoadWithPrefix(reader *base.ClassifierDeserializer, prefix string) error {

	pI := func(n string, i int) string {
		return reader.Prefix(prefix, reader.Prefix(n, fmt.Sprintf("%d", i)))
	}

	// Reload the instances
	fitOn, err := reader.GetInstancesForKey(reader.Prefix(prefix, "INSTANCE_STRUCTURE"))
	if err != nil {
		return base.DescribeError("Can't load INSTANCE_STRUCTURE", err)
	}
	m.Fit(fitOn)
	/*if err != nil {
		base.DescribeError("Could not fit reloaded classifier to the structure", err)
	}*/

	// Reload the filters
	numFiltersU64, err := reader.GetU64ForKey(reader.Prefix(prefix, "FILTER_COUNT"))
	if err != nil {
		return base.DescribeError("Can't load FILTER_COUNT", err)
	}
	m.filters = make([]*oneVsAllFilter, 0)
	numFilters := int(numFiltersU64)
	for i := 0; i < numFilters; i++ {
		f := oneVsAllFilter{}

		mapPrefix := pI("FILTER", i)
		mapCountKey := reader.Prefix(mapPrefix, "COUNT")
		numAttrsInMapU64, err := reader.GetU64ForKey(mapCountKey)
		if err != nil {
			return base.FormatError(err, "Unable to read %s", mapCountKey)
		}

		attrMap := make(map[base.Attribute]base.Attribute)

		for j := 0; j < int(numAttrsInMapU64); j++ {
			mapTupleKey := reader.Prefix(mapPrefix, fmt.Sprintf("%d", j))
			mapKeyKeyKey := reader.Prefix(mapTupleKey, "KEY")
			mapKeyValKey := reader.Prefix(mapTupleKey, "VAL")

			keyAttrRaw, err := reader.GetAttributeForKey(mapKeyKeyKey)
			if err != nil {
				return base.FormatError(err, "Unable to read Attr from %s", mapKeyKeyKey)
			}
			valAttrRaw, err := reader.GetAttributeForKey(mapKeyValKey)
			if err != nil {
				return base.FormatError(err, "Unable to read Attr from %s", mapKeyValKey)
			}

			keyAttr, err := base.ReplaceDeserializedAttributeWithVersionFromInstances(keyAttrRaw, m.fitOn)
			if err != nil {
				return base.FormatError(err, "Can't resolve this attribute: %s", keyAttrRaw)
			}

			attrMap[keyAttr] = valAttrRaw
		}
		f.attrs = attrMap
		mapClassKey := reader.Prefix(mapPrefix, "CLASS_ATTR")
		classAttrRaw, err := reader.GetAttributeForKey(mapClassKey)
		if err != nil {
			return base.FormatError(err, "Can't read from: %s", mapClassKey)
		}
		classAttr, err := base.ReplaceDeserializedAttributeWithVersionFromInstances(classAttrRaw, m.fitOn)
		if err != nil {
			return base.FormatError(err, "Can't resolve: %s", classAttr)
		}
		f.classAttr = classAttr

		classAttrValKey := reader.Prefix(mapPrefix, "CLASS_VAL")
		classVal, err := reader.GetU64ForKey(classAttrValKey)
		if err != nil {
			return base.FormatError(err, "Can't read from: %s", classAttrValKey)
		}
		f.classAttrVal = classVal
		m.filters = append(m.filters, &f)
	}
	// Reload the class values
	var classVals = make([]string, 0)
	err = reader.GetJSONForKey(reader.Prefix(prefix, "CLASS_VALUES"), &classVals)
	if err != nil {
		return base.DescribeError("Can't read CLASS_VALUES", err)
	}
	m.classValues = classVals

	// Reload the classifiers
	m.classifiers = make([]base.Classifier, 0)
	for i, c := range classVals {
		cls := m.NewClassifierFunction(c)
		clsPrefix := pI("CLASSIFIERS", i)

		err = cls.LoadWithPrefix(reader, clsPrefix)
		if err != nil {
			return base.FormatError(err, "Could not reload classifier at: %s", clsPrefix)
		}
		m.classifiers = append(m.classifiers, cls)
	}

	return nil
}

func (m *OneVsAllModel) GetMetadata() base.ClassifierMetadataV1 {
	return base.ClassifierMetadataV1{
		FormatVersion:      1,
		ClassifierName:     "OneVsAllModel",
		ClassifierVersion:  "1.0",
		ClassifierMetadata: nil,
	}
}

func (m *OneVsAllModel) Save(filePath string) error {
	writer, err := base.CreateSerializedClassifierStub(filePath, m.GetMetadata())
	if err != nil {
		return err
	}
	return m.SaveWithPrefix(writer, "")
}

func (m *OneVsAllModel) SaveWithPrefix(writer *base.ClassifierSerializer, prefix string) error {

	pI := func(n string, i int) string {
		return writer.Prefix(prefix, writer.Prefix(n, fmt.Sprintf("%d", i)))
	}

	// Save the instances
	err := writer.WriteInstancesForKey(writer.Prefix(prefix, "INSTANCE_STRUCTURE"), m.fitOn, false)
	if err != nil {
		return base.DescribeError("Unable to write INSTANCE_STRUCTURE", err)
	}

	// Write the class values
	err = writer.WriteJSONForKey(writer.Prefix(prefix, "CLASS_VALUES"), m.classValues)
	if err != nil {
		return base.DescribeError("Can't write CLASS_VALUES", err)
	}

	// Save the filters
	err = writer.WriteU64ForKey(writer.Prefix(prefix, "FILTER_COUNT"), uint64(len(m.filters)))
	if err != nil {
		return base.DescribeError("Unable to write FILTER_COUNT", err)
	}
	for i, f := range m.filters {
		mapPrefix := pI("FILTER", i)
		mapCountKey := writer.Prefix(mapPrefix, "COUNT")
		err := writer.WriteU64ForKey(mapCountKey, uint64(len(f.attrs)))
		if err != nil {
			return base.DescribeError("Unable to write the size of the filter map", err)
		}
		j := 0
		for key := range f.attrs {
			mapTupleKey := writer.Prefix(mapPrefix, fmt.Sprintf("%d", j))
			mapKeyKeyKey := writer.Prefix(mapTupleKey, "KEY")
			mapKeyValKey := writer.Prefix(mapTupleKey, "VAL")

			err = writer.WriteAttributeForKey(mapKeyKeyKey, key)
			if err != nil {
				return base.DescribeError("Unable to write filter map key", err)
			}
			err = writer.WriteAttributeForKey(mapKeyValKey, f.attrs[key])
			if err != nil {
				return base.DescribeError("Unable to write filter map value", err)
			}
			j++
		}
		mapClassKey := writer.Prefix(mapPrefix, "CLASS_ATTR")
		err = writer.WriteAttributeForKey(mapClassKey, f.classAttr)
		if err != nil {
			return base.DescribeError("Unable to write CLASS_ATTR", err)
		}
		classAttrValKey := writer.Prefix(mapPrefix, "CLASS_VAL")
		err = writer.WriteU64ForKey(classAttrValKey, f.classAttrVal)
		if err != nil {
			return base.DescribeError("Can't write CLASS_VAL", err)
		}
	}

	// Save the classifiers
	for i, c := range m.classifiers {
		clsPrefix := pI("CLASSIFIERS", i)
		err = c.SaveWithPrefix(writer, clsPrefix)
		if err != nil {
			return base.FormatError(err, "Can't save classifier for class %s", m.classValues[i])
		}
	}

	return writer.Close()
}

func (m *OneVsAllModel) generateAttributes(from base.FixedDataGrid) map[base.Attribute]base.Attribute {
	attrs := from.AllAttributes()
	classAttrs := from.AllClassAttributes()
	if len(classAttrs) != 1 {
		panic(fmt.Errorf("Only 1 class Attribute is supported, had %d", len(classAttrs)))
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
