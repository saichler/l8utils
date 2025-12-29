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
	"sync"
	"testing"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8utils/go/utils/aggregator"
)

type SendCall struct {
	method      ifs.VNicMethod
	destination string
	serviceName string
	serviceArea byte
	action      ifs.Action
	data        interface{}
}

type MockVNic struct {
	calls     []SendCall
	mtx       sync.Mutex
	resources ifs.IResources
}

func NewMockVNic() *MockVNic {
	return &MockVNic{
		calls:     make([]SendCall, 0),
		resources: globals,
	}
}

func (m *MockVNic) recordCall(method ifs.VNicMethod, destination, serviceName string, serviceArea byte, action ifs.Action, data interface{}) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.calls = append(m.calls, SendCall{
		method:      method,
		destination: destination,
		serviceName: serviceName,
		serviceArea: serviceArea,
		action:      action,
		data:        data,
	})
}

func (m *MockVNic) GetCalls() []SendCall {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	result := make([]SendCall, len(m.calls))
	copy(result, m.calls)
	return result
}

func (m *MockVNic) Start()                        {}
func (m *MockVNic) Shutdown()                     {}
func (m *MockVNic) Name() string                  { return "MockVNic" }
func (m *MockVNic) SendMessage([]byte) error      { return nil }
func (m *MockVNic) Reply(*ifs.Message, ifs.IElements) error { return nil }
func (m *MockVNic) RoundRobinRequest(string, byte, ifs.Action, interface{}, int, ...string) ifs.IElements {
	return nil
}
func (m *MockVNic) ProximityRequest(string, byte, ifs.Action, interface{}, int, ...string) ifs.IElements {
	return nil
}
func (m *MockVNic) LeaderRequest(string, byte, ifs.Action, interface{}, int, ...string) ifs.IElements {
	return nil
}
func (m *MockVNic) Local(string, byte, ifs.Action, interface{}) error { return nil }
func (m *MockVNic) LocalRequest(string, byte, ifs.Action, interface{}, int, ...string) ifs.IElements {
	return nil
}
func (m *MockVNic) Forward(*ifs.Message, string) ifs.IElements { return nil }
func (m *MockVNic) ServiceAPI(string, byte) ifs.ServiceAPI     { return nil }
func (m *MockVNic) Resources() ifs.IResources                  { return m.resources }
func (m *MockVNic) NotifyServiceAdded([]string, byte) error    { return nil }
func (m *MockVNic) NotifyServiceRemoved(string, byte) error    { return nil }
func (m *MockVNic) PropertyChangeNotification(set *l8notify.L8NotificationSet) {}
func (m *MockVNic) WaitForConnection()            {}
func (m *MockVNic) Running() bool                 { return true }
func (m *MockVNic) SetResponse(*ifs.Message, ifs.IElements) {}
func (m *MockVNic) IsVnet() bool                  { return false }

func (m *MockVNic) Unicast(destination, serviceName string, serviceArea byte, action ifs.Action, data interface{}) error {
	m.recordCall(ifs.Unicast, destination, serviceName, serviceArea, action, data)
	return nil
}

func (m *MockVNic) Request(destination, serviceName string, serviceArea byte, action ifs.Action, data interface{}, timeout int, aaa ...string) ifs.IElements {
	m.recordCall(ifs.Request, destination, serviceName, serviceArea, action, data)
	return nil
}

func (m *MockVNic) Multicast(serviceName string, serviceArea byte, action ifs.Action, data interface{}) error {
	m.recordCall(ifs.Multicast, "", serviceName, serviceArea, action, data)
	return nil
}

func (m *MockVNic) RoundRobin(serviceName string, serviceArea byte, action ifs.Action, data interface{}) error {
	m.recordCall(ifs.RoundRobin, "", serviceName, serviceArea, action, data)
	return nil
}

func (m *MockVNic) Proximity(serviceName string, serviceArea byte, action ifs.Action, data interface{}) error {
	m.recordCall(ifs.Proximity, "", serviceName, serviceArea, action, data)
	return nil
}

func (m *MockVNic) Leader(serviceName string, serviceArea byte, action ifs.Action, data interface{}) error {
	m.recordCall(ifs.Leader, "", serviceName, serviceArea, action, data)
	return nil
}

func TestAggregatorSingleElement(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("test-data", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)

	// Wait for the flush interval
	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 1 {
		Log.Fail(t, "Expected 1 call, got", len(calls))
		return
	}

	if calls[0].method != ifs.Unicast {
		Log.Fail(t, "Expected Unicast method")
		return
	}
	if calls[0].destination != "dest-1" {
		Log.Fail(t, "Expected destination dest-1")
		return
	}
	if calls[0].serviceName != "service-1" {
		Log.Fail(t, "Expected serviceName service-1")
		return
	}
}

func TestAggregatorMultipleElementsSameAttributes(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.Multicast, "", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Multicast, "", "service-1", 1, ifs.POST)
	agg.AddElement("data-3", ifs.Multicast, "", "service-1", 1, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 1 {
		Log.Fail(t, "Expected 1 call for aggregated elements, got", len(calls))
		return
	}

	data, ok := calls[0].data.([]interface{})
	if !ok {
		Log.Fail(t, "Expected data to be []interface{}")
		return
	}
	if len(data) != 3 {
		Log.Fail(t, "Expected 3 elements in aggregated data, got", len(data))
		return
	}
}

func TestAggregatorDifferentMethods(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Multicast, "", "service-1", 1, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 2 {
		Log.Fail(t, "Expected 2 calls for different methods, got", len(calls))
		return
	}

	if calls[0].method != ifs.Unicast {
		Log.Fail(t, "Expected first call to be Unicast")
		return
	}
	if calls[1].method != ifs.Multicast {
		Log.Fail(t, "Expected second call to be Multicast")
		return
	}
}

func TestAggregatorDifferentDestinations(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Unicast, "dest-2", "service-1", 1, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 2 {
		Log.Fail(t, "Expected 2 calls for different destinations, got", len(calls))
		return
	}

	if calls[0].destination != "dest-1" {
		Log.Fail(t, "Expected first call destination to be dest-1")
		return
	}
	if calls[1].destination != "dest-2" {
		Log.Fail(t, "Expected second call destination to be dest-2")
		return
	}
}

func TestAggregatorDifferentServices(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.RoundRobin, "", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.RoundRobin, "", "service-2", 1, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 2 {
		Log.Fail(t, "Expected 2 calls for different services, got", len(calls))
		return
	}

	if calls[0].serviceName != "service-1" {
		Log.Fail(t, "Expected first call serviceName to be service-1")
		return
	}
	if calls[1].serviceName != "service-2" {
		Log.Fail(t, "Expected second call serviceName to be service-2")
		return
	}
}

func TestAggregatorDifferentActions(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.Proximity, "", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Proximity, "", "service-1", 1, ifs.PUT)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 2 {
		Log.Fail(t, "Expected 2 calls for different actions, got", len(calls))
		return
	}

	if calls[0].action != ifs.POST {
		Log.Fail(t, "Expected first call action to be POST")
		return
	}
	if calls[1].action != ifs.PUT {
		Log.Fail(t, "Expected second call action to be PUT")
		return
	}
}

func TestAggregatorDifferentServiceAreas(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("data-1", ifs.Leader, "", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Leader, "", "service-1", 2, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 2 {
		Log.Fail(t, "Expected 2 calls for different service areas, got", len(calls))
		return
	}

	if calls[0].serviceArea != 1 {
		Log.Fail(t, "Expected first call serviceArea to be 1")
		return
	}
	if calls[1].serviceArea != 2 {
		Log.Fail(t, "Expected second call serviceArea to be 2")
		return
	}
}

func TestAggregatorRequest(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	agg.AddElement("request-data", ifs.Request, "dest-1", "service-1", 1, ifs.GET)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 1 {
		Log.Fail(t, "Expected 1 call, got", len(calls))
		return
	}

	if calls[0].method != ifs.Request {
		Log.Fail(t, "Expected Request method")
		return
	}
}

func TestAggregatorShutdown(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 10, 5)

	agg.AddElement("data-1", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
	agg.Shutdown()

	// Add after shutdown should not cause panic
	agg.AddElement("data-2", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
}

func TestAggregatorMixedBatching(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	// Add elements that should be batched together
	agg.AddElement("data-1", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
	agg.AddElement("data-2", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)
	// Different destination - new batch
	agg.AddElement("data-3", ifs.Unicast, "dest-2", "service-1", 1, ifs.POST)
	// Back to first destination - new batch (not merged with first)
	agg.AddElement("data-4", ifs.Unicast, "dest-1", "service-1", 1, ifs.POST)

	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 3 {
		Log.Fail(t, "Expected 3 calls for mixed batching, got", len(calls))
		return
	}

	// First batch should have 2 elements
	data1, ok := calls[0].data.([]interface{})
	if !ok || len(data1) != 2 {
		Log.Fail(t, "Expected first batch to have 2 elements")
		return
	}

	// Second batch should have 1 element
	data2, ok := calls[1].data.([]interface{})
	if !ok || len(data2) != 1 {
		Log.Fail(t, "Expected second batch to have 1 element")
		return
	}

	// Third batch should have 1 element
	data3, ok := calls[2].data.([]interface{})
	if !ok || len(data3) != 1 {
		Log.Fail(t, "Expected third batch to have 1 element")
		return
	}
}

func TestAggregatorEmptyFlush(t *testing.T) {
	mockVNic := NewMockVNic()
	agg := aggregator.NewAggregator(mockVNic, 1, 5)
	defer agg.Shutdown()

	// Wait for flush without adding elements
	time.Sleep(time.Millisecond * 1500)

	calls := mockVNic.GetCalls()
	if len(calls) != 0 {
		Log.Fail(t, "Expected 0 calls for empty flush, got", len(calls))
		return
	}
}
