package tests

import (
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8sysconfig"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

var globals ifs.IResources
var Log = logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

func init() {
	_log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	_log.SetLogLevel(ifs.Trace_Level)
	_resources := resources.NewResources(_log)
	_resources.Set(registry.NewRegistry())
	_security, err := ifs.LoadSecurityProvider(nil)
	if err != nil {
		panic("Failed to load security provider " + err.Error())
	}
	_resources.Set(_security)
	_config := &l8sysconfig.L8SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize: resources.DEFAULT_QUEUE_SIZE,
		VnetPort:    50000}
	_resources.Set(_config)
	globals = _resources

}
