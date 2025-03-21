package resources

import (
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

const (
	DEFAULT_MAX_DATA_SIZE = 1024 * 1024
	DEFAULT_QUEUE_SIZE    = 10000
)

type Resources struct {
	registry      common.IRegistry
	servicePoints common.IServicePoints
	security      common.ISecurityProvider
	logger        common.ILogger
	dataListener  common.IDatatListener
	serializers   map[common.SerializerMode]common.ISerializer
	config        *types.VNicConfig
	introspector  common.IIntrospector
}

func NewResources(registry common.IRegistry,
	security common.ISecurityProvider,
	servicePoints common.IServicePoints,
	logger common.ILogger,
	dataListener common.IDatatListener,
	serializer common.ISerializer,
	config *types.VNicConfig,
	introspector common.IIntrospector) common.IResources {
	r := &Resources{}
	r.registry = registry
	r.servicePoints = servicePoints
	r.security = security
	r.logger = logger
	r.dataListener = dataListener
	r.serializers = make(map[common.SerializerMode]common.ISerializer)
	if serializer != nil {
		r.serializers[serializer.Mode()] = serializer
	}
	r.config = config
	r.introspector = introspector
	return r
}

func (this *Resources) AddService(serviceName string, serviceArea int32) {
	common.AddService(this.config, serviceName, serviceArea)
}

func (this *Resources) Registry() common.IRegistry {
	return this.registry
}
func (this *Resources) ServicePoints() common.IServicePoints {
	return this.servicePoints
}
func (this *Resources) Security() common.ISecurityProvider {
	return this.security
}
func (this *Resources) DataListener() common.IDatatListener {
	return this.dataListener
}
func (this *Resources) Serializer(mode common.SerializerMode) common.ISerializer {
	return this.serializers[mode]
}
func (this *Resources) Logger() common.ILogger {
	return this.logger
}
func (this *Resources) Config() *types.VNicConfig {
	return this.config
}
func (this *Resources) Introspector() common.IIntrospector {
	return this.introspector
}
