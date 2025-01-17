package interfaces

import "github.com/saichler/shared/go/types"

type IResources interface {
	Registry() IRegistry
	ServicePoints() IServicePoints
	Security() ISecurityProvider
	DataListener() IDatatListener
	Serializer(SerializerMode) ISerializer
	Logger() ILogger
	Config() *types.VNicConfig
	SetDataListener(IDatatListener)
	SetSerializer(SerializerMode, ISerializer)
}
