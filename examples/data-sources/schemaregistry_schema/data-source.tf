data "schemaregistry_schema" "example" {
  subject = "example-subject"
}

output "schema_id" {
  value = data.schemaregistry_schema.example.schema_id
}

output "schema" {
  value = data.schemaregistry_schema.example.schema
}
