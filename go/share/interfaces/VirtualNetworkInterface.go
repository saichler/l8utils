package interfaces

import (
	"github.com/saichler/shared/go/types"
)

type IVirtualNetworkInterface interface {
	Start()
	Shutdown()
	Name() string
	SendMessage([]byte) error
	Unicast(types.Action, string, interface{}) error
	Multicast(types.Action, uint32, string, interface{}) error
	Resources() IResources
}

type IDatatListener interface {
	ShutdownVNic(IVirtualNetworkInterface)
	HandleData([]byte, IVirtualNetworkInterface)
	Failed([]byte, IVirtualNetworkInterface, string)
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
