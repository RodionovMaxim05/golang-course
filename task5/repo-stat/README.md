# Repository Statistics Microservices

## Description

A microservices-based application for fetching and managing GitHub repository statistics with subscription functionality. Users can subscribe to repositories and receive aggregated information about their subscribed repositories. The project follows Clean Architecture principles and uses gRPC for inter-service communication.

## Architecture

The system consists of 4 microservices:

- **API Gateway** - REST API server that accepts external HTTP requests and coordinates communication with other services
- **Processor** - Intermediate gRPC service that acts as a bridge between API Gateway and Collector
- **Collector** - gRPC service that handles GitHub API interactions and repository data fetching
- **Subscriber** - gRPC service that manages repository subscriptions in PostgreSQL database

## Requirements

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL (for subscription storage)
- Make

## Quick Start

### Run with Docker Compose

```bash
docker compose up --build
```

Stop:

```bash
docker compose stop

# or to stop and remove containers:

docker compose down
```

Services will start on:

- API Gateway: `http://localhost:28080`
- Processor: `grpc://localhost:8083`
- Collector: `grpc://localhost:8082`
- Subscriber: `grpc://localhost:8081`
- PostgreSQL: `localhost:5432`

### Environment Variables

The application uses the following environment variables (with defaults):

- `POSTGRES_USER=postgres`
- `POSTGRES_PASSWORD=password`
- `POSTGRES_DB=repo_stat`
- `POSTGRES_HOST=postgres`
- `POSTGRES_PORT=5432`
- `POSTGRES_SSLMODE=disable`

## Development

### Database

The Subscriber service uses PostgreSQL for storing repository subscriptions. Database migrations are automatically applied when the services start.

### Available Make Commands

```bash
make tools       # Install development tools
make protobuf    # Generate protobuf code
make swagger     # Generate swagger docs
make lint        # Run linter
make test        # Run tests in Docker
make build       # Build all services
```

### Testing

Run tests:

```bash
(cd ../ & make test)        # Full test suite in Docker
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

### Subscription Management

#### List Subscriptions

```
GET /subscriptions
```

Returns all subscribed repositories with creation timestamps:

**Example:**

```bash
curl "http://localhost:28080/subscriptions"
```

**Response:**

```json
[
  {
    "owner":"torvalds",
    "repo":"linux",
    "created_at":"2026-04-18T15:50:11.267046Z"
  },
  {
    "owner":"golang",
    "repo":"go",
    "created_at":"2026-04-18T15:49:34.514466Z"
  }
]
```

#### Subscribe to Repository

```
POST /subscriptions?url=<github_repo_url>
```

Subscribe to a GitHub repository for monitoring:

**Example:**

```bash
curl -X POST "http://localhost:28080/subscriptions?url=https://github.com/golang/go"
```

**Response:**

```json
{
  "owner": "golang",
  "repo": "go",
  "created_at":"2026-04-18T15:50:11.267046Z"
}
```

#### Unsubscribe from Repository

```
DELETE /subscriptions/{owner}/{repo}
```

Remove subscription for a specific repository:

**Example:**

```bash
curl -X DELETE "http://localhost:28080/subscriptions/golang/go"
```

#### Get Subscribed Repositories Info

```
GET /subscriptions/info
```

Retrieve aggregated information for all subscribed repositories:

**Example:**

```bash
curl "http://localhost:28080/subscriptions/info"
```

**Response:**

```json
[
  {
    "full_name":"torvalds/linux",
    "description":"Linux kernel source tree",
    "stars":229724,
    "forks":61710,
    "created_at":"2011-09-04T22:48:12Z"
  },
  {
    "full_name":"golang/go",
    "description":"The Go programming language",
    "stars":133520,
    "forks":18934,
    "created_at":"2014-08-19T04:33:40Z"
  }
]
```

### Swagger Documentation

API documentation available at:

```
http://localhost:28080/swagger/index.html
```
