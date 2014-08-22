package filters

import (
	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBinaryFilterClassPreservation(t *testing.T) {
	Convey("Given a contrived dataset...", t, func() {
		// Read the contrived dataset
		inst, err := base.ParseCSVToInstances("./binary_test.csv", true)
		So(err, ShouldEqual, nil)

		// Add all Attributes to the filter
		bFilt := NewBinaryConvertFilter()
		bAttrs := inst.AllAttributes()
		for _, a := range bAttrs {
			bFilt.AddAttribute(a)
		}
		bFilt.Train()

		// Construct a LazilyFilteredInstances to handle it
		instF := base.NewLazilyFilteredInstances(inst, bFilt)

		Convey("All the expected class Attributes should be present if discretised...", func() {
			attrMap := make(map[string]bool)
			attrMap["arbitraryClass_hi"] = false
			attrMap["arbitraryClass_there"] = false
			attrMap["arbitraryClass_world"] = false

			for _, a := range instF.AllClassAttributes() {
				attrMap[a.GetName()] = true
			}

			So(attrMap["arbitraryClass_hi"], ShouldEqual, true)
			So(attrMap["arbitraryClass_there"], ShouldEqual, true)
			So(attrMap["arbitraryClass_world"], ShouldEqual, true)
		})
	})
}

func TestBinaryFilter(t *testing.T) {

	Convey("Given a contrived dataset...", t, func() {

		// Read the contrived dataset
		inst, err := base.ParseCSVToInstances("./binary_test.csv", true)
		So(err, ShouldEqual, nil)

		// Add Attributes to the filter
		bFilt := NewBinaryConvertFilter()
		bAttrs := base.NonClassAttributes(inst)
		for _, a := range bAttrs {
			bFilt.AddAttribute(a)
		}
		bFilt.Train()

		// Construct a LazilyFilteredInstances to handle it
		instF := base.NewLazilyFilteredInstances(inst, bFilt)

		Convey("All the non-class Attributes should be binary...", func() {
			// Check that all the Attributes are the right type
			for _, a := range base.NonClassAttributes(instF) {
				_, ok := a.(*base.BinaryAttribute)
				So(ok, ShouldEqual, true)
			}
		})

		// Check that all the class Attributes made it
		Convey("All the class Attributes should have survived...", func() {
			origClassAttrs := inst.AllClassAttributes()
			newClassAttrs := instF.AllClassAttributes()
			intersectClassAttrs := base.AttributeIntersect(origClassAttrs, newClassAttrs)
			So(len(intersectClassAttrs), ShouldEqual, len(origClassAttrs))
		})
		// Check that the Attributes have the right names
		Convey("Attribute names should be correct...", func() {
			origNames := []string{"floatAttr", "shouldBe1Binary",
				"shouldBe3Binary_stoicism", "shouldBe3Binary_heroism",
				"shouldBe3Binary_romanticism", "arbitraryClass"}
			origMap := make(map[string]bool)
			for _, a := range origNames {
				origMap[a] = false
			}
			for _, a := range instF.AllAttributes() {
				name := a.GetName()
				_, ok := origMap[name]
				So(ok, ShouldBeTrue)
				origMap[name] = true
			}
			for a := range origMap {
				So(origMap[a], ShouldEqual, true)
			}
		})

		// Check that the Attributes have been discretised correctly
		Convey("Discretisation should have worked", func() {
			// Build Attribute map
			attrMap := make(map[string]base.Attribute)
			for _, a := range instF.AllAttributes() {
				attrMap[a.GetName()] = a
			}
			// For each attribute
			for name := range attrMap {
				So(name, ShouldBeIn, []string{
					"floatAttr",
					"shouldBe1Binary",
					"shouldBe3Binary_stoicism",
					"shouldBe3Binary_heroism",
					"shouldBe3Binary_romanticism",
					"arbitraryClass",
				})

				attr := attrMap[name]
				as, err := instF.GetAttribute(attr)
				So(err, ShouldEqual, nil)

				if name == "floatAttr" {
					So(instF.Get(as, 0), ShouldResemble, []byte{1})
					So(instF.Get(as, 1), ShouldResemble, []byte{1})
					So(instF.Get(as, 2), ShouldResemble, []byte{0})
				} else if name == "shouldBe1Binary" {
					So(instF.Get(as, 0), ShouldResemble, []byte{0})
					So(instF.Get(as, 1), ShouldResemble, []byte{1})
					So(instF.Get(as, 2), ShouldResemble, []byte{1})
				} else if name == "shouldBe3Binary_stoicism" {
					So(instF.Get(as, 0), ShouldResemble, []byte{1})
					So(instF.Get(as, 1), ShouldResemble, []byte{0})
					So(instF.Get(as, 2), ShouldResemble, []byte{0})
				} else if name == "shouldBe3Binary_heroism" {
					So(instF.Get(as, 0), ShouldResemble, []byte{0})
					So(instF.Get(as, 1), ShouldResemble, []byte{1})
					So(instF.Get(as, 2), ShouldResemble, []byte{0})
				} else if name == "shouldBe3Binary_romanticism" {
					So(instF.Get(as, 0), ShouldResemble, []byte{0})
					So(instF.Get(as, 1), ShouldResemble, []byte{0})
					So(instF.Get(as, 2), ShouldResemble, []byte{1})
				} else if name == "arbitraryClass" {
				}
			}
		})

	})

}
