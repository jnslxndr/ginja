package ginja

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func NewTestApi() *Api {
	return &Api{}
}

type TestItem struct {
	Name string `json:"name"`
}

var testItem = TestItem{
	Name: "A Name",
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

func TestNewDocument(t *testing.T) {
	Convey("Empty document has data:null", t, func() {
		d := NewDocument()

		payload, err := json.Marshal(&d)

		So(string(payload), ShouldEqual, `{"data":null}`)
		So(err, ShouldBeNil)

	})
}

func TestAddData(t *testing.T) {
	Convey("Document with simple data", t, func() {
		d := NewDocument()
		d.AddData(&ResourceObject{Id: "0", Object: &testItem})
		payload, err := json.Marshal(&d)

		So(string(payload), ShouldEqual, `{"data":{"type":"testitem","id":"0","attributes":{"name":"A Name"}}}`)
		So(err, ShouldBeNil)
	})
}

func BenchmarkNewDocument1000(b *testing.B) {
	ro := &ResourceObject{Id: "0", Object: &testItem}
	var d Document
	for n := 0; n < b.N; n++ {
		d = NewDocument()
		d.AddData(ro)
		d.MarshalJSON()
	}
}

func TestNewCollectionDocument(t *testing.T) {
	Convey("Empty collection document has data:[]", t, func() {
		d := NewCollectionDocument()

		// So(d, ShouldImplement, (*json.Marshaler)(nil))

		payload, err := json.Marshal(&d)

		So(string(payload), ShouldEqual, `{"data":[]}`)
		So(err, ShouldBeNil)
	})
}

func TestNewErrorDocument(t *testing.T) {
	Convey("Empty error document das no data, but empty errors field", t, func() {
		ed := NewErrorDocument()

		payload, err := json.Marshal(&ed)

		So(string(payload), ShouldEqual, `{"errors":[]}`)
		So(err, ShouldBeNil)
	})
}
