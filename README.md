# Repository Statistics Microservices

## Description

A microservices-based application for fetching and managing GitHub repository statistics with subscription functionality. Users can subscribe to repositories and receive aggregated information about their subscribed repositories. The project follows Clean Architecture principles, uses gRPC for synchronous inter-service communication, and Kafka for asynchronous event-driven interactions.

## Architecture

The system utilizes a distributed microservices pattern:

- **API Gateway** — REST API server that handles incoming HTTP traffic from the сlient. It uses **Redis** as a rate limiter to protect downstream services and as an edge cache to instantly serve frequent lookup results. It routes requests downstream using **gRPC**.
- **Processor** — The central orchestration engine of the system. It receives **gRPC** requests and maintains its own **Processor DB** to store cached repository metrics and workflow states. It acts as a producer for asynchronous background data operations via **Kafka** and fetches active user subscription lists through a **gRPC** link to the Subscriber.
- **Collector** — A stateless worker service dedicated to data integration and scraping. It processes tasks via the **Kafka Broker** and fetches repository details from the external **GitHub Cloud**. Additionally, it periodically independently queries the Subscriber via **gRPC** to fetch all active subscriptions and publishes update tasks back to Kafka for each one, ensuring that repository states in the Processor are regularly refreshed without direct user interaction.
- **Subscriber** — A dedicated domain service managing user repository subscriptions. It handles subscription lifecycle actions via incoming **gRPC** channels, independently communicates with the **GitHub Cloud** to validate repository existence when a new subscription is being added, and manages an isolated **Subscriber DB** to encapsulate subscription watchlists.

![UML Deploymnt Diagram](figures/deployment_diagram.svg)

## Requirements

### For Deployment & Quick Start

If you run the application via Docker Compose, you only need:

- **Docker & Docker Compose**
- **Make** (optional, for convenience commands)

### For Local Development

If you prefer to build and run the binary natively on your host machine:

- **Go 1.25+**
- **Apache Kafka**
- **PostgreSQL (separate instances for Subscriber and Processor)**

## Quick Start

### Run with Docker Compose

```bash
make up
```

To stop the services:

```bash
make down
```

To stop services and completely wipe database volumes:

```bash
make down-volumes
```

## Configuration

The application splits its configuration into environmental secrets (`.env`) and operational parameters (`<service>/config/config.yaml`).

Before running the application, you **must copy the example file and adjust them** for your environment.

### 1. Infrastructure & Secrets (`.env`)

Copy the template file to create your local environment file:

```bash
cp .env.example .env
```

Open .env and fill in your database credentials and infrastructure ports (to run via Docker, the pre-configured default values for Kafka and isolated PostgreSQL databases will work out of the box).

### 2. Services Settings (`<service>config/config.yaml`)

Edit each `<service>/config/config.yaml` file if you need to change the default settings for the corresponding service.

## Development

### Available Make Commands

```bash
# Infrastructure
make up                 # Up containers in background with docker compose
make down               # Stop and remove active containers
make down-volumes       # Stop containers and completely wipe all docker volumes

# Code Generation
make protobuf           # Compile all .proto contract files using module resolution
make swagger            # Generate Swagger documentation for API Gateway

# Linting & Formatting
make lint               # Run all active linters (protolint, gofmtcheck, golint)
make fix                # Auto-fix proto/go formatting and common lint issues

# Testing
make unit-test          # Run unit tests with race detector and generate HTML coverage
make integration-test   # Up isolated environment, run integration tests, and down env

# Tooling
make tools              # Install all necessary local development binaries and plugins
```

## API Endpoints

### Summary Table

| Method & Path | Description |
| :--- | :--- |
| `GET /api/ping` | **Health Check:** Returns the operational status of all downstream services. |
| `GET /api/repositories/info` | **Repository Information:** Fetches details for a specific repository directly from GitHub. |
| `GET /api/subscriptions` | **List Subscriptions:** Retrieves all currently monitored repositories with their creation timestamps. |
| `POST /api/subscriptions` | **Subscribe to Repository:** Adds a new GitHub repository to the monitoring system. |
| `DELETE /api/subscriptions/{owner}/{repo}` | **Unsubscribe from Repository:** Removes a specific repository from monitored subscriptions. |
| `GET /api/subscriptions/info` | **Get Subscribed Repositories Info:** Retrieves aggregated metrics and data for all subscribed repositories. |

### Health Check

```
GET /api/ping
```

Returns status of all services.

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

Fetch repository information from GitHub.

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
GET /api/subscriptions
```

Returns all subscribed repositories with creation timestamps.

**Example:**

```bash
curl "http://localhost:28080/api/subscriptions"
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
POST /api/subscriptions?url=<github_repo_url>
```

Subscribe to a GitHub repository for monitoring.

**Example:**

```bash
curl -X POST "http://localhost:28080/api/subscriptions?url=https://github.com/golang/go"
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
DELETE /api/subscriptions/{owner}/{repo}
```

Remove subscription for a specific repository.

**Example:**

```bash
curl -X DELETE "http://localhost:28080/api/subscriptions/golang/go"
```

#### Get Subscribed Repositories Info

```
GET /api/subscriptions/info
```

Retrieve aggregated information for all subscribed repositories.

**Example:**

```bash
curl "http://localhost:28080/api/subscriptions/info"
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

Interactive API playground is available at:

```
http://localhost:28080/swagger/index.html
```
