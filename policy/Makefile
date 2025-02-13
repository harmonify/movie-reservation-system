.PHONY: format
format: ## Format the code
	opa fmt -w .

.PHONY: build-examples
build-examples: ## Build the policy
	opa build --optimize=1 -e example -e example2 ./policies && opa inspect bundle.tar.gz

.PHONY: build
build: ## Build the policy
	opa build --optimize=1 ./policies

.PHONY: test
test: ## Test OPA policies
	if [ -z $(module) ]; then opa test -v .; else opa test -v $(module); fi

.PHONY: test-deployment
test-deployment: ## Test if the deployment is correct
	curl http://localhost:8181/v1/policies/policies/example/policy.rego &> /dev/null
	if [ $$? -ne 0 ]; then echo "Deployment failed"; exit 1; else echo "Deployment successful"; fi

.PHONY: test-eval
test-eval: ## Test if OPA successfully evaluate a request
	curl -X "POST" http://localhost:8181/v1/data/policies/example/policy -d @./example_input.json
