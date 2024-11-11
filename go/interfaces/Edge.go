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

var EdgeConfig *types.MessagingConfig
var EdgeSwitchConfig *types.MessagingConfig
var SwitchConfig *types.MessagingConfig

func NewMessageConfig(maxDataSize uint64,
	txQueueSize,
	rxQueueSize uint64,
	switchPort uint32,
	sendStatusInfo bool,
	sendStatusIntervalSeconds uint64) *types.MessagingConfig {
	mc := &types.MessagingConfig{
		MaxDataSize:               maxDataSize,
		TxQueueSize:               txQueueSize,
		RxQueueSize:               rxQueueSize,
		SwitchPort:                switchPort,
		SendStatusInfo:            sendStatusInfo,
		SendStatusIntervalSeconds: sendStatusIntervalSeconds,
	}
	return mc
}
