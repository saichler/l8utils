# Generic WebSocket Change Notifications From l8utils/cache

## Revision note

This plan was substantially redesigned after an earlier draft (Phases that added a
`vnic` field to `Cache` and changed `NewCache`'s signature) turned out to have far
larger blast radius than expected — every direct/indirect caller of `Cache`/`DCache`
across `l8services`, `l8orm`, `l8inventory`, and even unrelated projects (`l8vibe`,
`l8pollaris`) construct one, so changing that constructor rippled everywhere. That
work was started and then fully reverted (`git checkout --` on every touched file)
before landing anywhere, at the user's explicit instruction, once a much smaller,
correct design was found by re-reading how notifications are actually sent today.
This document reflects that smaller design; nothing below assumes the reverted
`vnic`-on-`Cache` approach.

## How this actually works today (verified against real code, not assumed)

`Cache` never sends anything itself — `Cache.Post/Put/Patch/Delete` only *build and
return* notification objects: `(n *l8notify.L8NotificationSet, cn *l8notify.L8NotificationSet, err error)`.

- **`n`** — the delta notification, for cross-node replication sync.
- **`cn`** — the "client notification" (`createClientNotification`/
  `createClientNotificationForPatch`, `notifications.go`), already the intended
  browser-facing half: it clones `n`'s fields (including `NotificationList`, the
  actual changed record) and attaches `AaaIds` from `this.subscriberAaaIds()`.

The **caller** — which already holds its own `vnic` on every write, via the
`IServiceHandler` method signature itself (`Post(pb ifs.IElements, vnic ifs.IVNic)`
etc.) — is responsible for sending `cn`. `base.BaseService.do()`
(`l8services/go/services/base/BaseServiceNotifications.go:82-87`) already does
exactly this:
```go
if cn != nil && this.vnic != nil {
    this.vnic.Multicast(WsServiceName, WsServiceArea, ifs.Action(cn.Type), cn)
}
```
This is real, wired code — it just never fires, because `cn` is always `nil` today,
for two independent reasons:

1. **Nothing calls `RegisterSubscription`** anywhere in the codebase except its own
   unit test (verified: `grep -rn "RegisterSubscription("` across every repo). So
   `this.subs` (the registry `HasSubscribers()`/`Subscribers()` read) is always
   empty, and `createClientNotification` hits its first line —
   `if delta == nil || !this.HasSubscribers() { return nil }` — every time.
2. **Even if it weren't empty, `createClientNotification`/`createClientNotificationForPatch`
   don't filter by query match at all** — `subscriberAaaIds()` just returns every
   currently-registered AAAId unconditionally, regardless of whether the changed
   record actually falls inside that subscriber's query.

Separately, `Cache.Fetch()` (`Fetch.go`) never reads `q.Register()` at all — the
`register=true` keyword `l8ql` parses is inspected nowhere downstream. And
`persist.OrmService` doesn't even reach this far: `OrmCache.go:45,55,66` hardcode
`this.cache.Post/Patch/Delete(element, false)` — since every function above returns
early (`if !createNotification { return n, nil, e }`) *before* computing `cn` at
all when the flag is `false`, ORM-backed writes never even build a `cn`, let alone
send one. `OrmService`'s own `cachePost`/`cachePatch`/`cacheDelete`
(`OrmCache.go:38-66`) discard whatever `Cache` returns entirely (`this.cache.Post(element, false)`,
no assignment; `cacheDelete` keeps only the error).

## Design

Four small, independent fixes — Phases 2 and 3 are inside `l8utils/cache`; Phases 1
and 1b are in the two real callers that currently discard `cn` entirely
(`l8orm`, `l8services/dcache`):

### Phase 1: `OrmService` sends `cn`, same pattern `base.BaseService` already uses

- `OrmCache.go:38-66` (`cachePost`/`cachePatch`/`cacheDelete`): add a `vnic ifs.IVNic`
  parameter to each (their caller, `cacheAction` in `OrmDoAction.go:57`, already
  receives `vnic` as a parameter — this is pure threading, one hop). Change
  `this.cache.Post(element, false)` → `true` (same for `Patch`/`Delete`), capture
  the returned `cn`, and if `cn != nil && vnic != nil`, call
  `vnic.Multicast(wsServiceName, wsServiceArea, ifs.Action(cn.Type), cn)`.
- Define local `wsServiceName = "websock"` / `wsServiceArea = byte(0)` constants in
  `l8orm/orm/persist` — same pattern `base`/`l8inventory` already use (each
  redefines its own copy rather than importing `l8web`), so `l8orm` gains no new
  dependency.
- Update `cachePost`/`cachePatch`/`cacheDelete`'s two-to-three call sites
  (`OrmDoAction.go`'s `cacheAction`, `OrmService.go:139,171`'s `Delete`) to pass
  `vnic` through.
- `n` (the other return value) is simply discarded, same as today — `OrmService`
  has no cross-node replication queue (no `nQueue` field, unlike `BaseService`), so
  there's nothing to do with it.

This alone does **not** yet produce any live notifications — `cn` will still always
be `nil`, because of Phases 2 and 3 below. It's a prerequisite, not the whole fix.

### Phase 1b: `DCache` sends `cn` too — literally has a `vnic` already, just typed narrower

`DCache` has no field declared as `ifs.IVNic` — only `listener ifs.IServiceCacheListener`
(`DCache.go:31`). But `IVNic` already satisfies `IServiceCacheListener`
(`VirtualNetworkInterface.PropertyChangeNotification`, `l8bus/go/overlay/vnic/Notifications.go:50-53`,
multicasts `n` back to the entity's own `ServiceName`/`ServiceArea` for cross-node
sync — a real, working, already-wired path, just unrelated to `"websock"`). Three of
the four known real callers pass their own `vnic` directly as that `listener`
argument (`ReplicationService.go:46-47`, `ProjectService.go:60`,
`PollarisCenter.go:72-73`); `l8inventory`'s `InventoryCenter.go:92-93` passes `nil`.

So the concrete object behind `this.listener` already *is* a `vnic`, for every
caller except `l8inventory` — recoverable with a type assertion, no new field, no
constructor change, no new caller updates:

- `DCachePost.go:23-28`, `DCachePatch.go`, `DCacheDelete.go` (`DCachePut.go` just
  delegates to `Post`, needs no separate change): capture `cn`, the second return
  value currently discarded (`n, _, e := this.cache.Post(v, createNotification)`),
  and add:
  ```go
  if v, ok := this.listener.(ifs.IVNic); ok && cn != nil {
      v.Multicast(wsServiceName, wsServiceArea, ifs.Action(cn.Type), cn)
  }
  ```
  For `l8inventory` (`listener == nil`), the type assertion fails harmlessly (`ok == false`)
  — no behavior change there, consistent with leaving `l8inventory`/`probler`
  untouched this pass.
- Same `wsServiceName`/`wsServiceArea` local constants as Phase 1, defined once in
  `l8services/go/services/dcache` (own local copy, same pattern as everywhere else).
- Also gated by Phase 1's requirement: `this.cache.Post/Patch/Delete(v, createNotification)`
  already passes `true` for `createNotification` on locally-sourced writes (see
  `DCachePost.go:24`, `!(sourceNotification != nil && ...)`) — no change needed
  there, unlike `OrmCache.go`. `cn` still won't be non-nil until Phases 2/3 land,
  same as Phase 1.

This makes `ReplicationService`, `l8vibe`'s `ProjectService`, and `l8pollaris`'s
`PollarisCenter` benefit automatically too, the moment a client fetches from any of
them with `register=true` — no changes needed in those three projects at all.

### Phase 2: `createClientNotification`/`createClientNotificationForPatch` actually match

- `Subscriptions.go`: change `Subscription.QueryText string` (and `QueryHash`) to
  `Subscription.Query ifs.IQuery` — store the live, already-parsed query object
  directly (Phase 3's caller already has it in hand at registration time; no need
  to re-parse text later, and no need to cross-reference into a separate map by
  hash).
- `notifications.go`: replace `subscriberAaaIds()` (all current subscribers,
  unconditionally) with a version that walks `this.subs.subscribers()` and calls
  `sub.Query.Match(newValue)` per subscriber, including only the AaaIds whose query
  actually matches the changed record. `createClientNotification`/
  `createClientNotificationForPatch` need the new/updated value passed in for this
  (today they only receive `delta`/`item`, which already carry what's needed —
  `delta.NotificationList` for the Post/Put/Delete path, `item` directly for the
  Patch path — no new parameter should be required, just reading what's already
  there instead of ignoring it in favor of `subscriberAaaIds()`).

Once Phase 1 and 2 both land, add/update/delete for ORM-backed services all flow
through the **same, already-existing** `cn`/`Multicast` mechanism uniformly — no
separate "unconditional broadcast" path is needed for add/delete the way an earlier
draft of this plan proposed. `createAddNotification`/`createReplaceNotification`/
`createDeleteNotification`/`createUpdateNotification` already embed the item in
`NotificationList`; `Match()` against an add's new value or a delete's old value
both work the same way as an update's new value — no special-casing required. This
also means `base.BaseService`'s own already-correct `Multicast(cn)` call
(`BaseServiceNotifications.go:85-87`) needs **no change or removal** — it starts
working correctly the moment `cn` stops being permanently nil/unfiltered, for any
`BaseService` too. (Confirmed no duplicate-send risk today: the old call is gated
on `cn != nil`, which stays false until this plan lands regardless.)

### Phase 3: Register on Fetch

- Make `RegisterSubscription` private (`registerSubscription`), and change its
  signature to take the whole `ifs.IQuery` rather than decomposed
  `(aaaId, queryHash int32, queryText string)` — feeds Phase 2's
  `Subscription.Query` field directly: `registerSubscription(aaaId string, q ifs.IQuery)`.
- `Cache.Fetch()` (`Fetch.go`) already holds `this.mtx.Lock()` for its whole body
  and already has `this.subs` as a field — no signature change to
  `internalCache.fetch()` needed at all. Add, right at the top:
  ```go
  if q.Register() && q.AAAId() != "" {
      this.registerSubscription(q.AAAId(), q)
  }
  ```
  Called **unconditionally on every matching `Fetch()` call**, not just when
  `internalCache.fetch()` creates a new `dq` entry — `subscriptions.register()`
  already refreshes `sub.lastSeen` every time, so a client that keeps re-issuing
  the same registered query (e.g. any kind of poll) naturally stays alive against
  TTL eviction without needing a separate "is this new" check.

## Known limitations, documented rather than fixed this pass

1. **`this.subs` is keyed by AAAId alone** (`Subscriptions.go`'s
   `subs map[string]*Subscription`) — a single Cache instance can hold only one
   active registered subscription per AAAId at a time. A client with two
   concurrently open registered queries against the *same model type* (e.g. two
   differently-filtered views of the same list open in two tabs under the same
   session) would have the second overwrite the first. Not fixed here — flagging
   as a real constraint, not silently working around it.
2. **No disconnect-triggered `UnregisterSubscription` call exists anywhere.** A
   client that registers once (one `Fetch()` with `register=true`) and then only
   listens via websocket — never re-fetching that query — relies purely on TTL
   eviction with no proactive cleanup on socket close. Whether that's acceptable
   depends on `l8web`'s `WebSocketManager` connection lifecycle, which this plan
   doesn't touch or resolve.
3. **`l8inventory`/`probler` are untouched**, per explicit instruction.
   `InventoryService`'s own hand-rolled `notifyWs` (`InventoryService.go:100-152`)
   is unaffected by Phase 1b: `InventoryCenter` passes `nil` as `DCache`'s
   `listener`, so the new `this.listener.(ifs.IVNic)` type assertion there always
   fails and `notifyWs` remains its only mechanism, unchanged. `ReplicationService`/
   `ProjectService`/`PollarisCenter`, which do pass their `vnic` as `listener`,
   benefit from Phase 1b automatically.

## Traceability Matrix

| # | Gap | Phase |
|---|---|---|
| 1 | `OrmCache.go` hardcodes `createNotification=false` — `cn` never even computed for ORM writes | Phase 1 |
| 2 | `OrmService`'s cache helpers never send `cn` even if computed | Phase 1 |
| 3 | `DCache` discards `cn` entirely in `Post`/`Patch`/`Delete` | Phase 1b |
| 4 | `createClientNotification`/`createClientNotificationForPatch` don't filter by `Match()` — blind broadcast to all current subscribers | Phase 2 |
| 5 | `Subscription` stores query text/hash, not a `Match()`-capable `ifs.IQuery` | Phase 2 |
| 6 | Nothing calls `RegisterSubscription` — `this.subs` always empty | Phase 3 |
| 7 | `Cache.Fetch()` never reads `q.Register()` | Phase 3 |
| 8 | One subscription per AAAId per Cache (no composite key) | Documented limitation, not fixed |
| 9 | No disconnect-triggered unregister | Documented limitation, not fixed |
| 10 | `l8inventory`'s `notifyWs` still hand-rolled, unaffected by Phase 1b (`nil` listener) | Out of scope (probler untouched) |

## Phase 5: End-to-End Verification

1. Extend `l8utils/go/tests/CacheSubscriptions_test.go` (already exercises
   `Subscription`/registration) to cover the real flow: `Fetch()` with
   `Register()==true` on a fake `IQuery`, then a `Post`/`Put`/`Patch`/`Delete` whose
   new/old value does/doesn't match that query, asserting `cn`'s `AaaIds` and
   `NotificationList` are correct in each case (matched, not-matched, no
   subscribers, `AAAId()==""`).
2. `go build ./...` / `go vet ./...` / `gofmt` clean in `l8utils`, `l8orm`, and
   `l8services` (the three repos touched by Phases 1/1b/2/3).
3. `probler`/`l8inventory`: explicitly out of scope, no verification here (see
   Known limitation 3). `l8vibe`/`l8pollaris` are touched by nothing (Phase 1b
   needs no changes in either — see its own text) but are worth a build check
   since they depend on `l8services/dcache`'s public API shape, which Phase 1b
   does not change.
4. The real end-to-end, live proof is downstream, in `l8secure-scan/plans/scanjob-live-progress.md`
   (its own Phase 4 depends on this plan), against `l8secure-scan`'s own
   `k8s/kind-start.sh` cluster — the first real `OrmService`-backed consumer of
   this mechanism.
