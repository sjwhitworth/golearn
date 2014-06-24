package edf

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"os"
	"testing"
)

func TestAllocFixed(t *testing.T) {
	Convey("Creating a non-existent file should succeed", t, func() {
		tempFile, err := ioutil.TempFile(os.TempDir(), "TestFileCreate")
		So(err, ShouldEqual, nil)
		Convey("Mapping the file should suceed", func() {
			mapping, err := EdfMap(tempFile, EDF_CREATE)
			So(err, ShouldEqual, nil)
			Convey("Allocation should suceed", func() {
				r, err := mapping.AllocPages(1, 2)
				So(err, ShouldEqual, nil)
				So(r.Start.Byte, ShouldEqual, 4*os.Getpagesize())
				So(r.Start.Segment, ShouldEqual, 0)
				Convey("Unmapping the file should suceed", func() {
					err = mapping.Unmap(EDF_UNMAP_SYNC)
					So(err, ShouldEqual, nil)
					Convey("Remapping the file should suceed", func() {
						mapping, err = EdfMap(tempFile, EDF_READ_ONLY)
						Convey("Should get the same allocations back", func() {
							rr, err := mapping.GetThreadBlocks(2)
							So(err, ShouldEqual, nil)
							So(len(rr), ShouldEqual, 1)
							So(rr[0], ShouldResemble, r)
						})
					})
				})
			})
		})
	})
}

func TestAllocWithExtraContentsBlock(t *testing.T) {
	Convey("Creating a non-existent file should succeed", t, func() {
		tempFile, err := ioutil.TempFile(os.TempDir(), "TestFileCreate")
		So(err, ShouldEqual, nil)
		Convey("Mapping the file should suceed", func() {
			mapping, err := EdfMap(tempFile, EDF_CREATE)
			So(err, ShouldEqual, nil)
			Convey("Allocation of 350 pages should suceed", func() {
				allocated := make([]EdfRange, 350)
				for i := 0; i < 350; i++ {
					r, err := mapping.AllocPages(1, 2)
					So(err, ShouldEqual, nil)
					allocated[i] = r
				}
				Convey("Unmapping the file should suceed", func() {
					err = mapping.Unmap(EDF_UNMAP_SYNC)
					So(err, ShouldEqual, nil)
					Convey("Remapping the file should suceed", func() {
						mapping, err = EdfMap(tempFile, EDF_READ_ONLY)
						Convey("Should get the same allocations back", func() {
							rr, err := mapping.GetThreadBlocks(2)
							So(err, ShouldEqual, nil)
							So(len(rr), ShouldEqual, 350)
							So(rr, ShouldResemble, allocated)
						})
					})
				})
			})
		})
	})
}
