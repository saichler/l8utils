package tests

import (
	"github.com/saichler/shared/go/share/logger"
	"github.com/saichler/shared/go/share/registry"
	"github.com/saichler/shared/go/share/resources"
	"github.com/saichler/shared/go/tests/infra"
	"github.com/saichler/types/go/common"
	"github.com/saichler/types/go/types"
)

var globals common.IResources
var log common.ILogger

func init() {
	log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	infra.Log = log
	config := &types.VNicConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  "tests"}
	secure, err := common.LoadSecurityProvider("security.so")
	if err != nil {
		panic(err)
	}
	globals = resources.NewResources(registry.NewRegistry(),
		secure, nil, log, nil, nil, config, nil)
	globals.Security().Init(globals)
}
