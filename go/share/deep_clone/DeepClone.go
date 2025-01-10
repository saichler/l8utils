package deep_clone

import (
	"reflect"
	"strconv"
	"strings"
)

type Cloner struct {
	cloners map[reflect.Kind]func(reflect.Value, string, map[string]reflect.Value) reflect.Value
}

func newCloner() *Cloner {
	cloner := &Cloner{}
	cloner.initCloners()
	return cloner
}

func (cloner *Cloner) initCloners() {
	cloner.cloners = make(map[reflect.Kind]func(reflect.Value, string, map[string]reflect.Value) reflect.Value)
	cloner.cloners[reflect.Int] = cloner.intCloner
	cloner.cloners[reflect.Int32] = cloner.int32Cloner
	cloner.cloners[reflect.Ptr] = cloner.ptrCloner
	cloner.cloners[reflect.Struct] = cloner.structCloner
	cloner.cloners[reflect.String] = cloner.stringCloner
	cloner.cloners[reflect.Slice] = cloner.sliceCloner
	cloner.cloners[reflect.Map] = cloner.mapCloner
	cloner.cloners[reflect.Bool] = cloner.boolCloner
	cloner.cloners[reflect.Int64] = cloner.int64Cloner
	cloner.cloners[reflect.Uint] = cloner.uintCloner
	cloner.cloners[reflect.Uint32] = cloner.uint32Cloner
	cloner.cloners[reflect.Uint64] = cloner.uint64Cloner
	cloner.cloners[reflect.Float32] = cloner.float32Cloner
	cloner.cloners[reflect.Float64] = cloner.float64Cloner
}

func (Cloner *Cloner) Clone(any interface{}) interface{} {
	value := reflect.ValueOf(any)
	stopLoop := make(map[string]reflect.Value)
	valueClone := Cloner.clone(value, "", stopLoop)
	return valueClone.Interface()
}

func (Cloner *Cloner) clone(value reflect.Value, fieldName string, stopLoop map[string]reflect.Value) reflect.Value {
	if !value.IsValid() {
		return value
	}
	kind := value.Kind()
	cloner := Cloner.cloners[kind]
	if cloner == nil {
		panic("No cloner for kind:" + kind.String() + ":" + fieldName)
	}
	return cloner(value, fieldName, stopLoop)
}

func (Cloner *Cloner) sliceCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}
	newSlice := reflect.MakeSlice(reflect.SliceOf(value.Index(0).Type()), value.Len(), value.Len())
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		elemClone := Cloner.clone(elem, name, stopLoop)
		newSlice.Index(i).Set(elemClone)
	}
	return newSlice
}

func (Cloner *Cloner) ptrCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}

	p := strconv.Itoa(int(value.Pointer()))
	exist, ok := stopLoop[p]
	if ok {
		return exist
	}

	newPtr := reflect.New(value.Elem().Type())
	stopLoop[p] = newPtr

	newPtr.Elem().Set(Cloner.clone(value.Elem(), name, stopLoop))

	return newPtr
}

func (Cloner *Cloner) structCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	cloneStruct := reflect.New(value.Type()).Elem()
	structType := value.Type()
	for i := 0; i < structType.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldName := structType.Field(i).Name
		if SkipFieldByName(fieldName) {
			continue
		}
		cloned := Cloner.clone(fieldValue, structType.Field(i).Name, stopLoop)
		if cloned.Kind() == reflect.Int32 {
			cloneStruct.Field(i).SetInt(cloned.Int())
		} else {
			cloneStruct.Field(i).Set(cloned)
		}
	}
	return cloneStruct
}

func (Cloner *Cloner) mapCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	if value.IsNil() {
		return value
	}
	mapKeys := value.MapKeys()
	mapClone := reflect.MakeMapWithSize(value.Type(), len(mapKeys))
	for _, key := range mapKeys {
		mapElem := value.MapIndex(key)
		mapElemClone := Cloner.clone(mapElem, name, stopLoop)
		mapClone.SetMapIndex(key, mapElemClone)
	}
	return mapClone
}

func (Cloner *Cloner) intCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int(i))
}

func (Cloner *Cloner) uintCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint(i))
}

func (Cloner *Cloner) uint32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint32(i))
}

func (Cloner *Cloner) uint64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Uint()
	return reflect.ValueOf(uint64(i))
}

func (Cloner *Cloner) float32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Float()
	return reflect.ValueOf(float32(i))
}

func (Cloner *Cloner) float64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Float()
	return reflect.ValueOf(float64(i))
}

func (Cloner *Cloner) boolCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	b := value.Bool()
	return reflect.ValueOf(b)
}

func (Cloner *Cloner) int32Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int32(i))
}

func (Cloner *Cloner) int64Cloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	i := value.Int()
	return reflect.ValueOf(int64(i))
}

func (Cloner *Cloner) stringCloner(value reflect.Value, name string, stopLoop map[string]reflect.Value) reflect.Value {
	s := value.String()
	return reflect.ValueOf(s)
}

func SkipFieldByName(fieldName string) bool {
	if fieldName == "DoNotCompare" {
		return true
	}
	if fieldName == "DoNotCopy" {
		return true
	}
	if len(fieldName) > 3 && fieldName[0:3] == "XXX" {
		return true
	}
	if fieldName[0:1] == strings.ToLower(fieldName[0:1]) {
		return true
	}
	return false
}
