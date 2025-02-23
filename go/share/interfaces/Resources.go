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
	Introspector() IIntrospector
	AddTopic(int32, string)
}

func AddTopic(config *types.VNicConfig, area int32, topic string) {
	if config == nil {
		return
	}
	if config.ServiceAreas == nil {
		config.ServiceAreas = &types.Areas{}
		config.ServiceAreas.AreasMap = make(map[int32]*types.Area)
	}
	_, ok := config.ServiceAreas.AreasMap[area]
	if !ok {
		config.ServiceAreas.AreasMap[area] = &types.Area{}
		config.ServiceAreas.AreasMap[area].Number = area
		config.ServiceAreas.AreasMap[area].Topics = make(map[string]bool)
	}
	config.ServiceAreas.AreasMap[area].Topics[topic] = true
}
