// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"testing"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/web"
)

func TestWebServiceNew(t *testing.T) {
	// Create a new web service
	ws := web.New("test-service", byte(1), 0)
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.POST, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.PUT, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.PATCH, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.DELETE, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.GET, &testtypes.TestProto{})

	// Test getters
	if ws.ServiceName() != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", ws.ServiceName())
	}

	if ws.ServiceArea() != byte(1) {
		t.Errorf("Expected service area 1, got %d", ws.ServiceArea())
	}

	if ws.Vnet() != 0 {
		t.Errorf("Expected vnet 0, got %d", ws.Vnet())
	}

	// Test that endpoints are registered
	l8ws := ws.Serialize()
	if len(l8ws.Endpoints) != 5 {
		t.Errorf("Expected 5 endpoints, got %d", len(l8ws.Endpoints))
	}

	// Verify POST endpoint
	postEndpoint := l8ws.Endpoints[int32(ifs.POST)]
	if postEndpoint == nil {
		t.Error("Expected POST endpoint to exist")
	} else {
		if postEndpoint.PrimaryBody != "TestProto" {
			t.Errorf("Expected POST primary body 'TestProto', got '%s'", postEndpoint.PrimaryBody)
		}
		if postEndpoint.Body2Response["TestProto"] != "TestProto" {
			t.Errorf("Expected POST response 'TestProto', got '%s'", postEndpoint.Body2Response["TestProto"])
		}
	}
}

func TestWebServiceNilValues(t *testing.T) {
	// Create a web service with nil body/response values
	ws := web.New("test-service", byte(1), 0)
	ws.AddEndpoint(nil, ifs.POST, nil)

	// Test that L8Empty is used for nil values
	l8ws := ws.Serialize()
	postEndpoint := l8ws.Endpoints[int32(ifs.POST)]
	if postEndpoint == nil {
		t.Error("Expected POST endpoint to exist")
	} else {
		if postEndpoint.PrimaryBody != "L8Empty" {
			t.Errorf("Expected POST primary body 'L8Empty', got '%s'", postEndpoint.PrimaryBody)
		}
		if postEndpoint.Body2Response["L8Empty"] != "L8Empty" {
			t.Errorf("Expected POST response 'L8Empty', got '%s'", postEndpoint.Body2Response["L8Empty"])
		}
	}
}

func TestWebServiceSerialize(t *testing.T) {
	// Create a web service
	ws := web.New("test-service", byte(1), 100)
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.POST, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.GET, &testtypes.TestProto{})

	// Serialize
	l8ws := ws.Serialize()

	// Verify serialization
	if l8ws.ServiceName != "test-service" {
		t.Errorf("Expected service name 'test-service', got '%s'", l8ws.ServiceName)
	}

	if l8ws.ServiceArea != int32(1) {
		t.Errorf("Expected service area 1, got %d", l8ws.ServiceArea)
	}

	if l8ws.Vnet != 100 {
		t.Errorf("Expected vnet 100, got %d", l8ws.Vnet)
	}

	if len(l8ws.Endpoints) != 2 {
		t.Errorf("Expected 2 endpoints, got %d", len(l8ws.Endpoints))
	}

	// Verify POST endpoint structure
	postEndpoint := l8ws.Endpoints[int32(ifs.POST)]
	if postEndpoint == nil {
		t.Error("Expected POST endpoint to exist")
	} else {
		if postEndpoint.PrimaryBody != "TestProto" {
			t.Errorf("Expected POST primary body 'TestProto', got '%s'", postEndpoint.PrimaryBody)
		}
	}
}

func TestWebServiceVnet(t *testing.T) {
	// Test vnet getter
	ws := web.New("test-service", byte(1), 12345)

	if ws.Vnet() != 12345 {
		t.Errorf("Expected vnet 12345, got %d", ws.Vnet())
	}
}

func TestWebServiceMultipleBodyTypes(t *testing.T) {
	// Test adding multiple body types to the same action
	ws := web.New("test-service", byte(1), 0)
	ws.AddEndpoint(&testtypes.TestProto{}, ifs.POST, &testtypes.TestProto{})
	ws.AddEndpoint(&testtypes.TestProtoSub{}, ifs.POST, &testtypes.TestProtoSub{})

	l8ws := ws.Serialize()
	postEndpoint := l8ws.Endpoints[int32(ifs.POST)]
	if postEndpoint == nil {
		t.Error("Expected POST endpoint to exist")
	} else {
		// Primary body should be the first one added
		if postEndpoint.PrimaryBody != "TestProto" {
			t.Errorf("Expected POST primary body 'TestProto', got '%s'", postEndpoint.PrimaryBody)
		}
		// Both body types should be in the map
		if len(postEndpoint.Body2Response) != 2 {
			t.Errorf("Expected 2 body types, got %d", len(postEndpoint.Body2Response))
		}
		if postEndpoint.Body2Response["TestProto"] != "TestProto" {
			t.Errorf("Expected TestProto response 'TestProto', got '%s'", postEndpoint.Body2Response["TestProto"])
		}
		if postEndpoint.Body2Response["TestProtoSub"] != "TestProtoSub" {
			t.Errorf("Expected TestProtoSub response 'TestProtoSub', got '%s'", postEndpoint.Body2Response["TestProtoSub"])
		}
	}
}
