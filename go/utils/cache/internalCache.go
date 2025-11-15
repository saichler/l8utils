package cache

import (
	"reflect"
	"strings"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

type internalCache struct {
	cache       map[string]interface{}
	addedOrder  []string
	key2order   map[string]int
	stamp       int64
	queries     map[string]*internalQuery
	countFunc   map[string]func(interface{}) (bool, string)
	globalCount *l8api.L8Counts
}

func newInternalCache() *internalCache {
	iq := &internalCache{}
	iq.cache = make(map[string]interface{})
	iq.addedOrder = make([]string, 0)
	iq.key2order = make(map[string]int)
	iq.queries = make(map[string]*internalQuery)
	iq.globalCount = newCounts()
	return iq
}

func newCounts() *l8api.L8Counts {
	counts := &l8api.L8Counts{}
	counts.Counts = &l8api.L8ValueCount{}
	counts.Counts.Counts = make(map[string]int32)
	counts.ValueCounts = make(map[string]*l8api.L8ValueCount)
	return counts
}

func (this *internalCache) removeFromCounts(key string) (interface{}, bool) {
	old, ok := this.cache[key]
	if ok && this.countFunc != nil {
		for count, f := range this.countFunc {
			ok1, v := f(old)
			if ok1 {
				this.globalCount.Counts.Counts[count]--
				if v != "" {
					vCount, ok2 := this.globalCount.ValueCounts[count]
					if !ok2 {
						vCount = &l8api.L8ValueCount{}
						vCount.Counts = make(map[string]int32)
						this.globalCount.ValueCounts[count] = vCount
					}
					vCount.Counts[v]--
				}
			}
		}
	}
	return old, ok
}

func (this *internalCache) addToCounts(value interface{}) {
	if this.countFunc != nil {
		for count, f := range this.countFunc {
			ok1, v := f(value)
			if ok1 {
				this.globalCount.Counts.Counts[count]++
				if v != "" {
					vCount, ok2 := this.globalCount.ValueCounts[count]
					if !ok2 {
						vCount = &l8api.L8ValueCount{}
						vCount.Counts = make(map[string]int32)
						this.globalCount.ValueCounts[count] = vCount
					}
					vCount.Counts[v]++
				}
			}
		}
	}
}

func addToCounts(value interface{}, countFunc map[string]func(interface{}) (bool, string), counts *l8api.L8Counts) {
	if countFunc != nil {
		for count, f := range countFunc {
			ok1, v := f(value)
			if ok1 {
				counts.Counts.Counts[count]++
				if v != "" {
					vCount, ok2 := counts.ValueCounts[count]
					if !ok2 {
						vCount = &l8api.L8ValueCount{}
						vCount.Counts = make(map[string]int32)
						counts.ValueCounts[count] = vCount
					}
					vCount.Counts[v]++
				}
			}
		}
	}
}

func (this *internalCache) put(key string, value interface{}) {
	_, ok := this.removeFromCounts(key)
	this.cache[key] = value
	if !ok {
		this.addedOrder = append(this.addedOrder, key)
		this.stamp = time.Now().Unix()
		this.key2order[key] = len(this.addedOrder) - 1
	}
	this.addToCounts(value)
}

func (this *internalCache) get(key string) (interface{}, bool) {
	item, ok := this.cache[key]
	return item, ok
}

func (this *internalCache) delete(key string) (interface{}, bool) {
	item, ok := this.removeFromCounts(key)
	delete(this.cache, key)
	this.stamp = time.Now().Unix()
	return item, ok
}

func (this *internalCache) size() int {
	return len(this.cache)
}

func (this *internalCache) fetch(start, blockSize int, q ifs.IQuery) ([]interface{}, *l8api.L8Counts) {

	noCriteriaOrSort := true

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
		noCriteriaOrSort = (!qrt.IsValid() || qrt.IsNil()) && strings.TrimSpace(q.SortBy()) == ""
		//Query does not have criteria so use the added order for keys
		if noCriteriaOrSort {
			dq.prepare(this.cache, this.addedOrder, this.stamp, this.countFunc)
		} else {
			//We need to create a dataset sorted by the sortby and filter by the criteria
			dq.prepare(this.cache, nil, this.stamp, this.countFunc)
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
	if !noCriteriaOrSort {
		return result, dq.counts
	}
	return result, this.globalCount
}

func (this *internalCache) addCountFunc(name string, f func(interface{}) (bool, string)) {
	if this.countFunc == nil {
		this.countFunc = make(map[string]func(interface{}) (bool, string))
	}
	this.countFunc[name] = f
	if len(this.cache) > 0 {
		for _, elem := range this.cache {
			ok1, v := f(elem)
			if ok1 {
				this.globalCount.Counts.Counts[name]++
				if v != "" {
					vCount, ok2 := this.globalCount.ValueCounts[name]
					if !ok2 {
						vCount = &l8api.L8ValueCount{}
						vCount.Counts = make(map[string]int32)
						this.globalCount.ValueCounts[name] = vCount
					}
					vCount.Counts[v]++
				}
			}
		}
	}
}
