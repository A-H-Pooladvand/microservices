.PHONY: build test lint proto swagger run clean help

# Build the application
build:
	go build -o build/app main.go

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Run linter
lint:
	golangci-lint run --timeout=5m

# Generate protobuf files
proto:
	buf generate proto/

# Generate swagger documentation
swagger:
	swag init

# Run the application
run:
	go run main.go serve

# Run database migrations
migrate:
	go run main.go migrate

# Seed the database
seed:
	go run main.go seed

# Clean build artifacts
clean:
	rm -rf build/
	rm -rf api/
	rm -rf docs/
	rm -f coverage.out coverage.html

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  proto          - Generate protobuf files"
	@echo "  swagger        - Generate swagger documentation"
	@echo "  run            - Run the application"
	@echo "  migrate        - Run database migrations"
	@echo "  seed           - Seed the database"
	@echo "  clean          - Clean build artifacts"