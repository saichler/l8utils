package logger

import (
	"errors"
	"fmt"
	"github.com/saichler/shared/go/string_utils"
	"testing"
	"time"
)

type FmtLogger struct{}

func format(t string, args ...interface{}) string {
	str := string_utils.New()
	str.Add(time.Now().String())
	str.Add(" ", t, ": ")
	if args != nil {
		for _, arg := range args {
			str.Add(str.StringOf(arg))
		}
	}
	return str.String()
}

func (fmtLog *FmtLogger) Trace(args ...interface{}) {
	fmt.Println(format("Trace", args...))
}
func (fmtLog *FmtLogger) Debug(args ...interface{}) {
	fmt.Println(format("Debug", args...))
}
func (fmtLog *FmtLogger) Info(args ...interface{}) {
	fmt.Println(format("Info", args...))
}
func (fmtLog *FmtLogger) Warning(args ...interface{}) {
	fmt.Println(format("Warning", args...))
}
func (fmtLog *FmtLogger) Error(args ...interface{}) error {
	msg := format("Error", args...)
	fmt.Println(msg)
	return errors.New(msg)
}
func (fmtLog *FmtLogger) Empty() bool {
	return false
}
func (fmtLog *FmtLogger) Fail(t interface{}, args ...interface{}) {
	fmtLog.Error(args...)
	ts, ok := t.(testing.T)
	if ok {
		ts.Fail()
	}
}
