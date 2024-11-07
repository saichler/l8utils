package interfaces

import (
	"google.golang.org/protobuf/proto"
	"reflect"
)

type IRegistry interface {
	RegisterStruct(interface{}) bool
	RegisterStructType(reflect.Type) bool
	NewProtobufInstance(string) (proto.Message, error)
	NewInstance(string) (interface{}, error)
	TypeByName(string) (reflect.Type, error)
}

var registry IRegistry

func Registry() IRegistry {
	return registry
}

func SetRegistry(r IRegistry) {
	registry = r
}
