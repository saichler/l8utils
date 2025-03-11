package tests

import (
	"errors"
	. "github.com/saichler/shared/go/tests/infra"
	"strings"
	"testing"
)

func TestFmtLogger(t *testing.T) {
	err := errors.New("Sample Error")
	Log.Trace("my trace message: ", err)
	Log.Debug("my debug message: ", err)
	Log.Info("my info message: ", err)
	Log.Warning("my warning message: ", err)
	err = Log.Error("my error message: ", err)
	if !strings.Contains(err.Error(), "(Error) - my error message: Sample Error") {
		t.Fail()
		Log.Error("Expected a formatted error message:", err.Error())
		return
	}
	Log.Empty()

	tt := &testing.T{}
	Log.Fail(tt, "my fail message", err)
}
