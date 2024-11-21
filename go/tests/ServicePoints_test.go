package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/tests/infra"
	"github.com/saichler/shared/go/types"
	"testing"
)

func TestServicePoints(t *testing.T) {
	testsp := infra.NewTestServicePointHandler("testsp")
	pb := &TestProto{}
	err := interfaces.ServicePoints().RegisterServicePoint(nil, testsp, interfaces.StructRegistry())
	if err == nil {
		interfaces.Fail("Expected an error")
		return
	}
	err = interfaces.ServicePoints().RegisterServicePoint(pb, nil, interfaces.StructRegistry())
	if err == nil {
		interfaces.Fail("Expected an error")
		return
	}
	err = interfaces.ServicePoints().RegisterServicePoint(pb, testsp, interfaces.StructRegistry())
	if err != nil {
		interfaces.Fail(t, err)
		return
	}
	sp, ok := interfaces.ServicePoints().ServicePointHandler("TestProto")
	if !ok {
		interfaces.Fail(t, "Service Point Not Found")
		return
	}
	sp.Topic()
	interfaces.ServicePoints().Handle(pb, types.Action_POST, nil)
	interfaces.ServicePoints().Handle(pb, types.Action_PUT, nil)
	interfaces.ServicePoints().Handle(pb, types.Action_DELETE, nil)
	interfaces.ServicePoints().Handle(pb, types.Action_GET, nil)
	interfaces.ServicePoints().Handle(pb, types.Action_PATCH, nil)
	interfaces.ServicePoints().Handle(pb, types.Action_Invalid_Action, nil)
	if testsp.PostNumber != 1 {
		interfaces.Fail(t, "Post is not 1")
	}
	if testsp.PutNumber != 1 {
		interfaces.Fail(t, "Put is not 1")
	}
	if testsp.DeleteNumber != 1 {
		interfaces.Fail(t, "Delete is not 1")
	}
	if testsp.PatchNumber != 1 {
		interfaces.Fail(t, "Patch is not 1")
	}
	if testsp.GetNumber != 1 {
		interfaces.Fail(t, "Get is not 1")
	}
}
