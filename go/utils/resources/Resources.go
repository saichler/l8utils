package resources

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
)

const (
	DEFAULT_MAX_DATA_SIZE = 1024 * 1024 * 50
	DEFAULT_QUEUE_SIZE    = 50000
)

type Resources struct {
	registry     ifs.IRegistry
	services     ifs.IServices
	security     ifs.ISecurityProvider
	logger       ifs.ILogger
	dataListener ifs.IDatatListener
	serializers  map[ifs.SerializerMode]ifs.ISerializer
	config       *types.SysConfig
	introspector ifs.IIntrospector
}

func NewResources(registry ifs.IRegistry,
	security ifs.ISecurityProvider,
	servicePoints ifs.IServices,
	logger ifs.ILogger,
	dataListener ifs.IDatatListener,
	serializer ifs.ISerializer,
	config *types.SysConfig,
	introspector ifs.IIntrospector) ifs.IResources {
	r := &Resources{}
	r.registry = registry
	r.services = servicePoints
	r.security = security
	r.logger = logger
	r.dataListener = dataListener
	r.serializers = make(map[ifs.SerializerMode]ifs.ISerializer)
	if serializer != nil {
		r.serializers[serializer.Mode()] = serializer
	}
	r.config = config
	r.introspector = introspector
	return r
}

func (this *Resources) AddService(serviceName string, serviceArea int32) {
	ifs.AddService(this.config, serviceName, serviceArea)
}

func (this *Resources) Registry() ifs.IRegistry {
	return this.registry
}
func (this *Resources) Services() ifs.IServices {
	return this.services
}
func (this *Resources) Security() ifs.ISecurityProvider {
	return this.security
}
func (this *Resources) DataListener() ifs.IDatatListener {
	return this.dataListener
}
func (this *Resources) Serializer(mode ifs.SerializerMode) ifs.ISerializer {
	return this.serializers[mode]
}
func (this *Resources) Logger() ifs.ILogger {
	return this.logger
}
func (this *Resources) SysConfig() *types.SysConfig {
	return this.config
}
func (this *Resources) Introspector() ifs.IIntrospector {
	return this.introspector
}
