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
	"github.com/saichler/l8types/go/types/l8notify"
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
//
// A paged subscription also carries a snapshot of the page the component
// displays, taken by the registering Fetch and replaced by each re-fetch, so a
// write notifies only the components whose displayed page it affects (see
// wants). An unpaged one (no block size, or an aggregate) matches by its
// WHERE clause alone.
type Subscription struct {
	AAAId      string
	RegisterId string
	Query      ifs.IQuery

	// Paged is true when the registering Fetch returned a page (block size > 0).
	Paged bool
	// PageKeys are the primary keys of the displayed page.
	PageKeys map[string]bool
	// FirstPage is true when the displayed page starts at the first element.
	FirstPage bool
	// Full is true when the page holds a whole block (more may follow it).
	Full bool
	// First and Last are the sort-order values (the sort-by value, or the
	// primary key when there is no sort-by) of the page's first and last
	// elements; nil for an empty page.
	First, Last interface{}
}

// wants reports whether a write of value (primary key key, notification type
// t) belongs to what this subscription's component displays: the value must
// match the query's WHERE clause and, for a paged subscription, an update or
// delete must be of an element on the displayed page, and an add must sort
// into the displayed page's range (or onto a page that is not full).
func (this *Subscription) wants(value interface{}, key string, t l8notify.L8NotificationType) bool {
	if this.Query == nil || !this.Query.Match(value) {
		return false
	}
	if !this.Paged {
		return true
	}
	if t != l8notify.L8NotificationType_Post {
		return this.PageKeys[key]
	}
	v := orderValue(this.Query, value, key)
	// Before the page's first element: it lands on this page only if this
	// is the first page (on a later page it belongs to an earlier one).
	if this.First != nil && !this.FirstPage && this.before(v, this.First) {
		return false
	}
	// After the page's last element: it lands on this page only if the page
	// is not full.
	if this.Last != nil && this.Full && this.before(this.Last, v) {
		return false
	}
	return true
}

// before reports whether a sorts before b in the query's order. It mirrors
// internalQuery.prepare: the sort-by value honors Descending, the primary
// key fallback is always ascending.
func (this *Subscription) before(a, b interface{}) bool {
	if this.Query.SortBy() != "" && this.Query.Descending() {
		return lessThan(b, a)
	}
	return lessThan(a, b)
}

// orderValue is the value an element is ordered by in q's results: its
// sort-by value, or its primary key when q has no sort-by (or the value is
// nil), matching internalQuery.prepare.
func orderValue(q ifs.IQuery, value interface{}, key string) interface{} {
	if q.SortBy() != "" {
		if v := q.SortByValue(value); v != nil {
			return v
		}
	}
	return key
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
