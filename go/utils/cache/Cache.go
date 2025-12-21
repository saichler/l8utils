package cache

import (
	"errors"
	"reflect"
	"sync"

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8reflect"
)

var cloner = cloning.NewCloner()

type Cache struct {
	iCache               *internalCache
	mtx                  *sync.RWMutex
	cond                 *sync.Cond
	store                ifs.IStorage
	modelType            string
	primaryKeyFieldNames []string
	uniqueKeyFieldNames  []string
	r                    ifs.IResources

	notifySequence uint32
	serviceName    string
	serviceArea    byte
	cleaner        *ttlCleaner
}

func NewCache(sampleElement interface{}, initElements []interface{}, store ifs.IStorage, r ifs.IResources) *Cache {
	this := &Cache{}
	this.iCache = newInternalCache()
	this.mtx = &sync.RWMutex{}
	this.cond = sync.NewCond(this.mtx)
	this.store = store
	this.r = r
	this.modelType = reflect.ValueOf(sampleElement).Elem().Type().Name()

	_, _, err := this.KeysFor(sampleElement)
	if err != nil {
		panic("Error in initialized elements " + err.Error())
	}

	loadedFromStore := false

	if this.store != nil {
		items := this.store.Collect(allElementsInCache)
		for _, v := range items {
			pk, uk, _ := this.KeysFor(v)
			this.iCache.put(pk, uk, v)
		}
		if len(items) > 0 {
			loadedFromStore = true
		}
	}

	if !loadedFromStore && this.store != nil {
		for _, item := range initElements {
			pk, _, er := this.KeysFor(item)
			if er != nil {
				r.Logger().Error(er.Error())
				continue
			}
			this.store.Put(pk, item)
		}
	}

	if !loadedFromStore && this.cacheEnabled() && initElements != nil {
		for _, item := range initElements {
			pk, uk, er := this.KeysFor(item)
			if er != nil {
				r.Logger().Error("#2 Init item", " error:", er.Error())
				continue
			}
			this.iCache.put(pk, uk, item)
		}
	}
	addTotalMetadata(this)

	// Start TTL cleaner for query cache
	this.cleaner = newTTLCleaner(this)
	this.cleaner.start()

	return this
}

func (this *Cache) SetNotificationsFor(serviceName string, serviceArea byte) {
	this.serviceName = serviceName
	this.serviceArea = serviceArea
}

func (this *Cache) cacheEnabled() bool {
	if this.store == nil {
		return true
	}
	return this.store.CacheEnabled()
}

func (this *Cache) Size() int {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	return this.iCache.size()
}

func (this *Cache) KeysFor(any interface{}) (string, string, error) {
	if any == nil {
		return "", "", errors.New("Cannot get keys for nil interface")
	}

	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Ptr {
		return "", "", errors.New("Cannot get keys for non-pointer interface")
	}
	v = v.Elem()

	if this.primaryKeyFieldNames == nil {
		node, _, err := this.r.Introspector().Decorators().NodeFor(any)
		if err != nil {
			return "", "", err
		}
		this.primaryKeyFieldNames, err = this.r.Introspector().Decorators().Fields(node, l8reflect.L8DecoratorType_Primary)
		if err != nil {
			return "", "", err
		}
		this.uniqueKeyFieldNames, err = this.r.Introspector().Decorators().Fields(node, l8reflect.L8DecoratorType_Unique)
	}

	pkValue, err := this.r.Introspector().Decorators().KeyForValue(this.primaryKeyFieldNames, v, this.modelType, true)
	ukValue, _ := this.r.Introspector().Decorators().KeyForValue(this.uniqueKeyFieldNames, v, this.modelType, false)
	return pkValue, ukValue, err
}

func allElementsInCache(i interface{}) (bool, interface{}) {
	return true, i
}

func (this *Cache) ServiceName() string {
	return this.serviceName
}

func (this *Cache) ServiceArea() byte {
	return this.serviceArea
}

func (this *Cache) Source() string {
	return this.r.SysConfig().LocalUuid
}

func (this *Cache) ModelType() string {
	return this.modelType
}

func (this *Cache) Close() {
	if this.cleaner != nil {
		this.cleaner.stop()
	}
}
