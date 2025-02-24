ifndef SERVICE_NAME
	SERVICE_NAME:=user-service
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## Setup the project
	echo "Copying .env.example"
	cp .env.example .env
	echo "Installing Air"
	go install github.com/air-verse/air@latest
	echo "Installing Mockery"
	go install github.com/vektra/mockery/v2@v2.42.1
	echo "Installing Migration Tools"
	go install github.com/rubenv/sql-migrate/...@latest
	echo "Installing Go Env"
	go install github.com/joho/godotenv/cmd/godotenv@latest

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
	mkdir -pv .dist && (air -c .air.toml | zap-pretty)

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
