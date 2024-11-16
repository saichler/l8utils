package defaults

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/service_points"
	"github.com/saichler/shared/go/share/shallow_security"
	"github.com/saichler/shared/go/share/struct_registry"
)

func LoadDefaultImplementations() {}

func init() {
	initLogger()
	initEdgeConfig()
	initRegistry()
	initServicePoints()
	initSecurityProvider()
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

func initSecurityProvider() {
	hash := md5.New()
	secret := "Default Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	interfaces.SetSecurityProvider(shallow_security.NewShallowSecurityProvider(k, secret))
}
