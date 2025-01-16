package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type ConfigType int

const (
	EdgeConfig       ConfigType = 1
	EdgeSwitchConfig ConfigType = 2
	SwitchConfig     ConfigType = 3
)

type IEdge interface {
	Start()
	Shutdown()
	Name() string
	Config() types.MessagingConfig
	Send([]byte) error
	Do(types.Action, string, proto.Message) error
	Resources() IResources
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
