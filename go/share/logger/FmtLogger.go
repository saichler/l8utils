package logger

import (
	"errors"
	"fmt"
	"testing"
)

type FmtLogger struct {
	isTraceEnabled   bool
	isDebugEnabled   bool
	isInfoEnabled    bool
	isWarningEnabled bool
}

func NewFmtLogger(isTraceEnabled, isDebugEnabled, isInfoEnabled, isWarningEnabled bool) *FmtLogger {
	fmtLogger := &FmtLogger{}
	fmtLogger.isTraceEnabled = isTraceEnabled
	fmtLogger.isDebugEnabled = isDebugEnabled
	fmtLogger.isInfoEnabled = isInfoEnabled
	fmtLogger.isWarningEnabled = isWarningEnabled
	return fmtLogger
}

func (fmtLog *FmtLogger) Trace(args ...interface{}) {
	if fmtLog.isTraceEnabled {
		fmt.Println(FormatLog(Trace, args...))
	}
}
func (fmtLog *FmtLogger) Debug(args ...interface{}) {
	if fmtLog.isDebugEnabled {
		fmt.Println(FormatLog(Debug, args...))
	}
}
func (fmtLog *FmtLogger) Info(args ...interface{}) {
	if fmtLog.isInfoEnabled {
		fmt.Println(FormatLog(Info, args...))
	}
}
func (fmtLog *FmtLogger) Warning(args ...interface{}) {
	if fmtLog.isWarningEnabled {
		fmt.Println(FormatLog(Warning, args...))
	}
}
func (fmtLog *FmtLogger) Error(args ...interface{}) error {
	msg := FormatLog(Error, args...)
	fmt.Println(msg)
	return errors.New(msg)
}
func (fmtLog *FmtLogger) Empty() bool {
	return false
}
func (fmtLog *FmtLogger) Fail(t interface{}, args ...interface{}) {
	args = append(args, FileAndLine("tests"))
	fmtLog.Error(args...)
	ts, ok := t.(*testing.T)
	if ok {
		ts.Fail()
	}
}

func (fmtLog *FmtLogger) IsTraceEnabled() bool {
	return fmtLog.isTraceEnabled
}
func (fmtLog *FmtLogger) IsDebugEnabled() bool {
	return fmtLog.isDebugEnabled
}
func (fmtLog *FmtLogger) IsInfoEnabled() bool {
	return fmtLog.isInfoEnabled
}
func (fmtLog *FmtLogger) IsWarningEnabled() bool {
	return fmtLog.isWarningEnabled
}
