package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/shared/go/share/shallow_security"
	"github.com/saichler/shared/go/tests/infra"
)

var globals interfaces.IResources
var log interfaces.ILogger

func init() {
	log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	infra.Log = log
	globals = resources.NewResources(registry.NewRegistry(),
		shallow_security.CreateShallowSecurityProvider(), nil, log, nil, nil, "tests")
}
