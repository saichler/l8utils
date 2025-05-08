package tests

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

var globals ifs.IResources
var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func init() {
	config := &types.SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:  "tests"}
	secure, err := ifs.LoadSecurityProvider()
	if err != nil {
		panic(err)
	}
	globals = resources.NewResources(registry.NewRegistry(),
		secure, nil, Log, nil, nil, config, nil)
}
