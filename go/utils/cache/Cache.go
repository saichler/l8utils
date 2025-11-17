package cache

import (
	"errors"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
)

var cloner = cloning.NewCloner()

type Cache struct {
	iCache        *internalCache
	mtx           *sync.RWMutex
	cond          *sync.Cond
	store         ifs.IStorage
	modelType     string
	keyFieldNames []string
	r             ifs.IResources

	notifySequence uint32
	serviceName    string
	serviceArea    byte
}

func NewCache(sampleElement interface{}, initElements []interface{}, store ifs.IStorage, r ifs.IResources) *Cache {
	this := &Cache{}
	this.iCache = newInternalCache()
	this.mtx = &sync.RWMutex{}
	this.cond = sync.NewCond(this.mtx)
	this.store = store
	this.r = r

	_, name, err := this.PrimaryKeyFor(sampleElement)
	if err != nil {
		panic("Error in initialized elements " + name + " " + err.Error())
	}

	loadedFromStore := false

	if this.store != nil {
		items := this.store.Collect(allElementsInCache)
		for k, v := range items {
			this.iCache.put(k, v)
		}
		if len(items) > 0 {
			loadedFromStore = true
		}
	}

	if !loadedFromStore && this.store != nil {
		for _, item := range initElements {
			k, n, er := this.PrimaryKeyFor(item)
			if er != nil {
				r.Logger().Error("#1 Init item", n, " error:", err.Error())
				continue
			}
			this.store.Put(k, item)
		}
	}

	if !loadedFromStore && this.cacheEnabled() && initElements != nil {
		for _, item := range initElements {
			k, n, er := this.PrimaryKeyFor(item)
			if er != nil {
				r.Logger().Error("#2 Init item", n, " error:", err.Error())
				continue
			}
			this.iCache.put(k, item)
		}
	}
	addTotalMetadata(this)
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

func (this *Cache) typeFor(any interface{}) (string, error) {
	if this.modelType != "" {
		return this.modelType, nil
	}
	if any == nil {
		debug.PrintStack()
		return "", errors.New("Cannot get type for nil interface")
	}
	v := reflect.ValueOf(any)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	this.modelType = v.Type().Name()
	return this.modelType, nil
}

func (this *Cache) PrimaryKeyFor(any interface{}) (string, string, error) {
	if any == nil {
		return "", "", errors.New("Cannot get key for nil interface")
	}

	v := reflect.ValueOf(any)
	if v.Kind() != reflect.Ptr {
		return "", "", errors.New("Cannot get key for non-pointer interface")
	}
	v = v.Elem()
	name := v.Type().Name()

	if this.keyFieldNames == nil {
		typ, err := this.typeFor(any)
		if err != nil {
			return "", name, err
		}
		node, ok := this.r.Introspector().Node(typ)
		if !ok {
			return "", name, errors.New("Could not find an interospector node for type " + typ)
		}
		pk := introspecting.PrimaryKeyDecorator(node)
		if pk == nil {
			return "", name, errors.New("No primary key decorator is defined for type " + typ)
		}
		this.keyFieldNames = pk.([]string)
	}

	if len(this.keyFieldNames) == 0 {
		return "", name, errors.New("Lost of keys is empty for type " + this.modelType)
	} else if len(this.keyFieldNames) == 1 {
		return strings.New(v.FieldByName(this.keyFieldNames[0]).Interface()).String(), name, nil
	} else if len(this.keyFieldNames) == 2 {
		return strings.New(v.FieldByName(this.keyFieldNames[0]).Interface(),
			v.FieldByName(this.keyFieldNames[1]).Interface()).String(), name, nil
	} else if len(this.keyFieldNames) == 3 {
		return strings.New(v.FieldByName(this.keyFieldNames[0]).Interface(),
			v.FieldByName(this.keyFieldNames[1]).Interface()).String(), name, nil
	}
	result := strings.New()
	for i := 0; i < len(this.keyFieldNames); i++ {
		result.Add(result.StringOf(v.FieldByName(this.keyFieldNames[0]).Interface()))
	}
	return result.String(), name, nil
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
