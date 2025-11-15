package registry

import (
	"errors"
	"reflect"
	"sync"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8health"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8types/go/types/l8reflect"
	"github.com/saichler/l8types/go/types/l8services"
	"github.com/saichler/l8types/go/types/l8sysconfig"
	"github.com/saichler/l8types/go/types/l8system"
	"github.com/saichler/l8types/go/types/l8web"
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
	registry.registerBase()
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

func (this *RegistryImpl) registerBase() {
	this.Register(&l8api.L8Query{})
	this.Register(&l8api.L8MetaData{})
	this.Register(&l8api.AuthToken{})
	this.Register(&l8api.AuthUser{})
	this.Register(&l8health.L8Health{})
	this.Register(&l8health.L8HealthList{})
	this.Register(&l8web.L8Empty{})
	this.Register(&l8notify.L8NotificationSet{})
	this.Register(&l8reflect.L8Node{})
	this.Register(&l8reflect.L8TableView{})
	this.Register(&l8web.L8WebService{})
	this.Register(&l8services.L8Services{})
	this.Register(&l8services.L8ReplicationIndex{})
	this.Register(&l8services.L8Transaction{})
	this.Register(&l8services.L8ServiceLink{})
	this.Register(&l8sysconfig.L8SysConfig{})
	this.Register(&l8system.L8SystemMessage{})
	this.Register(&l8web.L8Plugin{})
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

func (this *RegistryImpl) UnRegister(typeName string) (bool, error) {
	if typeName == "" {
		return false, errors.New("Cannot unregister a blank type")
	}
	return this.types.Del(typeName), nil
}

func (this *RegistryImpl) Info(name string) (ifs.IInfo, error) {
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

func (this *RegistryImpl) NewOf(any interface{}) interface{} {
	this.Register(any)
	typeName := reflect.ValueOf(any).Elem().Type().Name()
	info, _ := this.Info(typeName)
	r, _ := info.NewInstance()
	return r
}
