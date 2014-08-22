package edf

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestThreadDeserialize(t *testing.T) {
	bytes := []byte{0, 0, 0, 6, 83, 89, 83, 84, 69, 77, 0, 0, 0, 1}
	Convey("Given a byte slice", t, func() {
		var thread Thread
		size := thread.Deserialize(bytes)
		Convey("Decoded name should be SYSTEM", func() {
			So(thread.name, ShouldEqual, "SYSTEM")
		})
		Convey("Size should be the same as the array", func() {
			So(size, ShouldEqual, len(bytes))
		})
	})
}

func TestThreadSerialize(t *testing.T) {
	var thread Thread
	refBytes := []byte{0, 0, 0, 6, 83, 89, 83, 84, 69, 77, 0, 0, 0, 1}
	thread.name = "SYSTEM"
	thread.id = 1
	toBytes := make([]byte, len(refBytes))
	Convey("Should serialize correctly", t, func() {
		thread.Serialize(toBytes)
		So(toBytes, ShouldResemble, refBytes)
	})
}

func TestThreadFindAndWrite(t *testing.T) {
	Convey("Creating a non-existent file should succeed", t, func() {
		tempFile, err := os.OpenFile("hello.db", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0700) //ioutil.TempFile(os.TempDir(), "TestFileCreate")
		So(err, ShouldEqual, nil)
		Convey("Mapping the file should succeed", func() {
			mapping, err := EdfMap(tempFile, EDF_CREATE)
			So(err, ShouldEqual, nil)
			Convey("Writing the thread should succeed", func() {
				t := NewThread(mapping, "MyNameISWhat")
				Convey("Thread number should be 3", func() {
					So(t.id, ShouldEqual, 3)
				})
				Convey("Writing the thread should succeed", func() {
					err := mapping.WriteThread(t)
					So(err, ShouldEqual, nil)
					Convey("Should be able to find the thread again later", func() {
						id, err := mapping.FindThread("MyNameISWhat")
						So(err, ShouldEqual, nil)
						So(id, ShouldEqual, 3)
					})
				})
			})
		})
		os.Remove("hello.db")
	})
}
