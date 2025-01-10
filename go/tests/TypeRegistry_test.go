package tests

import (
	. "github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestTypeRegistry(t *testing.T) {
	protoName := "TestProto"
	unknowProtoName := "UnknowProto"

	TypeRegistry().Register(&TestProto{})
	ok := TypeRegistry().Register(TestProto{})
	if ok {
		Fail(t, "Type should have already been registered")
		return
	}
	typ, err := TypeRegistry().TypeInfo(protoName)
	if err != nil {
		Fail(t, "Failed to get type by name", err.Error())
		return
	}
	if typ.Name() != protoName {
		Fail(t, "Wrong type by name")
		return
	}
	_, err = TypeRegistry().TypeInfo(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}
	info, err := TypeRegistry().TypeInfo(protoName)
	if err != nil {
		Fail(t, "Failed to get type by name", err.Error())
		return
	}
	ins, err := info.NewInstance()
	if err != nil {
		Fail(t, "Failed to create instance", err.Error())
		return
	}
	_, ok = ins.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast instance")
		return
	}
	_, err = TypeRegistry().TypeInfo(unknowProtoName)
	if err == nil {
		Fail(t, "Expected an error")
		return
	}

	info, err = TypeRegistry().TypeInfo(protoName)
	if err != nil {
		Fail(t, "Failed to get type by name", err.Error())
		return
	}
	pb, err := info.NewInstance()
	if err != nil {
		Fail(t, "Failed to create protobuf instance", err.Error())
		return
	}
	_, ok = pb.(*TestProto)
	if !ok {
		Fail(t, "Failed to cast protobuf instance")
		return
	}
}
