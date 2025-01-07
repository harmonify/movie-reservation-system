ifndef SERVICE_NAME
	SERVICE_NAME:=mrs-cli
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

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

.PHONY: lint
lint: ## Run lint
	# go vet .
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix: ## Run lint and fix
	# go fmt .
	golangci-lint run --fix
