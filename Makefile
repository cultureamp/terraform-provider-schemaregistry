DOCKER_COMPOSE_RUN := docker compose run --rm

default: help

.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m

.PHONY: build ## Build the provider for all supported architectures
build: build-amd64 build-arm64

.PHONY: build-amd64
build-amd64: ## Build for amd64
	GOARCH=amd64 GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]') go build -o terraform-provider-schemaregistry-amd64

.PHONY: build-arm64
build-arm64: ## Build for arm64
	GOARCH=arm64 GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]') go build -o terraform-provider-schemaregistry-arm64

.PHONY: tidy
tidy: ## Run go mod tidy
	go mod tidy

.PHONY: download
download: ## Run go mod download
	go mod download

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
