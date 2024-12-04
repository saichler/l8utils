package struct_registry

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/maps"
	"reflect"
)

type TypesMap struct {
	impl *maps.SyncMap
}

type TypesMapEntry struct {
	typ        reflect.Type
	serializer interfaces.Serializer
}

func NewTypesMap() *TypesMap {
	s2t := &TypesMap{}
	s2t.impl = maps.NewSyncMap()
	return s2t
}

func (m *TypesMap) Put(key string, value reflect.Type, serializer interfaces.Serializer) bool {
	entry := &TypesMapEntry{}
	entry.typ = value
	entry.serializer = serializer
	return m.impl.Put(key, entry)
}

func (m *TypesMap) Get(key string) (reflect.Type, interfaces.Serializer, bool) {
	value, ok := m.impl.Get(key)
	if value != nil {
		entry := value.(*TypesMapEntry)
		return entry.typ, entry.serializer, ok
	}
	return nil, nil, ok
}

func (m *TypesMap) Contains(key string) bool {
	return m.impl.Contains(key)
}
