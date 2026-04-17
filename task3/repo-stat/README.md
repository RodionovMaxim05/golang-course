# Repository Statistics Microservices

## Description

A microservices-based application for fetching and managing GitHub repository statistics. The project follows Clean Architecture principles and uses gRPC for inter-service communication.

## Architecture

The system consists of 4 microservices:

- **API Gateway** - REST API server that accepts external HTTP requests and coordinates communication with other services
- **Processor** - Intermediate gRPC service that acts as a bridge between API Gateway and Collector
- **Collector** - gRPC service that handles GitHub API interactions and repository data fetching
- **Subscriber** - Reference implementation service for health checks

## Requirements

- Go 1.25+
- Docker & Docker Compose
- Make

## Quick Start

### Build everything

```bash
make protobuf  # Generate protobuf code
make build     # Build all services
```

### Run with Docker Compose

```bash
docker compose up --build
```

Services will start on:

- API Gateway: `http://localhost:28080`
- Processor: `grpc://localhost:8083`
- Collector: `grpc://localhost:8082`
- Subscriber: `grpc://localhost:8081`

## Development

### Available Make Commands

```bash
make protobuf    # Generate protobuf code
make lint        # Run linter
make gofmtcheck  # Check code formatting
make golint      # Run Go linter
make test        # Run tests in Docker
make build       # Build all services
```

### Testing

Run tests:

```bash
make test        # Full test suite in Docker
```

## API Endpoints

### Health Check

```
GET /api/ping
```

Returns status of all services:

```json
{
  "status": "ok",
  "services": [
    {
      "name": "processor",
      "status": "up"
    },
    {
      "name": "subscriber",
      "status": "up"
    }
  ]
}
```

Returns `503 Service Unavailable` if any service is down.

### Repository Information

```
GET /api/repositories/info?url=<github_repo_url>
```

Fetch repository information from GitHub:

**Example:**

```bash
curl "http://localhost:28080/api/repositories/info?url=https://github.com/golang/go"
```

**Response:**

```json
{
  "full_name": "golang/go",
  "description": "The Go programming language",
  "stars": 123456,
  "forks": 12345,
  "created_at": "2009-11-10T23:00:00Z"
}
```

### Swagger Documentation

API documentation available at:

```
http://localhost:28080/swagger/index.html
```
