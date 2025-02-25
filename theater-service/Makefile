ifndef SERVICE_NAME
	SERVICE_NAME:=theater-service
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

.PHONY: help
help: ## Display this help
	@echo "Usage: make <target> [VARIABLE=value]..."
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z0-9\-\_\\:]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	  sed -e 's/\\:/:/g' | \
	  sort | \
	  awk -F': *## ' '{printf "  \033[36m%-30s \033[0m%s\n", $$1, $$2}'

.PHONY: setup
setup: ## Setup the project
	echo "Copying .env.example"
	cp .env.example .env
	echo "Setup Workspace"
	go mod download
	echo "Done"

.PHONY: build
build: ## Build the binary file based on os
	go build -o .dist/${SERVICE_NAME}${EXT} .

.PHONY: start
start: ## Run server in normal mode
	chmod +x .dist/${SERVICE_NAME}${EXT}
	.dist/${SERVICE_NAME}${EXT} start --env dev

.PHONY: start-dev
start-dev: ## Run server in development mode (will restart if any changes)
	air -c .air.toml

.PHONY: test
test: ## Run test
	go test -v ./...

test\:e2e:
	INTEGRATION_TEST=true go test -v ./...

.PHONY: test-coverage
test-coverage: ## Run unit and integration testing
	mkdir -p coverage
	go test -v -coverprofile ./coverage/cover.out ./...
	go tool cover -html=./coverage/cover.out -o ./coverage/cover.html

.PHONY: lint
lint: ## Run lint
	# go vet .
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run lint and fix
	# go fmt .
	golangci-lint run --fix

.PHONY: mockery
mock: ## Run lint and fix
	mockery

.PHONY: migration\:mysql\:create
migration\:mysql\:create: ## Create new migration file
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration:mysql:create name=your-migration-name"; \
		exit 1; \
	fi
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
    echo "Detected ENV: $$ENV"; \
	godotenv -f .env sql-migrate new -config=./deploy/database/mysql/migration-config.yml -env=$$ENV $(name)

.PHONY: migration\:mysql\:up
migration\:mysql\:up: ## Run migrations up
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
	godotenv -f .env sql-migrate up -config=./deploy/database/mysql/migration-config.yml -env=$$ENV

.PHONY: migration\:mysql\:down
migration\:mysql\:down: ## Run migrations down
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
	godotenv -f .env sql-migrate down -config=./deploy/database/mysql/migration-config.yml -env=$$ENV
