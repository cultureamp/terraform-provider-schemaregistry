<!-- markdownlint-disable MD033 MD041 -->
<a href="https://terraform.io">
    <img src="https://www.svgrepo.com/show/448253/terraform.svg" alt="Terraform logo" title="Terraform" align="left" height="70" />
</a>

# Terraform Provider for Schema Registry

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

This has some limitations:

- Redpanda only supports `AVRO` and `PROTOBUF` encoding, cannot test `JSON` schemas
  [[1]](https://github.com/redpanda-data/redpanda/issues/6220)

## Using the Provider

If you're building the provider, follow the instructions to
[install it as a plugin](https://developer.hashicorp.com/terraform/cli/plugins#managing-plugin-installation).
After placing it into your plugins directory, run `terraform init` to initialize it.

```hcl
terraform {
  required_providers {
    schemaregistry = {
      source = "cultureamp/schemaregistry"
      version = "1.1.0"
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
