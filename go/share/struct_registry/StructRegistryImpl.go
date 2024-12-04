package struct_registry

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type StructRegistryImpl struct {
	types *TypesMap
}

func NewStructRegistry() *StructRegistryImpl {
	sr := &StructRegistryImpl{}
	sr.types = NewTypesMap()
	sr.registerPrimitives()
	return sr
}

func (r *StructRegistryImpl) registerPrimitives() {
	r.RegisterStructType(reflect.TypeOf(int8(0)), nil)
	r.RegisterStructType(reflect.TypeOf(int16(0)), nil)
	r.RegisterStructType(reflect.TypeOf(int32(0)), nil)
	r.RegisterStructType(reflect.TypeOf(int64(0)), nil)
	r.RegisterStructType(reflect.TypeOf(""), nil)
	r.RegisterStructType(reflect.TypeOf(true), nil)
	r.RegisterStructType(reflect.TypeOf(float32(0)), nil)
	r.RegisterStructType(reflect.TypeOf(float64(0)), nil)
}

func (r *StructRegistryImpl) RegisterStruct(any interface{}, serializer interfaces.Serializer) bool {
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		return r.RegisterStructType(v.Elem().Type(), serializer)
	}
	return r.RegisterStructType(v.Type(), serializer)
}

func (r *StructRegistryImpl) RegisterStructType(t reflect.Type, serializer interfaces.Serializer) bool {
	if t.Name() == "rtype" {
		return false
	}
	return r.types.Put(t.Name(), t, serializer)
}

func (r *StructRegistryImpl) TypeByName(name string) (reflect.Type, interfaces.Serializer, error) {
	typ, ser, ok := r.types.Get(name)
	if !ok {
		return nil, nil, errors.New("Unknown Struct Type: " + name)
	}
	return typ, ser, nil
}

func (r *StructRegistryImpl) NewInstance(typeName string) (interface{}, interfaces.Serializer, error) {
	if typeName == "" {
		return nil, nil, interfaces.Error("cannot create a new struct instance from blank type name")
	}
	typ, ser, ok := r.types.Get(typeName)
	if !ok {
		return nil, nil, interfaces.Error("Struct Type ", typeName, " is not registered")
	}
	n := reflect.New(typ)
	if !n.IsValid() {
		return nil, nil, interfaces.Error("Was not able to create new instance of type ", typeName)
	}
	return n.Interface(), ser, nil
}
