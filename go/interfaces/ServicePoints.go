package interfaces

import (
	"github.com/saichler/my.simple/go/net/model"
	"google.golang.org/protobuf/proto"
)

type IServicePoints interface {
	RegisterServicePoint(proto.Message, IServicePointHandler, IRegistry) error
	Handle(proto.Message, model.Action, IEdge) (proto.Message, error)
}

type IServicePointHandler interface {
	Post(proto.Message, IEdge) (proto.Message, error)
	Put(proto.Message, IEdge) (proto.Message, error)
	Patch(proto.Message, IEdge) (proto.Message, error)
	Delete(proto.Message, IEdge) (proto.Message, error)
	Get(proto.Message, IEdge) (proto.Message, error)
	EndPoint() string
}

var servicePoints IServicePoints

func ServicePoints() IServicePoints {
	return servicePoints
}

func SetServicePoints(sp IServicePoints) {
	servicePoints = sp
}
