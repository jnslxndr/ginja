package ginja

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuildUrl(t *testing.T) {
	Convey("buildUrl combines namespace and version of the API config", t, func() {
		c := Config{
			Namespace: "test-api",
			Version:   "1",
		}

		So(c.buildUrl(), ShouldEqual, "/test-api/v1")
	})

}

func TestApplyDefaults(t *testing.T) {
	Convey("Applay default values", t, func() {
		Convey("sets proper debug", func() {
			c := Config{Debug: false}
			c.ApplyDefaults()

			So(c.Debug, ShouldBeFalse)

			c = Config{Debug: true}
			c.ApplyDefaults()

			So(c.Debug, ShouldBeTrue)

		})
	})
}
