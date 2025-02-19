.PHONY: format
format: ## Format the code
	opa fmt -w .

.PHONY: build-examples
build-examples: ## Build the example policies
	opa build -o .dist/example.bundle.tar.gz --optimize=1 -e examples ./examples
	opa inspect .dist/example.bundle.tar.gz

.PHONY: build
build: ## Build the policies
	opa build -o .dist/bundle.tar.gz ./policies
	opa inspect .dist/bundle.tar.gz

.PHONY: test
test: ## Test OPA policies
	if [ -z $(module) ]; then opa test -v ./policies; else opa test -v $(module); fi

.PHONY: test-deploy
test-deploy: ## Test if the deployment is correct
	curl -s -X "GET" http://localhost:8181/v1/policies | jq '.result[].id'

.PHONY: test-eval
test-eval: ## Test if OPA successfully evaluate a request
	@if [ $$(curl -s -o /dev/null -w "%{http_code}" -X "POST" http://localhost:8181/v1/data/examples -d @./deploy/local/example_input.json) != 200 ]; then \
		echo "Failed"; \
		exit 1; \
	else \
		echo "Success"; \
	fi
