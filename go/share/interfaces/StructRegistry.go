package interfaces

import (
	"google.golang.org/protobuf/proto"
	"reflect"
)

type IStructRegistry interface {
	RegisterStruct(interface{}) bool
	RegisterStructType(reflect.Type) bool
	NewProtobufInstance(string) (proto.Message, error)
	NewInstance(string) (interface{}, error)
	TypeByName(string) (reflect.Type, error)
}

var registry IStructRegistry

func StructRegistry() IStructRegistry {
	return registry
}

func SetStructRegistry(r IStructRegistry) {
	registry = r
}
