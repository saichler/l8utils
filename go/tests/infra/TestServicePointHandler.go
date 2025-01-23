package infra

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

var Log interfaces.ILogger

type TestServicePointHandler struct {
	Name         string
	PostNumber   int
	PutNumber    int
	PatchNumber  int
	DeleteNumber int
	GetNumber    int
	FailedNumber int
}

const (
	TEST_TOPIC = "TestProto"
)

func NewTestServicePointHandler(name string) *TestServicePointHandler {
	tsp := &TestServicePointHandler{}
	tsp.Name = name
	return tsp
}

func (tsp *TestServicePointHandler) Post(pb proto.Message, edge interfaces.IVirtualNetworkInterface) (proto.Message, error) {
	Log.Debug("Post -", tsp.Name, "- Test callback")
	tsp.PostNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Put(pb proto.Message, edge interfaces.IVirtualNetworkInterface) (proto.Message, error) {
	Log.Debug("Put -", tsp.Name, "- Test callback")
	tsp.PutNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Patch(pb proto.Message, edge interfaces.IVirtualNetworkInterface) (proto.Message, error) {
	Log.Debug("Patch -", tsp.Name, "- Test callback")
	tsp.PatchNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Delete(pb proto.Message, edge interfaces.IVirtualNetworkInterface) (proto.Message, error) {
	Log.Debug("Delete -", tsp.Name, "- Test callback")
	tsp.DeleteNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Get(pb proto.Message, edge interfaces.IVirtualNetworkInterface) (proto.Message, error) {
	Log.Debug("Get -", tsp.Name, "- Test callback")
	tsp.GetNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Failed(pb proto.Message, edge interfaces.IVirtualNetworkInterface, info *types.Message) (proto.Message, error) {
	Log.Debug("Failed -", tsp.Name, "- Test callback")
	tsp.FailedNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) EndPoint() string {
	return "/Tests"
}
func (tsp *TestServicePointHandler) Topic() string {
	return TEST_TOPIC
}
