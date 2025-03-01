package interfaces

import (
	"github.com/saichler/shared/go/types"
	"time"
)
import "github.com/google/uuid"

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

func AddTopic(config *types.VNicConfig, vlan int32, topic string) {
	if config == nil {
		return
	}
	if config.Vlans == nil {
		config.Vlans = &types.Vlans{}
		config.Vlans.Vlans = make(map[int32]*types.Vlan)
	}
	_, ok := config.Vlans.Vlans[vlan]
	if !ok {
		config.Vlans.Vlans[vlan] = &types.Vlan{}
		config.Vlans.Vlans[vlan].Vlan = vlan
		config.Vlans.Vlans[vlan].Members = make(map[string]*types.VlanMembers)
	}
	_, ok = config.Vlans.Vlans[vlan].Members[topic]
	if !ok {
		config.Vlans.Vlans[vlan].Members[topic] = &types.VlanMembers{}
		config.Vlans.Vlans[vlan].Members[topic].MemberToJoinTime = make(map[string]int64)
	}
	if config.LocalUuid == "" {
		config.LocalUuid = NewUuid()
	}
	config.Vlans.Vlans[vlan].Members[topic].MemberToJoinTime[config.LocalUuid] = time.Now().UnixMilli()
}

func NewUuid() string {
	return uuid.New().String()
}
