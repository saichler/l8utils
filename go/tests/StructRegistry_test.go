package tests

import (
	. "github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestStructRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	StructRegistry().RegisterStruct(&TestProto{}, nil)
	ok := StructRegistry().RegisterStruct(TestProto{}, nil)
	if ok {
		Fail(t, "Type should have already been registered")
		return
	}
	typ, _, err := StructRegistry().TypeByName(protoName)
	if err != nil {
		Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Fail(t, "Wrong type by name")
		return
	}
	_, _, err = StructRegistry().TypeByName(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}

	ins, _, err := StructRegistry().NewInstance(protoName)
	if err != nil {
		Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast instance")
		return
	}
	_, _, err = StructRegistry().NewInstance(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}

	pb, _, err := StructRegistry().NewInstance(protoName)
	if err != nil {
		Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast protobuf instance")
		return
	}
	_, _, err = StructRegistry().NewInstance(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
	_, _, err = StructRegistry().NewInstance("")
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
}
