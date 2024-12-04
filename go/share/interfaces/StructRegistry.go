package interfaces

import (
	"reflect"
)

type IStructRegistry interface {
	RegisterStruct(interface{}, Serializer) bool
	RegisterStructType(reflect.Type, Serializer) bool
	NewInstance(string) (interface{}, Serializer, error)
	TypeByName(string) (reflect.Type, Serializer, error)
	Marshal(interface{}) ([]byte, error)
	Unmarshal(string, []byte) (interface{}, error)
}

var registry IStructRegistry

func StructRegistry() IStructRegistry {
	return registry
}

func SetStructRegistry(r IStructRegistry) {
	registry = r
}
