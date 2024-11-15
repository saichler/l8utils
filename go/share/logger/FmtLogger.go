package logger

import (
	"errors"
	"fmt"
	"testing"
)

type FmtLogger struct{}

func (fmtLog *FmtLogger) Trace(args ...interface{}) {
	fmt.Println(FormatLog(Trace, args...))
}
func (fmtLog *FmtLogger) Debug(args ...interface{}) {
	fmt.Println(FormatLog(Debug, args...))
}
func (fmtLog *FmtLogger) Info(args ...interface{}) {
	fmt.Println(FormatLog(Info, args...))
}
func (fmtLog *FmtLogger) Warning(args ...interface{}) {
	fmt.Println(FormatLog(Warning, args...))
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
	fmtLog.Error(args...)
	ts, ok := t.(*testing.T)
	if ok {
		ts.Fail()
	}
}
