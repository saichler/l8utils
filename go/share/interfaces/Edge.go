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

var edgeConfig *types.MessagingConfig
var edgeSwitchConfig *types.MessagingConfig
var switchConfig *types.MessagingConfig

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

func SetEdgeConfig(config *types.MessagingConfig) {
	edgeConfig = config
}

func SetEdgeSwitchConfig(config *types.MessagingConfig) {
	edgeSwitchConfig = config
}

func SetSwitchConfig(config *types.MessagingConfig) {
	switchConfig = config
}

func EdgeConfig() *types.MessagingConfig {
	return cloneConfig(*edgeConfig)
}

func EdgeSwitchConfig() *types.MessagingConfig {
	return cloneConfig(*edgeSwitchConfig)
}

func SwitchConfig() *types.MessagingConfig {
	return cloneConfig(*switchConfig)
}

func cloneConfig(config types.MessagingConfig) *types.MessagingConfig {
	return &config
}
