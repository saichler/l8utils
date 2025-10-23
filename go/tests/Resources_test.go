package tests

import (
	"testing"

	"github.com/saichler/l8utils/go/utils/logger"
	"github.com/saichler/l8utils/go/utils/resources"
)

func TestResources(t *testing.T) {
	// Create logger
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

	// Create resources
	r := resources.NewResources(log)

	// Test logger getter
	if r.Logger() == nil {
		t.Error("Logger should not be nil")
	}

	// Test other getters returning nil initially
	if r.Registry() != nil {
		t.Error("Registry should be nil initially")
	}
	if r.Services() != nil {
		t.Error("Services should be nil initially")
	}
	if r.Security() != nil {
		t.Error("Security should be nil initially")
	}
	if r.DataListener() != nil {
		t.Error("DataListener should be nil initially")
	}
	if r.Introspector() != nil {
		t.Error("Introspector should be nil initially")
	}
	if r.SysConfig() != nil {
		t.Error("SysConfig should be nil initially")
	}
}

func TestResourcesSet(t *testing.T) {
	// Create logger
	log := logger.NewLoggerDirectImpl(&logger.FmtLogMethod{})

	// Create resources
	r := resources.NewResources(log)

	// Test Set with nil (should do nothing)
	r.Set(nil)

	// Test Set with unknown type (should log error)
	r.Set("unknown type")
}
