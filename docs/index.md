---
layout: "schema registry"
page_title: "Provider: Kafka Schema Registry"
sidebar_current: "docs-schema registry-index"
description: |-
  The Kafka Schema Registry provider to interact with schemas
---

# Schema Registry Provider

The Schema Registry provider allows you to manage schema resources.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Schema Registry Provider
provider "schemaregistry" {
  schema_registry_url = "localhost:8081"
  username            = "example"
  password            = "example"
}

resource "schemaregistry_schema" "example" {
  subject             = "example"
  schema_type         = "avro"
  compatibility_level = "FORWARD_TRANSITIVE"
  schema              = "example"
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `schema_registry_url` - (Required) URI for Schema Registry API.
 May use `SCHEMA_REGISTRY_URL` environment variable.

* `username` - (Optional) Username for Schema Registry API.
 May use `SCHEMA_REGISTRY_USERNAME` environment variable.

* `password` - (Optional) Password for Schema Registry API.
 May use `SCHEMA_REGISTRY_PASSWORD` environment variable.
