package defaults

import (
	"github.com/saichler/shared/go/src/interfaces"
	"github.com/saichler/shared/go/src/logger"
	"github.com/saichler/shared/go/src/service_points"
	"github.com/saichler/shared/go/src/struct_registry"
)

func LoadDefaultImplementations() {}

func init() {
	initLogger()
	initEdgeConfig()
	initRegistry()
	initServicePoints()
}

func initLogger() {
	interfaces.SetLogger(&logger.FmtLogger{})
}

func initEdgeConfig() {
	interfaces.SetEdgeConfig(interfaces.NewMessageConfig(1024*1024, 1000, 1000, 50000, true, 30))
	interfaces.SetEdgeSwitchConfig(interfaces.NewMessageConfig(1024*1024, 1000, 1000, 50000, false, 0))
	interfaces.SetSwitchConfig(interfaces.NewMessageConfig(1024*1024, 5000, 5000, 50000, true, 30))
}

func initRegistry() {
	interfaces.SetStructRegistry(struct_registry.NewStructRegistry())
}

func initServicePoints() {
	interfaces.SetServicePoints(service_points.NewServicePoints())
}
