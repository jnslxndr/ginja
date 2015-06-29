package ginja

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func NewTestApi() *Api {
	api := &Api{}
	return api
}

type TestItem struct {
	Name string `json:"name"`
	Data int64  `json:"data"`
}

var testItem = TestItem{
	Name: "A Name",
	Data: math.MaxInt64,
}

var testItemPayload = map[string]interface{}{
	"data": map[string]interface{}{
		"type": "testitem",
		"id":   "0",
		"attributes": map[string]interface{}{
			"name": "A Name",
		},
	},
}

func TestStoreRegister(t *testing.T) {
	Convey("Api can register arbitrary types", t, func() {
		api := NewTestApi()
		api.Register(TestItem{})

		So(api.types[reflect.TypeOf(TestItem{})], ShouldResemble, [2]string{"testitem", "testitems"})
		So(api.NameFor(TestItem{}), ShouldEqual, "testitem")
	})
}

func TestDocumentMarshalJSON(t *testing.T) {
	Convey("Document is a json.Marshaler", t, func() {
		d := NewDocument()

		d.AddData(&ResourceObject{Id: "0", Object: testItem})
		// d.AddMeta(Fragment{"test": "some meta"})

		// So(d, ShouldImplement, (*json.Marshaler)(nil))

		payload, err := json.MarshalIndent(&d, "", "  ")

		So(string(payload), ShouldEqual, "lksdjflsdkjf")
		So(err, ShouldNotBeNil)

	})

}

// func TestSerialize(t *testing.T) {
// 	Convey("Serialize returns jsonapi compliant payload", t, func() {
// 		api := NewTestApi()

// 		api.Register(TestItem{})

// 		ti := testItem

// 		payload, err := api.Serialize(ti)

// 		So(err, ShouldBeNil)

// 		expected, _ := json.Marshal(testItemPayload)

// 		So(string(payload), ShouldEqual, string(expected))
// 	})
// }
