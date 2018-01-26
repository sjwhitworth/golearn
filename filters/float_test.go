package filters

import (
	"testing"

	"github.com/amclay/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFloatFilter(t *testing.T) {

	Convey("Given a contrived dataset...", t, func() {

		// Read the contrived dataset
		inst, err := base.ParseCSVToInstances("./binary_test.csv", true)
		So(err, ShouldEqual, nil)

		// Add Attributes to the filter
		bFilt := NewFloatConvertFilter()
		bAttrs := base.NonClassAttributes(inst)
		for _, a := range bAttrs {
			bFilt.AddAttribute(a)
		}
		bFilt.Train()

		// Construct a LazilyFilteredInstances to handle it
		instF := base.NewLazilyFilteredInstances(inst, bFilt)

		Convey("All the non-class Attributes should be floats...", func() {
			// Check that all the Attributes are the right type
			for _, a := range base.NonClassAttributes(instF) {
				_, ok := a.(*base.FloatAttribute)
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

		Convey("All Attributes should be the correct type...", func() {
			for _, a := range instF.AllAttributes() {
				if a.GetName() == "arbitraryClass" {
					_, ok := a.(*base.CategoricalAttribute)
					So(ok, ShouldEqual, true)
				} else {
					_, ok := a.(*base.FloatAttribute)
					So(ok, ShouldEqual, true)
				}
			}
		})

		// Check that the Attributes have been discretised correctly
		Convey("FloatConversion should have worked", func() {
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
					So(instF.Get(as, 0), ShouldResemble, base.PackFloatToBytes(1.0))
					So(instF.Get(as, 1), ShouldResemble, base.PackFloatToBytes(1.0))
					So(instF.Get(as, 2), ShouldResemble, base.PackFloatToBytes(0.0))
				} else if name == "shouldBe1Binary" {
					So(instF.Get(as, 0), ShouldResemble, base.PackFloatToBytes(0.0))
					So(instF.Get(as, 1), ShouldResemble, base.PackFloatToBytes(1.0))
					So(instF.Get(as, 2), ShouldResemble, base.PackFloatToBytes(1.0))
				} else if name == "shouldBe3Binary_stoicism" {
					So(instF.Get(as, 0), ShouldResemble, base.PackFloatToBytes(1.0))
					So(instF.Get(as, 1), ShouldResemble, base.PackFloatToBytes(0.0))
					So(instF.Get(as, 2), ShouldResemble, base.PackFloatToBytes(0.0))
				} else if name == "shouldBe3Binary_heroism" {
					So(instF.Get(as, 0), ShouldResemble, base.PackFloatToBytes(0.0))
					So(instF.Get(as, 1), ShouldResemble, base.PackFloatToBytes(1.0))
					So(instF.Get(as, 2), ShouldResemble, base.PackFloatToBytes(0.0))
				} else if name == "shouldBe3Binary_romanticism" {
					So(instF.Get(as, 0), ShouldResemble, base.PackFloatToBytes(0.0))
					So(instF.Get(as, 1), ShouldResemble, base.PackFloatToBytes(0.0))
					So(instF.Get(as, 2), ShouldResemble, base.PackFloatToBytes(1.0))
				} else if name == "arbitraryClass" {
				}
			}
		})
	})
}
