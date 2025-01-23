package interfaces

import (
	"reflect"
)

/*
IRegistry - Interface for encapsulating a Type registry service so a struct instance
can be instantiated based on the type name.
*/
type IRegistry interface {
	//Register - receive as an input an instance and register, extract the Type and register it.
	Register(interface{}) (bool, error)
	//RegisterType - receive a reflect.Type instance and register it.
	RegisterType(reflect.Type) (bool, error)
	//Info Retrieve the registered entry for this type.
	Info(string) (IInfo, error)
}

type IInfo interface {
	Type() reflect.Type
	Name() string
	Serializer(SerializerMode) ISerializer
	AddSerializer(ISerializer)
	NewInstance() (interface{}, error)
}
