resource "schemaregistry_schema" "test" {
  subject = "test-subject"
  schema  = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
  schema_type = "avro"
  compatibility_level = "FORWARD_TRANSITIVE"
}
