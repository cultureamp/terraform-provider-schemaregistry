package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaDataSource_basic(t *testing.T) {
	datasourceName := "data.schemaregistry_schema.test_01"
	subjectName := acctest.RandomWithPrefix("tf-acc-test-subject")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaDataSourceConfig_single(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(datasourceName, "schema", initialSchema),
					resource.TestCheckResourceAttr(datasourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(datasourceName, "compatibility_level", "NONE"),
				),
			},
		},
	})
}

func TestAccSchemaDataSource_multipleVersions(t *testing.T) {
	datasourceName := "data.schemaregistry_schema.test_01"
	subjectName := acctest.RandomWithPrefix("tf-acc-test-subject")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create initial version of the schema
			{
				Config: testAccSchemaDataSourceConfig_single(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(datasourceName, "schema", initialSchema),
					resource.TestCheckResourceAttr(datasourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(datasourceName, "compatibility_level", "NONE"),
				),
			},
			// Update schema to a new version
			{
				Config: testAccSchemaDataSourceConfig_update(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(datasourceName, "schema", updatedSchema),
					resource.TestCheckResourceAttr(datasourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(datasourceName, "compatibility_level", "BACKWARD"),
				),
			},
			// Validate updated version of the schema
			{
				Config: testAccSchemaDataSourceConfig_specificVersion(subjectName, 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttr(datasourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(datasourceName, "schema", updatedSchema),
					resource.TestCheckResourceAttr(datasourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(datasourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttr(datasourceName, "version", "2"),
				),
			},
		},
	})
}

func testAccSchemaDataSourceConfig_base() string {
	const baseTemplate = `
provider "schemaregistry" {
  schema_registry_url = "%s"
  username            = "%s"
  password            = "%s"
}
`
	return fmt.Sprintf(baseTemplate,
		getEnvOrDefault("SCHEMA_REGISTRY_URL", "localhost:9092"),
		getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", "superuser-1"),
		getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", "test"),
	)
}

func testAccSchemaDataSourceConfig_single(subject string) string {
	const singleTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = jsonencode({
  "type": "record",
  "name": "Test",
  "fields": [
    {
      "name": "f1",
      "type": "string"
    }
  ]
})
}

data "schemaregistry_schema" "test_01" {
  subject = schemaregistry_schema.test_01.subject
}

output "schema" {
  value = data.schemaregistry_schema.test_01.schema
}
`
	return ConfigCompose(testAccSchemaDataSourceConfig_base(),
		fmt.Sprintf(singleTemplate, subject))
}

func testAccSchemaDataSourceConfig_update(subject string) string {
	const updateTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject             = "%s"
  schema_type         = "AVRO"
  compatibility_level = "BACKWARD"
  hard_delete         = false
  schema              = jsonencode({
    type = "record",
    name = "TestUpdated",
    fields = [
      {
        name = "f1",
        type = "string"
      },
      {
        name = "f2",
        type = "int"
      }
    ]
  })
}

data "schemaregistry_schema" "test_01" {
  subject = schemaregistry_schema.test_01.subject
}

output "schema" {
  value = data.schemaregistry_schema.test_01.schema
}
`
	return ConfigCompose(testAccSchemaDataSourceConfig_base(),
		fmt.Sprintf(updateTemplate, subject))
}

func testAccSchemaDataSourceConfig_specificVersion(subject string, version int) string {
	const specificVersionTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject             = "%s"
  schema_type         = "AVRO"
  compatibility_level = "BACKWARD"
  schema              = jsonencode({
    type = "record",
    name = "TestUpdated",
    fields = [
      {
        name = "f1",
        type = "string"
      },
      {
        name = "f2",
        type = "int"
      }
    ]
  })
}

data "schemaregistry_schema" "test_01" {
  subject = schemaregistry_schema.test_01.subject
  version = %d
}

output "schema" {
  value = data.schemaregistry_schema.test_01.schema
}
`
	return ConfigCompose(testAccSchemaDataSourceConfig_base(),
		fmt.Sprintf(specificVersionTemplate, subject, version))
}
