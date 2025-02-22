package registry

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"reflect"
	"sync"
)

type RegistryImpl struct {
	types *TypesMap
	enums map[string]int32
	mtx   *sync.RWMutex
}

func NewRegistry() *RegistryImpl {
	registry := &RegistryImpl{}
	registry.types = NewTypesMap()
	registry.enums = make(map[string]int32)
	registry.mtx = new(sync.RWMutex)
	registry.registerPrimitives()
	return registry
}

func (this *RegistryImpl) registerPrimitives() {
	this.RegisterType(reflect.TypeOf(int8(0)))
	this.RegisterType(reflect.TypeOf(int16(0)))
	this.RegisterType(reflect.TypeOf(int32(0)))
	this.RegisterType(reflect.TypeOf(int64(0)))
	this.RegisterType(reflect.TypeOf(""))
	this.RegisterType(reflect.TypeOf(true))
	this.RegisterType(reflect.TypeOf(float32(0)))
	this.RegisterType(reflect.TypeOf(float64(0)))
}

func (this *RegistryImpl) Register(any interface{}) (bool, error) {
	v := reflect.ValueOf(any)
	if !v.IsValid() {
		return false, errors.New("invalid input")
	}
	if v.Kind() == reflect.Ptr {
		return this.RegisterType(v.Elem().Type())
	}
	return this.RegisterType(v.Type())
}

func (this *RegistryImpl) RegisterType(t reflect.Type) (bool, error) {
	if t == nil {
		return false, errors.New("Cannot register a nil type")
	}
	if t.Name() == "rtype" {
		return false, nil
	}
	return this.types.Put(t.Name(), t)
}

func (this *RegistryImpl) Info(name string) (interfaces.IInfo, error) {
	typ, ok := this.types.Get(name)
	if !ok {
		return nil, errors.New("Unknown Type: " + name)
	}
	return typ, nil
}

func (this *RegistryImpl) RegisterEnums(enums map[string]int32) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	for name, value := range enums {
		this.enums[name] = value
	}
}

func (this *RegistryImpl) Enum(name string) int32 {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	return this.enums[name]
}
