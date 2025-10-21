package main

import (
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/logger"
)

func main() {
	_log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})
	_log.SetLogLevel(ifs.Trace_Level)
	_log.Error("Just error")
	time.Sleep(time.Second)
	panic("panic")
}
