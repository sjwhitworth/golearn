package base

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestId3(t *testing.T) {
	Convey("Doing a error test", t, func() {
		var _gerr GoLearnError
		gerr := &_gerr
		gerr.attachFormattedStack()
		s := gerr.Error()
		So(s, ShouldNotBeNil)
		err := DescribeError("test", nil)
		So(err, ShouldNotBeNil)
		err = WrapError(nil)
		So(err, ShouldNotBeNil)
		s = wrapLinesWithTabPrefix("123\ntest\n")
		So(s, ShouldEqual, "\t123\n\ttest\n\t")
	})
}
