package resources

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/types"
)

const (
	DEFAULT_MAX_DATA_SIZE = 1024 * 1024
	DEFAULT_QUEUE_SIZE    = 10000
)

type Resources struct {
	registry      interfaces.IRegistry
	servicePoints interfaces.IServicePoints
	security      interfaces.ISecurityProvider
	logger        interfaces.ILogger
	dataListener  interfaces.IDatatListener
	serializers   map[interfaces.SerializerMode]interfaces.ISerializer
	config        *types.VNicConfig
	introspector  interfaces.IIntrospector
}

func NewResources(registry interfaces.IRegistry,
	security interfaces.ISecurityProvider,
	servicePoints interfaces.IServicePoints,
	logger interfaces.ILogger,
	dataListener interfaces.IDatatListener,
	serializer interfaces.ISerializer,
	config *types.VNicConfig,
	introspector interfaces.IIntrospector) interfaces.IResources {
	r := &Resources{}
	r.registry = registry
	r.servicePoints = servicePoints
	r.security = security
	r.logger = logger
	r.dataListener = dataListener
	r.serializers = make(map[interfaces.SerializerMode]interfaces.ISerializer)
	if serializer != nil {
		r.serializers[serializer.Mode()] = serializer
	}
	r.config = config
	r.introspector = introspector
	return r
}

func (this *Resources) AddTopic(area int32, topic string) {
	interfaces.AddTopic(this.config, area, topic)
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
func (this *Resources) Config() *types.VNicConfig {
	return this.config
}
func (this *Resources) Introspector() interfaces.IIntrospector {
	return this.introspector
}
