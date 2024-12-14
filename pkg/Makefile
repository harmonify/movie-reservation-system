ifndef SERVICE_NAME
	SERVICE_NAME:=auth-system
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

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
