package ginja

import (
	"errors"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	Convey("Error implements Error interface", t, func() {
		err := Error{Title: "Test error"}

		So(err, ShouldImplement, (*error)(nil))
		So(err.Error(), ShouldNotBeBlank)
		So(err.Error(), ShouldHaveSameTypeAs, "test")
		So(err.Error(), ShouldEqual, "Test error")
	})
}

type _testStringer struct{}

func (_ _testStringer) String() string {
	return "Stringer error"
}

func TestNewError(t *testing.T) {
	Convey("Error handles different sources", t, func() {
		err := NewError("Test error")
		So(err.Error(), ShouldEqual, "Test error")

		err = NewError(errors.New("Test error"))
		So(err.Error(), ShouldEqual, "Test error")

		err = NewError(os.ErrInvalid)
		So(err.Error(), ShouldEqual, "invalid argument")

		err = NewError(_testStringer{})
		So(err.Error(), ShouldEqual, "Stringer error")

		err = NewError(struct{}{})
		So(err.Error(), ShouldEqual, "Unknown error occurred")
	})
}
