# Layer 8 Utils

[![Go Version](https://img.shields.io/badge/Go-1.23.8-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)]()

Common shared utilities and building blocks for Layer 8 microservices applications.

## Overview

Layer 8 Utils provides a comprehensive collection of utilities, interfaces, and default implementations designed to address common requirements in microservices architecture. The library emphasizes modularity, allowing projects to use default implementations or provide custom implementations that adhere to defined interfaces.

## Recent Updates

### Latest Changes (October 2025)
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
  - Query support with pagination and filtering
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
    "github.com/saichler/l8utils/go/utils/logger"
    "github.com/saichler/l8utils/go/utils/queues"
    "github.com/saichler/l8types/go/ifs"
)

func main() {
    // Create a logger with console output
    log := logger.NewLoggerImpl(logger.NewFmtLogMethod())
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
logger := logger.NewLoggerImpl(
    logger.NewFileLogMethod("app.log"),
    logger.NewFmtLogMethod(),
)
logger.SetLogLevel(ifs.LOG_INFO)
logger.Info("Application ready")
```

### Web Services
RESTful service utilities with protocol buffer support:

```go
service := web.NewWebService("user-service", serviceArea)
service.HandlePost(userCreateHandler)
service.HandleGet(userGetHandler)
```

### Cache
High-performance in-memory cache with optional storage persistence:

```go
// Create cache with storage backend
cache := cache.NewCache(&MyModel{}, storage, nil, resources)

// CRUD operations
cache.Post(item, true)  // Add with notification
item := cache.Get(key)
cache.Put(key, updatedItem, true)
cache.Patch(key, changes, true)
cache.Delete(key, true)

// Query with pagination
results := cache.Fetch(0, 25, query)

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
- **l8reflect** (v0.0.0-20251012150625-41304187d527): Reflection utilities for dynamic type handling
- **l8srlz** (v0.0.0-20251010183545-1dc2ad85aec0): Serialization framework for data exchange
- **l8types** (v0.0.0-20251020135133-5dc074d0670c): Core type definitions and interfaces
- **Protocol Buffers** (v1.36.10): Message serialization and data exchange

### Indirect Dependencies
- **Google UUID** (v1.6.0): Unique identifier generation and management
- **l8bus** (v0.0.0-20251020135521-8892f5ac8f9c): Event bus for distributed messaging
- **l8ql** (v0.0.0-20250927164348-155ee588c3cb): Query language support
- **l8services** (v0.0.0-20251019174910-451b61827826): Services framework
- **l8test** (v0.0.0-20251019130747-4b37b734925d): Testing utilities

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

### October 2025
- `d1faf0e` - Add NewOf function for dynamic instance creation
- `e5abe4b` - Add TotalStats feature for automatic cache statistics
- `f854626` - Remove sync dependencies
- `bc797de` - Fix tests
- `0a66b52` - Add model type support
- `e928f7d` - Add Collect functionality for cache aggregation
- `0f9e63f` - Add comprehensive cache & notification systems
- `92ff836` - Add tests for notification framework
- `ceb655e` - Update README documentation

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For questions and support, please open an issue in the GitHub repository.