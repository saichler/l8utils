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

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8utils/go/utils/cache"
)

// fetchWithRegister sends a query through Fetch the way a client's does: the
// register action and register id are fields of the query body (not its
// text), carried onto the parsed query by object.NewFromQuery.
func fetchWithRegister(c *cache.Cache, aaaId, registerId string, action l8api.L8Register,
	gsql string, res ifs.IResources) ([]interface{}, ifs.IQuery) {
	elems, e := object.NewFromQuery(&l8api.L8Query{Text: gsql, AaaId: aaaId,
		Register: action, RegisterId: registerId}, res)
	if e != nil {
		panic(e)
	}
	q, e := elems.Query(res)
	if e != nil {
		panic(e)
	}
	values, _ := c.Fetch(0, 10, q)
	return values, q
}

// registerViaFetch registers the client component aaaId+"main" -- Fetch() is
// the only entry point for registration since Phase 3
// (l8utils/plans/generic-websocket-change-notifications.md).
func registerViaFetch(c *cache.Cache, aaaId, gsql string, res ifs.IResources) ifs.IQuery {
	_, q := fetchWithRegister(c, aaaId, "main", l8api.L8Register_L8_REGISTER_REGISTER, gsql, res)
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
	if subs[0].Query.Register() != l8api.L8Register_L8_REGISTER_REGISTER {
		t.Error("Expected the stored query to have Register()==Register")
	}
	if subs[0].RegisterId != "main" {
		t.Errorf("Expected RegisterId 'main', got '%s'", subs[0].RegisterId)
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

	fetchWithRegister(c, "", "main", l8api.L8Register_L8_REGISTER_REGISTER, "select * from TestProto", res)

	if c.HasSubscribers() {
		t.Error("Expected no subscription from a registered query with no AAAId")
	}
}

func TestSubscriptionReplacesSameComponent(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	registerViaFetch(c, "aaa-1", "select * from TestProto where MyString='hello'", res)

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber after replace, got %d", len(subs))
	}
	if subs[0].Query.Text() != "select * from TestProto where MyString='hello'" {
		t.Errorf("Expected replaced query text, got '%s'", subs[0].Query.Text())
	}
}

func TestSubscriptionUnregister(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	registerViaFetch(c, "aaa-1", "select * from TestProto", res)
	registerViaFetch(c, "aaa-2", "select * from TestProto", res)
	c.UnregisterSubscription("aaa-1", "main")

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
	c.UnregisterSubscription("aaa-999", "main")

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

	c.UnregisterSubscription("aaa-1", "main")
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
				c.UnregisterSubscription(aaaId, "main")
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

// One user holds one registration per client component: two components on
// the same model keep their own queries, and unregistering one leaves the
// other untouched.
func TestSubscriptionOnePerComponent(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	fetchWithRegister(c, "aaa-1", "table", l8api.L8Register_L8_REGISTER_REGISTER,
		"select * from TestProto where MyString='a'", res)
	fetchWithRegister(c, "aaa-1", "chart", l8api.L8Register_L8_REGISTER_REGISTER,
		"select * from TestProto where MyString='b'", res)
	if n := len(c.Subscribers()); n != 2 {
		t.Fatalf("Expected 2 subscriptions for two components, got %d", n)
	}

	fetchWithRegister(c, "aaa-1", "table", l8api.L8Register_L8_REGISTER_UNREGISTER,
		"select * from TestProto where MyString='a'", res)
	subs := c.Subscribers()
	if len(subs) != 1 || subs[0].RegisterId != "chart" {
		t.Fatalf("Expected only 'chart' to remain, got %v", subs)
	}
}

// An Unregister removes the registration and fetches nothing: the client
// only sends it to drop interest, never for data.
func TestUnregisterFetchesNothing(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()
	c.Post(createModel(1), false)

	values, _ := fetchWithRegister(c, "aaa-1", "main", l8api.L8Register_L8_REGISTER_REGISTER,
		"select * from TestProto", res)
	if len(values) != 1 {
		t.Fatalf("Expected the registering fetch to return the row, got %d", len(values))
	}
	values, _ = fetchWithRegister(c, "aaa-1", "main", l8api.L8Register_L8_REGISTER_UNREGISTER,
		"select * from TestProto", res)
	if len(values) != 0 {
		t.Errorf("Expected Unregister to return no rows, got %d", len(values))
	}
	if c.HasSubscribers() {
		t.Error("Expected Unregister to remove the subscription")
	}
}

// Old clients put "register=true" in the query text; it still registers,
// under an empty register id.
func TestRegisterFromQueryText(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	q := createIQuery("select * from TestProto register=true", res)
	q.SetAAAId("aaa-1")
	c.Fetch(0, 10, q)
	subs := c.Subscribers()
	if len(subs) != 1 || subs[0].RegisterId != "" {
		t.Fatalf("Expected one subscription with an empty register id, got %v", subs)
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
// Fetch() itself is what creates/replaces the subscription, on every call
// with Register, not just the first.
func TestFetchRegistersAndRefreshes(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	fetchWithRegister(c, "aaa-1", "main", l8api.L8Register_L8_REGISTER_REGISTER, "select * from TestProto", res)
	if !c.HasSubscribers() {
		t.Fatal("Expected Fetch with Register to create a subscription")
	}

	// A second Fetch of the same query must not error or duplicate --
	// still exactly one subscriber for this component.
	fetchWithRegister(c, "aaa-1", "main", l8api.L8Register_L8_REGISTER_REGISTER, "select * from TestProto", res)
	if len(c.Subscribers()) != 1 {
		t.Errorf("Expected exactly 1 subscriber after two fetches, got %d", len(c.Subscribers()))
	}
}
