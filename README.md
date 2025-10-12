# Layer 8 Utils

[![Go Version](https://img.shields.io/badge/Go-1.23.8-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)]()

Common shared utilities and building blocks for Layer 8 microservices applications.

## Overview

Layer 8 Utils provides a comprehensive collection of utilities, interfaces, and default implementations designed to address common requirements in microservices architecture. The library emphasizes modularity, allowing projects to use default implementations or provide custom implementations that adhere to defined interfaces.

## Recent Updates

- **Cache System**: New high-performance in-memory cache with storage integration, CRUD operations, notifications, and query support (78.1% test coverage)
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
  - Statistics tracking with functional filters
  - Clone-based isolation for concurrent access
  - Query support with pagination and filtering

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

// Statistics tracking
stats := cache.Stats(func(item interface{}) bool {
    return item.(*MyModel).Status == "active"
})
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

## Testing

The library includes comprehensive test suites with high code coverage:

- **Cache Tests** (`go/tests/Cache_test.go`): 78.1% coverage
  - CRUD operations testing
  - Storage integration tests
  - Notification system tests
  - Statistics tracking validation
  - Concurrent access and isolation tests

- **Notification Tests** (`go/tests/Notifications_test.go`): 87.8% coverage
  - All notification type tests (Add, Delete, Update, Replace, Sync)
  - Serialization/deserialization validation
  - ItemOf extraction tests
  - Property-level change tracking tests
  - Sequence and service area validation

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

- **l8types** (v0.0.0-20250926135209-1d316857fdcf): Core type definitions and interfaces
- **Protocol Buffers** (v1.36.9): Message serialization and data exchange
- **Google UUID** (v1.6.0): Unique identifier generation and management

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

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For questions and support, please open an issue in the GitHub repository.