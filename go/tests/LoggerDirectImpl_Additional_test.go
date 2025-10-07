package tests

import (
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/logger"
)

func TestLoggerDirectImplThreeMethods(t *testing.T) {
	logMethod1 := &TestLogMethod{logs: make([]string, 0)}
	logMethod2 := &TestLogMethod{logs: make([]string, 0)}
	logMethod3 := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerDirectImpl(logMethod1, logMethod2, logMethod3)

	log.Info("test message")

	// All three log methods should have the message
	if len(logMethod1.logs) != 1 {
		t.Errorf("Expected 1 log in method1, got %d", len(logMethod1.logs))
	}
	if len(logMethod2.logs) != 1 {
		t.Errorf("Expected 1 log in method2, got %d", len(logMethod2.logs))
	}
	if len(logMethod3.logs) != 1 {
		t.Errorf("Expected 1 log in method3, got %d", len(logMethod3.logs))
	}
}

func TestLoggerDirectImplSetLogLevel(t *testing.T) {
	logMethod := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerDirectImpl(logMethod)

	// Set log level to ERROR
	log.SetLogLevel(ifs.Error_Level)

	// Test that trace is not logged
	log.Trace("should not be logged")
	if len(logMethod.logs) != 0 {
		t.Errorf("Expected no logs, got %d", len(logMethod.logs))
	}

	// Test that debug is not logged
	log.Debug("should not be logged")
	if len(logMethod.logs) != 0 {
		t.Errorf("Expected no logs, got %d", len(logMethod.logs))
	}

	// Test that info is not logged
	log.Info("should not be logged")
	if len(logMethod.logs) != 0 {
		t.Errorf("Expected no logs, got %d", len(logMethod.logs))
	}

	// Test that warning is not logged
	log.Warning("should not be logged")
	if len(logMethod.logs) != 0 {
		t.Errorf("Expected no logs, got %d", len(logMethod.logs))
	}

	// Error should be logged
	log.Error("this should be logged")
	if len(logMethod.logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logMethod.logs))
	}
}
