# Layer 8 Utils

[![Go Version](https://img.shields.io/badge/Go-1.24.9-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)]()

Common shared utilities and building blocks for Layer 8 microservices applications.

## Overview

Layer 8 Utils provides a comprehensive collection of utilities, interfaces, and default implementations designed to address common requirements in microservices architecture. The library emphasizes modularity, allowing projects to use default implementations or provide custom implementations that adhere to defined interfaces.

## Recent Updates

### Latest Changes (December 2025)
- **Query TTL Support**: Added automatic TTL-based cleanup for cached queries with configurable expiration (30s default)
- **Query Sorting Fix**: Fixed sorting behavior in cache queries for consistent ordering
- **Web Service Refactoring**: Improved web service architecture with better code organization
- **Performance Analysis**: Comprehensive performance analysis document added with scalability recommendations
- **Bug Fixes**: Fixed missing root key handling, zero value handling, nil support improvements

### November 2025 Updates
- **VNet Support**: Added VNet (Virtual Network) support to WebService for enhanced networking capabilities
- **Logging Enhancements**: Fixed and improved file logging functionality with better error handling
- **Shared Resources**: Added new shared resource utilities (`NewResources`) for centralized resource management
- **Service Integration**: Enhanced integration with l8services for improved microservices support

### October 2025 Updates
- **Registry Enhancement**: Added `NewOf()` function for dynamic instance creation from registered types
- **Cache Statistics**: Introduced `TotalStats` feature with automatic total counting for all cache items
- **Collection Support**: Added `Collect()` functionality for advanced data aggregation in cache
- **Model Type Integration**: Enhanced type handling and model type support across the framework
- **Source Management**: Improved source tracking and management in components

### Core Improvements
- **Cache System**: High-performance in-memory cache with storage integration, CRUD operations, notifications, and query support (78.1% test coverage)
- **Notification Framework**: Comprehensive notification system for distributed state management with support for Add, Delete, Update, Replace, and Sync operations (87.8% test coverage)
- **Test Suite Expansion**: Added comprehensive test utilities (createModel, newResources, createChanges) for building robust test suites
- **Enhanced Security**: Improved certificate management with self-signed certificate support
- **Interface Improvements**: Fixed interface implementations for better type safety
- **Token Validation**: Added robust token validation mechanisms
- **Performance Optimizations**: Enhanced byte queue performance and enlarged queue sizes

## Features

### 🚀 Core Utilities

- **Cache**: High-performance in-memory cache with storage integration
  - CRUD operations (Post, Get, Put, Patch, Delete)
  - Storage layer integration with persistence support
  - Built-in notification system for change tracking
  - Enhanced statistics tracking with named stat functions and automatic totals
  - Collection operations with `Collect()` for data aggregation
  - Clone-based isolation for concurrent access
  - Query support with pagination, filtering, and sorting
  - Query TTL with automatic cleanup (configurable, 30s default)
  - Dynamic stat functions registration with `AddStatFunc()`

- **Notifications**: Comprehensive notification system for distributed state management
  - Support for Add, Delete, Update, Replace, and Sync notification types
  - Serialization/deserialization with protocol buffers
  - Change tracking with property-level granularity
  - Sequence numbering for ordering guarantees
  - Service-area based routing support

- **Queues**: High-performance thread-safe queues with priority support
  - `ByteQueue`: Optimized byte queue with priority handling
  - `Queue`: Generic queue implementation
  - Support for concurrent operations and backpressure

- **Logging**: Flexible logging framework with multiple output methods
  - File-based logging with rotation
  - Console/fmt logging
  - Configurable log levels
  - Asynchronous logging with queue-based processing

- **String Utilities**: Comprehensive string manipulation and conversion
  - Type-safe string to primitive conversions
  - String formatting and parsing utilities
  - Validation helpers

### 🔐 Security & Infrastructure

- **Certificate Management**: TLS/SSL certificate utilities
- **Shallow Security**: Basic security providers and authentication utilities
- **Maps**: Thread-safe map implementations with sync.Map optimizations

### 🌐 Web Services

- **Web Service Framework**: RESTful service utilities
  - HTTP method handlers (GET, POST, PUT, PATCH, DELETE)
  - Request/response marshaling
  - Protocol buffer integration
  - VNet (Virtual Network) support for distributed networking
  - Enhanced service-to-service communication

### 📊 System Management

- **Registry**: Centralized resource and configuration management
  - Type registration and lookup
  - Dynamic instance creation with `NewOf()` function
  - Enum registration and management
  - Thread-safe operations
- **Resources**: Resource loading and management utilities
- **Workers**: Worker pool implementations for concurrent processing

## Installation

```bash
go get github.com/saichler/l8utils/go
```

## Quick Start

```go
package main

import (
    "github.com/saichler/l8utils/go/utils"
    "github.com/saichler/l8utils/go/utils/logger"
    "github.com/saichler/l8utils/go/utils/queues"
    "github.com/saichler/l8types/go/ifs"
)

func main() {
    // Create shared resources with VNet support
    resources := utils.NewResources("my-service", 8080, 30)

    // Get logger from resources
    log := resources.Logger()
    log.Info("Application started")

    // Create a high-performance byte queue
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
│   ├── cache/          # High-performance cache with storage integration
│   ├── certs/          # Certificate management
│   ├── logger/         # Logging framework
│   ├── maps/           # Thread-safe map implementations
│   ├── notify/         # Notification system for state changes
│   ├── queues/         # High-performance queue implementations
│   ├── registry/       # Resource registry
│   ├── resources/      # Resource management
│   ├── shallow_security/ # Basic security utilities
│   ├── strings/        # String manipulation utilities
│   ├── web/           # Web service framework
│   └── workers/       # Worker pool implementations
└── tests/
    ├── Cache_test.go         # Cache tests (78.1% coverage)
    └── Notifications_test.go # Notification tests (87.8% coverage)
```

## Key Components

### ByteQueue
High-performance, thread-safe queue optimized for byte operations with priority support:

```go
queue := queues.NewByteQueue("processor", 5000)
queue.Add(data, ifs.PRIORITY_HIGH)
result := queue.Poll() // Non-blocking
result := queue.Next() // Blocking until data available
```

### Logger
Asynchronous logging with multiple output methods:

```go
// Console logging
logger := logger.NewLoggerImpl(logger.NewFmtLogMethod())

// File logging with automatic file management
fileLog := logger.NewFileLogMethod("app.log")
logger := logger.NewLoggerImpl(fileLog)

// Combined logging
logger := logger.NewLoggerImpl(
    logger.NewFileLogMethod("app.log"),
    logger.NewFmtLogMethod(),
)
logger.SetLogLevel(ifs.Error_Level)
logger.Info("Application ready")
```

### Web Services
RESTful service utilities with protocol buffer and VNet support:

```go
// Create web service with VNet support
service := web.NewWebService("user-service", serviceArea)
service.SetVnet(8080) // Set VNet port

// Register handlers
service.HandlePost(userCreateHandler)
service.HandleGet(userGetHandler)
service.HandlePut(userUpdateHandler)
service.HandleDelete(userDeleteHandler)
```

### Shared Resources
Centralized resource management with integrated components:

```go
// Create resources with alias, VNet port, and keep-alive settings
resources := utils.NewResources("my-service", 8080, 30)

// Resources automatically include:
// - Logger with error-level default
// - Registry for type management
// - Security provider
// - System configuration
// - Introspection capabilities
// - Service manager
```

### Cache
High-performance in-memory cache with optional storage persistence and query TTL:

```go
// Create cache with storage backend
cache := cache.NewCache(&MyModel{}, initElements, storage, resources)

// CRUD operations
cache.Post(item, true)  // Add with notification
item := cache.Get(key)
cache.Put(key, updatedItem, true)
cache.Patch(key, changes, true)
cache.Delete(key, true)

// Query with pagination (queries cached with 30s TTL by default)
results := cache.Fetch(0, 25, query)

// Query cache management
queryCount := cache.QueryCount()           // Monitor cached queries
cache.CleanupQueriesNow(60)               // Manual cleanup with custom TTL
defer cache.Close()                        // Stop TTL cleaner on shutdown

// Enhanced statistics tracking
cache.AddStatFunc("active", func(item interface{}) bool {
    return item.(*MyModel).Status == "active"
})
cache.AddStatFunc("pending", func(item interface{}) bool {
    return item.(*MyModel).Status == "pending"
})
stats := cache.Stats() // Returns map with counts for "Total", "active", "pending"

// Collection operations
collection := cache.Collect(predicate)
```

### Notifications
Create and manage notifications for distributed state synchronization:

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
// Create registry
registry := registry.NewRegistry()

// Register a type
registry.Register(&MyModel{})

// Create new instance dynamically
newInstance := registry.NewOf(&MyModel{})

// Get type information
info, err := registry.Info("MyModel")

// Register enums
registry.RegisterEnum("Status", []string{"active", "pending", "completed"})
```

## Testing

The library includes comprehensive test suites with high code coverage:

### Coverage Reports
- **Cache Package**: 78.1% coverage
- **Notifications Package**: 87.8% coverage
- **Overall**: Comprehensive test coverage across all packages

### Test Files
- **Cache Tests** (`go/tests/Cache_test.go`):
  - CRUD operations testing
  - Storage integration tests
  - Notification system tests
  - Statistics tracking validation including new TotalStats feature
  - Collection operations testing
  - Concurrent access and isolation tests

- **Notification Tests** (`go/tests/Notifications_test.go`):
  - All notification type tests (Add, Delete, Update, Replace, Sync)
  - Serialization/deserialization validation
  - ItemOf extraction tests
  - Property-level change tracking tests
  - Sequence and service area validation

### Running Tests
```bash
cd go
./test.sh  # Runs all tests with coverage reporting
```

### Coverage Visualization
HTML coverage reports are generated at:
- `go/cover-report.html` - Full coverage report
- `go/coverage_notify.html` - Notification-specific coverage

Test utilities available:
```go
// Create test model instance
model := createModel(i)

// Create resources with registry and introspection
resources := newResources()

// Generate changes between models
changes := createChanges(oldModel, newModel, resources)
```

## Dependencies

### Direct Dependencies
- **l8reflect** (v0.0.0-20251020202633-feaa244d0a2b): Reflection utilities for dynamic type handling
- **l8services** (v0.0.0-20251031163521-852f7c020c80): Services framework for microservices management
- **l8srlz** (v0.0.0-20251027151455-5149a019bed7): Serialization framework for data exchange
- **l8types** (v0.0.0-20251103131334-82d3444d09d8): Core type definitions and interfaces
- **Protocol Buffers** (v1.36.10): Message serialization and data exchange

### Indirect Dependencies
- **Google UUID** (v1.6.0): Unique identifier generation and management
- **l8bus** (v0.0.0-20251031141311-e67190ca68dc): Event bus for distributed messaging
- **l8ql** (v0.0.0-20251030150208-8a58a1d7ac8a): Query language support
- **l8test** (v0.0.0-20251030140121-4de54523fc40): Testing utilities

## Performance Features

- Lock-free operations where possible
- Memory-efficient byte handling
- Priority-based queue processing
- Asynchronous logging to prevent I/O blocking
- Connection pooling and reuse

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Recent Commits

### December 2025
- `038ab24` - Fix sorting in queries - Fixed query result sorting
- `65c752f` - Add TTL to cache query - Automatic query cache cleanup with TTL support
- `c7794a3` - Fix missing root key - Improved key handling
- `8296546` - Refactor web - Improved web service architecture
- `0b87c4e` - Fix zeroValue - Fixed zero value handling
- `156e414` - Add support for nil - Improved nil handling

### November 2025
- `2119a29` - Add vnet to web - Added VNet support to WebService
- `64ab2f9` - Fix log to file - Fixed file logging functionality
- `5ca22f1` - move loader - Reorganized loader components
- `27fa796` - add shared - Added shared resource utilities
- `cb17025` - Log to files - Enhanced file logging capabilities

### October 2025
- `d7f9d0d` - Wait for signal - Added signal handling
- `fe2c916` - remove default user - Security enhancement
- `74eaae4` - Add log to files - Initial file logging implementation
- `046924f` - Add file logs - File logging support
- `06f7f2d` - update readme - Documentation updates

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For questions and support, please open an issue in the GitHub repository.