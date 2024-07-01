default: build

# Run acceptance tests (these tests interact with the API through the terraform test framework)
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build the provider
.PHONY: build
build: build-amd64 build-arm64

# Build for amd64
.PHONY: build-amd64
build-amd64:
	GOARCH=amd64 GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]') go build -o terraform-provider-schemaregistry-amd64

# Build for arm64
.PHONY: build-arm64
build-arm64:
	GOARCH=arm64 GOOS=$(shell uname -s | tr '[:upper:]' '[:lower:]') go build -o terraform-provider-schemaregistry-arm64

# Tidy the module dependencies
.PHONY: tidy
tidy:
	go mod tidy

# Download the module dependencies
.PHONY: download
download:
	go mod download

# Clean the build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-schemaregistry

# Run go vet to check the code
.PHONY: vet
vet:
	go vet ./...

# Format the code
.PHONY: fmt
fmt:
	go fmt ./...

# Generate provider documentation for the Terraform registry
# Requires: go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
.PHONY: docs
docs:
	tfplugindocs
