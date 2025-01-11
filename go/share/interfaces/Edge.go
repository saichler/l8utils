package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type IEdge interface {
	Start()
	Config() types.MessagingConfig
	Send([]byte) error
	Name() string
	Do(types.Action, string, proto.Message) error
	Shutdown()
	CreatedAt() int64
	PublishState()
	RegisterTopic(string)
}

type IDatatListener interface {
	PortShutdown(IEdge)
	HandleData([]byte, IEdge)
}

func NewMessageConfig(maxDataSize uint64,
	txQueueSize,
	rxQueueSize uint64,
	switchPort uint32,
	sendStateInfo bool,
	sendStateIntervalSeconds int64) *types.MessagingConfig {
	mc := &types.MessagingConfig{
		MaxDataSize:              maxDataSize,
		TxQueueSize:              txQueueSize,
		RxQueueSize:              rxQueueSize,
		SwitchPort:               switchPort,
		SendStateInfo:            sendStateInfo,
		SendStateIntervalSeconds: sendStateIntervalSeconds,
	}
	return mc
}
