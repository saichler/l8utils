package tests

import (
	"crypto/md5"
	"encoding/base64"
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/service_points"
	"github.com/saichler/shared/go/share/shallow_security"
	"github.com/saichler/shared/go/share/type_registry"
)

var providers *interfaces.Providers

const (
	DEFAULT_MAX_DATA_SIZE     = 1024 * 1024
	DEFAULT_EDGE_QUEUE_SIZE   = 10000
	DEFAULT_SWITCH_QUEUE_SIZE = 50000
	DEFAULT_SWITCH_PORT       = 50000
)

func init() {
	providers = interfaces.NewProviders(
		type_registry.NewTypeRegistry(),
		createSecurityProvider(),
		service_points.NewServicePoints(),
		logger.NewLoggerImpl(&logger.FmtLogMethod{}))
	a := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE,
		DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30)
	b := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_EDGE_QUEUE_SIZE,
		DEFAULT_EDGE_QUEUE_SIZE, DEFAULT_SWITCH_PORT, false, 0)
	c := interfaces.NewMessageConfig(DEFAULT_MAX_DATA_SIZE, DEFAULT_SWITCH_QUEUE_SIZE,
		DEFAULT_SWITCH_QUEUE_SIZE, DEFAULT_SWITCH_PORT, true, 30)
	providers.SetDefaultMessageConfig(a, c, b)
}

func createSecurityProvider() interfaces.ISecurityProvider {
	hash := md5.New()
	secret := "Default Security Provider"
	hash.Write([]byte(secret))
	kHash := hash.Sum(nil)
	k := base64.StdEncoding.EncodeToString(kHash)
	return shallow_security.NewShallowSecurityProvider(k, secret)
}
