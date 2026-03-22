# Plan: Cache In-Memory Aggregate Support (with GROUP BY and HAVING)

## Goal
When `Fetch` receives a query with aggregates (Count, Sum, Avg, Min, Max), compute the results in-memory over cached objects — including GROUP BY grouping and HAVING filtering — and return them in `L8MetaData.KeyCount.Counts`. Return an empty `[]interface{}` slice — no raw records.

## Key Discovery
The l8ql `Query` struct already implements the full aggregate pipeline:
- `IsAggregate()` — checks if query has aggregate functions
- `Filter(items, onlySelectedColumns)` — filters by WHERE criteria
- `Aggregate(items)` — groups by GROUP BY fields, computes all aggregate functions per group, returns `[]map[string]interface{}`

These methods are on the concrete `*Query` struct but **not exposed** via the `IQuery` interface. The cache only sees `IQuery`.

## Changes

### 1. l8types: Add 3 methods to `IQuery` interface

**File:** `l8types/go/ifs/API.go`

```go
// IsAggregate returns true if this query has aggregate functions.
IsAggregate() bool
// Filter returns the subset of items matching the WHERE clause.
Filter([]interface{}, bool) []interface{}
// Aggregate groups items by GROUP BY fields and computes aggregate functions.
// Returns one map per group with aggregate alias keys and group-by field values.
Aggregate([]interface{}) []map[string]interface{}
```

### 2. l8utils cache: New file `Aggregate.go`

**File:** `l8utils/go/utils/cache/Aggregate.go`

Handles the aggregate computation path. Called from `internalCache.fetch` when the query is an aggregate query.

```go
func (this *internalCache) fetchAggregate(q ifs.IQuery) ([]interface{}, *l8api.L8MetaData) {
    // 1. Collect all cached objects
    items := make([]interface{}, 0, len(this.cache))
    for _, v := range this.cache {
        items = append(items, v)
    }

    // 2. Filter by WHERE criteria
    filtered := q.Filter(items, false)

    // 3. Compute aggregates (handles GROUP BY internally)
    groups := q.Aggregate(filtered)

    // 4. Filter by HAVING
    groups = filterByHaving(groups, q.Having())

    // 5. Pack results into metadata
    metadata := newMetadata()
    packAggregateResults(groups, q.Aggregates(), q.GroupBy(), metadata)

    return []interface{}{}, metadata
}
```

### 3. l8utils cache: New file `GroupBy.go`

**File:** `l8utils/go/utils/cache/GroupBy.go`

Handles packing grouped aggregate results into `L8MetaData.KeyCount.Counts`.

**Encoding scheme for Counts map:**
- **No GROUP BY** (single group): key = aggregate alias (e.g., `"count"`, `"sumMyInt32"`)
- **With GROUP BY** (multiple groups): key = `"alias:groupFieldValue"` (e.g., `"count:GroupA"`, `"sumMyInt32:GroupA"`)
  - For multi-field GROUP BY: `"alias:val1|val2"`

```go
func packAggregateResults(groups []map[string]interface{}, aggregates []*l8api.L8AggregateFunction, groupByFields []string, metadata *l8api.L8MetaData) {
    if len(groupByFields) == 0 {
        // Single group — flat keys
        if len(groups) == 1 {
            for _, agg := range aggregates {
                if val, ok := toFloat64(groups[0][agg.Alias]); ok {
                    metadata.KeyCount.Counts[agg.Alias] = val
                }
            }
        }
        return
    }

    // Multiple groups — prefixed keys
    for _, group := range groups {
        groupKey := buildGroupKeyString(group, groupByFields)
        for _, agg := range aggregates {
            key := agg.Alias + ":" + groupKey
            if val, ok := toFloat64(group[agg.Alias]); ok {
                metadata.KeyCount.Counts[key] = val
            }
        }
    }
}

func toFloat64(v interface{}) (float64, bool) {
    switch n := v.(type) {
    case int:
        return float64(n), true
    case int32:
        return float64(n), true
    case int64:
        return float64(n), true
    case float64:
        return n, true
    case float32:
        return float64(n), true
    case uint32:
        return float64(n), true
    case uint64:
        return float64(n), true
    }
    // Handle protobuf enum types via reflection
    rv := reflect.ValueOf(v)
    switch rv.Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return float64(rv.Int()), true
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return float64(rv.Uint()), true
    case reflect.Float32, reflect.Float64:
        return rv.Float(), true
    }
    return 0, false
}

func buildGroupKeyString(group map[string]interface{}, groupByFields []string) string {
    // Concatenate group-by field values with "|" separator
    // e.g., group["myString"] = "GroupA" → "GroupA"
    // e.g., group["dept"] = "Sales", group["region"] = "West" → "Sales|West"
}
```

### 4. l8utils cache: New file `Having.go`

**File:** `l8utils/go/utils/cache/Having.go`

Filters grouped aggregate results by the HAVING clause.

The HAVING expression's comparators reference aggregate aliases as their left operand. Since aggregate results are `map[string]interface{}` (not protobuf structs), the existing `Expression.Match()` which uses `Property.Get()` cannot be reused. Instead, implement a simple map-based expression evaluator:

```go
func filterByHaving(groups []map[string]interface{}, having ifs.IExpression) []map[string]interface{} {
    if having == nil {
        return groups
    }
    result := make([]map[string]interface{}, 0)
    for _, group := range groups {
        if matchHaving(group, having) {
            result = append(result, group)
        }
    }
    return result
}

func matchHaving(group map[string]interface{}, expr ifs.IExpression) bool {
    // Evaluate condition: look up comparator.Left() in group map, compare to comparator.Right()
    // Support operators: =, !=, >, <, >=, <=
    // Chain with AND/OR via expr.Operator() and expr.Next()
    // Recurse into expr.Child() for nested expressions
}

func compareValues(left interface{}, right string, operator string) bool {
    // Convert left (from map) and right (literal string) to float64, then compare
}
```

### 5. l8utils cache: Wire into `internalCache.fetch`

**File:** `l8utils/go/utils/cache/internalCache.go`

At the top of `fetch()`, add:

```go
if q.IsAggregate() {
    return this.fetchAggregate(q)
}
```

### 6. l8utils cache: Update `Fetch.go` (public method)

**File:** `l8utils/go/utils/cache/Fetch.go`

When `q.IsAggregate()`, the `[]interface{}` from `iCache.fetch` will be empty — skip cloning loop and just clone+return the metadata.

## Full Aggregate Pipeline in Cache

```
Fetch(start, blockSize, query)
  │
  ├── query.IsAggregate()? ──YES──► fetchAggregate(query)
  │                                    │
  │                                    ├── 1. Collect all cached objects
  │                                    ├── 2. query.Filter(items) — apply WHERE
  │                                    ├── 3. query.Aggregate(filtered) — GROUP BY + compute aggs
  │                                    ├── 4. filterByHaving(groups, query.Having()) — apply HAVING
  │                                    ├── 5. packAggregateResults(groups) — write to metadata.Counts
  │                                    └── 6. return []interface{}{}, metadata
  │
  └── NO ──► existing fetch path (pagination, sorting, etc.)
```

## Files Changed (summary)

| # | Repo | File | Change |
|---|------|------|--------|
| 1 | l8types | `go/ifs/API.go` | Add `IsAggregate()`, `Filter()`, `Aggregate()` to `IQuery` |
| 2 | l8utils | `go/utils/cache/Aggregate.go` | **New** — `fetchAggregate()`, orchestrates the pipeline |
| 3 | l8utils | `go/utils/cache/GroupBy.go` | **New** — `packAggregateResults()`, `toFloat64()`, group key encoding |
| 4 | l8utils | `go/utils/cache/Having.go` | **New** — `filterByHaving()`, `matchHaving()`, map-based expression eval |
| 5 | l8utils | `go/utils/cache/internalCache.go` | Add `IsAggregate()` check at top of `fetch()` |
| 6 | l8utils | `go/utils/cache/Fetch.go` | Skip clone loop for aggregate results |

## Notes
- No changes needed in l8ql — `Query` already implements `IsAggregate()`, `Filter()`, and `Aggregate()`, they just need to be added to the interface
- The `Aggregate()` method in l8ql already handles GROUP BY grouping and uses its internal `Accumulator` for count/sum/avg/min/max
- HAVING is evaluated on the cache side against `[]map[string]interface{}` because the existing `Expression.Match()` uses `Property.Get()` which expects struct objects, not maps
