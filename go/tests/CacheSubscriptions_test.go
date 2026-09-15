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
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/cache"
)

// registerViaFetch is the only real entry point for subscription
// registration since Phase 3 (l8utils/plans/generic-websocket-change-notifications.md):
// Cache.RegisterSubscription is now private, called by Fetch() itself
// whenever the query has register=true and an AAAId. "register" is a plain
// substring match in the query text (l8ql's getBoolTag), no "=true" needed.
func registerViaFetch(c *cache.Cache, aaaId, gsql string, res ifs.IResources) ifs.IQuery {
	q := createIQuery(gsql+" register", res)
	q.SetAAAId(aaaId)
	c.Fetch(0, 10, q)
	return q
}

func TestSubscriptionRegister(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber, got %d", len(subs))
	}
	if subs[0].AAAId != "aaa-1" {
		t.Errorf("Expected AAAId 'aaa-1', got '%s'", subs[0].AAAId)
	}
	if subs[0].Query == nil {
		t.Fatal("Expected a non-nil Query")
	}
	if !subs[0].Query.Register() {
		t.Error("Expected the stored query to have Register()==true")
	}
}

func TestSubscriptionIgnoredWithoutRegisterTag(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	q := createIQuery("select * from TestProto", res)
	q.SetAAAId("aaa-1")
	c.Fetch(0, 10, q)

	if c.HasSubscribers() {
		t.Error("Expected no subscription from a query without register")
	}
}

func TestSubscriptionIgnoredWithoutAAAId(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	q := createIQuery("select * from TestProto register", res)
	c.Fetch(0, 10, q)

	if c.HasSubscribers() {
		t.Error("Expected no subscription from a registered query with no AAAId")
	}
}

func TestSubscriptionReplacesSameAAAId(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	registerViaFetch(c, "aaa-1", "select * from TestProto where MyString='hello'", res)

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber after replace, got %d", len(subs))
	}
	if subs[0].Query.Text() != "select * from TestProto where MyString='hello' register" {
		t.Errorf("Expected replaced query text, got '%s'", subs[0].Query.Text())
	}
}

func TestSubscriptionUnregister(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	registerViaFetch(c, "aaa-2", "select * from TestProto", res)
	c.UnregisterSubscription("aaa-1")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber after unregister, got %d", len(subs))
	}
	if subs[0].AAAId != "aaa-2" {
		t.Errorf("Expected remaining AAAId 'aaa-2', got '%s'", subs[0].AAAId)
	}
}

func TestSubscriptionUnregisterNonExistent(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	c.UnregisterSubscription("aaa-999")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber unchanged, got %d", len(subs))
	}
}

func TestSubscribersReturnsCopy(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)

	subs := c.Subscribers()
	subs[0].AAAId = "mutated"

	fresh := c.Subscribers()
	if fresh[0].AAAId != "aaa-1" {
		t.Errorf("Mutation leaked into cache: got '%s'", fresh[0].AAAId)
	}
}

func TestHasSubscribers(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	if c.HasSubscribers() {
		t.Error("Expected no subscribers on empty cache")
	}

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	if !c.HasSubscribers() {
		t.Error("Expected HasSubscribers true after register")
	}

	c.UnregisterSubscription("aaa-1")
	if c.HasSubscribers() {
		t.Error("Expected HasSubscribers false after unregister all")
	}
}

func TestSubscriptionMultipleAAAIds(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	for i := 0; i < 10; i++ {
		registerViaFetch(c, fmt.Sprintf("aaa-%d", i), "select * from TestProto", res)
	}

	subs := c.Subscribers()
	if len(subs) != 10 {
		t.Fatalf("Expected 10 subscribers, got %d", len(subs))
	}

	ids := make(map[string]bool)
	for _, s := range subs {
		ids[s.AAAId] = true
	}
	for i := 0; i < 10; i++ {
		if !ids[fmt.Sprintf("aaa-%d", i)] {
			t.Errorf("Missing aaa-%d", i)
		}
	}
}

func TestSubscriptionConcurrentSafety(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			aaaId := fmt.Sprintf("aaa-%d", idx)
			registerViaFetch(c, aaaId, "select * from TestProto", res)
			c.Subscribers()
			c.HasSubscribers()
			if idx%3 == 0 {
				c.UnregisterSubscription(aaaId)
			}
		}(i)
	}
	wg.Wait()

	subs := c.Subscribers()
	// 100 registered, every 3rd unregistered (indices 0,3,6,...,99 = 34 removed)
	// Remaining: 100 - 34 = 66
	expected := 66
	if len(subs) != expected {
		t.Errorf("Expected %d subscribers after concurrent ops, got %d", expected, len(subs))
	}
}

func TestSubscriptionEvictStale(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-old", "select * from TestProto", res)
	time.Sleep(2 * time.Second)
	registerViaFetch(c, "aaa-fresh", "select * from TestProto", res)

	evicted := c.EvictStaleSubscriptions(1)
	if evicted != 1 {
		t.Errorf("Expected 1 evicted, got %d", evicted)
	}

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 remaining, got %d", len(subs))
	}
	if subs[0].AAAId != "aaa-fresh" {
		t.Errorf("Expected 'aaa-fresh' to survive, got '%s'", subs[0].AAAId)
	}
}

func TestSubscriptionRefreshPreventsEviction(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	time.Sleep(2 * time.Second)
	// Re-fetching the same registered query refreshes lastSeen.
	registerViaFetch(c, "aaa-1", "select * from TestProto", res)

	evicted := c.EvictStaleSubscriptions(1)
	if evicted != 0 {
		t.Errorf("Expected 0 evicted after refresh, got %d", evicted)
	}
	if !c.HasSubscribers() {
		t.Error("Expected subscriber to survive after refresh")
	}
}

// TestClientNotificationMatchesQuery exercises the real end-to-end flow added
// in Phases 2/3: Fetch with register=true creates the subscription, and a
// later write only notifies a subscriber whose query actually matches the
// changed record, carrying the full record
// (l8utils/plans/generic-websocket-change-notifications.md).
func TestClientNotificationMatchesQuery(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-match", "select * from TestProto where MyString='wanted'", res)
	registerViaFetch(c, "aaa-nomatch", "select * from TestProto where MyString='other'", res)

	t1 := createModel(1)
	t1.MyString = "wanted"
	_, cn, err := c.Post(t1, true)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if cn == nil {
		t.Fatal("Expected a non-nil client notification for a matching subscriber")
	}
	if len(cn.AaaIds) != 1 || !cn.AaaIds["aaa-match"] {
		t.Errorf("Expected AaaIds to contain only 'aaa-match', got %v", cn.AaaIds)
	}
	if cn.NotificationList == nil {
		t.Error("Expected NotificationList (the record) to be populated")
	}
}

func TestClientNotificationNoMatchIsNil(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-nomatch", "select * from TestProto where MyString='other'", res)

	t1 := createModel(1)
	t1.MyString = "wanted"
	_, cn, err := c.Post(t1, true)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if cn != nil {
		t.Errorf("Expected nil client notification when no subscriber matches, got %v", cn)
	}
}

func TestClientNotificationNoSubscribersIsNil(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	t1 := createModel(1)
	_, cn, err := c.Post(t1, true)
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if cn != nil {
		t.Errorf("Expected nil client notification with no subscribers, got %v", cn)
	}
}

// TestFetchRegistersAndRefreshes is the direct test of Phase 3's own change:
// Fetch() itself is what creates/refreshes the subscription, unconditionally
// on every call with register=true, not just the first.
func TestFetchRegistersAndRefreshes(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	q := createIQuery("select * from TestProto register", res)
	q.SetAAAId("aaa-1")

	c.Fetch(0, 10, q)
	if !c.HasSubscribers() {
		t.Fatal("Expected Fetch with register=true to create a subscription")
	}

	// A second Fetch of the same query must not error or duplicate --
	// still exactly one subscriber for this AAAId.
	c.Fetch(0, 10, q)
	if len(c.Subscribers()) != 1 {
		t.Errorf("Expected exactly 1 subscriber after two fetches, got %d", len(c.Subscribers()))
	}
}
