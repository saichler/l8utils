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

// Page-scoped client notifications: a registered, paged view is notified only
// of writes that belong to the page it displays
// (probler/plans/websocket-page-scoped-notifications.md Phase 1).

import (
	"fmt"
	"testing"

	"github.com/saichler/l8srlz/go/serialize/object"
	"github.com/saichler/l8types/go/ifs"
	"github.com/saichler/l8types/go/testtypes"
	"github.com/saichler/l8types/go/types/l8api"
	"github.com/saichler/l8types/go/types/l8notify"
	"github.com/saichler/l8utils/go/utils/cache"
)

const pageSize = 15

// pageKey is the primary key of item i; zero-padded so key order is numeric order.
func pageKey(i int) string { return fmt.Sprintf("k-%05d", i) }

func pageItem(key string, i int) *testtypes.TestProto {
	return &testtypes.TestProto{MyString: key, MyInt32: int32(i), MyFloat64: 1}
}

// newPagedCache returns a cache holding n light items k-00000..k-<n-1>.
func newPagedCache(n int, res ifs.IResources) *cache.Cache {
	c := cache.NewCache(&testtypes.TestProto{}, nil, nil, res)
	for i := 0; i < n; i++ {
		c.Post(pageItem(pageKey(i), i), false)
	}
	return c
}

// fetchPage registers aaaId+"table" with the page starting at start, the way
// a paged table's fetch does.
func fetchPage(c *cache.Cache, aaaId, gsql string, start, block int, res ifs.IResources) []interface{} {
	elems, e := object.NewFromQuery(&l8api.L8Query{Text: gsql, AaaId: aaaId,
		Register: l8api.L8Register_L8_REGISTER_REGISTER, RegisterId: "table"}, res)
	if e != nil {
		panic(e)
	}
	q, e := elems.Query(res)
	if e != nil {
		panic(e)
	}
	values, _ := c.Fetch(start, block, q)
	return values
}

// patchKey changes item key's MyFloat64 and returns the client notification.
func patchKey(t *testing.T, c *cache.Cache, key string, v float64) *l8notify.L8NotificationSet {
	t.Helper()
	_, cn, err := c.Patch(&testtypes.TestProto{MyString: key, MyFloat64: v}, true)
	if err != nil {
		t.Fatalf("Patch %s: %v", key, err)
	}
	return cn
}

func postKey(t *testing.T, c *cache.Cache, key string, i int, v float64) *l8notify.L8NotificationSet {
	t.Helper()
	item := pageItem(key, i)
	item.MyFloat64 = v
	_, cn, err := c.Post(item, true)
	if err != nil {
		t.Fatalf("Post %s: %v", key, err)
	}
	return cn
}

func deleteKey(t *testing.T, c *cache.Cache, key string) *l8notify.L8NotificationSet {
	t.Helper()
	_, cn, err := c.Delete(&testtypes.TestProto{MyString: key}, true)
	if err != nil {
		t.Fatalf("Delete %s: %v", key, err)
	}
	return cn
}

func expectNotified(t *testing.T, what string, cn *l8notify.L8NotificationSet, aaaId string) {
	t.Helper()
	if cn == nil || !cn.AaaIds[aaaId] {
		t.Errorf("%s: expected a notification for %s, got %v", what, aaaId, cn)
	}
}

func expectSilent(t *testing.T, what string, cn *l8notify.L8NotificationSet) {
	t.Helper()
	if cn != nil {
		t.Errorf("%s: expected no notification, got one for %v", what, cn.AaaIds)
	}
}

// The probler case: 30K items, an unfiltered first page. Only changes to the
// 15 displayed items reach the subscriber.
func TestPageScopedUpdatesAtScale(t *testing.T) {
	res := newResources()
	c := newPagedCache(30000, res)
	defer c.Close()

	page := fetchPage(c, "aaa-1", "select * from TestProto", 0, pageSize, res)
	if len(page) != pageSize {
		t.Fatalf("expected a page of %d, got %d", pageSize, len(page))
	}

	off := 0
	for i := pageSize; i < 30000; i += 97 {
		if cn := patchKey(t, c, pageKey(i), float64(i)+0.5); cn != nil {
			off++
		}
	}
	if off != 0 {
		t.Errorf("expected no notifications for off-page patches, got %d", off)
	}
	expectNotified(t, "patch on-page", patchKey(t, c, pageKey(3), 42), "aaa-1")
	expectSilent(t, "replace off-page", postKey(t, c, pageKey(20000), 20000, 7))
	expectNotified(t, "replace on-page", postKey(t, c, pageKey(4), 4, 7), "aaa-1")
}

func TestPageScopedDelete(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto", 0, pageSize, res)

	expectSilent(t, "delete off-page", deleteKey(t, c, pageKey(50)))
	expectNotified(t, "delete on-page", deleteKey(t, c, pageKey(5)), "aaa-1")
}

// An add notifies only when it sorts into the displayed page's range.
func TestPageScopedAddFirstPage(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto", 0, pageSize, res)

	expectNotified(t, "add inside the page", postKey(t, c, "k-00003a", 1000, 1), "aaa-1")
	expectNotified(t, "add before the first page", postKey(t, c, "a-00000", 1001, 1), "aaa-1")
	expectSilent(t, "add after a full first page", postKey(t, c, "k-99999", 1002, 1))
}

func TestPageScopedAddLaterPage(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto", pageSize, pageSize, res) // k-00015..k-00029

	expectSilent(t, "add before page 2", postKey(t, c, "k-00003a", 1000, 1))
	expectNotified(t, "add inside page 2", postKey(t, c, "k-00020a", 1001, 1), "aaa-1")
	expectSilent(t, "add after page 2", postKey(t, c, "k-00050a", 1002, 1))
	expectSilent(t, "patch on page 1", patchKey(t, c, pageKey(3), 9))
	expectNotified(t, "patch on page 2", patchKey(t, c, pageKey(16), 9), "aaa-1")
}

// The last page is not full: an add after its last element lands on it.
func TestPageScopedAddLastPartialPage(t *testing.T) {
	res := newResources()
	c := newPagedCache(20, res)
	defer c.Close()
	page := fetchPage(c, "aaa-1", "select * from TestProto", pageSize, pageSize, res)
	if len(page) != 5 {
		t.Fatalf("expected a partial last page of 5, got %d", len(page))
	}
	expectNotified(t, "add after the partial last page", postKey(t, c, "k-99999", 1000, 1), "aaa-1")
}

// A re-fetch replaces the page snapshot.
func TestPageScopedRefetchMovesSnapshot(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto", 0, pageSize, res)
	fetchPage(c, "aaa-1", "select * from TestProto", pageSize, pageSize, res)

	if n := len(c.Subscribers()); n != 1 {
		t.Fatalf("expected the re-fetch to replace the subscription, got %d", n)
	}
	expectSilent(t, "patch on the old page", patchKey(t, c, pageKey(3), 9))
	expectNotified(t, "patch on the new page", patchKey(t, c, pageKey(16), 9), "aaa-1")
}

// The WHERE clause still applies on top of the page.
func TestPageScopedWhereStillApplies(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto where MyString='k-00003'", 0, pageSize, res)

	expectNotified(t, "patch of the matching item", patchKey(t, c, pageKey(3), 9), "aaa-1")
	expectSilent(t, "patch of a non-matching item", patchKey(t, c, pageKey(4), 9))
	expectSilent(t, "add not matching WHERE", postKey(t, c, "k-00003a", 1000, 1))
}

// Each subscriber gets only the changes on its own page.
func TestPageScopedPerSubscriber(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-p1", "select * from TestProto", 0, pageSize, res)
	fetchPage(c, "aaa-p2", "select * from TestProto", pageSize, pageSize, res)

	cn := patchKey(t, c, pageKey(16), 9)
	if cn == nil || len(cn.AaaIds) != 1 || !cn.AaaIds["aaa-p2"] {
		t.Errorf("expected only aaa-p2 to be notified, got %v", cn)
	}
}

// Descending sort-by: the page range follows the sort order, not the key.
func TestPageScopedSortByDescending(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	// MyInt32 descending: page 2 holds MyInt32 84..70.
	page := fetchPage(c, "aaa-1", "select * from TestProto sort-by MyInt32 descending", pageSize, pageSize, res)
	if first := page[0].(*testtypes.TestProto).MyInt32; first != 84 {
		t.Fatalf("expected page 2 to start at MyInt32 84, got %d", first)
	}

	expectNotified(t, "add with MyInt32 inside page 2", postKey(t, c, "new-1", 77, 1), "aaa-1")
	expectSilent(t, "add with MyInt32 above page 2 (earlier page)", postKey(t, c, "new-2", 95, 1))
	expectSilent(t, "add with MyInt32 below page 2 (later page)", postKey(t, c, "new-3", 10, 1))
	expectNotified(t, "patch on page 2", patchKey(t, c, pageKey(80), 9), "aaa-1")
	expectSilent(t, "patch on page 1", patchKey(t, c, pageKey(95), 9))
}

// An unpaged fetch (block size 0) keeps WHERE-only matching.
func TestUnpagedSubscriptionMatchesAll(t *testing.T) {
	res := newResources()
	c := newPagedCache(100, res)
	defer c.Close()
	fetchPage(c, "aaa-1", "select * from TestProto", 0, 0, res)

	expectNotified(t, "patch of any item", patchKey(t, c, pageKey(77), 9), "aaa-1")
}
