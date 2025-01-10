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
