<!-- markdownlint-disable MD033 MD041 -->
<a href="https://terraform.io">
    <img src="https://www.svgrepo.com/show/448253/terraform.svg" alt="Terraform logo" title="Terraform" align="left" height="70" />
</a>

# Terraform Provider for Schema Registry

 <p align="center">
    <a href="https://goreportcard.com/report/github.com/cultureamp/terraform-provider-schemaregistry">
      <img src="https://goreportcard.com/badge/github.com/cultureamp/terraform-provider-schemaregistry" alt="Go report badge">
    </a>
    <a href="https://github.com/cultureamp/terraform-provider-schemaregistry/releases/latest">
      <img src="https://img.shields.io/github/release/cultureamp/terraform-provider-schemaregistry.svg" alt="Release badge">
    </a>
    <a href="https://securityscorecards.dev/viewer/?uri=github.com/cultureamp/terraform-provider-schemaregistry">
      <img src="https://api.securityscorecards.dev/projects/github.com/cultureamp/terraform-provider-schemaregistry/badge" alt="OpenSSF Scorecard">
    </a>
</p>

This provider is built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.0+
- Go 1.16+ (if you plan to build the provider from the source)

## Building the Provider

If you want to build the provider from source, follow these steps:

1. Clone the repository
2. Run `make tidy` to install dependencies
3. Build the provider using `make build`
4. Run tests with `make testacc`

## Testing the Provider

The acceptance tests rely on [Testcontainers for Go (Redpanda)](https://golang.testcontainers.org/modules/redpanda/) to
provide a Schema Registry API.

This is protocol compliant but may have some limitations [[test-containers][01] / [redpanda][02]].

## Using the Provider

If you're building the provider, follow the instructions to
[install it as a plugin](https://developer.hashicorp.com/terraform/cli/plugins#managing-plugin-installation).
After placing it into your plugins directory, run `terraform init` to initialize it.

```hcl
terraform {
  required_providers {
    schemaregistry = {
      source = "cultureamp/schemaregistry"
      version = "1.5.0"
    }
  }
}

provider "schemaregistry" {
  schema_registry_url = "https://schemaregistry.example.com"
  username            = "your-username"
  password            = "your-password"
}

resource "schemaregistry_schema" "example" {
  subject              = "example"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = file("path/to/your/schema.avsc")

  # optional list of schema references
  references = [
    {
      name    = "ref-schema-name-1"
      subject = schemaregistry_schema.ref_schema_1.subject
      version = schemaregistry_schema.ref_schema_1.version
    },
  ]
}
```

[01]:https://github.com/testcontainers/testcontainers-go/issues?q=sort%3Aupdated-desc%20is%3Aissue%20redpanda
[02]:https://github.com/redpanda-data/redpanda/issues?q=sort%3Aupdated-desc%20is%3Aissue%20schema%20registry
