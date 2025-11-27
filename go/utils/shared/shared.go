package shared

import (
	"time"

	"github.com/saichler/l8reflect/go/reflect/introspecting"
	"github.com/saichler/l8services/go/services/manager"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8sysconfig"
	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/registry"
	"github.com/saichler/l8utils/go/utils/resources"
)

func ResourcesOf(alias string, vnetPort, keepAlive uint32, logToFile bool, others ...ifs.IResources) ifs.IResources {
	if logToFile {
		logger.SetLogToFile(alias)
	}
	log := logger.NewLoggerImpl(&logger.FmtLogMethod{})
	log.SetLogLevel(ifs.Error_Level)
	res := resources.NewResources(log)

	res.Set(registry.NewRegistry())

	if others != nil {
		for _, other := range others {
			res.Set(other.Security())
		}
	} else {
		sec, err := ifs.LoadSecurityProvider(res)
		if err != nil {
			time.Sleep(time.Second * 10)
			panic(err.Error())
		}
		res.Set(sec)
	}

	conf := &l8sysconfig.L8SysConfig{MaxDataSize: resources.DEFAULT_MAX_DATA_SIZE,
		RxQueueSize:              resources.DEFAULT_QUEUE_SIZE,
		TxQueueSize:              resources.DEFAULT_QUEUE_SIZE,
		LocalAlias:               alias,
		VnetPort:                 vnetPort,
		KeepAliveIntervalSeconds: int64(keepAlive)}
	res.Set(conf)

	res.Set(introspecting.NewIntrospect(res.Registry()))
	res.Set(manager.NewServices(res))

	return res
}
