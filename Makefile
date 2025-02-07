ifndef SERVICE_NAME
	SERVICE_NAME:=auth-system
endif

ifeq ($(OS), Windows_NT)
	SERVICE_NAME :=${SERVICE_NAME}
	EXT :=.exe
endif

.PHONY: setup
setup: ## Setup the project
	echo "Installing Air"
	go install github.com/cosmtrek/air@latest
	echo "Installing Protoc"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	echo "Installing Mockery"
	go install github.com/vektra/mockery/v2@v2.42.1
	echo "Installing Go Env"
	go install github.com/joho/godotenv/cmd/godotenv@latest
	echo "Setup Workspace"
	go mod download
	echo "Sync Submodules"
	git submodule foreach git pull origin main
	echo "Installing Docker plugin"
	docker plugin install grafana/loki-docker-driver:2.9.2 --alias loki --grant-all-permissions
	echo "Done..."

.PHONY: start
start: ## Deploy infrastructure
	docker-compose up -d

.PHONY: build-cli
build-cli: ## Build MRS CLI
	bin/build_cli.sh

.PHONY: gen-proto
gen-proto: ## Generate protobuf code for Golang in specified output directory
	@if [ -z "$(output)" ]; then \
		echo "Example usage: make gen-proto output=user-service/internal/driven/proto"; \
		exit 1; \
	fi
	bin/gen_proto.sh "$(output)"

.PHONY: export-puml
export-puml: ## ex: make export-puml darkmode=true
	bin/export_puml.sh "$(darkmode)"
