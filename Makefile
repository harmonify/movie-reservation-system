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
	echo "Sync Submodules"
	git submodule foreach git pull origin main
	echo "Installing tools..."
	echo "Installing Air"
	go install github.com/air-verse/air@latest
	echo "Installing Mockery"
	go install github.com/vektra/mockery/v2@v2.42.1
	echo "Installing Migration Tools"
	go install github.com/rubenv/sql-migrate/...@latest
	go install -tags 'postgres,mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	echo "Installing Go Env"
	go install github.com/joho/godotenv/cmd/godotenv@latest
	echo "Installing Pretty Logging"
	go install github.com/maoueh/zap-pretty/cmd/zap-pretty@latest
	echo "Installing Protoc"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.2
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	# echo "Installing Docker plugin"
	# docker plugin install grafana/loki-docker-driver:2.9.2 --alias loki --grant-all-permissions
	echo "Done..."

.PHONY: build_puml
build_puml: ## Generate PlantUML diagrams
	bin/export_puml.sh "$(darkmode)"

.PHONY: build_proto
build_proto: ## Generate protobuf code for Golang
	bin/gen_proto.sh

.PHONY: build_cli
build_cli: ## Build MRS CLI
	bin/build_cli.sh

.PHONY: deploy
deploy: ## Deploy infrastructure
	docker-compose up -d
