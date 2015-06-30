package ginja

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetType(t *testing.T) {
	Convey("GetType returns the reflected type of the underlying Object", t, func() {
		Convey("on values", func() {
			ro := ResourceObject{Object: testItem}
			So(ro.getType(), ShouldEqual, "testitem")
		})
		Convey("as well as on pointers", func() {
			ro := ResourceObject{Object: &testItem}
			So(ro.getType(), ShouldEqual, "testitem")
		})
	})
}
