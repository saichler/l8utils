package tests

import (
	"reflect"
	"testing"

	"github.com/saichler/l8utils/go/utils/registry"
)

type TestStruct struct {
	Name string
}

func TestTypesMapDelAndContains(t *testing.T) {
	tm := registry.NewTypesMap()

	// Test Contains on empty map
	if tm.Contains("TestStruct") {
		t.Error("Should not contain TestStruct initially")
	}

	// Add a value
	testType := reflect.TypeOf(TestStruct{})
	tm.Put("TestStruct", testType)

	// Test Contains
	if !tm.Contains("TestStruct") {
		t.Error("Should contain TestStruct after adding")
	}

	// Test Del
	tm.Del("TestStruct")

	// Test Contains after delete
	if tm.Contains("TestStruct") {
		t.Error("Should not contain TestStruct after deletion")
	}
}
