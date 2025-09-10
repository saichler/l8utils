# L8Utils

[![Go Version](https://img.shields.io/badge/Go-1.23.8-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Common shared utilities and building blocks for Layer8 microservices applications.

## Overview

L8Utils provides a comprehensive collection of utilities, interfaces, and default implementations designed to address common requirements in microservices architecture. The library emphasizes modularity, allowing projects to use default implementations or provide custom implementations that adhere to defined interfaces.

## Features

### 🚀 Core Utilities

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
└── utils/
    ├── certs/          # Certificate management
    ├── logger/         # Logging framework
    ├── maps/           # Thread-safe map implementations
    ├── queues/         # High-performance queue implementations
    ├── registry/       # Resource registry
    ├── resources/      # Resource management
    ├── shallow_security/ # Basic security utilities
    ├── strings/        # String manipulation utilities
    ├── web/           # Web service framework
    └── workers/       # Worker pool implementations
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

## Dependencies

- **l8types**: Core type definitions and interfaces
- **Protocol Buffers**: Message serialization
- **UUID**: Unique identifier generation

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