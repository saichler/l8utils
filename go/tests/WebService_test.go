package tests

import (
	"testing"

	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/web"
)

func TestWebServiceNew(t *testing.T) {
	// Create a new web service
	ws := web.New(
		"test-service",
		byte(1),
		&testtypes.TestProto{}, &testtypes.TestProto{}, // POST
		&testtypes.TestProto{}, &testtypes.TestProto{}, // PUT
		&testtypes.TestProto{}, &testtypes.TestProto{}, // PATCH
		&testtypes.TestProto{}, &testtypes.TestProto{}, // DELETE
		&testtypes.TestProto{}, &testtypes.TestProto{}, // GET
	)

	// Test getters
	if ws.ServiceName() != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", ws.ServiceName())
	}

	if ws.ServiceArea() != byte(1) {
		t.Errorf("Expected service area 1, got %d", ws.ServiceArea())
	}

	// Test that type names are set
	if ws.PostBody() != "TestProto" {
		t.Errorf("Expected PostBody 'TestProto', got '%s'", ws.PostBody())
	}

	if ws.PostResp() != "TestProto" {
		t.Errorf("Expected PostResp 'TestProto', got '%s'", ws.PostResp())
	}

	if ws.PutBody() != "TestProto" {
		t.Errorf("Expected PutBody 'TestProto', got '%s'", ws.PutBody())
	}

	if ws.PutResp() != "TestProto" {
		t.Errorf("Expected PutResp 'TestProto', got '%s'", ws.PutResp())
	}

	if ws.PatchBody() != "TestProto" {
		t.Errorf("Expected PatchBody 'TestProto', got '%s'", ws.PatchBody())
	}

	if ws.PatchResp() != "TestProto" {
		t.Errorf("Expected PatchResp 'TestProto', got '%s'", ws.PatchResp())
	}

	if ws.DeleteBody() != "TestProto" {
		t.Errorf("Expected DeleteBody 'TestProto', got '%s'", ws.DeleteBody())
	}

	if ws.DeleteResp() != "TestProto" {
		t.Errorf("Expected DeleteResp 'TestProto', got '%s'", ws.DeleteResp())
	}

	if ws.GetBody() != "TestProto" {
		t.Errorf("Expected GetBody 'TestProto', got '%s'", ws.GetBody())
	}

	if ws.GetResp() != "TestProto" {
		t.Errorf("Expected GetResp 'TestProto', got '%s'", ws.GetResp())
	}
}

func TestWebServiceNilValues(t *testing.T) {
	// Create a web service with nil values
	ws := web.New(
		"test-service",
		byte(1),
		nil, nil, // POST
		nil, nil, // PUT
		nil, nil, // PATCH
		nil, nil, // DELETE
		nil, nil, // GET
	)

	// Test that empty strings are returned for nil values
	if ws.PostBody() != "" {
		t.Errorf("Expected empty PostBody, got '%s'", ws.PostBody())
	}

	if ws.PostResp() != "" {
		t.Errorf("Expected empty PostResp, got '%s'", ws.PostResp())
	}
}

func TestWebServiceSerialize(t *testing.T) {
	// Create a web service
	ws := web.New(
		"test-service",
		byte(1),
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
	)

	// Serialize
	l8ws := ws.Serialize()

	// Verify serialization
	if l8ws.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", l8ws.ServiceName)
	}

	if l8ws.ServiceArea != int32(1) {
		t.Errorf("Expected service area 1, got %d", l8ws.ServiceArea)
	}

	if l8ws.PostBodyType != "TestProto" {
		t.Errorf("Expected PostBodyType 'TestProto', got '%s'", l8ws.PostBodyType)
	}

	if l8ws.PostRespType != "TestProto" {
		t.Errorf("Expected PostRespType 'TestProto', got '%s'", l8ws.PostRespType)
	}
}

func TestWebServiceDeSerialize(t *testing.T) {
	// Create a web service
	ws := web.New(
		"original-service",
		byte(1),
		nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
	)

	// Create a serialized version manually
	ws2 := web.New(
		"test-service",
		byte(2),
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
		&testtypes.TestProto{}, &testtypes.TestProto{},
	)

	l8ws := ws2.Serialize()

	// Deserialize into the first web service
	ws.DeSerialize(l8ws)

	// Verify deserialization
	if ws.ServiceName() != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", ws.ServiceName())
	}

	if ws.ServiceArea() != byte(2) {
		t.Errorf("Expected service area 2, got %d", ws.ServiceArea())
	}

	if ws.PostBody() != "TestProto" {
		t.Errorf("Expected PostBody 'TestProto', got '%s'", ws.PostBody())
	}
}
