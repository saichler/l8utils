package tests

import (
	"github.com/saichler/shared/go/share/interfaces"
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/shared/go/tests/infra"
	"github.com/saichler/shared/go/types"
)

var globals interfaces.IResources
var log interfaces.ILogger

func init() {
	log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	infra.Log = log
	config := &types.VNicConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  "tests"}
	secure, err := interfaces.LoadSecurityProvider("../share/shallow_security/security.so")
	if err != nil {
		panic(err)
	}
	globals = resources.NewResources(registry.NewRegistry(),
		secure, nil, log, nil, nil, config, nil)
	globals.Security().Init(globals)
}
