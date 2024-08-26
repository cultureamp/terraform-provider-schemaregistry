resource "schemaregistry_schema" "ref_01" {
  subject             = "%s"
  schema_type         = "AVRO"
  compatibility_level = "NONE"
  schema              = <<EOF
{
  "type": "record",
  "name": "Test",
  "fields": [
    {
      "name": "f1",
      "type": "string"
    }
  ]
}
EOF
}

resource "schemaregistry_schema" "example_01" {
  subject             = "%s"
  schema_type         = "AVRO"
  compatibility_level = "NONE"
  hard_delete         = true
  schema = jsonencode({
    "type" : "record",
    "name" : "Example",
    "fields" : [
      {
        "name" : "f1",
        "type" : "string"
      }
    ]
  })
  references = [
    {
      name    = "TestRef01"
      subject = schemaregistry_schema.ref_01.subject
      version = schemaregistry_schema.ref_01.version
    },
  ]
}
