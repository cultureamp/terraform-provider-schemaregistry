package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource_CreateReadImportUpdate(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSchemaRegistryConfig_base(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", rName),
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
				Config: testAccSchemaRegistryConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", rName),
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
func testAccSchemaRegistryConfig_base(rName string) string {
	return fmt.Sprintf(`
provider "schemaregistry" {
  schema_registry_url = "%s"
  username            = "%s"
  password            = "%s"
}

resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "NONE"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}
`, getEnvOrDefault("SCHEMA_REGISTRY_URL", "localhost:9092"),
		getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", "superuser-1"),
		getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", "test"),
		rName)
}

func testAccSchemaRegistryConfig_update(rName string) string {
	return fmt.Sprintf(`
provider "schemaregistry" {
  schema_registry_url = "%s"
  username            = "%s"
  password            = "%s"
}

resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "BACKWARD"
  schema               = "{\"type\":\"record\",\"name\":\"TestUpdated\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"
}
`, getEnvOrDefault("SCHEMA_REGISTRY_URL", "localhost:9092"),
		getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", "superuser-1"),
		getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", "test"),
		rName)
}
