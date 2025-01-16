package interfaces

import (
	"reflect"
)

type IInfo interface {
	Type() reflect.Type
	Name() string
	Serializer(SerializerMode) ISerializer
	AddSerializer(ISerializer)
	NewInstance() (interface{}, error)
}

type IRegistry interface {
	Register(interface{}) (bool, error)
	RegisterType(reflect.Type) (bool, error)
	Info(string) (IInfo, error)
}
