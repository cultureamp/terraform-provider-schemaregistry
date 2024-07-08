resource "schemaregistry_schema" "example" {
  subject             = "example-subject"
  schema              = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
  schema_type         = "avro"
  compatibility_level = "FORWARD_TRANSITIVE"
}
