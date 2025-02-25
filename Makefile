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

.PHONY: build-puml
build-puml: ## ex: make build-puml darkmode=true
	bin/export_puml.sh "$(darkmode)"

.PHONY: build-proto
build-proto: ## Generate protobuf code for Golang
	bin/gen_proto.sh

.PHONY: build-cli
build-cli: ## Build MRS CLI
	bin/build_cli.sh

.PHONY: deploy
deploy: ## Deploy infrastructure
	docker-compose up -d
