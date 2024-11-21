package tests

import (
	. "github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestStructRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	StructRegistry().RegisterStruct(&TestProto{})
	ok := StructRegistry().RegisterStruct(TestProto{})
	if ok {
		Fail(t, "Type should have already been registered")
		return
	}
	typ, err := StructRegistry().TypeByName(protoName)
	if err != nil {
		Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Fail(t, "Wrong type by name")
		return
	}
	_, err = StructRegistry().TypeByName(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}

	ins, err := StructRegistry().NewInstance(protoName)
	if err != nil {
		Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast instance")
		return
	}
	_, err = StructRegistry().NewInstance(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}

	pb, err := StructRegistry().NewProtobufInstance(protoName)
	if err != nil {
		Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast protobuf instance")
		return
	}
	_, err = StructRegistry().NewProtobufInstance(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
	_, err = StructRegistry().NewProtobufInstance("")
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
	type TT struct{}
	StructRegistry().RegisterStruct(TT{})
	_, err = StructRegistry().NewProtobufInstance("TT")
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
}
