export POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

.PHONY: run
run:
	@echo "Running the program..."
	@go run cmd/main.go

## Create a new migration file: Usage `make migration/create name=<migration_name>`
.PHONY: migration/create
migration/create:
	@echo "Creating a new migration..."
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir postgres/migrations -seq $(name)

.PHONY: dev/start
dev/start:
	@echo "Starting the development server..."
	@docker-compose up -d

.PHONY: dev/stop
dev/stop:
	@echo "Stopping the development server..."
	@docker-compose down

.PHONY: integration-test
integration-test:
	@echo "Running integration tests..."
	@go test -v ./...
