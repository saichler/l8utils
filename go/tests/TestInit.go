package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/shared/go/tests/infra"
)

var globals interfaces.IResources
var log interfaces.ILogger

func init() {
	log = logger.NewLoggerImpl(&logger.FmtLogMethod{})
	infra.Log = log
	globals = resources.NewDefaultResources(log)
}
