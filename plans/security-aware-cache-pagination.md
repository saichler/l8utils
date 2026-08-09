# Security-Aware Cache Pagination

## Problem 1: Inaccurate Pagination for Unprivileged Users

When a user with limited permissions fetches paginated data, the cache returns pages based on the full unfiltered dataset. Security filtering (ScopeView) happens downstream in the service manager — after pagination has already sliced. This causes:

1. Pages with fewer records than expected (or completely empty pages)
2. Total page count reflecting unfiltered data (e.g., "Page 1 of 5" when only 2 pages worth of data is visible)
3. Poor user experience — the user clicks through empty pages

## Problem 2: Ghost Missing Entry After Delete (Bug)

After deleting one entity, a subsequent `select *` query returns N-2 items instead of N-1. One unrelated entity silently disappears. The entity still exists in the cache and can be found via a filtered query, but the unfiltered list skips it.

**Root cause:** `internalQuery.prepare()` has a special case for queries with no criteria — it assigns `data = addedOrder` directly (a shared slice reference, not a copy). When `internalCache.delete()` writes a tombstone (`""`) into `addedOrder`, the tombstone is visible in the query's `data` too. The subsequent `sort.Slice` sorts the empty tombstone to the beginning, shifting indices while `key2order` retains the old indices. This causes one real entry to be unreachable.

**Reproduced in:** `l8vendingmachine` (VendRoute), `l8physio` (therapist list).

See `/home/saichler/Downloads/bug-cache-delete-ghost.md` for full reproduction steps.

## Solution

Both problems are solved together by:

1. **Removing the `addedOrder` special case in `prepare()` entirely** — all queries now iterate the cache map and apply `Match()`, eliminating the shared slice bug and giving a single code path for security filtering
2. **Including the AAA ID in the query hash** so each user gets their own `internalQuery`
3. **Applying `ScopeItem` during `prepare()`** to filter out denied records before pagination
4. **Computing metadata only over visible records**

Removing the special case also eliminates the need for `addedOrder`, `key2order`, `deleteCount`, and `cleanupOrder()` on `internalCache` — they existed solely to support the no-criteria shortcut path.

## Current Architecture

### Query hash and internalQuery caching

- `internalCache.fetch()` (`internalCache.go:174`) computes `hash := q.Hash()` and looks up `this.queries[hash]`
- `q.Hash()` (`l8ql/.../Query.go:372`) computes an MD5 of the normalized query text, then hex-encodes it to a 32-character string — **not user-specific**
- The MD5+hex approach is heavyweight for a cache lookup key: it allocates an `md5.New()` hasher, performs a cryptographic digest, then hex-encodes into a string — all unnecessary for a non-security hash map key
- Two users with the same WHERE clause share the same `internalQuery`, which holds a pre-sorted key list (`data []string`) of ALL matching records
- The key list is rebuilt only when the cache's mutation stamp changes (`dq.stamp != this.stamp`)

### internalQuery.prepare() — the two paths

- `internalQuery.prepare()` (`internalQuery.go:44`) builds the sorted key list by either:
  - **No-criteria path:** Using `addedOrder` directly (`data = addedOrder`) — shared slice reference, source of the ghost bug
  - **Criteria path:** Iterating all cache entries, calling `this.query.Match(v)` to filter, then sorting
- Metadata (counts) is computed during the criteria path only; the no-criteria path uses `globalMetadata`
- The key list is then sliced by `start` and `blockSize` in `fetch()`

### The addedOrder machinery

- `addedOrder []string` — insertion-ordered list of primary keys, maintained by `put()` and `delete()`
- `key2order map[string]int` — reverse index from key to position in `addedOrder`
- `deleteCount int` — tombstone counter for triggering `cleanupOrder()`
- `cleanupOrder()` — compacts `addedOrder` by removing tombstones when threshold is exceeded (100 tombstones or 25% of slice)

All of this exists solely to support the no-criteria shortcut. With the shortcut removed, none of it is needed.

### Security filtering today

- `ScopeView` is called in `ServiceManager.go:127` AFTER the cache returns paginated results
- `ScopeView` takes `IVNic` — not available in the cache
- New `ScopeItem` (`Security.go:55-56`) takes `IResources` (available via `Cache.r`) and filters a single item:
  ```
  ScopeItem(IResources, interface{}, string, string, ...string) interface{}
  ```
  Returns the item if allowed, nil if denied.

### AAA ID availability

- `IQuery` has `AAAId() string` and `SetAAAId(string)` (`API.go:118-121`)
- The AAA ID is set on the query before it reaches the cache
- The cache has `r ifs.IResources` which provides `r.Security()` and `r.SysConfig().LocalUuid`

## Changes

### Change 0: Replace MD5 string hash with int32 arithmetic hash (l8ql + l8types + cache)

The current `IQuery.Hash()` returns a `string` computed via MD5+hex encoding. This is unnecessarily expensive for a cache map key — MD5 is a cryptographic hash, and the hex encoding allocates a 32-character string. Replace it with a Java-style `hashCode` that returns `int32`: fast arithmetic, no allocations, good distribution for hash map keys.

**File:** `l8types/go/ifs/API.go` (IQuery interface)
**Line:** 100

Change the interface method return type:
```go
// Before
Hash() string

// After
Hash() int32
```

**File:** `l8ql/go/gsql/interpreter/Query.go`
**Method:** `Hash()` (line 372)

Replace the MD5 implementation with a Java-style arithmetic hash:
```go
func (this *Query) Hash() int32 {
    text := strings.TrimSpace(strings.ToLower(this.Text()))
    var h int32
    for _, c := range text {
        h = 31*h + int32(c)
    }
    return h
}
```

Remove the `crypto/md5` and `encoding/hex` imports from the file (if no longer used elsewhere).

**File:** `go/utils/cache/internalCache.go`

The `queries` map key type changes from `string` to `int64`. We use `int64` (not `int32`) because the cache combines the query hash with the AAA ID hash for per-user isolation — two `int32` values packed into one `int64`:

```go
// Before
queries map[string]*internalQuery

// After
queries map[int64]*internalQuery
```

In `fetch()`, the hash combination with AAA ID changes from string concatenation to arithmetic:
```go
// Before (from Change 3)
hash := q.Hash()
if aaaId != "" {
    hash = hash + ":" + aaaId
}

// After
hash := int64(q.Hash())
if aaaId != "" {
    hash = hash<<32 | int64(hashString(aaaId))
}
```

Where `hashString` is a local helper using the same Java-style algorithm:
```go
func hashString(s string) int32 {
    var h int32
    for _, c := range s {
        h = 31*h + int32(c)
    }
    return h
}
```

**File:** `go/utils/cache/internalQuery.go`

The `hash` field type changes from `string` to `int64`:
```go
type internalQuery struct {
    query    ifs.IQuery
    data     []string
    stamp    int64
    hash     int64
    metadata *l8api.L8MetaData
    lastUsed int64
}
```

In `newInternalQuery`, the assignment changes accordingly:
```go
iq.hash = int64(query.Hash())
```

**File:** `go/utils/cache/Subscriptions.go`

The `QueryHash` field on `Subscription` changes from `string` to `int32`:
```go
type Subscription struct {
    AAAId     string
    QueryHash int32
    QueryText string
    lastSeen  int64
}
```

**File:** `go/utils/cache/Cache.go`

The `RegisterSubscription` signature changes:
```go
// Before
func (this *Cache) RegisterSubscription(aaaId, queryHash, queryText string)

// After
func (this *Cache) RegisterSubscription(aaaId string, queryHash int32, queryText string)
```

**Ripple to callers:** `RegisterSubscription` is called from `l8services` (service handler after Fetch). That call site passes `q.Hash()` — since the return type changes, the call site compiles without modification. The `Subscription.QueryHash` field is only used internally for subscription tracking, not as a map key or external API.

**Note:** This change is in `l8ql` and `l8types` (sibling projects), not in `l8utils`. The user manages pushing and re-vendoring. After the changes are committed in those projects, `l8utils` re-vendors to pick them up.

**Framework interface change justification (per `framework-interface-boundaries.md`):** This modifies an existing method signature on `IQuery` in `l8types/go/ifs/API.go`. This is appropriate because: (a) the change is initiated by the framework owner, not a consuming project; (b) it improves the fundamental model — replacing a heavyweight cryptographic hash with a lightweight arithmetic hash for all cache consumers; (c) it is not driven by a single consumer's feature need but benefits the entire ecosystem.

### Change 1: Remove the addedOrder special case and simplify prepare() (internalQuery.go)

**File:** `go/utils/cache/internalQuery.go`
**Method:** `prepare()`

Remove the `addedOrder` parameter entirely. All queries now take the same path: iterate the cache map, apply `Match()`, build the key list, sort, compute metadata.

Current signature:
```go
func (this *internalQuery) prepare(cache map[string]interface{}, addedOrder []string, stamp int64, descending bool, metadataFunc map[string]func(interface{}) (bool, string))
```

New signature:
```go
func (this *internalQuery) prepare(cache map[string]interface{}, stamp int64, descending bool, metadataFunc map[string]func(interface{}) (bool, string), r ifs.IResources, aaaId string)
```

New body:
```go
func (this *internalQuery) prepare(cache map[string]interface{}, stamp int64, descending bool, metadataFunc map[string]func(interface{}) (bool, string), r ifs.IResources, aaaId string) {
    this.stamp = stamp
    this.metadata = newMetadata()

    data := make([]string, 0)
    hasScopeFilter := r != nil && r.Security() != nil && aaaId != ""
    uuid := ""
    if r != nil && r.SysConfig() != nil {
        uuid = r.SysConfig().LocalUuid
    }

    for k, v := range cache {
        if !this.query.Match(v) {
            continue
        }
        if hasScopeFilter {
            if r.Security().ScopeItem(r, v, uuid, aaaId) == nil {
                continue
            }
        }
        data = append(data, k)
        addToMetadata(v, metadataFunc, this.metadata)
    }

    sort.Slice(data, func(i, j int) bool {
        if this.query.SortBy() != "" {
            v1 := this.query.SortByValue(cache[data[i]])
            v2 := this.query.SortByValue(cache[data[j]])
            if v1 != nil && v2 != nil {
                result := lessThan(v1, v2)
                if descending {
                    return !result
                }
                return result
            }
        }
        return lessThan(data[i], data[j])
    })
    this.data = data
}
```

When `Match()` has no criteria, it returns true for all items — so the behavior is equivalent to the old `addedOrder` path, but without the shared slice bug. Sorting falls back to key-based ordering when there's no `SortBy`, which is deterministic.

### Change 2: Remove addedOrder machinery from internalCache (internalCache.go)

**File:** `go/utils/cache/internalCache.go`

Remove the following fields from `internalCache`:
- `addedOrder []string`
- `key2order map[string]int`
- `deleteCount int`

Remove the following methods/code:
- `cleanupOrder()` method (lines 146-168)
- All `addedOrder`/`key2order`/`deleteCount` maintenance in `put()` (lines 108-110) and `delete()` (lines 135-139, 142)
- `newInternalCache()` — remove `addedOrder` and `key2order` initialization (lines 43-44); change `queries` map initialization to `map[int64]*internalQuery`

Remove the `globalMetadata` field and related logic:
- `globalMetadata *l8api.L8MetaData` field
- `addToMetadata` / `removeFromMetadata` calls in `put()` and `delete()` that update `globalMetadata`
- The `globalMetadata` initialization in `newInternalCache()`

Since every `internalQuery` now computes its own metadata during `prepare()`, `globalMetadata` is no longer needed.

### Change 3: Simplify fetch() in internalCache (internalCache.go)

**File:** `go/utils/cache/internalCache.go`
**Method:** `fetch()`

Update to pass `IResources` through to `prepare()`, remove the two-path logic, and always return `dq.metadata`.

Current:
```go
func (this *internalCache) fetch(start, blockSize int, q ifs.IQuery) ([]interface{}, *l8api.L8MetaData) {
    // ...
    noCriteriaOrSort := true
    hash := q.Hash()
    dq, ok := this.queries[hash]
    if !ok {
        dq = newInternalQuery(q)
        this.queries[hash] = dq
    }
    atomic.StoreInt64(&dq.lastUsed, time.Now().Unix())
    if dq.stamp != this.stamp {
        qrt := reflect.ValueOf(q.Criteria())
        noCriteriaOrSort = (...)
        if noCriteriaOrSort {
            dq.prepare(this.cache, this.addedOrder, this.stamp, q.Descending(), this.metadataFunc)
        } else {
            dq.prepare(this.cache, nil, this.stamp, q.Descending(), this.metadataFunc)
        }
    }
    // ... slice result ...
    if !noCriteriaOrSort {
        return result, dq.metadata
    }
    return result, this.globalMetadata
}
```

New:
```go
func (this *internalCache) fetch(start, blockSize int, q ifs.IQuery, r ifs.IResources) ([]interface{}, *l8api.L8MetaData) {
    if q.IsAggregate() {
        return this.fetchAggregate(q)
    }

    aaaId := q.AAAId()
    hash := int64(q.Hash())
    if aaaId != "" {
        hash = hash<<32 | int64(hashString(aaaId))
    }

    dq, ok := this.queries[hash]
    if !ok {
        dq = newInternalQuery(q)
        this.queries[hash] = dq
    }

    atomic.StoreInt64(&dq.lastUsed, time.Now().Unix())

    if dq.stamp != this.stamp {
        dq.prepare(this.cache, this.stamp, q.Descending(), this.metadataFunc, r, aaaId)
    }

    result := make([]interface{}, 0)
    for i := start; i < len(dq.data); i++ {
        key := dq.data[i]
        value, ok := this.cache[key]
        if ok {
            result = append(result, value)
        }
        if blockSize == 0 {
            continue
        }
        if len(result) >= blockSize {
            break
        }
    }
    return result, dq.metadata
}
```

Key changes:
- AAA ID appended to hash for per-user isolation
- Single `prepare()` call — no branching on criteria presence
- Always returns `dq.metadata` — no `globalMetadata` fallback
- `IResources` passed as parameter

### Change 4: Update Fetch.go to pass IResources (Fetch.go)

**File:** `go/utils/cache/Fetch.go`
**Line 28**

Change:
```go
values, metadata := this.iCache.fetch(start, blockSize, q)
```

To:
```go
values, metadata := this.iCache.fetch(start, blockSize, q, this.r)
```

### Change 5: Update addTotalMetadata and metadata functions (internalCache.go)

The `addTotalMetadata` function in `Cache.go:118` calls into the metadata machinery. Since `globalMetadata` is removed, this function needs to be updated to only register the metadata functions (stored in `metadataFunc`) without pre-computing global counts. The metadata functions are still used by `prepare()` to compute per-query metadata.

Review `addTotalMetadata` and adjust: it should still register `metadataFunc` entries on `internalCache`, but skip the initial global count computation that populates `globalMetadata`.

### Change 6: Clean up put() and delete() (internalCache.go)

**`put()` method** — remove:
- `this.addedOrder = append(this.addedOrder, pk)` and `this.key2order[pk] = ...`
- `this.addToMetadata(value)` and `this.removeFromMetadata(pk)` calls that update `globalMetadata`

Keep:
- `this.cache[pk] = value`
- `this.putUnique(pk, uk)`
- `this.stamp = time.Now().Unix()` (still needed to trigger `prepare()` rebuild)

**`delete()` method** — remove:
- Tombstone logic (`this.addedOrder[idx] = ""`, `this.key2order`, `this.deleteCount++`)
- `this.cleanupOrder()` call
- `this.removeFromMetadata(pk)` call

Keep:
- `delete(this.cache, pk)`
- `this.deleteUnique(pk, uk)`
- `this.stamp = time.Now().Unix()`
- Return the deleted item

## Traceability Matrix

| # | Concern | Change |
|---|---------|--------|
| 1 | Lightweight query hash (replace MD5 string with int32) | Change 0: l8ql Hash() returns int32, cache uses int64 map keys |
| 2 | Ghost missing entry bug (addedOrder shared slice) | Change 1 + Change 2: Remove addedOrder entirely |
| 3 | Per-user query isolation | Change 3: AAA ID combined into int64 hash key |
| 4 | IResources flow to prepare() | Change 3 + Change 4: Pass through fetch() |
| 5 | Security filtering before pagination | Change 1: ScopeItem in unified prepare() |
| 6 | Accurate metadata per user | Change 1 + Change 5: Per-query metadata, no globalMetadata |
| 7 | Backward compatibility (no AAA ID) | Change 1 + Change 3: Empty AAA ID = no ScopeItem, shared query as before |
| 8 | Memory overhead from per-user queries | Managed by existing TTL cleanup — per-user internalQuery entries expire after 30s of inactivity |
| 9 | Dead code removal (addedOrder machinery) | Change 2 + Change 6: Remove addedOrder, key2order, deleteCount, cleanupOrder, globalMetadata |

## Files Modified

| File | Project | What Changes |
|------|---------|-------------|
| `go/ifs/API.go` | l8types | `IQuery.Hash()` return type changes from `string` to `int32` |
| `go/gsql/interpreter/Query.go` | l8ql | `Hash()` implementation: MD5+hex replaced with Java-style int32 arithmetic hash; remove `crypto/md5` and `encoding/hex` imports |
| `go/utils/cache/internalQuery.go` | l8utils | `hash` field type changes to `int64`; `prepare()` — remove `addedOrder` parameter, single code path with Match + ScopeItem filtering, always compute per-query metadata |
| `go/utils/cache/internalCache.go` | l8utils | `queries` map key changes from `string` to `int64`; add `hashString()` helper; remove `addedOrder`, `key2order`, `deleteCount`, `globalMetadata` fields; remove `cleanupOrder()`; simplify `put()` and `delete()`; update `fetch()` signature and logic |
| `go/utils/cache/Fetch.go` | l8utils | Pass `this.r` to `iCache.fetch()` |
| `go/utils/cache/Subscriptions.go` | l8utils | `QueryHash` field type changes from `string` to `int32` |
| `go/utils/cache/Cache.go` | l8utils | `RegisterSubscription` parameter type changes for `queryHash` |

## Backward Compatibility

- **Hash type change:** `IQuery.Hash()` changes from `string` to `int32`. All callers in l8utils (cache) and l8services (subscription registration) pass the return value directly — no string formatting or parsing. The change is source-compatible at all call sites. After re-vendoring l8types and l8ql, all consuming projects recompile without modification.
- When `q.AAAId()` is empty (no authenticated user, system queries, tests), the hash is unchanged and no ScopeItem filtering is applied. Behavior is identical to today minus the ghost bug.
- When `r.Security()` is nil (ShallowSecurityProvider scenarios), no filtering is applied.
- ShallowSecurityProvider's `ScopeItem` should return the item as-is (permissive), matching its existing `ScopeView` behavior.
- Sorting for no-criteria queries changes from insertion order to key-based order. This is acceptable — users always see sorted or paginated data through the UI, and insertion order was never guaranteed to be meaningful.

## Impact on Downstream (l8services)

After this change, the cache returns already-filtered, correctly-paginated data. The existing `ScopeView` call in `ServiceManager.go:127` becomes redundant for GET requests served from cache. It should remain in place as a safety net (ScopeView on already-filtered data is a no-op), but could be removed in a future cleanup.

## Testing

All test files live in `go/tests/` per `test-location-and-approach.md`. Tests exercise the cache through its public API (`NewCache`, `Fetch`, `Post`, `Delete`, etc.), not internal functions.

1. **Test:** Cache with a mock security provider that denies specific items. Verify that Fetch returns only allowed items, and metadata reflects the filtered count.
2. **Test:** Same query from two different AAA IDs with different permissions. Verify each user sees their own filtered result set with correct pagination.
3. **Test:** Empty AAA ID. Verify behavior is unchanged from current (no filtering, shared internalQuery).
4. **Test:** TTL cleanup. Verify per-user internalQuery entries are cleaned up after expiry.
5. **Test:** Delete an item from a 5-item cache, then fetch with `select *`. Verify exactly 4 items returned (ghost bug regression test).
6. **Test:** Delete an item, then fetch with criteria that would have matched the deleted item. Verify it is not returned.
7. **Test (l8ql, in `../l8ql/go/tests/`):** Verify `Hash()` returns consistent `int32` values for the same query text, and different values for different queries. Verify case-insensitivity (same hash for `SELECT * FROM Foo` and `select * from foo`).
8. **Existing tests:** All existing cache tests must pass — the behavior change is limited to removing insertion-order preservation for unsorted queries.

## Final Verification Phase

After all changes are implemented across l8types, l8ql, and l8utils:

1. Re-vendor l8types and l8ql into l8utils
2. `cd go && go build ./...` — verify all packages compile
3. `cd go && go vet ./...` — verify no static analysis issues
4. Run full test suite: `cd go && bash test.sh` — all existing tests pass
5. Run the new security-aware pagination tests (tests 1-6 above) — all pass
6. Verify no file in the cache package exceeds 500 lines after the changes (per `maintainability.md`)
7. Re-vendor l8types into l8services — verify `go build ./...` compiles (RegisterSubscription call site, ScopeView call site)
8. In a project with security (e.g., l8physio or l8erp), re-vendor all updated deps, start the system locally, and verify:
   - [ ] Admin user sees all records with correct pagination
   - [ ] Restricted user sees only permitted records with correct page count
   - [ ] Deleting a record shows N-1 records on next fetch (ghost bug regression)
   - [ ] Metadata total count matches visible record count for both users

## Implementation Order

1. **l8types** — Change `IQuery.Hash()` return type from `string` to `int32` in the interface
2. **l8ql** — Replace MD5 implementation with int32 arithmetic hash
3. **l8utils** — Re-vendor l8types and l8ql, then implement Changes 1-6 in the cache package
4. **l8services** — Re-vendor l8types (for the interface change); call sites compile without modification
