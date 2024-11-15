package tests

import (
	. "github.com/saichler/shared/go/interfaces"
	"testing"
)

func TestStructRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	ok := StructRegistry().RegisterStruct(&TestProto{})
	if !ok {
		Fail(t, "#1 Failed to register struct")
		return
	}
	ok = StructRegistry().RegisterStruct(TestProto{})
	if ok {
		Fail(t, "#2 Type should have already been registered")
		return
	}
	typ, err := StructRegistry().TypeByName(protoName)
	if err != nil {
		Fail(t, "#3 Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Fail(t, "#4 Wrong type by name")
		return
	}
	_, err = StructRegistry().TypeByName(unknowProtoName)
	if err == nil {
		Fail(t, "#4.1 Expected an error")
		return
	}

	ins, err := StructRegistry().NewInstance(protoName)
	if err != nil {
		Fail(t, "#5 Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Fail(t, "#6 Failed to cast instance")
		return
	}
	_, err = StructRegistry().NewInstance(unknowProtoName)
	if err == nil {
		Fail(t, "#6.1 Expected an error")
		return
	}

	pb, err := StructRegistry().NewProtobufInstance(protoName)
	if err != nil {
		Fail(t, "#7 Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Fail(t, "#8 Failed to cast protobuf instance")
		return
	}
	_, err = StructRegistry().NewProtobufInstance(unknowProtoName)
	if err == nil {
		Fail(t, "#8.1 Expected an error")
		return
	}
	_, err = StructRegistry().NewProtobufInstance("")
	if err == nil {
		Fail(t, "#8.2 Expected an error")
		return
	}
	type TT struct{}
	StructRegistry().RegisterStruct(TT{})
	_, err = StructRegistry().NewProtobufInstance("TT")
	if err == nil {
		Fail(t, "#8.3 Expected an error")
		return
	}
}
