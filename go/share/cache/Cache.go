package cache

import (
	"errors"
	"github.com/saichler/reflect/go/reflect/clone"
	"github.com/saichler/reflect/go/reflect/common"
	"github.com/saichler/reflect/go/reflect/updater"
	"sync"
)

type Cache struct {
	cache        map[string]interface{}
	mtx          *sync.RWMutex
	cond         *sync.Cond
	listener     ICacheListener
	cloner       *clone.Cloner
	introspector common.IIntrospect
}

type ICacheListener interface {
	ModelItemAdded(interface{})
	ModelItemDeleted(interface{})
	PropertyChangeNotification(interface{}, string, interface{}, interface{})
}

func NewModelCache(listener ICacheListener, introspector common.IIntrospect) *Cache {
	this := &Cache{}
	this.cache = make(map[string]interface{})
	this.mtx = &sync.RWMutex{}
	this.cond = sync.NewCond(this.mtx)
	this.listener = listener
	this.cloner = clone.NewCloner()
	this.introspector = introspector
	return this
}

func (this *Cache) Get(k string) interface{} {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	item, ok := this.cache[k]
	if ok {
		itemClone := this.cloner.Clone(item)
		return itemClone
	}
	return nil
}

func (this *Cache) Put(k string, v interface{}) error {
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

	go func() {
		if this.listener != nil {
			for _, change := range changes {
				this.listener.PropertyChangeNotification(itemClone, change.PropertyId(), change.OldValue(), change.NewValue())
			}
		}
	}()
	return nil
}

func (this *Cache) Patch(k string, v interface{}) error {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	item, ok := this.cache[k]
	//If the item does not exist in the cache
	if !ok {
		return errors.New("Key " + k + " not found")
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

func (this *Cache) Delete(k string) {
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

func (this *Cache) Collect(f func(interface{}) interface{}) map[string]interface{} {
	result := map[string]interface{}{}
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	for k, v := range this.cache {
		result[k] = f(v)
	}
	return result
}
