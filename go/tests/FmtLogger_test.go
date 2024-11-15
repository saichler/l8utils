package tests

import (
	"errors"
	"github.com/saichler/shared/go/src/interfaces"
	"strings"
	"testing"
)

func TestFmtLogger(t *testing.T) {
	err := errors.New("Sample Error")
	interfaces.Trace("my trace message: ", err)
	interfaces.Debug("my debug message: ", err)
	interfaces.Info("my info message: ", err)
	interfaces.Warning("my warning message: ", err)
	err = interfaces.Error("my error message: ", err)
	if !strings.Contains(err.Error(), "(  Error) - my error message: Sample Error") {
		t.Fail()
		interfaces.Error("Expected a formatted error message")
		return
	}
	interfaces.Empty()

	tt := &testing.T{}
	interfaces.Fail(tt, "my fail message", err)
}
