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

# Create a new migration file
migrate-new:
	docker run --rm -v $(MIGRATION_PATH):/migrations --network=host $(MIGRATE_IMAGE) create -ext sql -dir /migrations -seq $(name)

