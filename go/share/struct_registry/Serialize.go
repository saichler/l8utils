package struct_registry

import (
	"github.com/saichler/shared/go/share/interfaces"
	"google.golang.org/protobuf/proto"
	"reflect"
	"sync"
)

var empty = make([]byte, 0)
var pbMtx = &sync.Mutex{}
var UseProtoBuiltInSerializer = true

func (r *StructRegistryImpl) Marshal(any interface{}) ([]byte, error) {
	if any == nil {
		return empty, nil
	}

	pb, ok := any.(proto.Message)
	if ok && UseProtoBuiltInSerializer {
		pbMtx.Lock()
		defer pbMtx.Unlock()
		return proto.Marshal(pb)
	}

	val := reflect.ValueOf(any)
	if !val.IsValid() {
		return empty, nil
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return empty, nil
		}
		val = val.Elem()
	}

	_, ser, _ := r.types.Get(val.Type().Name())
	if ser != nil {
		b, _ := ser.Add(any, r)
		return b, nil
	}
	return nil, interfaces.Error("serializer not found in struct registry for ", val.Type().Name())
}

func (r *StructRegistryImpl) Unmarshal(name string, b []byte) (interface{}, error) {
	ins, ser, err := r.NewInstance(name)
	if err == nil && UseProtoBuiltInSerializer {
		pb := ins.(proto.Message)
		pbMtx.Lock()
		defer pbMtx.Unlock()
		err = proto.Unmarshal(b, pb)
		if err != nil {
			return nil, err
		}
		return pb, nil
	}

	if err == nil && ser != nil {
		dins, _ := ser.Get(b, 0, r)
		return dins, nil
	}
	return nil, interfaces.Error("serializer not found in struct registry for ", name)
}
