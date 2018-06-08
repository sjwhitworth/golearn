package linear_models

import (
	//"github.com/sjwhitworth/golearn/base"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLinearSVC(t *testing.T) {
	Convey("Doing a LinearSVC test", t, func(){
		_, err := NewLinearSVC("l1", "l1", false, 1.0, -1e6)
		So(err, ShouldNotBeNil)
		_, err = NewLinearSVC("l0", "l1", false, 1.0, -1e6)
		So(err, ShouldNotBeNil)
		_, err = NewLinearSVC("l1", "l0", false, 1.0, -1e6)
		So(err, ShouldNotBeNil)
		_, err = NewLinearSVC("l1", "l2", false, 1.0, -1e6)
		So(err, ShouldNotBeNil)
		SVC, err := NewLinearSVC("l1", "l2", true, 1.0, -1e6)
		So(SVC, ShouldNotBeNil)
		So(err, ShouldBeNil)
		_, err = NewLinearSVC("l2", "l2", false, 1.0, -1e6)
		So(err, ShouldBeNil)
		_, err = NewLinearSVC("l2", "l2", true, 1.0, -1e6)
		So(err, ShouldBeNil)
		_, err = NewLinearSVC("l2", "l1", false, 1.0, -1e6)
		So(err, ShouldBeNil)
		_, err = NewLinearSVC("l2", "l1", true, 1.0, -1e6)
		So(err, ShouldNotBeNil)
		So(func(){ SVC.GetMetadata() } , ShouldNotPanic)
		//err = SVC.Save("tmp")

		params := &LinearSVCParams{0, []float64{0.0}, 1.0, -1e6, false, false}
		params = params.Copy()
		So(params, ShouldNotBeNil)

		g := [][]float64{{1, 2}, {1, 2}, {1, 2}}
		v := []float64{1, 2}
		var bias float64
		problem := NewProblem(g[:], v[:], bias)
		param := NewParameter(0, 1.0, -1e6)
		model := Train(problem, param)
		So(model, ShouldNotBeNil)
		err = Export(model, "tmp")
		So(err, ShouldBeNil)
		err = Load(model, "tmp")
		So(err, ShouldBeNil)

		SVC.model = model
		err = SVC.Save("tmp")
		So(err, ShouldBeNil)
		err = SVC.Load("tmp")
		So(err, ShouldBeNil)
		s := SVC.String()
		So(s, ShouldEqual, "LogisticSVC")
	})
}
