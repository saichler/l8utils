package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type IServicePoints interface {
	RegisterServicePoint(int32, proto.Message, IServicePointHandler) error
	Handle(proto.Message, types.Action, IVirtualNetworkInterface, *types.Message) (proto.Message, error)
	Notify(proto.Message, types.Action, IVirtualNetworkInterface, *types.Message) (proto.Message, error)
	ServicePointHandler(string) (IServicePointHandler, bool)
	Areas() *types.Areas
}

type IServicePointHandler interface {
	Post(proto.Message, IVirtualNetworkInterface) (proto.Message, error)
	Put(proto.Message, IVirtualNetworkInterface) (proto.Message, error)
	Patch(proto.Message, IVirtualNetworkInterface) (proto.Message, error)
	Delete(proto.Message, IVirtualNetworkInterface) (proto.Message, error)
	Get(proto.Message, IVirtualNetworkInterface) (proto.Message, error)
	Failed(proto.Message, IVirtualNetworkInterface, *types.Message) (proto.Message, error)
	EndPoint() string
	Topic() string
}
