provider "schemaregistry" {
  schema_registry_url = "localhost:8081"
  username            = "test-user"
  password            = "test-pass"
}

data "schemaregistry_schema" "example" {
  subject = "example-subject"
}

output "schema_id" {
  value = data.schemaregistry_schema.example.schema_id
}

output "schema" {
  value = data.schemaregistry_schema.example.schema
}
