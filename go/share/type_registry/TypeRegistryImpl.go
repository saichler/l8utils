package type_registry

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
)

type TypeRegistryImpl struct {
	types *TypesMap
}

func NewTypeRegistry() *TypeRegistryImpl {
	tri := &TypeRegistryImpl{}
	tri.types = NewTypesMap()
	tri.registerPrimitives()
	return tri
}

func (this *TypeRegistryImpl) registerPrimitives() {
	this.RegisterType(reflect.TypeOf(int8(0)))
	this.RegisterType(reflect.TypeOf(int16(0)))
	this.RegisterType(reflect.TypeOf(int32(0)))
	this.RegisterType(reflect.TypeOf(int64(0)))
	this.RegisterType(reflect.TypeOf(""))
	this.RegisterType(reflect.TypeOf(true))
	this.RegisterType(reflect.TypeOf(float32(0)))
	this.RegisterType(reflect.TypeOf(float64(0)))
}

func (this *TypeRegistryImpl) Register(any interface{}) bool {
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		return this.RegisterType(v.Elem().Type())
	}
	return this.RegisterType(v.Type())
}

func (this *TypeRegistryImpl) RegisterType(t reflect.Type) bool {
	if t.Name() == "rtype" {
		return false
	}
	return this.types.Put(t.Name(), t)
}

func (this *TypeRegistryImpl) TypeInfo(name string) (interfaces.ITypeInfo, error) {
	typ, ok := this.types.Get(name)
	if !ok {
		return nil, errors.New("Unknown Type: " + name)
	}
	return typ, nil
}
