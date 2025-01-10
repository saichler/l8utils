package defaults

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/service_points"
	"github.com/saichler/shared/go/share/shallow_security"
	"github.com/saichler/shared/go/share/type_registry"
)

func LoadDefaultImplementations() {}

const (
	DEFAULT_MAX_DATA_SIZE     = 1024 * 1024
	DEFAULT_EDGE_QUEUE_SIZE   = 10000
	DEFAULT_SWITCH_QUEUE_SIZE = 50000
	DEFAULT_SWITCH_PORT       = 50000
)

func init() {
	initLogger()
	initEdgeConfig()
	initRegistry()
	initServicePoints()
	initSecurityProvider()
}

func initLogger() {
	interfaces.SetLogger(logger.NewLoggerImpl(&logger.FmtLogMethod{}))
}

func initEdgeConfig() {
	interfaces.SetEdgeConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30))
	interfaces.SetEdgeSwitchConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, false, 0))
	interfaces.SetSwitchConfig(interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30))
}

func initRegistry() {
	interfaces.SetTypeRegistry(type_registry.NewTypeRegistry())
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
