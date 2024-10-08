---
# generated by https://github.com/hashicorp/terraform-plugin-docs
layout: "schemaregistry"
page_title: "Schema Registry: schemaregistry_schema"
sidebar_current: "docs-schemaregistry-resource-schema"
description: |-
  Provides a Schema Registry Schema
---
# schemaregistry_schema Resource

Provides a Schema Registry Schema resource.

## Example Usage

```hcl
resource "schemaregistry_schema" "example" {
  subject             = "example"
  schema_type         = "AVRO"
  compatibility_level = "FORWARD_TRANSITIVE"
  hard_delete         = true
  schema              = "example"

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

## Argument Reference

The following arguments are supported:

* `subject` - (Required) The name of the subject under which the schema will be registered.
* `schema_type` - (Required) The schema format. Accepted values are: `AVRO`, `PROTOBUF` and `JSON`.
* `compatibility_level` - (Required) The schema compatibility level. Accepted values are `BACKWARD`,
`BACKWARD_TRANSITIVE`, `FORWARD`, `FORWARD_TRANSITIVE`, `FULL`, `FULL_TRANSITIVE` and `NONE`.
* `schema` - (Optional) The schema string.
* `references` - (Optional) The referenced schema list.
* `hard_delete` - (Optional) Controls whether the subject is soft or hard deleted. Must not be set when importing.
Defaults to `false` (soft delete).

## Attributes Reference

* `schema_id` - The ID of the schema.
* `version` - The schema version.

## Import

Schemas can be imported using their `subject` ID, e.g.

```sh
terraform import schemaregistry_schema.example subject
```
