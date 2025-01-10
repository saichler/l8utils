package type_registry

import (
	"github.com/saichler/shared/go/share/maps"
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

func (m *TypesMap) Put(key string, value reflect.Type) bool {
	return m.impl.Put(key, NewTypeInfo(value))
}

func (m *TypesMap) Get(key string) (*TypeInfo, bool) {
	value, ok := m.impl.Get(key)
	if value != nil {
		info := value.(*TypeInfo)
		return info, ok
	}
	return nil, ok
}

func (m *TypesMap) Contains(key string) bool {
	return m.impl.Contains(key)
}
