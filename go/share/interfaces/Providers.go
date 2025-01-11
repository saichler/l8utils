package interfaces

import "github.com/saichler/shared/go/types"

type Providers struct {
	registry      ITypeRegistry
	security      ISecurityProvider
	servicePoints IServicePoints
	logger        ILogger
	edgeConfig    *types.MessagingConfig
	edgeSwitch    *types.MessagingConfig
	swithConfig   *types.MessagingConfig
}

func NewProviders(registry ITypeRegistry,
	security ISecurityProvider,
	servicePoints IServicePoints,
	logger ILogger) *Providers {
	p := &Providers{
		registry:      registry,
		security:      security,
		servicePoints: servicePoints,
		logger:        logger,
	}
	return p
}

func (this *Providers) Registry() ITypeRegistry {
	return this.registry
}
func (this *Providers) Security() ISecurityProvider {
	return this.security
}
func (this *Providers) ServicePoints() IServicePoints {
	return this.servicePoints
}
func (this *Providers) Logger() ILogger {
	return this.logger
}
func (this *Providers) SetDefaultMessageConfig(edgeConfig, switchConfig, edgeSwitch *types.MessagingConfig) {
	this.edgeSwitch = edgeSwitch
	this.edgeConfig = edgeConfig
	this.swithConfig = switchConfig
}
func (this *Providers) EdgeConfig() types.MessagingConfig {
	return *this.edgeConfig
}
func (this *Providers) EdgeSwitch() types.MessagingConfig {
	return *this.edgeSwitch
}
func (this *Providers) Switch() types.MessagingConfig {
	return *this.swithConfig
}
