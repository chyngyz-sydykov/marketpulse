# Load environment variables from .env
include .env
export

# Database migration variables
MIGRATE_IMAGE=migrate/migrate
MIGRATION_PATH=$(PWD)/migrations
DATABASE_URL=postgres://$(DB_USERNAME):$(DB_PASSWORD)@localhost:$(DB_PORT_EXTERNAL)/$(DB_DATABASE)?sslmode=disable

# Build and start Docker containers
up:
	docker-compose up -d

# Stop and remove Docker containers
down:
	docker-compose down

# Restart the application
restart:
	docker-compose restart web

# View logs of the application container
logs:
	docker-compose logs -f web

# Run migrations inside a temporary container
migrate-up:
	docker run --rm -v $(MIGRATION_PATH):/migrations --network=host $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" up

# Rollback the last migration
migrate-down:
	docker run --rm -v $(MIGRATION_PATH):/migrations --network=host $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" down 1

# Fix dirty migrations
migrate-fix:
	docker run --rm -v $(MIGRATION_PATH):/migrations --network=host $(MIGRATE_IMAGE) -path=/migrations -database "$(DATABASE_URL)" force $(version)

# Create a new migration file
migrate-new:
	docker run --rm -v $(MIGRATION_PATH):/migrations --network=host $(MIGRATE_IMAGE) create -ext sql -dir /migrations -seq $(name)

# Run Tests
test:
	go test ./... -v -cover  # Run all tests with verbose output and coverage

# Run Code Linting
lint:
	golangci-lint run ./...  # Run linting using golangci-lint

# Format Code
fmt:
	go fmt ./...  # Format all Go files

# Clean Docker Cache
clean-docker:
	docker system prune -f  # Clean unused Docker objects

# Build & Run the Application Locally (Outside Docker)
run-local:
	go run cmd/main.go

# Build Go Binary (For Deployment)
build:
	go build -o marketpulse ./cmd/main.go

# Remove Built Binary
clean:
	rm -f marketpulse

# Help Command
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  up              Start Docker containers"
	@echo "  down            Stop and remove containers"
	@echo "  restart         Restart the web container"
	@echo "  logs            Show application logs"
	@echo "  migrate-up      Apply database migrations"
	@echo "  migrate-down    Rollback last migration"
	@echo "  migrate-fix     Force dirty migration"
	@echo "  migrate-new   Create a new migration file (use name=<name>)"
	@echo "  test            Run all tests"
	@echo "  lint            Run code linting"
	@echo "  fmt             Format Go code"
	@echo "  clean-docker    Clean up unused Docker objects"
	@echo "  run-local       Run application locally (outside Docker)"
	@echo "  build           Compile Go application binary"
	@echo "  clean           Remove compiled binary"