package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type IServicePoints interface {
	RegisterServicePoint(proto.Message, IServicePointHandler) error
	Handle(proto.Message, types.Action, IEdge) (proto.Message, error)
	ServicePointHandler(string) (IServicePointHandler, bool)
}

type IServicePointHandler interface {
	Post(proto.Message, IEdge) (proto.Message, error)
	Put(proto.Message, IEdge) (proto.Message, error)
	Patch(proto.Message, IEdge) (proto.Message, error)
	Delete(proto.Message, IEdge) (proto.Message, error)
	Get(proto.Message, IEdge) (proto.Message, error)
	EndPoint() string
	Topic() string
}
