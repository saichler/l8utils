package tests

import (
	"strings"
	"testing"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8utils/go/utils/logger"
)

type TestLogMethod struct {
	logs []string
}

func (tlm *TestLogMethod) Log(level ifs.LogLevel, msg string) {
	tlm.logs = append(tlm.logs, msg)
}

func TestLoggerImpl(t *testing.T) {
	logMethod := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerImpl(logMethod)

	// Test Trace
	log.Trace("trace message")
	time.Sleep(100 * time.Millisecond)

	// Test Debug
	log.Debug("debug message")
	time.Sleep(100 * time.Millisecond)

	// Test Info
	log.Info("info message")
	time.Sleep(100 * time.Millisecond)

	// Test Warning
	log.Warning("warning message")
	time.Sleep(100 * time.Millisecond)

	// Test Error
	err := log.Error("error message")
	if err == nil {
		t.Error("Error should return an error")
	}
	time.Sleep(100 * time.Millisecond)

	// Wait for queue to process
	for i := 0; i < 10; i++ {
		if log.Empty() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Check that we have logs
	if len(logMethod.logs) == 0 {
		t.Error("Expected to have logs")
	}

	// Verify log content
	found := make(map[string]bool)
	for _, logStr := range logMethod.logs {
		if strings.Contains(logStr, "trace message") {
			found["trace"] = true
		}
		if strings.Contains(logStr, "debug message") {
			found["debug"] = true
		}
		if strings.Contains(logStr, "info message") {
			found["info"] = true
		}
		if strings.Contains(logStr, "warning message") {
			found["warning"] = true
		}
		if strings.Contains(logStr, "error message") {
			found["error"] = true
		}
	}

	if !found["trace"] {
		t.Error("Expected to find trace message in logs")
	}
	if !found["debug"] {
		t.Error("Expected to find debug message in logs")
	}
	if !found["info"] {
		t.Error("Expected to find info message in logs")
	}
	if !found["warning"] {
		t.Error("Expected to find warning message in logs")
	}
	if !found["error"] {
		t.Error("Expected to find error message in logs")
	}
}

func TestLoggerImplSetLogLevel(t *testing.T) {
	logMethod := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerImpl(logMethod)

	// Set log level to ERROR
	log.SetLogLevel(ifs.Error_Level)

	// Test that debug/info/warning are not logged
	log.Debug("should not be logged")
	log.Info("should not be logged")
	log.Warning("should not be logged")
	time.Sleep(100 * time.Millisecond)

	// Wait for queue to process
	for i := 0; i < 10; i++ {
		if log.Empty() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Should have no logs since we set level to ERROR
	if len(logMethod.logs) != 0 {
		t.Errorf("Expected no logs, got %d", len(logMethod.logs))
	}

	// Now test Error level
	log.Error("this should be logged")
	time.Sleep(100 * time.Millisecond)

	// Wait for queue to process
	for i := 0; i < 10; i++ {
		if log.Empty() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if len(logMethod.logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(logMethod.logs))
	}
}

func TestLoggerImplMultipleMethods(t *testing.T) {
	logMethod1 := &TestLogMethod{logs: make([]string, 0)}
	logMethod2 := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerImpl(logMethod1, logMethod2)

	log.Info("test message")
	time.Sleep(100 * time.Millisecond)

	// Wait for queue to process
	for i := 0; i < 10; i++ {
		if log.Empty() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	// Both log methods should have the message
	if len(logMethod1.logs) != 1 {
		t.Errorf("Expected 1 log in method1, got %d", len(logMethod1.logs))
	}
	if len(logMethod2.logs) != 1 {
		t.Errorf("Expected 1 log in method2, got %d", len(logMethod2.logs))
	}
}

func TestLoggerImplFail(t *testing.T) {
	logMethod := &TestLogMethod{logs: make([]string, 0)}
	log := logger.NewLoggerImpl(logMethod)

	// Create a mock testing.T - note: this won't actually fail the test
	// but it will test the Fail method code path
	mockT := &testing.T{}
	log.Fail(mockT, "fail message")
	time.Sleep(100 * time.Millisecond)

	// Wait for queue to process
	for i := 0; i < 10; i++ {
		if log.Empty() {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if len(logMethod.logs) == 0 {
		t.Error("Expected to have logs from Fail")
	}
}
