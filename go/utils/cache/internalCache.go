package cache

import (
	"reflect"
	"strings"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8api"
)

type internalCache struct {
	cache          map[string]interface{}
	addedOrder     []string
	key2order      map[string]int
	stamp          int64
	queries        map[string]*internalQuery
	metadataFunc   map[string]func(interface{}) (bool, string)
	globalMetadata *l8api.L8MetaData
}

func newInternalCache() *internalCache {
	iq := &internalCache{}
	iq.cache = make(map[string]interface{})
	iq.addedOrder = make([]string, 0)
	iq.key2order = make(map[string]int)
	iq.queries = make(map[string]*internalQuery)
	iq.globalMetadata = newMetadata()
	return iq
}

func newMetadata() *l8api.L8MetaData {
	metadata := &l8api.L8MetaData{}
	metadata.KeyCount = &l8api.L8Count{}
	metadata.KeyCount.Counts = make(map[string]int32)
	metadata.ValueCount = make(map[string]*l8api.L8Count)
	return metadata
}

func (this *internalCache) removeFromMetadata(key string) (interface{}, bool) {
	old, ok := this.cache[key]
	if ok && this.metadataFunc != nil {
		for count, f := range this.metadataFunc {
			ok1, v := f(old)
			if ok1 {
				this.globalMetadata.KeyCount.Counts[count]--
				if v != "" {
					vCount, ok2 := this.globalMetadata.ValueCount[count]
					if !ok2 {
						vCount = &l8api.L8Count{}
						vCount.Counts = make(map[string]int32)
						this.globalMetadata.ValueCount[count] = vCount
					}
					vCount.Counts[v]--
				}
			}
		}
	}
	return old, ok
}

func (this *internalCache) addToMetadata(value interface{}) {
	addToMetadata(value, this.metadataFunc, this.globalMetadata)
}

func addToMetadata(value interface{}, metadataFunc map[string]func(interface{}) (bool, string), metadata *l8api.L8MetaData) {
	if metadataFunc != nil {
		for count, f := range metadataFunc {
			ok1, v := f(value)
			if ok1 {
				metadata.KeyCount.Counts[count]++
				if v != "" {
					vCount, ok2 := metadata.ValueCount[count]
					if !ok2 {
						vCount = &l8api.L8Count{}
						vCount.Counts = make(map[string]int32)
						metadata.ValueCount[count] = vCount
					}
					vCount.Counts[v]++
				}
			}
		}
	}
}

func (this *internalCache) put(key string, value interface{}) {
	_, ok := this.removeFromMetadata(key)
	this.cache[key] = value
	if !ok {
		this.addedOrder = append(this.addedOrder, key)
		this.stamp = time.Now().Unix()
		this.key2order[key] = len(this.addedOrder) - 1
	}
	this.addToMetadata(value)
}

func (this *internalCache) get(key string) (interface{}, bool) {
	item, ok := this.cache[key]
	return item, ok
}

func (this *internalCache) delete(key string) (interface{}, bool) {
	item, ok := this.removeFromMetadata(key)
	delete(this.cache, key)
	this.stamp = time.Now().Unix()
	return item, ok
}

func (this *internalCache) size() int {
	return len(this.cache)
}

func (this *internalCache) fetch(start, blockSize int, q ifs.IQuery) ([]interface{}, *l8api.L8MetaData) {

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
			dq.prepare(this.cache, this.addedOrder, this.stamp, this.metadataFunc)
		} else {
			//We need to create a dataset sorted by the sortby and filter by the criteria
			dq.prepare(this.cache, nil, this.stamp, this.metadataFunc)
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
		return result, dq.metadata
	}
	return result, this.globalMetadata
}

func (this *internalCache) addMetadataFunc(name string, f func(interface{}) (bool, string)) {
	if this.metadataFunc == nil {
		this.metadataFunc = make(map[string]func(interface{}) (bool, string))
	}
	this.metadataFunc[name] = f
	if len(this.cache) > 0 {
		for _, elem := range this.cache {
			ok1, v := f(elem)
			if ok1 {
				this.globalMetadata.KeyCount.Counts[name]++
				if v != "" {
					vCount, ok2 := this.globalMetadata.ValueCount[name]
					if !ok2 {
						vCount = &l8api.L8Count{}
						vCount.Counts = make(map[string]int32)
						this.globalMetadata.ValueCount[name] = vCount
					}
					vCount.Counts[v]++
				}
			}
		}
	}
}
