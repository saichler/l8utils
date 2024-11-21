package testsp

import (
	"github.com/saichler/shared/go/share/interfaces"
	"google.golang.org/protobuf/proto"
)

type TestServicePointHandler struct {
	Name         string
	PostNumber   int
	PutNumber    int
	PatchNumber  int
	DeleteNumber int
	GetNumber    int
}

const (
	TEST_TOPIC = "Tests"
)

func NewTestServicePointHandler(name string) *TestServicePointHandler {
	tsp := &TestServicePointHandler{}
	tsp.Name = name
	return tsp
}

func (tsp *TestServicePointHandler) Post(pb proto.Message, edge interfaces.IEdge) (proto.Message, error) {
	interfaces.Logger().Debug("Post Test callback")
	tsp.PostNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Put(pb proto.Message, edge interfaces.IEdge) (proto.Message, error) {
	interfaces.Logger().Debug("Put Test callback")
	tsp.PutNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Patch(pb proto.Message, edge interfaces.IEdge) (proto.Message, error) {
	interfaces.Logger().Debug("Patch Test callback")
	tsp.PatchNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Delete(pb proto.Message, edge interfaces.IEdge) (proto.Message, error) {
	interfaces.Logger().Debug("Delete Test callback")
	tsp.DeleteNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) Get(pb proto.Message, edge interfaces.IEdge) (proto.Message, error) {
	interfaces.Logger().Debug("Get Test callback")
	tsp.GetNumber++
	return nil, nil
}
func (tsp *TestServicePointHandler) EndPoint() string {
	return "/Tests"
}
func (tsp *TestServicePointHandler) Topic() string {
	return TEST_TOPIC
}
