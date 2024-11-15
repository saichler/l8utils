package tests

import (
	. "github.com/saichler/shared/go/interfaces"
	"github.com/saichler/shared/go/maps"
	"reflect"

	"testing"
)

func TestSyncMap(t *testing.T) {
	key := "key"
	val := "val"
	m := maps.NewSyncMap()
	m.Put(key, val)
	v, ok := m.Get(key)
	if !ok {
		Fail(t, "Expected key to exist")
		return
	}
	if v != val {
		Fail(t, "Expected value to be '"+val+"'")
		return
	}

	if m.Size() != 1 {
		Fail(t, "Expected size to be 1")
		return
	}

	m.Clean()
	if m.Size() != 0 {
		Fail(t, "Expected size to be 0")
		return
	}

	m.Put(key, val)

	v, ok = m.Delete(key)
	if !ok {
		Fail(t, "Expected key to exist")
		return
	}
	if v != val {
		Fail(t, "Expected value to be '"+val+"'")
		return
	}

	if m.Contains(key) {
		Fail(t, "Expected key '"+key+" to NOT exist")
	}

	m.Put("a", "b")
	m.Put("c", "d")
	m.Put("e", "f")

	vFilter := func(filter interface{}) bool {
		k := filter.(string)
		if k == "d" {
			return false
		}
		return true
	}

	l := m.ValuesAsList(reflect.ValueOf(val).Type(), vFilter)
	list := l.([]string)

	if len(list) != 2 {
		Fail(t, "Expected length of list to be 2, but it is:", len(list))
		return
	}

	if !m.Contains("a") || !m.Contains("e") {
		Fail(t, "Expected 'a' & 'e' keys to exist")
		return
	}

	l = m.ValuesAsList(reflect.ValueOf(val).Type(), nil)
	list = l.([]string)

	if len(list) != 3 {
		Fail(t, "Expected length of list to be 3, but it is:", len(list))
		return
	}

	if !m.Contains("a") || !m.Contains("e") || !m.Contains("c") {
		Fail(t, "Expected 'a', 'c' & 'e' keys to exist")
		return
	}

	l = m.KeysAsList(reflect.ValueOf(val).Type(), nil)
	list = l.([]string)
	if len(list) != 3 {
		Fail(t, "Expected length of list to be 3, but it is:", len(list))
		return
	}

	kFilter := func(filter interface{}) bool {
		k := filter.(string)
		if k == "c" {
			return false
		}
		return true
	}

	l = m.KeysAsList(reflect.ValueOf(val).Type(), kFilter)
	list = l.([]string)
	if len(list) != 2 {
		Fail(t, "Expected length of list to be 2, but it is:", len(list))
		return
	}

	itf := func(key interface{}, val interface{}) {
		Debug("key:", key, " val:", val)
	}

	m.Iterate(itf)
}
