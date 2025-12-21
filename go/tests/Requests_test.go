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
	"github.com/saichler/l8utils/go/utils/requests"
)

// Mock IElements implementation for testing
type MockElement struct {
	data string
}

func (m *MockElement) Element() any {
	return m.data
}

func (m *MockElement) Elements() []any {
	return []any{m.data}
}

func (m *MockElement) Keys() []any {
	return nil
}

func (m *MockElement) Errors() []error {
	return nil
}

func (m *MockElement) Query(ifs.IResources) (ifs.IQuery, error) {
	return nil, nil
}

func (m *MockElement) Key() any {
	return nil
}

func (m *MockElement) Error() error {
	return nil
}

func (m *MockElement) Serialize() ([]byte, error) {
	return nil, nil
}

func (m *MockElement) Deserialize([]byte, ifs.IRegistry) error {
	return nil
}

func (m *MockElement) Notification() bool {
	return false
}

func (m *MockElement) Append(ifs.IElements) {
}

func (m *MockElement) AsList(ifs.IRegistry) (any, error) {
	return nil, nil
}

func (m *MockElement) IsFilterMode() bool {
	return false
}

func (m *MockElement) IsReplica() bool {
	return false
}

func (m *MockElement) Replica() byte {
	return 0
}

// Mock IElements for transaction testing
type MockTransaction struct {
	End uint32
}

func (m *MockTransaction) Element() any {
	return &MockTransactionData{End: m.End}
}

func (m *MockTransaction) Elements() []any {
	return []any{&MockTransactionData{End: m.End}}
}

func (m *MockTransaction) Keys() []any {
	return nil
}

func (m *MockTransaction) Errors() []error {
	return nil
}

func (m *MockTransaction) Query(ifs.IResources) (ifs.IQuery, error) {
	return nil, nil
}

func (m *MockTransaction) Key() any {
	return nil
}

func (m *MockTransaction) Error() error {
	return nil
}

func (m *MockTransaction) Serialize() ([]byte, error) {
	return nil, nil
}

func (m *MockTransaction) Deserialize([]byte, ifs.IRegistry) error {
	return nil
}

func (m *MockTransaction) Notification() bool {
	return false
}

func (m *MockTransaction) Append(ifs.IElements) {
}

func (m *MockTransaction) AsList(ifs.IRegistry) (any, error) {
	return nil, nil
}

func (m *MockTransaction) IsFilterMode() bool {
	return false
}

func (m *MockTransaction) IsReplica() bool {
	return false
}

func (m *MockTransaction) Replica() byte {
	return 0
}

type MockTransactionData struct {
	End uint32
}

// TestNewRequests tests the creation of Requests manager
func TestNewRequests(t *testing.T) {
	reqs := requests.NewRequests()
	if reqs == nil {
		Log.Fail(t, "NewRequests should not return nil")
		return
	}
}

// TestNewRequest tests creating a new request
func TestNewRequest(t *testing.T) {
	reqs := requests.NewRequests()
	req, err := reqs.NewRequest(1, "source1", 5, Log)
	if err != nil {
		Log.Fail(t, "NewRequest should not fail:", err)
		return
	}
	if req == nil {
		Log.Fail(t, "NewRequest should return a request")
		return
	}
	if req.MsgNum() != 1 {
		Log.Fail(t, "Expected message number to be 1")
		return
	}
	if req.MsgSource() != "source1" {
		Log.Fail(t, "Expected message source to be 'source1'")
		return
	}
}

// TestDuplicateRequest tests that duplicate requests return an error
func TestDuplicateRequest(t *testing.T) {
	reqs := requests.NewRequests()
	_, err := reqs.NewRequest(1, "source1", 5, Log)
	if err != nil {
		Log.Fail(t, "First request should succeed")
		return
	}

	// Try to create duplicate
	_, err = reqs.NewRequest(1, "source1", 5, Log)
	if err == nil {
		Log.Fail(t, "Duplicate request should return an error")
		return
	}
	if err != requests.ErrDuplicateRequest {
		Log.Fail(t, "Error should be ErrDuplicateRequest")
		return
	}
}

// TestGetRequest tests retrieving a request
func TestGetRequest(t *testing.T) {
	reqs := requests.NewRequests()
	req1, _ := reqs.NewRequest(1, "source1", 5, Log)

	req2 := reqs.GetRequest(1, "source1")
	if req2 == nil {
		Log.Fail(t, "GetRequest should return the request")
		return
	}
	if req1 != req2 {
		Log.Fail(t, "GetRequest should return the same request instance")
		return
	}

	// Test non-existent request
	req3 := reqs.GetRequest(999, "nonexistent")
	if req3 != nil {
		Log.Fail(t, "GetRequest should return nil for non-existent request")
		return
	}
}

// TestDelRequest tests manual deletion of a request
func TestDelRequest(t *testing.T) {
	reqs := requests.NewRequests()
	reqs.NewRequest(1, "source1", 5, Log)

	// Verify it exists
	req := reqs.GetRequest(1, "source1")
	if req == nil {
		Log.Fail(t, "Request should exist before deletion")
		return
	}

	// Delete it
	reqs.DelRequest(1, "source1")

	// Verify it's gone
	req = reqs.GetRequest(1, "source1")
	if req != nil {
		Log.Fail(t, "Request should be nil after deletion")
		return
	}
}

// TestRequestTimeout tests that requests timeout correctly
func TestRequestTimeout(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 1, Log) // 1 second timeout

	start := time.Now()
	resp := req.Wait()
	elapsed := time.Since(start)

	// Should have timed out after ~1 second
	if elapsed < 900*time.Millisecond || elapsed > 1200*time.Millisecond {
		Log.Fail(t, "Request should timeout after ~1 second, got:", elapsed)
		return
	}

	if resp.Error() == nil {
		// Check if it's an error response
		if resp.Element() == nil {
			Log.Fail(t, "Timeout should return an error response")
			return
		}
	}

	// Verify the request was auto-cleaned up
	req2 := reqs.GetRequest(1, "source1")
	if req2 != nil {
		Log.Fail(t, "Request should be auto-cleaned up after timeout")
		return
	}
}

// TestRequestResponse tests that responses are delivered correctly
func TestRequestResponse(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 10, Log)

	mockResp := &MockElement{data: "test response"}

	// Send response in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		req.SetResponse(mockResp)
	}()

	start := time.Now()
	resp := req.Wait()
	elapsed := time.Since(start)

	// Should return quickly (within 200ms, not wait for full timeout)
	if elapsed > 300*time.Millisecond {
		Log.Fail(t, "Request should return quickly when response is set, got:", elapsed)
		return
	}

	if resp == nil || resp.Element() != mockResp.Element() {
		Log.Fail(t, "Response should match the mock response")
		return
	}

	// Verify auto-cleanup
	req2 := reqs.GetRequest(1, "source1")
	if req2 != nil {
		Log.Fail(t, "Request should be auto-cleaned up after response")
		return
	}
}

// TestRequestResponseMethod tests the Response() method
func TestRequestResponseMethod(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 10, Log)

	mockResp := &MockElement{data: "test response"}

	// Send response
	go func() {
		time.Sleep(100 * time.Millisecond)
		req.SetResponse(mockResp)
	}()

	req.Wait()

	// Test Response() method
	resp := req.Response()
	if resp == nil || resp.Element() != mockResp.Element() {
		Log.Fail(t, "Response() should return the set response")
		return
	}
}

// TestConcurrentRequests tests multiple concurrent requests
func TestConcurrentRequests(t *testing.T) {
	reqs := requests.NewRequests()
	numRequests := 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			defer wg.Done()
			req, err := reqs.NewRequest(uint32(idx), "source", 5, Log)
			if err != nil {
				Log.Fail(t, "Failed to create request:", idx, err)
				return
			}

			// Simulate response
			go func() {
				time.Sleep(10 * time.Millisecond)
				req.SetResponse(&MockElement{data: "response"})
			}()

			resp := req.Wait()
			if resp == nil {
				Log.Fail(t, "Response should not be nil for request:", idx)
				return
			}
		}(i)
	}

	wg.Wait()

	// Verify all requests were cleaned up
	for i := 0; i < numRequests; i++ {
		req := reqs.GetRequest(uint32(i), "source")
		if req != nil {
			Log.Fail(t, "Request", i, "should be cleaned up")
			return
		}
	}
}

// TestMultipleResponsesOnlyFirstCounts tests that only the first response is processed
func TestMultipleResponsesOnlyFirstCounts(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 10, Log)

	mockResp1 := &MockElement{data: "first response"}
	mockResp2 := &MockElement{data: "second response"}

	// Send multiple responses
	go func() {
		time.Sleep(50 * time.Millisecond)
		req.SetResponse(mockResp1)
		req.SetResponse(mockResp2) // This should be ignored
	}()

	resp := req.Wait()

	// Should get the first response
	if resp == nil || resp.Element() != mockResp1.Element() {
		Log.Fail(t, "Should receive the first response")
		return
	}
}

// TestNilResponseIgnored tests that nil responses are ignored
func TestNilResponseIgnored(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 2, Log)

	// Send nil response
	go func() {
		time.Sleep(100 * time.Millisecond)
		req.SetResponse(nil) // Should be ignored
	}()

	start := time.Now()
	resp := req.Wait()
	elapsed := time.Since(start)

	// Should timeout since nil is ignored
	if elapsed < 1800*time.Millisecond {
		Log.Fail(t, "Should have timed out since nil response is ignored")
		return
	}

	// Should get timeout error
	if resp.Error() == nil {
		if resp.Element() == nil {
			Log.Fail(t, "Should get timeout error")
			return
		}
	}
}

// TestMemoryCleanup verifies no goroutine leaks by testing rapid request creation/completion
func TestMemoryCleanup(t *testing.T) {
	reqs := requests.NewRequests()
	numIterations := 1000

	for i := 0; i < numIterations; i++ {
		req, _ := reqs.NewRequest(uint32(i), "source", 10, Log)

		// Immediately send response
		req.SetResponse(&MockElement{data: "response"})
		req.Wait()

		// Verify cleanup
		if reqs.GetRequest(uint32(i), "source") != nil {
			Log.Fail(t, "Request should be cleaned up immediately")
			return
		}
	}
}

// TestRaceConditions tests concurrent read/write access
func TestRaceConditions(t *testing.T) {
	reqs := requests.NewRequests()
	req, _ := reqs.NewRequest(1, "source1", 10, Log)

	var wg sync.WaitGroup
	wg.Add(3)

	// Concurrent SetResponse calls
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			req.SetResponse(&MockElement{data: "response"})
		}
	}()

	// Concurrent Response() calls
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = req.Response()
			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Wait for response
	go func() {
		defer wg.Done()
		req.Wait()
	}()

	wg.Wait()
}

// TestDifferentSourcesAndNumbers tests requests with different sources and numbers
func TestDifferentSourcesAndNumbers(t *testing.T) {
	reqs := requests.NewRequests()

	// Create requests with same number, different sources
	req1, _ := reqs.NewRequest(1, "source1", 10, Log)
	req2, _ := reqs.NewRequest(1, "source2", 10, Log)

	if req1 == req2 {
		Log.Fail(t, "Requests with different sources should be different")
		return
	}

	// Create requests with same source, different numbers
	req3, _ := reqs.NewRequest(2, "source1", 10, Log)

	if req1 == req3 {
		Log.Fail(t, "Requests with different numbers should be different")
		return
	}

	// Cleanup
	req1.SetResponse(&MockElement{data: "r1"})
	req2.SetResponse(&MockElement{data: "r2"})
	req3.SetResponse(&MockElement{data: "r3"})

	req1.Wait()
	req2.Wait()
	req3.Wait()
}

// TestFastResponseBeforeTimeout tests that timeout goroutine is cleaned up properly
func TestFastResponseBeforeTimeout(t *testing.T) {
	reqs := requests.NewRequests()

	// Create many requests with long timeouts but fast responses
	numRequests := 100
	var wg sync.WaitGroup
	wg.Add(numRequests)

	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			defer wg.Done()
			req, _ := reqs.NewRequest(uint32(idx), "source", 60, Log) // Long timeout

			// Send response immediately
			req.SetResponse(&MockElement{data: "fast"})

			// Should return immediately, not wait for 60 seconds
			start := time.Now()
			req.Wait()
			elapsed := time.Since(start)

			if elapsed > 100*time.Millisecond {
				Log.Fail(t, "Request should return immediately, not wait for timeout")
				return
			}
		}(i)
	}

	wg.Wait()
}
