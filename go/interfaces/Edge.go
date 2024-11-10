package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type IEdge interface {
	Start()
	Addr() string
	Uuid() string
	Send([]byte) error
	Name() string
	Do(*types.Request, string, proto.Message) error
	Shutdown()
	CreatedAt() int64
}

type IDatatListener interface {
	PortShutdown(IEdge)
	HandleData([]byte, IEdge)
}
