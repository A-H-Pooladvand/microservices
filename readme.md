# Microservices Application

A production-ready microservices application built with Go, following the Hexagonal Architecture (Ports and Adapters) pattern. This project provides a solid foundation for building scalable, maintainable microservices.

## Architecture

This project follows the **Hexagonal Architecture** pattern:

```
internal/
├── core/                    # Business Logic (Domain)
│   ├── domain/             # Domain entities and errors
│   ├── port/               # Interfaces (ports)
│   └── service/            # Business logic implementation
├── adapter/                 # Adapters (Infrastructure)
│   ├── inbound/            # Driving adapters (HTTP, gRPC handlers)
│   │   ├── http/           # HTTP handlers
│   │   └── grpc/           # gRPC server handlers
│   └── outbound/           # Driven adapters (Repositories, Clients)
│       ├── repository/     # Database repositories
│       └── grpcclient/     # gRPC client implementations
```

### Key Principles

- **Domain-Centric**: Business logic is isolated in the core layer
- **Dependency Inversion**: Core depends on interfaces, not implementations
- **Testability**: Each layer can be tested independently with mocks
- **Flexibility**: Easy to swap infrastructure components

## Features

- [x] **Echo** - High-performance HTTP framework
- [x] **gRPC** - Both client and server implementations
- [x] **PostgreSQL** - Database with GORM ORM
- [x] **Redis** - Caching and distributed locks
- [x] **Vault** - Secrets management
- [x] **RabbitMQ** - Message queue support
- [x] **Swagger** - API documentation
- [x] **Jaeger/OTLP** - Distributed tracing
- [x] **Prometheus** - Metrics collection
- [x] **Uber Fx** - Dependency injection
- [x] **Zap** - Structured logging
- [x] **Cobra** - CLI commands

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL
- Redis (optional)
- Docker (optional, for running dependencies)

### Installation

1. Clone the repository
2. Copy and configure the config file:
   ```bash
   cp config.example.yml config.yml
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```

### Configuration

Edit `config.yml` to match your environment:

```yaml
app:
  Name: my-service
  Env: local  # Options: local, dev, staging, production
  Debug: true
  Port: 8000

postgres:
  Host: 127.0.0.1
  Port: 5432
  Username: postgres
  Password: postgres
  Db: microservices

grpc:
  addr: "9000"  # Leave empty to disable gRPC server
```

### Running the Application

```bash
# Build the application
make build

# Or run directly
go run main.go serve

# Run database migrations
go run main.go migrate

# Seed the database
go run main.go seed
```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Generating Protobuf Files

```bash
# Install buf tool
go install github.com/bufbuild/buf/cmd/buf@latest

# Generate protobuf files
buf generate proto/
```

### Generating Swagger Documentation

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init
```

## API Endpoints

### HTTP Endpoints

- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get a user by ID
- `PUT /api/v1/users/:id` - Update a user
- `DELETE /api/v1/users/:id` - Delete a user
- `GET /health` - Health check
- `GET /metric` - Prometheus metrics
- `GET /swagger/*` - Swagger documentation

### gRPC Services

The User service provides the following RPCs:
- `CreateUser`
- `GetUser`
- `ListUsers`
- `UpdateUser`
- `DeleteUser`

## Local Development

1. Set `app.Env` to `local` in your config
2. The application will skip certain production-only features (e.g., Vault)
3. Debug mode enables verbose logging

## Production Deployment

1. Set `app.Env` to `production`
2. Set `app.Debug` to `false`
3. Configure proper database connection pooling
4. Enable TLS for gRPC if needed
5. Configure proper tracing sampling ratio

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).
