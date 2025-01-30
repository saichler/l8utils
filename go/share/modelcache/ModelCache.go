package modelcache

import (
	"github.com/saichler/reflect/go/reflect/clone"
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/reflect/updater"
	"sync"
)

type ModelCache struct {
	cache        map[string]interface{}
	mtx          *sync.RWMutex
	cond         *sync.Cond
	listener     IModelCacheListener
	cloner       *clone.Cloner
	introspector common.IIntrospect
}

type IModelCacheListener interface {
	ModelItemAdded(interface{})
	ModelItemDeleted(interface{})
	PropertyChangeNotification(interface{}, string, interface{}, interface{})
}

func NewModelCache(listener IModelCacheListener, introspector common.IIntrospect) *ModelCache {
	mc := &ModelCache{}
	mc.cache = make(map[string]interface{})
	mc.mtx = &sync.RWMutex{}
	mc.cond = sync.NewCond(mc.mtx)
	mc.listener = listener
	mc.cloner = clone.NewCloner()
	mc.introspector = introspector
	return mc
}

func (this *ModelCache) Get(k string) interface{} {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	item, ok := this.cache[k]
	if ok {
		itemClone := this.cloner.Clone(item)
		return itemClone
	}
	return nil
}

func (this *ModelCache) Put(k string, v interface{}) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	item, ok := this.cache[k]
	//If the item does not exist in the cache
	if !ok {
		//First clone the value so we can use it in the notification.
		itemClone := this.cloner.Clone(v)
		//Place the value in the cache
		this.cache[k] = v
		//Send the notification using the clone outside the current go routine
		if this.listener != nil {
			go this.listener.ModelItemAdded(itemClone)
		}
		return nil
	}
	//Clone the existing item
	itemClone := this.cloner.Clone(item)
	//Create a new updater
	putUpdater := updater.NewUpdater(this.introspector, true)
	//update the item clone with the new element where nil is valid
	err := putUpdater.Update(itemClone, v)
	if err != nil {
		return err
	}

	//if there are changes, then nothing to do
	changes := putUpdater.Changes()
	if changes == nil {
		return nil
	}

	//Apply the changes to the existing item
	for _, change := range changes {
		change.Apply(item)
	}
	defer func() {
		if this.listener != nil {
			for _, change := range changes {
				this.listener.PropertyChangeNotification(itemClone, change.PropertyId(), change.OldValue(), change.NewValue())
			}
		}
	}()
	return nil
}

func (this *ModelCache) Patch(k string, v interface{}) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	item, ok := this.cache[k]
	//If the item does not exist in the cache
	if !ok {
		return nil
	}
	//Clone the existing item
	itemClone := this.cloner.Clone(item)
	//Create a new updater
	putUpdater := updater.NewUpdater(this.introspector, false)
	//update the item clone with the new element where nil is valid
	err := putUpdater.Update(itemClone, v)
	if err != nil {
		return err
	}

	//if there are changes, then nothing to do
	changes := putUpdater.Changes()
	if changes == nil {
		return nil
	}

	//Apply the changes to the existing item
	for _, change := range changes {
		change.Apply(item)
	}
	go func() {
		if this.listener != nil {
			for _, change := range changes {
				this.listener.PropertyChangeNotification(itemClone, change.PropertyId(), change.OldValue(), change.NewValue())
			}
		}
	}()
	return nil
}

func (this *ModelCache) Delete(k string) {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	item, ok := this.cache[k]
	if !ok {
		return
	}
	delete(this.cache, k)
	if this.listener != nil {
		go this.listener.ModelItemDeleted(item)
	}
}

func (this *ModelCache) Attributes(f func(interface{}) interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	for k, v := range this.cache {
		result[k] = f(v)
	}
	return result
}
