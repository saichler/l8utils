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

	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8utils/go/utils/cache"
)

func TestSubscriptionRegister(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "hash-abc", "select * from TestProto")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber, got %d", len(subs))
	}
	if subs[0].Token != "token-1" {
		t.Errorf("Expected token 'token-1', got '%s'", subs[0].Token)
	}
	if subs[0].QueryHash != "hash-abc" {
		t.Errorf("Expected hash 'hash-abc', got '%s'", subs[0].QueryHash)
	}
	if subs[0].QueryText != "select * from TestProto" {
		t.Errorf("Expected query text, got '%s'", subs[0].QueryText)
	}
}

func TestSubscriptionReplacesSameToken(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "hash-old", "select * from TestProto")
	c.RegisterSubscription("token-1", "hash-new", "select * from TestProto where MyString=hello")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber after replace, got %d", len(subs))
	}
	if subs[0].QueryHash != "hash-new" {
		t.Errorf("Expected replaced hash 'hash-new', got '%s'", subs[0].QueryHash)
	}
}

func TestSubscriptionUnregister(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "h1", "q1")
	c.RegisterSubscription("token-2", "h2", "q2")
	c.UnregisterSubscription("token-1")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber after unregister, got %d", len(subs))
	}
	if subs[0].Token != "token-2" {
		t.Errorf("Expected remaining token 'token-2', got '%s'", subs[0].Token)
	}
}

func TestSubscriptionUnregisterNonExistent(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "h1", "q1")
	c.UnregisterSubscription("token-999")

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 subscriber unchanged, got %d", len(subs))
	}
}

func TestSubscribersReturnsCopy(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "h1", "q1")

	subs := c.Subscribers()
	subs[0].Token = "mutated"

	fresh := c.Subscribers()
	if fresh[0].Token != "token-1" {
		t.Errorf("Mutation leaked into cache: got '%s'", fresh[0].Token)
	}
}

func TestHasSubscribers(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	if c.HasSubscribers() {
		t.Error("Expected no subscribers on empty cache")
	}

	c.RegisterSubscription("token-1", "h1", "q1")
	if !c.HasSubscribers() {
		t.Error("Expected HasSubscribers true after register")
	}

	c.UnregisterSubscription("token-1")
	if c.HasSubscribers() {
		t.Error("Expected HasSubscribers false after unregister all")
	}
}

func TestSubscriptionMultipleTokens(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	for i := 0; i < 10; i++ {
		c.RegisterSubscription(fmt.Sprintf("token-%d", i), fmt.Sprintf("h%d", i), fmt.Sprintf("q%d", i))
	}

	subs := c.Subscribers()
	if len(subs) != 10 {
		t.Fatalf("Expected 10 subscribers, got %d", len(subs))
	}

	tokens := make(map[string]bool)
	for _, s := range subs {
		tokens[s.Token] = true
	}
	for i := 0; i < 10; i++ {
		if !tokens[fmt.Sprintf("token-%d", i)] {
			t.Errorf("Missing token-%d", i)
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
			token := fmt.Sprintf("token-%d", idx)
			c.RegisterSubscription(token, "h", "q")
			c.Subscribers()
			c.HasSubscribers()
			if idx%3 == 0 {
				c.UnregisterSubscription(token)
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

	c.RegisterSubscription("token-old", "h1", "q1")
	time.Sleep(2 * time.Second)
	c.RegisterSubscription("token-fresh", "h2", "q2")

	evicted := c.EvictStaleSubscriptions(1)
	if evicted != 1 {
		t.Errorf("Expected 1 evicted, got %d", evicted)
	}

	subs := c.Subscribers()
	if len(subs) != 1 {
		t.Fatalf("Expected 1 remaining, got %d", len(subs))
	}
	if subs[0].Token != "token-fresh" {
		t.Errorf("Expected 'token-fresh' to survive, got '%s'", subs[0].Token)
	}
}

func TestSubscriptionRefreshPreventsEviction(t *testing.T) {
	res := newResources()
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	defer c.Close()

	c.RegisterSubscription("token-1", "h1", "q1")
	time.Sleep(2 * time.Second)
	// Re-register refreshes lastSeen
	c.RegisterSubscription("token-1", "h1", "q1")

	evicted := c.EvictStaleSubscriptions(1)
	if evicted != 0 {
		t.Errorf("Expected 0 evicted after refresh, got %d", evicted)
	}
	if !c.HasSubscribers() {
		t.Error("Expected subscriber to survive after refresh")
	}
}
