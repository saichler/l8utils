package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/tests/infra"
	"reflect"
	"testing"
	"time"
)

func TestRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	ok, err := globals.Registry().Register(nil)
	if err == nil {
		log.Fail("Expected an error for nil type")
	}

	ok, err = globals.Registry().Register(&TestProto{})
	if !ok || err != nil {
		log.Fail("Expected to register a proto successfully")
	}

	ok, err = globals.Registry().Register(TestProto{})
	if ok {
		log.Fail(t, "Type should have already been registered")
		return
	}
	typ, err := globals.Registry().Info(protoName)
	if err != nil {
		log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		log.Fail(t, "Wrong type by name")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		log.Fail(t, "Expected an error")
		return
	}
	info, err := globals.Registry().Info(protoName)
	if err != nil {
		log.Fail(t, "Failed to get type by name", err.Error())
		return
	}
	ins, err := info.NewInstance()
	if err != nil {
		log.Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		log.Fail(t, "Failed to cast instance")
		return
	}
	_, err = globals.Registry().Info(unknowProtoName)
	if err == nil {
		log.Fail(t, "Expected an error")
		return
	}

	info, err = globals.Registry().Info(protoName)
	if err != nil {
		log.Fail(t, "Failed to get type by name", err.Error())
		return
	}

	if info.Type() == nil || info.Type().Name() != protoName {
		log.Fail(t, "Wrong type by name")
		return
	}

	info.AddSerializer(&infra.TestSerializer{})
	ser := info.Serializer(interfaces.BINARY)

	if ser == nil {
		log.Fail(t, "Failed to create serializer")
		return
	}

	if reflect.ValueOf(ser).Elem().Type().Name() != "TestSerializer" {
		log.Fail(t, "Wrong type by name")
		return
	}

	pb, err := info.NewInstance()
	if err != nil {
		log.Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		log.Fail(t, "Failed to cast protobuf instance")
		return
	}

	i, e := registry.NewInfo(nil)
	defer time.Sleep(time.Second)

	if e == nil {
		log.Fail(t, "Expected an error")
		return
	}

	if i != nil {
		log.Fail(t, "Expected nil instance")
		return
	}

	b, e := globals.Registry().RegisterType(nil)
	if e == nil {
		log.Fail(t, "Expected an error")
		return
	}
	if b {
		log.Fail(t, "Expected a false")
		return
	}

	b, e = globals.Registry().RegisterType(reflect.ValueOf(reflect.TypeOf(5)).Type().Elem()szc)
	if e != nil {
		log.Fail(t, "Did not expect an error")
		return
	}
	if b {
		log.Fail(t, "Expected a false")
		return
	}
}
