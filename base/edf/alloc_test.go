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
		Convey("Mapping the file should succeed", func() {
			mapping, err := edfMap(tempFile, EDF_CREATE)
			So(err, ShouldEqual, nil)
			Convey("Allocation should succeed", func() {
				r, err := mapping.AllocPages(1, 2)
				So(err, ShouldEqual, nil)
				So(r.Start.Byte, ShouldEqual, 4*os.Getpagesize())
				So(r.Start.Segment, ShouldEqual, 0)
				Convey("Unmapping the file should succeed", func() {
					err = mapping.unmap(EDF_UNMAP_SYNC)
					So(err, ShouldEqual, nil)
					Convey("Remapping the file should succeed", func() {
						mapping, err = edfMap(tempFile, EDF_READ_ONLY)
						Convey("Should get the same allocations back", func() {
							rr, err := mapping.getThreadBlocks(2)
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
		Convey("Mapping the file should succeed", func() {
			mapping, err := edfMap(tempFile, EDF_CREATE)
			So(err, ShouldEqual, nil)
			Convey("Allocation of 10 pages should succeed", func() {
				allocated := make([]edfRange, 10)
				for i := 0; i < 10; i++ {
					r, err := mapping.AllocPages(1, 2)
					So(err, ShouldEqual, nil)
					allocated[i] = r
				}
				Convey("Unmapping the file should succeed", func() {
					err = mapping.unmap(EDF_UNMAP_SYNC)
					So(err, ShouldEqual, nil)
					Convey("Remapping the file should succeed", func() {
						mapping, err = edfMap(tempFile, EDF_READ_ONLY)
						Convey("Should get the same allocations back", func() {
							rr, err := mapping.getThreadBlocks(2)
							So(err, ShouldEqual, nil)
							So(len(rr), ShouldEqual, 10)
							So(rr, ShouldResemble, allocated)
						})
					})
				})
			})
		})
	})
}
