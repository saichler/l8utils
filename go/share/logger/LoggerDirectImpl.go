package logger

import (
	"errors"
	"github.com/saichler/shared/go/share/interfaces"
	"testing"
	"time"
)

type LoggerDirectImpl struct {
	logMethods []ILogMethod
	logLevel   interfaces.LogLevel
}

func NewLoggerDirectImpl(logMethods ...ILogMethod) *LoggerDirectImpl {
	logImpl := &LoggerDirectImpl{}
	logImpl.logMethods = logMethods
	return logImpl
}

func (loggerImpl *LoggerDirectImpl) processEntry(entry *LoggerEntry) {
	str := FormatLog(entry.l, entry.t, entry.anys...)
	if len(loggerImpl.logMethods) == 1 {
		loggerImpl.logMethods[0].Log(entry.l, str)
	} else if len(loggerImpl.logMethods) == 2 {
		loggerImpl.logMethods[0].Log(entry.l, str)
		loggerImpl.logMethods[1].Log(entry.l, str)
	} else if len(loggerImpl.logMethods) == 3 {
		loggerImpl.logMethods[0].Log(entry.l, str)
		loggerImpl.logMethods[1].Log(entry.l, str)
		loggerImpl.logMethods[2].Log(entry.l, str)
	}
}

func (loggerImpl *LoggerDirectImpl) Empty() bool {
	return true
}

func (loggerImpl *LoggerDirectImpl) Trace(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Trace_Level {
		return
	}
	loggerImpl.processEntry(newEntry(interfaces.Trace_Level, anys...))
}

func (loggerImpl *LoggerDirectImpl) Debug(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Debug_Level {
		return
	}
	loggerImpl.processEntry(newEntry(interfaces.Debug_Level, anys...))
}

func (loggerImpl *LoggerDirectImpl) Info(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Info_Level {
		return
	}
	loggerImpl.processEntry(newEntry(interfaces.Info_Level, anys...))
}

func (loggerImpl *LoggerDirectImpl) Warning(anys ...interface{}) {
	if loggerImpl.logLevel > interfaces.Warning_Level {
		return
	}
	loggerImpl.processEntry(newEntry(interfaces.Warning_Level, anys...))
}

func (loggerImpl *LoggerDirectImpl) Error(anys ...interface{}) error {
	anys = append(anys, FileAndLine(".go", false))
	loggerImpl.processEntry(newEntry(interfaces.Error_Level, anys...))
	err := FormatLog(interfaces.Error_Level, time.Now().Unix(), anys...)
	return errors.New(err)
}

func (loggerImpl *LoggerDirectImpl) Fail(t interface{}, args ...interface{}) {
	args = append(args, FileAndLine("tests", true))
	loggerImpl.Error(args...)
	ts, ok := t.(*testing.T)
	if ok {
		ts.Fail()
	}
}

func (loggerImpl *LoggerDirectImpl) SetLogLevel(level interfaces.LogLevel) {
	loggerImpl.logLevel = level
}
