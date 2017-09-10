package base

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"io/ioutil"
	"testing"
)

func TestSerializeToCSV(t *testing.T) {
	Convey("Reading some instances...", t, func() {
		inst, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("Saving the instances to CSV...", func() {
			f, err := ioutil.TempFile("", "instTmpCSV")
			So(err, ShouldBeNil)
			err = SerializeInstancesToCSV(inst, f.Name())
			So(err, ShouldBeNil)
			Convey("What's written out should match what's read in", func() {
				dinst, err := ParseCSVToInstances(f.Name(), true)
				So(err, ShouldBeNil)
				So(InstancesAreEqual(inst, dinst), ShouldBeTrue)
			})
		})
	})
}

func TestCreateAndReadClassifierStub(t *testing.T) {
	Convey("Creating a classifier stub...", t, func() {

		exampleClassifierMetadata := make(map[string]interface{})
		exampleClassifierMetadata["num_trees"] = 4

		metadata := ClassifierMetadataV1{
			FormatVersion:      1,
			ClassifierName:     "test",
			ClassifierVersion:  "1",
			ClassifierMetadata: exampleClassifierMetadata,
		}

		Convey("Saving the classifier...", func() {
			f, err := ioutil.TempFile("", "classTmpF")
			So(err, ShouldBeNil)
			serializer, err := CreateSerializedClassifierStub(f.Name(), metadata)
			So(err, ShouldBeNil)
			err = serializer.Close()
			So(err, ShouldBeNil)
			Convey("Should be able to read the information back...", func() {
				reader, err := ReadSerializedClassifierStub(f.Name())
				So(err, ShouldBeNil)
				So(reader, ShouldNotBeNil)
				So(reader.Metadata.FormatVersion, ShouldEqual, 1)
				So(reader.Metadata.ClassifierName, ShouldEqual, "test")
				So(reader.Metadata.ClassifierVersion, ShouldEqual, "1")
				So(reader.Metadata.ClassifierMetadata["num_trees"], ShouldEqual, 4)
			})
		})
	})
}

func TestSerializeToFile(t *testing.T) {
	Convey("Reading some instances...", t, func() {
		inst, err := ParseCSVToInstances("../examples/datasets/iris_headers.csv", true)
		So(err, ShouldBeNil)

		Convey("Dumping to file...", func() {
			f, err := ioutil.TempFile("", "instTmpF")
			So(err, ShouldBeNil)
			err = SerializeInstances(inst, f)
			So(err, ShouldBeNil)
			f.Seek(0, 0)
			Convey("Contents of the archive should be right...", func() {
				gzr, err := gzip.NewReader(f)
				So(err, ShouldBeNil)
				tr := tar.NewReader(gzr)
				classAttrsPresent := false
				manifestPresent := false
				regularAttrsPresent := false
				dataPresent := false
				dimsPresent := false
				readBytes := make([]byte, len([]byte(SerializationFormatVersion)))
				for {
					hdr, err := tr.Next()
					if err == io.EOF {
						break
					}
					So(err, ShouldBeNil)
					switch hdr.Name {
					case "MANIFEST":
						tr.Read(readBytes)
						manifestPresent = true
						break
					case "CATTRS":
						classAttrsPresent = true
						break
					case "ATTRS":
						regularAttrsPresent = true
						break
					case "DATA":
						dataPresent = true
						break
					case "DIMS":
						dimsPresent = true
						break
					default:
						fmt.Printf("Unknown file: %s\n", hdr.Name)
					}
				}
				Convey("MANIFEST should be present", func() {
					So(manifestPresent, ShouldBeTrue)
					Convey("MANIFEST should be right...", func() {
						So(readBytes, ShouldResemble, []byte(SerializationFormatVersion))
					})

					Convey("DATA should be present", func() {
						So(dataPresent, ShouldBeTrue)
					})
					Convey("ATTRS should be present", func() {
						So(regularAttrsPresent, ShouldBeTrue)
					})
					Convey("CATTRS should be present", func() {
						So(classAttrsPresent, ShouldBeTrue)
					})
					Convey("DIMS should be present", func() {
						So(dimsPresent, ShouldBeTrue)
					})
				})
				Convey("Should be able to reconstruct...", func() {
					f.Seek(0, 0)
					dinst, err := DeserializeInstances(f)
					So(err, ShouldBeNil)
					So(InstancesAreEqual(inst, dinst), ShouldBeTrue)
				})
			})
		})
	})
}
