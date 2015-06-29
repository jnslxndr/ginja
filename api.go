package ginja

import (
	"log"
	"strings"

	"reflect"
)

type Api struct {
	Config
	Store
}

func (a *Api) Register(i interface{}) {
	t := a.registerType(i)
	log.Println("Registering type", t.String())
	log.Println(a.types)
	// a.registerSerializerFor(t)
}

type Store struct {
	types map[reflect.Type][2]string
}

func (s *Store) registerType(i interface{}) reflect.Type {
	if s.types == nil {
		s.types = make(map[reflect.Type][2]string)
	}
	t := reflect.TypeOf(i)
	name := strings.ToLower(t.Name())

	s.types[t] = [2]string{name, name + "s"}
	return t
}

func (s *Store) NameFor(i interface{}) string {
	return s.types[reflect.TypeOf(i)][0]
}
