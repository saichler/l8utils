package service_points

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
	"reflect"
)

type ServicePointsImpl struct {
	structName2ServicePoint *String2ServicePointMap
}

func NewServicePoints() interfaces.IServicePoints {
	sp := &ServicePointsImpl{}
	sp.structName2ServicePoint = NewString2ServicePointMap()
	return sp
}

func (servicePoints *ServicePointsImpl) RegisterServicePoint(pb proto.Message, handler interfaces.IServicePointHandler, registry interfaces.ITypeRegistry) error {
	if pb == nil {
		return interfaces.Error("cannot register handler with nil proto")
	}
	typ := reflect.ValueOf(pb).Elem().Type()
	if handler == nil {
		return interfaces.Error("cannot register nil handler for type ", typ.Name())
	}
	registry.Register(typ)
	servicePoints.structName2ServicePoint.Put(typ.Name(), handler)
	return nil
}

func (servicePoints *ServicePointsImpl) Handle(pb proto.Message, action types.Action, edge interfaces.IEdge) (proto.Message, error) {
	tName := reflect.ValueOf(pb).Elem().Type().Name()
	h, ok := servicePoints.structName2ServicePoint.Get(tName)
	if !ok {
		return nil, interfaces.Error("Cannot find handler for type ", tName)
	}
	switch action {
	case types.Action_POST:
		return h.Post(pb, edge)
	case types.Action_PUT:
		return h.Put(pb, edge)
	case types.Action_PATCH:
		return h.Patch(pb, edge)
	case types.Action_DELETE:
		return h.Delete(pb, edge)
	case types.Action_GET:
		return h.Get(pb, edge)
	case types.Action_Invalid_Action:
		return nil, interfaces.Error("Invalid Action, ignoring")
	}
	panic("Unknown Action:" + action.String())
}

func (servicePoints *ServicePointsImpl) ServicePointHandler(topic string) (interfaces.IServicePointHandler, bool) {
	return servicePoints.structName2ServicePoint.Get(topic)
}
