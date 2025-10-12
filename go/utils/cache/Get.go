package cache

import (
	"errors"

	"github.com/saichler/l8utils/go/utils/strings"
)

func (this *Cache) Get(v interface{}) (interface{}, error) {
	var item interface{}
	var e error
	var k string
	var ok bool

	k, e = this.PrimaryKeyFor(v)
	if e != nil {
		return item, e
	}

	if k == "" {
		e = errors.New("Interface does not contain the Key attributes")
		return item, e
	}

	this.mtx.RLock()
	defer this.mtx.RUnlock()

	if this.cacheEnabled() {
		item, ok = this.iCache.get(k)
		if ok {
			itemClone := cloner.Clone(item)
			return itemClone, e
		}
	} else {
		item, e = this.store.Get(k)
		if e == nil {
			return item, e
		}
		e = errors.New(strings.New("Cache:", this.serviceName, ":", this.serviceArea, " ", e.Error()).String())
		return item, e
	}
	e = errors.New("Not found in the cache")
	return item, e
}
