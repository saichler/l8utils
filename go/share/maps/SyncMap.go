package maps

import (
	"reflect"
	"sync"
)

type SyncMap struct {
	m map[interface{}]interface{}
	s *sync.RWMutex
}

func NewSyncMap() *SyncMap {
	mm := &SyncMap{}
	mm.m = make(map[interface{}]interface{})
	mm.s = &sync.RWMutex{}
	return mm
}

func (this *SyncMap) Put(key, value interface{}) bool {
	if this == nil {
		return false
	}
	this.s.Lock()
	defer this.s.Unlock()
	_, ok := this.m[key]
	this.m[key] = value
	return !ok
}

func (this *SyncMap) Get(key interface{}) (interface{}, bool) {
	if this == nil {
		return nil, false
	}
	this.s.RLock()
	defer this.s.RUnlock()
	v, ok := this.m[key]
	return v, ok
}

func (this *SyncMap) Contains(key interface{}) bool {
	if this == nil {
		return false
	}
	this.s.RLock()
	defer this.s.RUnlock()
	_, ok := this.m[key]
	return ok
}

func (this *SyncMap) Delete(key interface{}) (interface{}, bool) {
	if this == nil {
		return nil, false
	}
	this.s.Lock()
	defer this.s.Unlock()
	v, ok := this.m[key]
	delete(this.m, key)
	return v, ok
}

func (this *SyncMap) Size() int {
	if this == nil {
		return 0
	}
	this.s.RLock()
	defer this.s.RUnlock()
	return len(this.m)
}

func (this *SyncMap) Clean() map[interface{}]interface{} {
	if this == nil {
		return nil
	}
	this.s.Lock()
	defer this.s.Unlock()
	result := this.m
	this.m = make(map[interface{}]interface{})
	return result
}

func (this *SyncMap) ValuesAsList(typ reflect.Type, filter func(interface{}) bool) interface{} {
	if this == nil {
		return false
	}
	this.s.RLock()
	defer this.s.RUnlock()
	newSlice := reflect.MakeSlice(reflect.SliceOf(typ), len(this.m), len(this.m))
	index := 0
	for _, v := range this.m {
		if filter != nil && !filter(v) {
			continue
		}
		newSlice.Index(index).Set(reflect.ValueOf(v))
		index++
	}

	if index < len(this.m) {
		filterSlice := reflect.MakeSlice(reflect.SliceOf(typ), index, index)
		for i := 0; i < index; i++ {
			filterSlice.Index(i).Set(newSlice.Index(i))
		}
		return filterSlice.Interface()
	}

	return newSlice.Interface()
}

func (this *SyncMap) KeysAsList(typ reflect.Type, filter func(interface{}) bool) interface{} {
	if this == nil {
		return false
	}
	this.s.RLock()
	defer this.s.RUnlock()
	newSlice := reflect.MakeSlice(reflect.SliceOf(typ), len(this.m), len(this.m))
	index := 0
	for v, _ := range this.m {
		if filter != nil && !filter(v) {
			continue
		}
		newSlice.Index(index).Set(reflect.ValueOf(v))
		index++
	}

	if index < len(this.m) {
		filterSlice := reflect.MakeSlice(reflect.SliceOf(typ), index, index)
		for i := 0; i < index; i++ {
			filterSlice.Index(i).Set(newSlice.Index(i))
		}
		return filterSlice.Interface()
	}

	return newSlice.Interface()
}

func (this *SyncMap) Iterate(do func(k, v interface{})) {
	if this == nil {
		return
	}
	this.s.RLock()
	defer this.s.RUnlock()
	for k, v := range this.m {
		do(k, v)
	}
}
