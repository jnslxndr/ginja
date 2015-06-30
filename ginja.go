package ginja

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type GinApi struct {
	*gin.RouterGroup
	Api
}

// New returns a new ginja.Api struct
func New(server *gin.Engine, config Config, middleware ...gin.HandlerFunc) *GinApi {
	config.ApplyDefaults()
	api := &GinApi{
		RouterGroup: server.Group(config.buildUrl()),
		Api:         Api{Config: config},
	}

	api.init()

	api.Use(middleware...)

	// TODO extract!
	if api.MountStats {
		api.GET(api.StatsURL, func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"config": api.Config})
		})
	}

	return api
}

func contentTypeSetter(isDebug bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isDebug {
			c.Header("Content-type", "application/json")
		} else {
			c.Header("Content-type", "application/vnd.api+json")
		}
		c.Next()
	}
}

func (a *GinApi) init() *GinApi {
	if a.Debug == false {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	a.Use(contentTypeSetter(a.Debug))

	if a.Debug {
		config, _ := json.MarshalIndent(a.Config, "  ", "  ")
		log.Println("ginja API Config:")
		log.Println("  " + string(config))
	}

	return a
}

type Resource interface {
	GetId() string
	Attributes() interface{}
	// setType(interface{})
	getType() string
}

type ResourceObject struct {
	Id     string
	Object interface{}
}

func (r ResourceObject) GetId() string {
	return r.Id
}

func (r ResourceObject) getType() string {
	if reflect.TypeOf(r.Object).Kind() == reflect.Ptr {
		return strings.ToLower(reflect.TypeOf(r.Object).Elem().Name())
	} else {
		return strings.ToLower(reflect.TypeOf(r.Object).Name())
	}
}

func (r ResourceObject) Attributes() interface{} {
	return &r.Object
}

type Fragment struct {
	Type       string      `json:"type"`
	Id         string      `json:"id"`
	Attributes interface{} `json:"attributes,omitempty"`
}

type Document struct {
	Data    interface{}
	Meta    map[string]interface{}
	Errors  []Error
	isError bool
}

type metaObject struct {
	Meta map[string]interface{} `json:"meta,omitempty"`
}

type document struct {
	Data interface{} `json:"data"`
	metaObject
}

type errorDocument struct {
	Errors []Error `json:"errors"`
	metaObject
}

type ErrorDocument Document

type collection []interface{}

// func (ic collection) String() string {
// 	return "a slice of interfaces"
// }

func NewDocument() Document {
	return Document{}
}

func NewCollectionDocument() Document {
	return Document{Data: []Fragment{}}
}

func NewErrorDocument() Document {
	return Document{Errors: []Error{}, isError: true}
}

func (d *Document) AddMeta(meta map[string]interface{}) *Document {
	d.Meta = meta
	return d
}

func (d *Document) AddData(data Resource) {
	d.Data = Fragment{
		Type:       data.getType(),
		Id:         data.GetId(),
		Attributes: data.Attributes(),
	}
}

func (d *Document) AddError(err error) *Document {
	d.Errors = append(d.Errors, NewError(err))
	return d
}

var documentPool = sync.Pool{
	New: func() interface{} {
		return &document{}
	},
}

var errorDocumentPool = sync.Pool{
	New: func() interface{} {
		return &errorDocument{}
	},
}

func (d Document) MarshalJSON() ([]byte, error) {
	if len(d.Errors) > 0 || d.isError {
		payload := errorDocumentPool.Get().(*errorDocument)
		payload.Errors = d.Errors
		payload.Meta = nil
		if len(d.Meta) > 0 {
			payload.Meta = d.Meta
		}
		defer errorDocumentPool.Put(payload)
		return json.Marshal(&payload)
	} else {
		payload := documentPool.Get().(*document)
		payload.Data = d.Data
		payload.Meta = nil
		if len(d.Meta) > 0 {
			payload.Meta = d.Meta
		}
		defer documentPool.Put(payload)
		return json.Marshal(&payload)
	}
}

// func (d Document) UnmarshalJSON(data []byte) error {
// 	log.Println("testing unmarshalling")
// 	return nil
// }
