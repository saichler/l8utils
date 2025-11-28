package cache

import (
	"errors"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"

	"github.com/saichler/l8reflect/go/reflect/cloning"
	"github.com/saichler/l8reflect/go/reflect/helping"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/strings"
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
}

func NewCache(sampleElement interface{}, initElements []interface{}, store ifs.IStorage, r ifs.IResources) *Cache {
	this := &Cache{}
	this.iCache = newInternalCache()
	this.mtx = &sync.RWMutex{}
	this.cond = sync.NewCond(this.mtx)
	this.store = store
	this.r = r

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
		fmt.Println("Stack trace:")
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
		typ, err := this.typeFor(any)
		if err != nil {
			return "", "", err
		}
		node, ok := this.r.Introspector().Node(typ)
		if !ok {
			return "", "", errors.New("Could not find an interospector node for type " + typ)
		}
		pk := helping.PrimaryKeyDecorator(node)
		uk := helping.UniqueKeyDecorator(node)
		if pk == nil {
			return "", "", errors.New("No primary key decorator is defined for type " + typ)
		}
		this.primaryKeyFieldNames = pk.([]string)
		this.uniqueKeyFieldNames, _ = uk.([]string)
	}

	pkValue, err := keyFor(this.primaryKeyFieldNames, v, this.modelType, true)
	ukValue, _ := keyFor(this.uniqueKeyFieldNames, v, this.modelType, false)
	return pkValue, ukValue, err
}

func keyFor(names []string, v reflect.Value, modelType string, returnError bool) (string, error) {
	if names == nil || len(names) == 0 {
		if returnError {
			return "", errors.New("Primary Key Decorator  is empty for type " + modelType)
		}
		return "", nil
	}
	switch len(names) {
	case 0:
		if returnError {
			return "", errors.New("Primary Key Decorator  is empty for type " + modelType)
		}
		return "", nil
	case 1:
		return strings.New(v.FieldByName(names[0]).Interface()).String(), nil
	case 2:
		strings.New(v.FieldByName(names[0]).Interface(), v.FieldByName(names[1]).Interface()).String()
	case 3:
		strings.New(v.FieldByName(names[0]).Interface(),
			v.FieldByName(names[1]).Interface(),
			v.FieldByName(names[2]).Interface()).String()
	default:
		result := strings.New()
		for i := 0; i < len(names); i++ {
			result.Add(result.StringOf(v.FieldByName(names[i]).Interface()))
		}
		return result.String(), nil
	}
	return "", errors.New("Unexpected code")
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
