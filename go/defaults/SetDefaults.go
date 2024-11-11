package defaults

import (
	"github.com/saichler/shared/go/interfaces"
	"github.com/saichler/shared/go/service_points"
	"github.com/saichler/shared/go/struct_registry"
)

func LoadDefaults() {}

func init() {
	initEdge()
	initRegistry()
	initServicePoints()
}

func initEdge() {
	interfaces.EdgeConfig = interfaces.NewMessageConfig(1024*1024, 1000, 1000, 50000, true, 30)
	interfaces.EdgeSwitchConfig = interfaces.NewMessageConfig(1024*1024, 1000, 1000, 50000, false, 0)
	interfaces.SwitchConfig = interfaces.NewMessageConfig(1024*1024, 5000, 5000, 50000, true, 30)
}

func initRegistry() {
	interfaces.SetStructRegistry(struct_registry.NewStructRegistry())
}

func initServicePoints() {
	interfaces.SetServicePoints(service_points.NewServicePoints())
}
