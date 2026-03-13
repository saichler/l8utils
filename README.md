# Layer 8 Utils

[![Go Version](https://img.shields.io/badge/Go-1.25.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)]()

Common shared utilities and building blocks for Layer 8 microservices applications.

## Overview

Layer 8 Utils provides a comprehensive collection of utilities, interfaces, and default implementations designed to address common requirements in microservices architecture. The library emphasizes modularity, allowing projects to use default implementations or provide custom implementations that adhere to defined interfaces.

## Recent Updates

### Latest Changes (March 2026)
- **TSDB Notifications**: Added TSDB (Time Series Database) notification support for time-series data change tracking
- **CPU Profiling**: Added CPU usage analysis and memory dump capabilities via pprof integration
- **Cache Optimization**: Moved Patch operations to use dry run for improved performance
- **Windows Compatibility**: Fixed platform compatibility issues

### February 2026
- **Aggregator**: Completed data aggregation utilities with full test coverage
- **Logging Improvements**: Fixed concurrent panic in logger, corrupted log argument handling, configurable log directory
- **Shared Resources**: Renamed and consolidated shared resource utilities, uses shallow security as default
- **Escaping Fixes**: Fixed string escaping issues

### January 2026
- **Memory Leak Fixes**: Fixed memory leaks in cache and resource management
- **pprof Integration**: Added pprof heap dump support for memory profiling
- **Host Name Support**: Added hostname resolution utilities

### December 2025
- **Query TTL Support**: Added automatic TTL-based cleanup for cached queries with configurable expiration (30s default)
- **Query Sorting Fix**: Fixed sorting behavior in cache queries for consistent ordering
- **Web Service Refactoring**: Improved web service architecture with better code organization

## Features

### Core Utilities

- **Cache**: High-performance in-memory cache with storage integration
  - CRUD operations (Post, Get, Put, Patch, Delete)
  - Storage layer integration with persistence support
  - Built-in notification system for change tracking
  - Statistics tracking with named stat functions and automatic totals
  - Collection operations with `Collect()` for data aggregation
  - Clone-based isolation for concurrent access
  - Query support with pagination, filtering, and sorting
  - Query TTL with automatic cleanup (configurable, 30s default)
  - Dry-run Patch support for optimized updates

- **Notifications**: Distributed state change notification system
  - Support for Add, Delete, Update, Replace, and Sync notification types
  - TSDB notification support for time-series data
  - Serialization/deserialization with protocol buffers
  - Change tracking with property-level granularity
  - Sequence numbering for ordering guarantees
  - Service-area based routing support

- **Queues**: High-performance thread-safe queues
  - `ByteQueue`: Optimized byte queue with 8 priority levels and O(1) bit operations
  - `Queue`: Generic queue implementation for any type
  - Support for concurrent operations, blocking/non-blocking dequeue

- **Logging**: Flexible async logging framework
  - File-based logging with configurable directory
  - Console/fmt logging
  - Configurable log levels (Trace, Debug, Info, Warning, Error)
  - Asynchronous queue-based processing (50k entry limit)
  - pprof integration for CPU and heap profiling
  - Direct logger implementation for synchronous use
  - Platform-specific support (Unix and Windows)

- **String Utilities**: Comprehensive string manipulation and conversion
  - Type-safe string to/from primitive conversions
  - String formatting and parsing utilities
  - Escape/unescape handling

### Security & Infrastructure

- **Certificate Management**: TLS/SSL certificate utilities with self-signed certificate support
- **Shallow Security**: Default security provider with token validation
- **Maps**: Thread-safe `SyncMap` with reflection-based type-safe value/key lists

### Web Services

- **Web Service Framework**: RESTful service utilities
  - HTTP method handlers (GET, POST, PUT, PATCH, DELETE)
  - Request/response marshaling with protocol buffer integration
  - VNet (Virtual Network) support for distributed networking

### System Management

- **Registry**: Type registration and management
  - Dynamic instance creation with `NewOf()`
  - Enum registration and management
  - Layer 8 core types pre-registered
  - Thread-safe operations
- **Resources**: Centralized dependency injection container (logger, registry, security, serializers, config, services)
- **Workers**: Worker pool with configurable concurrency limits and condition-variable coordination
- **Aggregator**: Data aggregation utilities for collection processing
- **Tasks**: Task queue management
- **IP Segment**: IP segment parsing and management

## Installation

```bash
go get github.com/saichler/l8utils/go
```

## Quick Start

```go
package main

import (
    "github.com/saichler/l8utils/go/utils/shared"
    "github.com/saichler/l8utils/go/utils/queues"
    "github.com/saichler/l8types/go/ifs"
)

func main() {
    // Create shared resources
    resources := shared.NewResources("my-service", 8080, 30)

    // Get logger from resources
    log := resources.Logger()
    log.Info("Application started")

    // Create a high-performance byte queue with 8 priority levels
    queue := queues.NewByteQueue("main-queue", 10000)
    queue.Add([]byte("Hello World"), ifs.PRIORITY_MEDIUM)

    // Process queue items
    if data := queue.Poll(); data != nil {
        log.Info("Processed:", string(data))
    }
}
```

## Project Structure

```
go/
├── utils/
│   ├── aggregator/        # Data aggregation utilities
│   ├── cache/             # High-performance cache with storage integration
│   ├── certs/             # TLS/SSL certificate management
│   ├── ipsegment/         # IP segment parsing and management
│   ├── logger/            # Async logging framework with pprof support
│   ├── maps/              # Thread-safe SyncMap implementation
│   ├── notify/            # Notification system (including TSDB support)
│   ├── queues/            # Priority queue implementations
│   ├── registry/          # Type registration and dynamic instantiation
│   ├── requests/          # HTTP request utilities
│   ├── resources/         # Centralized resource container
│   ├── shallow_security/  # Default security provider
│   ├── shared/            # Pre-configured Resources factory
│   ├── strings/           # String conversion and escaping
│   ├── tasks/             # Task queue management
│   ├── web/               # RESTful web service framework
│   └── workers/           # Worker pool implementations
├── tests/                 # All test files (25 files)
└── vendor/                # Vendored dependencies
```

## Key Components

### Cache
High-performance in-memory cache with optional storage persistence and query TTL:

```go
// Create cache with storage backend
cache := cache.NewCache(&MyModel{}, initElements, storage, resources)

// CRUD operations
cache.Post(item, true)  // Add with notification
item := cache.Get(key)
cache.Put(key, updatedItem, true)
cache.Patch(key, changes, true)  // Uses dry run for optimization
cache.Delete(key, true)

// Query with pagination (queries cached with 30s TTL by default)
results := cache.Fetch(0, 25, query)

// Query cache management
queryCount := cache.QueryCount()
cache.CleanupQueriesNow(60)       // Manual cleanup with custom TTL
defer cache.Close()                // Stop TTL cleaner on shutdown

// Statistics tracking
cache.AddStatFunc("active", func(item interface{}) bool {
    return item.(*MyModel).Status == "active"
})
stats := cache.Stats() // Returns map with counts for "Total", "active", etc.

// Collection operations
collection := cache.Collect(predicate)
```

### ByteQueue
High-performance, thread-safe queue with 8 priority levels and O(1) bit operations:

```go
queue := queues.NewByteQueue("processor", 5000)
queue.Add(data, ifs.PRIORITY_HIGH)
result := queue.Poll() // Non-blocking
result := queue.Next() // Blocking until data available
```

### Logger
Asynchronous logging with multiple output methods and pprof integration:

```go
// Console logging
log := logger.NewLoggerImpl(logger.NewFmtLogMethod())

// File logging with configurable directory
fileLog := logger.NewFileLogMethod("app.log")
log := logger.NewLoggerImpl(fileLog)

// Combined logging
log := logger.NewLoggerImpl(
    logger.NewFileLogMethod("app.log"),
    logger.NewFmtLogMethod(),
)
log.SetLogLevel(ifs.Error_Level)
log.Info("Application ready")
```

### Web Services
RESTful service utilities with protocol buffer and VNet support:

```go
service := web.NewWebService("user-service", serviceArea)
service.SetVnet(8080)

service.HandlePost(userCreateHandler)
service.HandleGet(userGetHandler)
service.HandlePut(userUpdateHandler)
service.HandleDelete(userDeleteHandler)
```

### Notifications
Distributed state change notifications with TSDB support:

```go
// Create Add notification
notSet, err := notify.CreateAddNotification(
    model, "service-name", "model-key",
    serviceArea, "ModelType", "source", 1, sequence,
)

// Create Update notification from changes
notSet, err := notify.CreateUpdateNotification(
    changes, "service-name", "model-key",
    serviceArea, "ModelType", "source", len(changes), sequence,
)

// Extract item from notification
item, err := notify.ItemOf(notSet, resources)
```

### Registry
Type registration and dynamic instance creation:

```go
registry := registry.NewRegistry()
registry.Register(&MyModel{})
newInstance := registry.NewOf(&MyModel{})
info, err := registry.Info("MyModel")
registry.RegisterEnum("Status", []string{"active", "pending", "completed"})
```

### Resources
Centralized dependency injection container:

```go
resources := shared.NewResources("my-service", 8080, 30)

// Includes: Logger, Registry, Security provider,
// System configuration, Introspection, Service manager
```

## Testing

### Test Coverage
The library includes 25 test files covering all major packages:

- Cache (CRUD, cloning, queries, notifications, statistics, collections)
- Queues (blocking/non-blocking, priority ordering, edge cases)
- Logger (async processing, log levels, direct implementation)
- Registry (type registration, lookup, dynamic instantiation, enums)
- Notifications (all types, serialization, property tracking)
- Resources (container setup, component storage/retrieval)
- Workers (concurrency, limit enforcement)
- Certificates (generation, validation)
- Strings (conversion, parsing, escaping)
- Aggregator (data aggregation operations)
- SyncMap (thread-safe operations)
- Security (shallow security provider, encryption/decryption)

### Running Tests
```bash
cd go
./test.sh  # Runs all tests with coverage reporting
```

## Dependencies

### Direct Dependencies
- **l8types** (v0.0.0-20260313): Core type definitions and interfaces
- **l8reflect** (v0.0.0-20260306): Reflection utilities for dynamic type handling
- **l8services** (v0.0.0-20260309): Services framework for microservices management
- **l8srlz** (v0.0.0-20251226): Serialization framework for data exchange
- **Protocol Buffers** (v1.36.11): Message serialization
- **golang.org/x/sys** (v0.42.0): System call support

### Indirect Dependencies
- **Google UUID** (v1.6.0): Unique identifier generation
- **l8bus** (v0.0.0-20260310): Event bus for distributed messaging
- **l8ql** (v0.0.0-20260228): Query language support
- **l8test** (v0.0.0-20260307): Testing utilities

## Design Principles

- **Interface-driven**: All major components expose interfaces from `l8types`
- **Thread-safe**: All packages use sync.RWMutex, sync.Cond, or sync.Map
- **Clone-based isolation**: Cache returns clones to prevent external mutation
- **Async processing**: Logger and cache use background goroutines with condition variables
- **Dependency injection**: Resources container enables passing components without global state
- **Vendored dependencies**: All external deps vendored in `go/vendor/`
- **Maintainability**: No file exceeds 500 lines; single responsibility per file

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For questions and support, please open an issue in the GitHub repository.
