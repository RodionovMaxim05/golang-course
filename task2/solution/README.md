# GitHub Repository Info - Distributed System

## Description

A distributed system for fetching GitHub repository information, consisting of two microservices:

- **Collector** - gRPC server that communicates with the GitHub API
- **Gateway** - REST server and gRPC client that exposes a public HTTP API with Swagger UI

## Architecture

`Client` → **HTTP** → `Gateway` → **gRPC** → `Collector` → **GitHub API**

Both services follow **Clean Architecture** principles with the following layers:

```bash
adapters/   - incoming (HTTP/gRPC handlers) and outgoing (GitHub API, gRPC clients)
usecases/   - business logic orchestration
services/   - domain-specific logic
domain/     - domain models and errors
```

## Project Structure

```bash
.
├── api/                       # Protobuf schemas and generated code
│   ├── collector.proto
│   ├── gen/
│   └── Makefile
├── collector/                 # gRPC server
│   ├── cmd/main.go
│   └── internal/
│       ├── adapters/
│       │   ├── clients/       # GitHub API client
│       │   └── controllers/   # gRPC handler
│       ├── usecases/
│       ├── services/
│       └── domain/
├── gateway/                   # REST server
│   ├── cmd/main.go
│   ├── docs/                  # Generated Swagger docs
│   └── internal/
│       ├── adapters/
│       │   ├── clients/       # gRPC client to Collector
│       │   └── controllers/   # HTTP handler + Swagger
│       ├── usecases/
│       ├── services/
│       └── domain/
└── infrastructure/
    └── local/
        ├── docker-compose.yaml
        ├── collector.dockerfile
        └── gateway.dockerfile
```

## Requirements

- Go 1.25+
- Docker + Docker Compose v2

## Running with Docker

```bash
docker compose -f infrastructure/local/docker-compose.yaml up --build
```

Stop:

```bash
docker compose -f infrastructure/local/docker-compose.yaml stop

# or to stop and remove containers:

docker compose -f infrastructure/local/docker-compose.yaml down
```

## Running locally

**Collector** (terminal 1):

```bash
cd collector
go run cmd/main.go
```

**Gateway** (terminal 2):

```bash
cd gateway
go run cmd/main.go
```

## Usage

Once both services are running, open `Swagger UI`:

```bash
http://localhost:8080/swagger/index.html
```

Or use `curl` directly:

```bash
curl "http://localhost:8080/api/v1/repo?url=https://github.com/torvalds/linux"
```

**Example response:**

```json
{
  "Name":"linux",
  "Description":"Linux kernel source tree",
  "StargazersCount":223280,
  "ForksCount":61037,
  "CreatedAt":"2011-09-04T22:48:12Z"
}
```

## Environment Variables

| Variable         | Service   | Default          | Description                  |
|------------------|-----------|------------------|------------------------------|
| `COLLECTOR_PORT` | Collector | `:50051`         | gRPC server port             |
| `COLLECTOR_ADDR` | Gateway   | `localhost:50051`| Collector gRPC address       |
| `GATEWAY_PORT`   | Gateway   | `:8080`          | HTTP server port             |
