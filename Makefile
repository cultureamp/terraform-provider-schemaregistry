DOCKER_COMPOSE_RUN := docker compose run --rm

default: help

.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 TESTCONTAINERS_RYUK_DISABLED=true go test ./... -v $(TESTARGS) -timeout 120m

# Build the provider
.PHONY: build
build:
	go build -o terraform-provider-schemaregistry

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

# Clean the build artifacts
.PHONY: clean
clean: ## Clean up the build artifacts
	rm -f terraform-provider-schemaregistry

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: fmt
fmt: ## Run go fmt
	go fmt ./...

.PHONY: docs
docs: ## Generate provider documentation for the Terraform registry
	tfplugindocs

# Add double hash '##' plus the help text you would like to display for "make help" against that command
help: ### Show help for documented commands
	@echo "-------------------------------"
	@echo "cultureamp/terraform-provider-kafka-schema-registry"
	@echo "-------------------------------"
	@grep --no-filename -E '^[-a-zA-Z_:]+.*[^#]## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | \
		sort
	@echo "-------------------------------"
