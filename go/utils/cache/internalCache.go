package cache

import (
	"reflect"
	"strings"
	"time"

	"github.com/saichler/l8types/go/ifs"
)

type internalCache struct {
	cache      map[string]interface{}
	addedOrder []string
	key2order  map[string]int
	stamp      int64
	queries    map[string]*internalQuery
	stats      map[string]int32
	statsFunc  map[string]func(interface{}) bool
}

func newInternalCache() *internalCache {
	iq := &internalCache{}
	iq.cache = make(map[string]interface{})
	iq.addedOrder = make([]string, 0)
	iq.key2order = make(map[string]int)
	iq.queries = make(map[string]*internalQuery)
	iq.stats = make(map[string]int32)
	iq.statsFunc = make(map[string]func(interface{}) bool)
	return iq
}

func (this *internalCache) removeFromStats(key string) (interface{}, bool) {
	old, ok := this.cache[key]
	if ok && this.statsFunc != nil {
		for stat, f := range this.statsFunc {
			if f(old) {
				this.stats[stat]--
			}
		}
	}
	return old, ok
}

func (this *internalCache) addToStats(value interface{}) {
	if this.statsFunc != nil {
		for stat, f := range this.statsFunc {
			if f(value) {
				this.stats[stat]++
			}
		}
	}
}

func (this *internalCache) put(key string, value interface{}) {
	_, ok := this.removeFromStats(key)
	this.cache[key] = value
	if !ok {
		this.addedOrder = append(this.addedOrder, key)
		this.stamp = time.Now().Unix()
		this.key2order[key] = len(this.addedOrder) - 1
	}
	this.addToStats(value)
}

func (this *internalCache) get(key string) (interface{}, bool) {
	item, ok := this.cache[key]
	return item, ok
}

func (this *internalCache) delete(key string) (interface{}, bool) {
	item, ok := this.removeFromStats(key)
	delete(this.cache, key)
	this.stamp = time.Now().Unix()
	return item, ok
}

func (this *internalCache) size() int {
	return len(this.cache)
}

func (this *internalCache) fetch(start, blockSize int, q ifs.IQuery) []interface{} {

	dq, ok := this.queries[q.Hash()]
	//This is a new query, so create it
	if !ok {
		dq = newInternalQuery(q)
		this.queries[q.Hash()] = dq
	}

	//If the query timestamp has changed, it means elements were added/removed
	//so we need re-cresate the sorted set
	if dq.stamp != this.stamp {
		qrt := reflect.ValueOf(q.Criteria())
		//Query does not have criteria so use the added order for keys
		if (!qrt.IsValid() || qrt.IsNil()) && strings.TrimSpace(q.SortBy()) == "" {
			dq.prepare(this.cache, this.addedOrder, this.stamp)
		} else {
			//We need to create a dataset sorted by the sortby and filter by the criteria
			dq.prepare(this.cache, nil, this.stamp)
		}
	}

	//return just the subset of rows requested
	result := make([]interface{}, 0)
	for i := start; i < len(dq.data); i++ {
		key := dq.data[i]
		value, ok := this.cache[key]
		if ok {
			result = append(result, value)
		}
		if blockSize == 0 {
			continue
		}
		if len(result) >= blockSize {
			break
		}
	}
	return result
}

func (this *internalCache) addStatsFunc(name string, f func(interface{}) bool) {
	if this.statsFunc == nil {
		this.statsFunc = make(map[string]func(interface{}) bool)
		this.stats = make(map[string]int32)
	}
	this.statsFunc[name] = f
	if len(this.cache) > 0 {
		for _, elem := range this.cache {
			if f(elem) {
				this.stats[name]++
			}
		}
	}
}
