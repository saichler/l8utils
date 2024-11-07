package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type OverlayEdge interface {
	Start()
	Addr() string
	Uuid() string
	Send([]byte) error
	Name() string
	Do(types.Action, string, proto.Message) error
	Shutdown()
	CreatedAt() int64
}

type DatatListener interface {
	PortShutdown(OverlayEdge)
	HandleData([]byte, OverlayEdge)
}
