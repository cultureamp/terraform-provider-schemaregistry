<a href="https://terraform.io">
    <img src="https://www.svgrepo.com/show/376353/terraform.svg" alt="Terraform logo" title="Terraform" align="left" height="50" />
</a>

# Terraform Provider for Kafka Schema Registry

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

## Using the Provider

If you're building the provider, follow the instructions to [install it as a plugin](https://developer.hashicorp.com/terraform/cli/plugins#managing-plugin-installation). After placing it into your plugins directory, run `terraform init` to initialize it.

```hcl
terraform {
  required_providers {
    schemaregistry = {
      source = "../schemaregistry"
      version = "0.1.0"
    }
  }
}

provider "schemaregistry" {
  schema_registry_url = "https://schema-registry.example.com"
  username            = "your-username"
  password            = "your-password"
}

resource "schemaregistry_schema" "example" {
  subject              = "example"
  schema_type          = "avro"
  compatibility_level  = "NONE"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}
```
