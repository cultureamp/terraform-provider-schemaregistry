DOCKER_COMPOSE_RUN := docker compose run --rm
VERSION ?= 1.0.0
PROVIDER_NAME = terraform-provider-schemaregistry
ARCHS = amd64 arm64 arm 386
PLATFORMS = linux darwin windows freebsd
UNSUPPORTED_COMBOS = darwin/arm darwin/386

default: help

.PHONY: all
all: ## Build and zip the provider for all supported platforms and architectures
all: $(foreach platform,$(PLATFORMS),$(foreach arch,$(ARCHS),build-zip-$(platform)-$(arch)))

define build_zip_target
.PHONY: build-zip-$(1)-$(2)
build-zip-$(1)-$(2):
	@echo "Building and zipping for GOOS=$(1) GOARCH=$(2)"
	@if [ "$(filter $(1)/$(2), $(UNSUPPORTED_COMBOS))" ]; then \
		echo "Skipping unsupported GOOS/GOARCH pair: $(1)/$(2)"; \
	else \
		GOOS=$(1) GOARCH=$(2) go build -o $(PROVIDER_NAME)_$(1)_$(2)$(if $(findstring windows,$(1)),.exe,); \
		zip $(PROVIDER_NAME)_$(VERSION)_$(1)_$(2).zip $(PROVIDER_NAME)_$(1)_$(2)$(if $(findstring windows,$(1)),.exe,); \
		rm $(PROVIDER_NAME)_$(1)_$(2)$(if $(findstring windows,$(1)),.exe,); \
	fi
endef

$(foreach platform,$(PLATFORMS),$(foreach arch,$(ARCHS),$(eval $(call build_zip_target,$(platform),$(arch)))))

.PHONY: build
build: ## Build the provider for the current architecture
	GOARCH=$(shell go env GOARCH) GOOS=$(shell go env GOOS) go build -o $(PROVIDER_NAME)_$(shell go env GOOS)_$(shell go env GOARCH)

.PHONY: testacc
testacc: ## Run acceptance tests
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m

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
	@echo "cultureamp/terraform-provider-schemaregistry"
	@echo "-------------------------------"
	@grep --no-filename -E '^[-a-zA-Z_:]+.*[^#]## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-28s\033[0m %s\n", $$1, $$2}' | \
		sort
	@echo "-------------------------------"
