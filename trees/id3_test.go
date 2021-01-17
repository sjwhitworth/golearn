package trees

import (
	"fmt"
	"testing"

	"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
)

func TestId3(t *testing.T) {
	Convey("Doing a id3 test", t, func() {
		var rule DecisionTreeRule
		s := rule.String()
		So(s, ShouldNotBeNil)
		rule.SplitAttr = nil
		s = rule.String()
		So(s, ShouldNotBeNil)

		instances, err := base.ParseCSVToInstances("onerow.csv", true)
		So(err, ShouldBeNil)

		trainData, _ := base.InstancesTrainTestSplit(instances, 0.6)
		gRuleGen := new(GiniCoefficientRuleGenerator)
		root := InferID3Tree(trainData, gRuleGen)
		s = root.getNestedString(3)
		So(s, ShouldNotBeNil)
		s = root.String()
		So(s, ShouldNotBeNil)

		//var proba1 ClassProba
		//var proba2 ClassProba
		//probas := ClassesProba{proba1, proba2}
		_, rc := trainData.Size()
		fmt.Println(rc)
		id3tree := NewID3DecisionTree(0.1)
		So(id3tree, ShouldNotBeNil)
		id3tree.Root = root
		probas, err := id3tree.PredictProba(trainData)
		So(err, ShouldBeNil)
		var proba1, proba2 ClassProba
		probas = ClassesProba{proba1, proba2}
		L := probas.Len()
		So(L, ShouldEqual, 2)
		probas.Swap(0, 1)
		less := probas.Less(0, 1)
		So(less, ShouldEqual, false)

		data := id3tree.GetMetadata()
		So(data, ShouldNotBeNil)
		s = id3tree.String()
		So(s, ShouldNotBeNil)
		_, err = id3tree.Predict(trainData)
		So(err, ShouldBeNil)

		// Test save and load model
		err = id3tree.Save("tmp")
		So(err, ShouldBeNil)

		id3tree = NewID3DecisionTree(0.1)
		err = id3tree.Load("tmp")
		So(err, ShouldBeNil)

		_, err = id3tree.Predict(trainData)
		So(err, ShouldBeNil)
	})
}
