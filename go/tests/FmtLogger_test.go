package tests

import (
	"errors"
	"strings"
	"testing"
)

func TestFmtLogger(t *testing.T) {
	err := errors.New("Sample Error")
	log.Trace("my trace message: ", err)
	log.Debug("my debug message: ", err)
	log.Info("my info message: ", err)
	log.Warning("my warning message: ", err)
	err = log.Error("my error message: ", err)
	if !strings.Contains(err.Error(), "(Error) - my error message: Sample Error") {
		t.Fail()
		log.Error("Expected a formatted error message:", err.Error())
		return
	}
	log.Empty()

	tt := &testing.T{}
	log.Fail(tt, "my fail message", err)
}
