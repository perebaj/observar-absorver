export POSTGRES_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-z\/]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: migration/create
migration/create: ## create a new migration file. Usage `make migration/create name=<migration_name>`
	@echo "Creating a new migration..."
	@go run github.com/golang-migrate/migrate/v4/cmd/migrate create -ext sql -dir postgres/migrations -seq $(name)

.PHONY: dev/start
dev/start: ## Start the development server
	@echo "Starting the development server..."
	@docker-compose up -d

.PHONY: dev/stop
dev/stop: ## Stop the development server
	@echo "Stopping the development server..."
	@docker-compose down

.PHONY: integration-test
integration-test: ## Run integration tests
	@echo "Running integration tests..."
	@go test -v ./...
