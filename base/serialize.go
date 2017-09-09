package base

import (
	"archive/tar"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
)

const (
	SerializationFormatVersion = "golearn 0.5"
)

func SerializeInstancesToFile(inst FixedDataGrid, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	err = SerializeInstances(inst, f)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return fmt.Errorf("Couldn't flush file: %s", err)
	}
	f.Close()
	return nil
}

func SerializeInstancesToCSV(inst FixedDataGrid, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	defer func() {
		f.Sync()
		f.Close()
	}()

	return SerializeInstancesToCSVStream(inst, f)
}

func SerializeInstancesToCSVStream(inst FixedDataGrid, f io.Writer) error {
	// Create the CSV writer
	w := csv.NewWriter(f)

	colCount, _ := inst.Size()

	// Write out Attribute headers
	// Start with the regular Attributes
	normalAttrs := NonClassAttributes(inst)
	classAttrs := inst.AllClassAttributes()
	allAttrs := make([]Attribute, colCount)
	n := copy(allAttrs, normalAttrs)
	copy(allAttrs[n:], classAttrs)
	headerRow := make([]string, colCount)
	for i, v := range allAttrs {
		headerRow[i] = v.GetName()
	}
	w.Write(headerRow)

	specs := ResolveAttributes(inst, allAttrs)
	curRow := make([]string, colCount)
	inst.MapOverRows(specs, func(row [][]byte, rowNo int) (bool, error) {
		for i, v := range row {
			attr := allAttrs[i]
			curRow[i] = attr.GetStringFromSysVal(v)
		}
		w.Write(curRow)
		return true, nil
	})

	w.Flush()
	return nil
}

func writeAttributesToFilePart(attrs []Attribute, f *tar.Writer, name string) error {
	// Get the marshaled Attribute array
	body, err := json.Marshal(attrs)
	if err != nil {
		return err
	}

	// Write a header
	hdr := &tar.Header{
		Name: name,
		Size: int64(len(body)),
	}
	if err := f.WriteHeader(hdr); err != nil {
		return err
	}

	// Write the marshaled data
	if _, err := f.Write([]byte(body)); err != nil {
		return err
	}

	return nil
}

func getTarContent(tr *tar.Reader, name string) []byte {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		if hdr.Name == name {
			ret := make([]byte, hdr.Size)
			n, err := tr.Read(ret)
			if int64(n) != hdr.Size {
				panic("Size mismatch")
			}
			if err != nil {
				panic(err)
			}
			return ret
		}
	}
	panic("File not found!")
}

func MarshalAttribute(a Attribute) (map[string]interface{}, error) {
	ret := make(map[string]interface{})
	marshaledAttrRaw, err := a.MarshalJSON()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(marshaledAttrRaw, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func DeserializeAttribute(data []byte) (Attribute, error) {
	type JSONAttribute struct {
		Type string          `json:"type"`
		Name string          `json:"name"`
		Attr json.RawMessage `json:"attr"`
	}

	var rawAttr JSONAttribute
	err := json.Unmarshal(data, &rawAttr)
	if err != nil {
		return nil, err
	}
	var attr Attribute

	switch rawAttr.Type {
	case "binary":
		attr = new(BinaryAttribute)
		break
	case "float":
		attr = new(FloatAttribute)
		break
	case "categorical":
		attr = new(CategoricalAttribute)
		break
	default:
		return nil, fmt.Errorf("Unrecognised Attribute format: %s", rawAttr.Type)
	}

	err = attr.UnmarshalJSON(rawAttr.Attr)
	if err != nil {
		return nil, fmt.Errorf("Can't deserialize: %s (error: %s)", rawAttr, err)
	}
	attr.SetName(rawAttr.Name)
	return attr, nil
}

// DeserializeAttributes constructs a ve
func DeserializeAttributes(data []byte) ([]Attribute, error) {

	// Define a JSON shim Attribute
	var attrs []json.RawMessage
	err := json.Unmarshal(data, &attrs)
	if err != nil {
		return nil, fmt.Errorf("Failed to deserialize attributes: %v", err)
	}

	ret := make([]Attribute, len(attrs))
	for i, v := range attrs {
		ret[i], err = DeserializeAttribute(v)
		if err != nil {
			return nil, err
		}
	}

	return ret, nil
}

func DeserializeInstancesFromTarReader(tr *tar.Reader, prefix string) (ret *DenseInstances, err error) {

	p := func(n string) string {
		return fmt.Sprintf("%s%s", prefix, n)
	}

	// Retrieve the MANIFEST and verify
	manifestBytes := getTarContent(tr, p("MANIFEST"))
	if !reflect.DeepEqual(manifestBytes, []byte(SerializationFormatVersion)) {
		return nil, fmt.Errorf("Unsupported MANIFEST: %s", string(manifestBytes))
	}

	// Get the size
	sizeBytes := getTarContent(tr, p("DIMS"))
	attrCount := int(UnpackBytesToU64(sizeBytes[0:8]))
	rowCount := int(UnpackBytesToU64(sizeBytes[8:]))

	// Unmarshal the Attributes
	attrBytes := getTarContent(tr, p("CATTRS"))
	cAttrs, err := DeserializeAttributes(attrBytes)
	if err != nil {
		return nil, err
	}
	attrBytes = getTarContent(tr, p("ATTRS"))
	normalAttrs, err := DeserializeAttributes(attrBytes)
	if err != nil {
		return nil, err
	}

	// Create the return instances
	ret = NewDenseInstances()

	// Normal Attributes first, class Attributes on the end
	allAttributes := make([]Attribute, attrCount)
	for i, v := range normalAttrs {
		ret.AddAttribute(v)
		allAttributes[i] = v
	}
	for i, v := range cAttrs {
		ret.AddAttribute(v)
		err = ret.AddClassAttribute(v)
		if err != nil {
			return nil, fmt.Errorf("Could not set Attribute as class Attribute: %s", err)
		}
		allAttributes[i+len(normalAttrs)] = v
	}
	// Allocate memory
	err = ret.Extend(int(rowCount))
	if err != nil {
		return nil, fmt.Errorf("Could not allocate memory")
	}

	// Seek through the TAR file until we get to the DATA section
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil, fmt.Errorf("DATA section missing!")
		} else if err != nil {
			return nil, fmt.Errorf("Error seeking to DATA section: %s", err)
		}
		if hdr.Name == p("DATA") {
			break
		}
	}

	// Resolve AttributeSpecs
	specs := ResolveAttributes(ret, allAttributes)

	// Finally, read the values out of the data section
	for i := 0; i < rowCount; i++ {
		for _, s := range specs {
			r := ret.Get(s, i)
			n, err := tr.Read(r)
			if n != len(r) {
				return nil, fmt.Errorf("Expected %d bytes (read %d) on row %d", len(r), n, i)
			}
			if err != nil {
				return nil, fmt.Errorf("Read error: %s", err)
			}
			ret.Set(s, i, r)
		}
	}

	return ret, nil
}

func DeserializeInstances(f io.Reader) (ret *DenseInstances, err error) {

	// Recovery function
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
		}
	}()

	// Open the .gz layer
	gzReader, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("Can't open: %s", err)
	}
	// Open the .tar layer
	tr := tar.NewReader(gzReader)

	ret, deSerializeErr := DeserializeInstancesFromTarReader(tr, "")

	if err = gzReader.Close(); err != nil {
		return ret, fmt.Errorf("Error closing gzip stream: %s", err)
	}

	return ret, deSerializeErr
}

// ClassifierMetadataV1 is what gets written into METADATA
// in a classification file format.
type ClassifierMetadataV1 struct {
	// FormatVersion should always be 1 for this structure
	FormatVersion int `json:"format_version"`
	// Uses the classifier name (provided by the classifier)
	ClassifierName string `json:"classifier"`
	// ClassifierVersion is also provided by the classifier
	// and checks whether this version of GoLearn can read what's
	// be written.
	ClassifierVersion string `json"classifier_version"`
	// This is a custom metadata field, provided by the classifier
	ClassifierMetadata map[string]interface{} `json:"classifier_metadata"`
}

type ClassifierDeserializer struct {
	gzipReader io.Reader
	fileReader io.ReadCloser
	tarReader  *tar.Reader
	Metadata   *ClassifierMetadataV1
}

// ReadSerializedClassifierStub is the counterpart of CreateSerializedClassifierStub.
// It's used inside SaveableClassifiers to read information from a perviously saved
// model file.
func ReadSerializedClassifierStub(filePath string) (*ClassifierDeserializer, error) {

	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Can't open file: %s", err)
	}

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("Can't decompress file: %s", err)
	}

	tz := tar.NewReader(gzr)

	// Check that the serialization format is right
	// Retrieve the MANIFEST and verify
	manifestBytes := getTarContent(tz, "MANIFEST")
	if !reflect.DeepEqual(manifestBytes, []byte(SerializationFormatVersion)) {
		return nil, fmt.Errorf("Unsupported MANIFEST: %s", string(manifestBytes))
	}

	//
	// Parse METADATA
	//
	var metadata ClassifierMetadataV1
	ret := &ClassifierDeserializer{
		f,
		gzr,
		tz,
		&metadata,
	}

	err = ret.GetJSONForKey("METADATA", ret.Metadata)
	if err != nil {
		return nil, fmt.Errorf("Error whilst reading METADATA: %s", err)
	}

	// Check that we can understand this archive
	if metadata.FormatVersion != 1 {
		return nil, fmt.Errorf("METADATA: wrong format_version for this version of golearn")
	}

	return ret, nil
}

// GetBytesForKey returns the bytes at a given location in the output.
func (c *ClassifierDeserializer) GetBytesForKey(key string) ([]byte, error) {
	return getTarContent(c.tarReader, key), nil
}

// GetJSONForKey deserializes a JSON key in the output file.
func (c *ClassifierDeserializer) GetJSONForKey(key string, v interface{}) error {
	b, err := c.GetBytesForKey(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}

// GetUInt64ForKey returns a int64 stored at a given key
func (c *ClassifierDeserializer) GetU64ForKey(key string) (uint64, error) {
	b, err := c.GetBytesForKey(key)
	if err != nil {
		return 0, err
	}
	return UnpackBytesToU64(b), nil
}

// Close cleans up everything.
func (c *ClassifierDeserializer) Close() {
	c.fileReader.Close()
}

// ClassifierSerializer is an object used by SaveableClassifiers.
type ClassifierSerializer struct {
	gzipWriter *gzip.Writer
	fileWriter io.WriteCloser
	tarWriter  *tar.Writer
}

// Close finalizes the Classifier serialization session.
func (c *ClassifierSerializer) Close() error {

	// Finally, close and flush the various levels
	if err := c.tarWriter.Flush(); err != nil {
		return fmt.Errorf("Could not flush tar: %s", err)
	}

	if err := c.tarWriter.Close(); err != nil {
		return fmt.Errorf("Could not close tar: %s", err)
	}

	if err := c.gzipWriter.Flush(); err != nil {
		return fmt.Errorf("Could not flush gz: %s", err)
	}

	if err := c.gzipWriter.Close(); err != nil {
		return fmt.Errorf("Could not close gz: %s", err)
	}

	if err := c.fileWriter.Close(); err != nil {
		return fmt.Errorf("Could not close file: %s", err)
	}

	return nil
}

// WriteBytesForKey creates a new entry in the serializer file with some user-defined bytes.
func (c *ClassifierSerializer) WriteBytesForKey(key string, b []byte) error {

	//
	// Write header for key
	//
	hdr := &tar.Header{
		Name: key,
		Size: int64(len(b)),
	}

	if err := c.tarWriter.WriteHeader(hdr); err != nil {
		return fmt.Errorf("Could not write header for '%s': %s", key, err)
	}
	//
	// Write data
	//
	if _, err := c.tarWriter.Write(b); err != nil {
		return fmt.Errorf("Could not write data for '%s': %s", key, err)
	}

	return nil
}

// WriteJSONForKey creates a new entry in the file with an interface serialized as JSON.
func (c *ClassifierSerializer) WriteJSONForKey(key string, v interface{}) error {

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return c.WriteBytesForKey(key, b)

}

func (c *ClassifierSerializer) WriteInstancesForKey(key string, g FixedDataGrid, includeData bool) error {
	return SerializeInstancesToTarWriter(g, c.tarWriter, fmt.Sprintf("%s/"), includeData)
}

// CreateSerializedClassifierStub generates a file to serialize into
// and writes the METADATA header.
func CreateSerializedClassifierStub(filePath string, metadata ClassifierMetadataV1) (*ClassifierSerializer, error) {

	// Open the filePath
	f, err := os.OpenFile(filePath, os.O_RDWR | os.O_TRUNC, 0600)
	if err != nil {
		return nil, nil
	}

	var hdr *tar.Header
	gzWriter := gzip.NewWriter(f)
	tw := tar.NewWriter(gzWriter)

	ret := &ClassifierSerializer{
		gzipWriter: gzWriter,
		fileWriter: f,
		tarWriter: tw,
	}

	//
	// Write the MANIFEST entry
	//
	hdr = &tar.Header{
		Name: "MANIFEST",
		Size: int64(len(SerializationFormatVersion)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, fmt.Errorf("Could not write MANIFEST header: %s", err)
	}

	if _, err := tw.Write([]byte(SerializationFormatVersion)); err != nil {
		return nil, fmt.Errorf("Could not write MANIFEST contents: %s", err)
	}

	//
	// Write the METADATA entry
	//

	// Marshal the classifier information (TODO: split this into another method)
	err = ret.WriteJSONForKey("METADATA", &metadata)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %s", err)
	}

	return ret, nil

}

func SerializeInstances(inst FixedDataGrid, f io.Writer) error {
	// Create a .tar.gz container
	gzWriter := gzip.NewWriter(f)
	tw := tar.NewWriter(gzWriter)

	serializeErr := SerializeInstancesToTarWriter(inst, tw, "", true)
	// Finally, close and flush the various levels
	if err := tw.Flush(); err != nil {
		return fmt.Errorf("Could not flush tar: %s", err)
	}

	if err := tw.Close(); err != nil {
		return fmt.Errorf("Could not close tar: %s", err)
	}

	if err := gzWriter.Flush(); err != nil {
		return fmt.Errorf("Could not flush gz: %s", err)
	}

	if err := gzWriter.Close(); err != nil {
		return fmt.Errorf("Could not close gz: %s", err)
	}

	return serializeErr
}

func SerializeInstancesToTarWriter(inst FixedDataGrid, tw *tar.Writer, prefix string, includeData bool) error {
	var hdr *tar.Header

	p := func(n string) string {
		return fmt.Sprintf("%s%s", prefix, n)
	}

	// Write the MANIFEST entry
	hdr = &tar.Header{
		Name: p("MANIFEST"),
		Size: int64(len(SerializationFormatVersion)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("Could not write MANIFEST header: %s", err)
	}

	if _, err := tw.Write([]byte(SerializationFormatVersion)); err != nil {
		return fmt.Errorf("Could not write MANIFEST contents: %s", err)
	}

	// Now write the dimensions of the dataset
	attrCount, rowCount := inst.Size()
	hdr = &tar.Header{
		Name: p("DIMS"),
		Size: 16,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("Could not write DIMS header: %s", err)
	}

	if _, err := tw.Write(PackU64ToBytes(uint64(attrCount))); err != nil {
		return fmt.Errorf("Could not write DIMS (attrCount): %s", err)
	}
	if _, err := tw.Write(PackU64ToBytes(uint64(rowCount))); err != nil {
		return fmt.Errorf("Could not write DIMS (rowCount): %s", err)
	}

	// Write the ATTRIBUTES files
	classAttrs := inst.AllClassAttributes()
	normalAttrs := NonClassAttributes(inst)
	if err := writeAttributesToFilePart(classAttrs, tw, p("CATTRS")); err != nil {
		return fmt.Errorf("Could not write CATTRS: %s", err)
	}
	if err := writeAttributesToFilePart(normalAttrs, tw, p("ATTRS")); err != nil {
		return fmt.Errorf("Could not write ATTRS: %s", err)
	}

	if !includeData {
		return nil
	}

	// Data must be written out in the same order as the Attributes
	allAttrs := make([]Attribute, attrCount)
	normCount := copy(allAttrs, normalAttrs)
	for i, v := range classAttrs {
		allAttrs[normCount+i] = v
	}

	allSpecs := ResolveAttributes(inst, allAttrs)

	// First, estimate the amount of data we'll need...
	dataLength := int64(0)
	inst.MapOverRows(allSpecs, func(val [][]byte, row int) (bool, error) {
		for _, v := range val {
			dataLength += int64(len(v))
		}
		return true, nil
	})

	// Then write the header
	hdr = &tar.Header{
		Name: p("DATA"),
		Size: dataLength,
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return fmt.Errorf("Could not write DATA: %s", err)
	}

	// Then write the actual data
	writtenLength := int64(0)
	if err := inst.MapOverRows(allSpecs, func(val [][]byte, row int) (bool, error) {
		for _, v := range val {
			wl, err := tw.Write(v)
			writtenLength += int64(wl)
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}); err != nil {
		return err
	}

	if writtenLength != dataLength {
		return fmt.Errorf("Could not write DATA: changed size from %v to %v", dataLength, writtenLength)
	}

	return nil
}
