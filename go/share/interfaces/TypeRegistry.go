package interfaces

import (
	"reflect"
)

type ITypeInfo interface {
	Type() reflect.Type
	Name() string
	Serializer(SerializerMode) Serializer
	AddSerializer(Serializer)
	NewInstance() (interface{}, error)
}

type ITypeRegistry interface {
	Register(interface{}) bool
	RegisterType(reflect.Type) bool
	TypeInfo(string) (ITypeInfo, error)
}

var registry ITypeRegistry

func TypeRegistry() ITypeRegistry {
	return registry
}

func SetTypeRegistry(r ITypeRegistry) {
	registry = r
}
