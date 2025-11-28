package cache

import (
	"errors"

	"github.com/saichler/l8types/go/types/l8notify"
)

func (this *Cache) Delete(v interface{}, createNotification bool) (*l8notify.L8NotificationSet, error) {
	pk, uk, err := this.KeysFor(v)
	if err != nil {
		return nil, err
	}
	if pk == "" {
		return nil, errors.New("Interface does not contain the Key attributes")
	}

	this.mtx.Lock()
	defer this.mtx.Unlock()

	var n *l8notify.L8NotificationSet
	var e error
	var item interface{}
	var ok bool

	if this.cacheEnabled() {
		item, ok = this.iCache.delete(pk, uk)
		if !ok {
			return n, errors.New("Delete Key " + pk + " not found")
		}
	}

	if this.store != nil {
		item, e = this.store.Delete(pk)
		if e != nil {
			return n, e
		}
	}

	if !createNotification {
		return n, e
	}

	n, e = this.createDeleteNotification(item, pk)
	return n, e
}
