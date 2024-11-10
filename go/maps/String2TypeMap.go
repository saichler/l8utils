package maps

import (
	"reflect"
)

type String2TypeMap struct {
	impl *SyncMap
}

func NewString2TypeMao() *String2TypeMap {
	s2t := &String2TypeMap{}
	s2t.impl = NewSyncMap()
	return s2t
}

func (s2t *String2TypeMap) Put(key string, value reflect.Type) bool {
	return s2t.impl.Put(key, value)
}

func (s2t *String2TypeMap) Get(key string) (reflect.Type, bool) {
	value, ok := s2t.impl.Get(key)
	if value != nil {
		return value.(reflect.Type), ok
	}
	return nil, ok
}

func (s2t *String2TypeMap) Contains(key string) bool {
	return s2t.impl.Contains(key)
}
