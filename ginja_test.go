package ginja

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("New ginja router object generator", t, func() {
		server := gin.New()
		c := defaultConfig
		c.MountStats = true
		g := New(server, c)

		So(g.Namespace, ShouldEqual, "api")
		So(g.Version, ShouldEqual, "1")

		Convey("server reponse to _stats with application/json in debug mode", func() {

			req, _ := http.NewRequest("GET", "/api/v1/_stats", nil)
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Header().Get("Content-type"), ShouldEqual, "application/json")
			So(resp.Body, ShouldNotBeEmpty)
		})

		Convey("server reponse to _stats with application/vnd.api+json in release mode", func() {

			c.Debug = false
			server := gin.New()
			g = New(server, c)

			req, _ := http.NewRequest("GET", "/api/v1/_stats", nil)
			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Header().Get("Content-type"), ShouldEqual, "application/vnd.api+json")
			So(resp.Body, ShouldNotBeEmpty)

		})

	})
}
