package tests

import (
	"testing"
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
}
