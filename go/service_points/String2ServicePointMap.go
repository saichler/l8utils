package service_points

import (
	"github.com/saichler/shared/go/interfaces"
	"github.com/saichler/shared/go/maps"
)

type String2ServicePointMap struct {
	impl *maps.SyncMap
}

func NewString2ServicePointMap() *String2ServicePointMap {
	newMap := &String2ServicePointMap{}
	newMap.impl = maps.NewSyncMap()
	return newMap
}

func (mp *String2ServicePointMap) Put(key string, value interfaces.IServicePointHandler) bool {
	return mp.impl.Put(key, value)
}

func (mp *String2ServicePointMap) Get(key string) (interfaces.IServicePointHandler, bool) {
	value, ok := mp.impl.Get(key)
	if value != nil {
		return value.(interfaces.IServicePointHandler), ok
	}
	return nil, ok
}

func (mp *String2ServicePointMap) Contains(key string) bool {
	return mp.impl.Contains(key)
}
