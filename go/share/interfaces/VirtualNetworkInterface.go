package interfaces

import (
	"github.com/saichler/shared/go/types"
	"google.golang.org/protobuf/proto"
)

type IVirtualNetworkInterface interface {
	Start()
	Shutdown()
	Name() string
	Send([]byte) error
	Do(types.Action, string, proto.Message) error
	Resources() IResources
}

type IDatatListener interface {
	ShutdownVNic(IVirtualNetworkInterface)
	HandleData([]byte, IVirtualNetworkInterface)
}

func NewVNicConfig(maxDataSize uint64, txQueueSize, rxQueueSize uint64, switchPort uint32) *types.VNicConfig {
	mc := &types.VNicConfig{
		MaxDataSize: maxDataSize,
		TxQueueSize: txQueueSize,
		RxQueueSize: rxQueueSize,
		SwitchPort:  switchPort,
	}
	return mc
}
