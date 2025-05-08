package registry

import (
	"github.com/saichler/l8utils/go/utils/maps"
	"reflect"
)

type TypesMap struct {
	impl *maps.SyncMap
}

func NewTypesMap() *TypesMap {
	s2t := &TypesMap{}
	s2t.impl = maps.NewSyncMap()
	return s2t
}

func (m *TypesMap) Put(key string, value reflect.Type) (bool, error) {
	info, err := NewInfo(value)
	if err != nil {
		return false, err
	}
	return m.impl.Put(key, info), err
}

func (m *TypesMap) Get(key string) (*Info, bool) {
	value, ok := m.impl.Get(key)
	if value != nil {
		info := value.(*Info)
		return info, ok
	}
	return nil, ok
}

func (m *TypesMap) Contains(key string) bool {
	return m.impl.Contains(key)
}
