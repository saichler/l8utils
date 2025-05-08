package tests

import (
	"github.com/saichler/l8types/go/ifs"
	. "github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/registry"
	"reflect"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	ok, err := globals.Registry().Register(nil)
	if err == nil {
		Log.Fail("Expected an error for nil type")
	}

	ok, err = globals.Registry().Register(&TestProto{})
	if !ok || err != nil {
		Log.Fail("Expected to register a proto successfully")
	}

	ok, err = globals.Registry().Register(TestProto{})
	if ok {
		Log.Fail(t, "Type should have already been registered")
		return
	}
	typ, err := globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Log.Fail(t, "Wrong type by name")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		Log.Fail(t, "Expected an error")
		return
	}
	info, err := globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	ins, err := info.NewInstance()
	if err != nil {
		Log.Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Log.Fail(t, "Failed to cast instance")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		Log.Fail(t, "Expected an error")
		return
	}

	info, err = globals.Registry().Info(protoName)
	if err != nil {
		Log.Fail(t, "Failed to get type by name", err.Error())
		return
	}

	if info.Type() == nil || info.Type().Name() != protoName {
		Log.Fail(t, "Wrong type by name")
		return
	}

	info.AddSerializer(&TestSerializer{})
	ser := info.Serializer(ifs.BINARY)

	if ser == nil {
		Log.Fail(t, "Failed to create serializer")
		return
	}

	if reflect.ValueOf(ser).Elem().Type().Name() != "TestSerializer" {
		Log.Fail(t, "Wrong type by name")
		return
	}

	pb, err := info.NewInstance()
	if err != nil {
		Log.Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Log.Fail(t, "Failed to cast protobuf instance")
		return
	}

	i, e := registry.NewInfo(nil)
	defer time.Sleep(time.Second)

	if e == nil {
		Log.Fail(t, "Expected an error")
		return
	}

	if i != nil {
		Log.Fail(t, "Expected nil instance")
		return
	}

	b, e := globals.Registry().RegisterType(nil)
	if e == nil {
		Log.Fail(t, "Expected an error")
		return
	}
	if b {
		Log.Fail(t, "Expected a false")
		return
	}

	b, e = globals.Registry().RegisterType(reflect.ValueOf(reflect.TypeOf(5)).Type().Elem())
	if e != nil {
		Log.Fail(t, "Did not expect an error")
		return
	}
	if b {
		Log.Fail(t, "Expected a false")
		return
	}
}
