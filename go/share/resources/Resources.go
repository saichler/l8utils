package resources

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/service_points"
	"github.com/saichler/shared/go/share/shallow_security"
	"github.com/saichler/shared/go/types"
)

const (
	DEFAULT_MAX_DATA_SIZE     = 1024 * 1024
	DEFAULT_EDGE_QUEUE_SIZE   = 10000
	DEFAULT_SWITCH_QUEUE_SIZE = 50000
	DEFAULT_SWITCH_PORT       = 50000
)

type Resources struct {
	registry      interfaces.IRegistry
	servicePoints interfaces.IServicePoints
	security      interfaces.ISecurityProvider
	dataListener  interfaces.IDatatListener
	logger        interfaces.ILogger
	serializers   map[interfaces.SerializerMode]interfaces.ISerializer
	configs       map[interfaces.ConfigType]*types.MessagingConfig
}

func NewDefaultResources(logger interfaces.ILogger) interfaces.IResources {
	return NewResources(registry.NewRegistry(), createSecurityProvider(), logger)
}

func NewResources(registry interfaces.IRegistry,
	security interfaces.ISecurityProvider,
	logger interfaces.ILogger) interfaces.IResources {
	r := &Resources{}
	r.registry = registry
	r.servicePoints = service_points.NewServicePoints(r)
	r.security = security
	r.logger = logger
	r.serializers = make(map[interfaces.SerializerMode]interfaces.ISerializer)
	r.configs = make(map[interfaces.ConfigType]*types.MessagingConfig)
	r.createConfigs()
	return r
}

func (this *Resources) Registry() interfaces.IRegistry {
	return this.registry
}
func (this *Resources) ServicePoints() interfaces.IServicePoints {
	return this.servicePoints
}
func (this *Resources) Security() interfaces.ISecurityProvider {
	return this.security
}
func (this *Resources) DataListener() interfaces.IDatatListener {
	return this.dataListener
}
func (this *Resources) Serializer(mode interfaces.SerializerMode) interfaces.ISerializer {
	return this.serializers[mode]
}
func (this *Resources) Logger() interfaces.ILogger {
	return this.logger
}
func (this *Resources) Config(configType interfaces.ConfigType) types.MessagingConfig {
	return *this.configs[configType]
}

func (this *Resources) createConfigs() {
	edgeConfig := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE,
		DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30)
	edgeSwitchConfig := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE,
		DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, false, 0)
	switchConfig := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_SWITCH_QUEUE_SIZE,
		DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30)
	this.configs[interfaces.EdgeConfig] = edgeConfig
	this.configs[interfaces.SwitchConfig] = switchConfig
	this.configs[interfaces.EdgeSwitchConfig] = edgeSwitchConfig
}

func createSecurityProvider() interfaces.ISecurityProvider {
	hash := md5.New()
	secret := "Default Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	return shallow_security.NewShallowSecurityProvider(k, secret)
}
