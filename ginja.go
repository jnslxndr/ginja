package ginja

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/imdario/mergo"
)

type GinApi struct {
	*gin.RouterGroup
	Api
}

// New returns a new ginja.Api struct
func New(server *gin.Engine, config Config, middleware ...gin.HandlerFunc) *GinApi {
	c := defaultConfig

	setDebug := config.Debug
	shouldDebug := false
	if setDebug {
		shouldDebug = config.Debug
	}

	mergo.MergeWithOverwrite(&c, config)

	if setDebug {
		c.Debug = shouldDebug
	}

	api := &GinApi{
		RouterGroup: server.Group(c.buildUrl()),
		Api:         Api{Config: c},
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

func (a *GinApi) init() *GinApi {
	if !a.Debug {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	a.Use(func(c *gin.Context) {
		if a.Debug {
			c.Header("Content-type", "application/json")
		} else {
			c.Header("Content-type", "application/vnd.api+json")
		}

		c.Next()
	})

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
	setType(interface{})
	getType() reflect.Type
}

type ResourceObject struct {
	Id           string
	resourceType reflect.Type `jons:"-"`
	Object       interface{}
}

func (r ResourceObject) GetId() string {
	return r.Id
}

func (r ResourceObject) setType(i interface{}) {
	r.resourceType = reflect.TypeOf(i)
}

func (r ResourceObject) getType() reflect.Type {
	return reflect.TypeOf(r.Object)
}

func (r *ResourceObject) Attributes() interface{} {
	return &r.Object
}

type member *struct{}

type Fragment map[string]interface{}

type Document Fragment

type ErrorDocument Document

type collection []interface{}

func (ic collection) String() string {
	return "a slice of interfaces"
}

func NewDocument() Document {
	return Document{"data": nil}
}

func NewCollectionDocument() Document {
	return Document{"data": []interface{}{}}
}

func NewErrorDocument() Document {
	return Document{"errors": []interface{}{}}
}

func (d *Document) AddMeta(meta map[string]interface{}) *Document {
	(*(*map[string]interface{})(d))["meta"] = meta
	return d
}

func (d *Document) AddData(data Resource) *Document {
	data.setType(reflect.TypeOf(data))
	(*(*map[string]interface{})(d))["data"] = Fragment{
		"type":       strings.ToLower(data.getType().Name()),
		"id":         data.GetId(),
		"attributes": data.Attributes(),
	}
	return d
}

func (d *Document) AddError(err error) *Document {
	errors := (*map[string]interface{})(d)
	if (*errors)["errors"] == nil {
		(*errors)["errors"] = []Error{}
	}
	(*errors)["errors"] = append((*errors)["errors"].([]interface{}), NewError(err))
	return d
}

// func (d Document) MarshalJSON() ([]byte, error) {
// 	payload := make(map[string]interface{})

// 	return json.Marshal(&payload)
// }

// func (d Document) UnmarshalJSON(data []byte) error {
// 	log.Println("testing unmarshalling")
// 	return nil
// }
