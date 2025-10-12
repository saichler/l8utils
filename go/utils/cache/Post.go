package cache

import (
	"errors"

	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8types/go/types/l8notify"
)

func (this *Cache) Post(v interface{}, createNotification bool) (*l8notify.L8NotificationSet, error) {
	k, err := this.PrimaryKeyFor(v)
	if err != nil {
		return nil, err
	}
	if k == "" {
		return nil, errors.New("Interface does not contain the Key attributes")
	}

	//Make sure we clone the input value, so the caller don't have a reference to the cache element
	v = cloner.Clone(v)

	var n *l8notify.L8NotificationSet
	var e error
	var item interface{}
	var ok bool

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if this.cacheEnabled() {
		item, ok = this.iCache.get(k)
	} else {
		item, e = this.store.Get(k)
		ok = e == nil
	}

	//If the item does not exist in the cache
	if !ok {
		//First clone the value so we can use it in the notification.
		itemClone := cloner.Clone(v)
		if this.cacheEnabled() {
			//Place the value in the cache
			this.iCache.put(k, v)
		}
		if this.store != nil {
			e = this.store.Put(k, v)
			if e != nil {
				return n, e
			}
		}
		//Create the notification using the clone outside the current go routine
		if createNotification {
			n, e = this.createAddNotification(itemClone, k)
			return n, e
		}
		return n, e
	}

	//From this point onwards, it the Put implementation, e.g. item exist and is full
	//Clone the instance so it won't be able to be updated outside the scope
	vClone := cloner.Clone(v)

	if this.cacheEnabled() {
		//Place the value in the cache
		this.iCache.put(k, vClone)
	}

	if this.store != nil {
		e = this.store.Put(k, vClone)
		if e != nil {
			return n, e
		}
	}

	if !createNotification {
		return n, e
	}

	//From this point onward, the item is no longer in the cache
	//so we don't need to clone it

	//Create a new updater
	putUpdater := updating.NewUpdater(this.r, true, true)

	//update the item clone with the new element where nil is valid
	e = putUpdater.Update(item, vClone)
	if e != nil {
		return n, e
	}

	//if there are changes, then nothing to do
	changes := putUpdater.Changes()
	if changes == nil || len(changes) == 0 {
		return nil, nil
	}

	n, e = this.createReplaceNotification(item, v, k)
	return n, e
}
