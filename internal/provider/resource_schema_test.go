package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource_CreateReadImport(t *testing.T) {
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
`, getEnvOrDefault("SCHEMA_REGISTRY_URL", "localhost:8081"),
		getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", "test-user"),
		getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", "test-pass"),
		rName)
}
