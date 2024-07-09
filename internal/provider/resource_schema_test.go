package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource_CreateReadImportUpdate(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSchemaResourceConfig_single(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema", "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "avro"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "NONE"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
			// ImportState testing
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// No attributes to ignore during import
				ImportStateVerifyIgnore: []string{},
			},
			// Update and Read testing
			{
				Config: testAccSchemaResourceConfig_singleUpdate(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema", "{\"type\":\"record\",\"name\":\"TestUpdated\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "avro"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
		},
	})
}

func testAccSchemaResourceConfig_base() string {
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

func testAccSchemaResourceConfig_single(subject string) string {
	const singleTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "NONE"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(singleTemplate, subject))
}

func testAccSchemaResourceConfig_singleUpdate(subject string) string {
	const updateTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "BACKWARD"
  schema               = "{\"type\":\"record\",\"name\":\"TestUpdated\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(updateTemplate, subject))
}
