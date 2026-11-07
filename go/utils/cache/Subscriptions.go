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

package cache

import (
	"sync"

	"github.com/saichler/l8types/go/ifs"
)

// Subscription represents one client component's interest in change
// notifications for a specific model type. A user (AAAId) can hold several,
// one per component, told apart by RegisterId; the pair is the key.
// Query is the live, already-parsed query object -- stored directly (not
// re-derived from text) so a write can call Query.Match() against the
// changed record without re-parsing
// (l8utils/plans/generic-websocket-change-notifications.md Phase 2).
//
// A subscription lives until its owner unregisters it (an evicted client
// component, or the web server when the user's websocket closes); there is
// no time-based expiry, since a live view applies pushes in place and does
// not re-fetch to stay registered.
type Subscription struct {
	AAAId      string
	RegisterId string
	Query      ifs.IQuery
}

func subscriptionKey(aaaId, registerId string) string {
	return aaaId + "/" + registerId
}

// subscriptions tracks the registered client components for a single Cache
// instance (one model type). Thread-safe for concurrent access.
type subscriptions struct {
	mu   sync.RWMutex
	subs map[string]*Subscription
}

func newSubscriptions() *subscriptions {
	return &subscriptions{
		subs: make(map[string]*Subscription),
	}
}

func (this *subscriptions) register(sub *Subscription) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.subs[subscriptionKey(sub.AAAId, sub.RegisterId)] = sub
}

func (this *subscriptions) unregister(aaaId, registerId string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.subs, subscriptionKey(aaaId, registerId))
}

func (this *subscriptions) subscribers() []*Subscription {
	this.mu.RLock()
	defer this.mu.RUnlock()
	result := make([]*Subscription, 0, len(this.subs))
	for _, sub := range this.subs {
		cp := *sub
		result = append(result, &cp)
	}
	return result
}

func (this *subscriptions) hasSubscribers() bool {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return len(this.subs) > 0
}
