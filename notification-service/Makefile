ifndef SERVICE_NAME
	SERVICE_NAME:=user-service
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

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

.PHONY: migration\:postgresql\:create
migration\:postgresql\:create: ## Create new migration file
	@if [ -z "$(name)" ]; then \
		echo "Usage: make migration:postgresql:create name=your-migration-name"; \
		exit 1; \
	fi
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
    echo "Detected ENV: $$ENV"; \
	godotenv -f .env sql-migrate new -config=./cmd/database/postgresql/migration-config.yml -env=$$ENV $(name)

.PHONY: migration\:postgresql\:up
migration\:postgresql\:up: ## Run migrations up
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
	godotenv -f .env sql-migrate up -config=./cmd/database/postgresql/migration-config.yml -env=$$ENV

.PHONY: migration\:postgresql\:down
migration\:postgresql\:down: ## Run migrations down
	@ENV=$$(grep '^ENV=' .env | cut -d '=' -f2); \
	godotenv -f .env sql-migrate down -config=./cmd/database/postgresql/migration-config.yml -env=$$ENV

.PHONY: debezium\:register-postgresql-source-connector
debezium\:register-postgresql-source-connector: ## Register PostgreSQL source connector
	curl -i -X POST -H "Accept:application/json" -H  "Content-Type:application/json" http://localhost:8083/connectors/ -d @./deploy/local/run/config/postgresql-source-connector.json

.PHONY: debezium\:deregister-postgresql-source-connector
debezium\:deregister-postgresql-source-connector: ## Deregister PostgreSQL source connector
	curl -i -X DELETE -H "Accept:application/json" -H "Content-Type:application/json" http://localhost:8083/connectors/users-connector

.PHONY: debezium\:register-elasticsearch-sink-connector
debezium\:register-elasticsearch-sink-connector: ## Register Elasticsearch sink connector
	curl -i -X POST -H "Accept:application/json" -H  "Content-Type:application/json" http://localhost:8083/connectors/ -d @./deploy/local/run/config/elasticsearch-sink-connector.json
