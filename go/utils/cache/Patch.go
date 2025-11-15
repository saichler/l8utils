package cache

import (
	"errors"

	"github.com/saichler/l8reflect/go/reflect/updating"
	"github.com/saichler/l8types/go/types/l8notify"
)

func (this *Cache) Patch(v interface{}, createNotification bool) (*l8notify.L8NotificationSet, error) {
	k, err := this.PrimaryKeyFor(v)
	if err != nil {
		return nil, err
	}
	if k == "" {
		return nil, errors.New("Interface does not contain the Key attributes")
	}

	this.mtx.Lock()
	defer this.mtx.Unlock()
	var n *l8notify.L8NotificationSet
	var e error
	var item interface{}
	var ok bool

	if this.cacheEnabled() {
		item, ok = this.iCache.get(k)
	} else {
		item, e = this.store.Get(k)
		ok = e == nil
	}

	//If the item does not exist in the cache
	if !ok {
		//Clone the value for the cache
		vClone := cloner.Clone(v)

		if this.cacheEnabled() {
			//Place the new Item clone in the cache
			this.iCache.put(k, vClone)
		}

		if this.store != nil {
			//place the new item clone in the store
			e = this.store.Put(k, vClone)
		}

		if !createNotification {
			return n, e
		}

		//Clone the value for the notification
		itemClone := cloner.Clone(v)
		n, e = this.createAddNotification(itemClone, k)
		return n, e
	}

	//Clone the existing item
	itemClone := cloner.Clone(item)
	//Create a new updater
	patchUpdater := updating.NewUpdater(this.r, false, false)
	//update the item clone with the new element where nil is valid
	e = patchUpdater.Update(itemClone, v)
	if e != nil {
		return n, e
	}

	//if there are no changes, then nothing to do
	changes := patchUpdater.Changes()
	if changes == nil {
		return n, e
	}

	//Remove the item from the stats to make sure if one of the attributes
	//that are going to change affect the stats
	this.iCache.removeFromCounts(k)

	//Apply the changes to the existing item in the cache
	for _, change := range changes {
		change.Apply(item)
	}

	//Re-Add the item to the stats
	this.iCache.addToCounts(item)

	if this.store != nil {
		e = this.store.Put(k, item)
	}

	if !createNotification {
		return n, e
	}

	n, e = this.createUpdateNotification(changes, k)
	return n, e
}
