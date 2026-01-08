# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`serve-grpc` is a gRPC server that exposes the Senzing SDK API. It translates gRPC requests into Senzing Go SDK calls using `senzing/sz-sdk-go-core` and returns responses to gRPC clients. The server supports TLS (server-side and mutual TLS) and optional HTTP/gRPC-web connectivity.

## Prerequisites

The Senzing C library must be installed:

- `/opt/senzing/er/lib` - shared objects
- `/opt/senzing/er/sdk/c` - SDK header files
- `/etc/opt/senzing` - configuration

Set `LD_LIBRARY_PATH=/opt/senzing/er/lib` when running.

## Build Commands

```bash
# Install development dependencies (one-time)
make dependencies-for-development

# Install Go dependencies
make dependencies

# Build
make clean build

# Build with system SQLite (Linux)
make build-with-libsqlite3

# Lint (runs golangci-lint, govulncheck, cspell)
make lint

# Run all lint auto-fixers
make fix
```

## Test Commands

```bash
# Setup test database and run tests
make clean setup test

# Run tests with HTTP enabled
make clean setup test-http

# Run a single test (Linux)
go test -tags "libsqlite3 linux" -v -run TestFunctionName ./package/...

# Coverage
make clean setup coverage

# Check coverage thresholds
make check-coverage
```

**Important**: On Linux, tests require the build tags `-tags "libsqlite3 linux"`.

## Running the Server

```bash
# Run with all services enabled
make run

# Run with HTTP enabled
make run-http

# Run with TLS
make clean setup run-server-side-tls

# Run with mutual TLS
make clean setup run-mutual-tls
```

Default ports: gRPC on 8261, HTTP on 8260.

## Architecture

### Package Structure

- `cmd/` - CLI using Cobra/Viper. Entry point is `root.go` which configures and starts the gRPC server.
- `grpcserver/` - Core gRPC server implementation. `BasicGrpcServer` initializes and registers all service handlers.
- `httpserver/` - Optional HTTP server wrapping gRPC for gRPC-web support.
- `sz*server/` - Individual gRPC service implementations:
  - `szconfigserver/` - Configuration management
  - `szconfigmanagerserver/` - Configuration manager
  - `szdiagnosticserver/` - Diagnostics
  - `szengineserver/` - Core entity resolution engine
  - `szproductserver/` - Product info

### Request Flow

1. CLI (`cmd/root.go`) parses configuration via Viper (environment variables or flags)
2. `BasicGrpcServer.Initialize()` sets up the gRPC server and registers enabled services
3. Each `sz*server` wraps the corresponding Senzing SDK component from `sz-sdk-go-core`
4. Requests come in via gRPC, get translated to SDK calls, and responses are returned

### Configuration

Configuration is via environment variables (`SENZING_TOOLS_*`) or CLI flags. Key options:

- `SENZING_TOOLS_DATABASE_URL` - Database connection string
- `SENZING_TOOLS_ENABLE_ALL` - Enable all services
- `SENZING_TOOLS_GRPC_PORT` / `SENZING_TOOLS_HTTP_PORT` - Server ports
- TLS options: `SENZING_TOOLS_SERVER_CERTIFICATE_FILE`, `SENZING_TOOLS_SERVER_KEY_FILE`, `SENZING_TOOLS_CLIENT_CA_CERTIFICATE_FILE`

## Docker

```bash
# Build Docker image
make docker-build

# Run container
docker run -it -p 8261:8261 --rm senzing/serve-grpc
```

## Code Style

- Uses `golangci-lint` with extensive linter set (see `.github/linters/.golangci.yaml`)
- Max line length: 120 characters
- Max function complexity: 11
- JSON tags use `upperSnake` case
- Run `make fix` to auto-fix common lint issues
