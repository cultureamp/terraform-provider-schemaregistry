terraform {
  required_providers {
    schemaregistry = {
      source = "local/schemaregistry"
      version = "0.1.0"
    }
  }
}

provider "schemaregistry" {
  schema_registry_url = "localhost:8081"
  username            = "test-user"
  password            = "test-pass"
}
