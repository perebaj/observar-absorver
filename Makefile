GOLANGCI_LINT_VERSION = v1.57.2
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
	if [ -n "$(testcase)" ]; then \
		go test ./... -timeout 5s -v -run="^$(testcase)$$" ; \
	else \
		go test ./... -timeout 5s; \
	fi

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run -v

## Enrich the chunk of text with the GPT model (create embeddings over it)
.PHONY: run/gpt
run/gpt:
	@echo "Running embedding example..."
	@go run cmd/gpt/main.go
