package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "schemaregistry_schema" "test" {
  subject = "test-subject"
  schema  = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
  schema_type = "avro"
  compatibility_level = "FULL"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify schema attributes
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", "test-subject"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema_type", "json"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "compatibility_level", "FULL"),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "version"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "schemaregistry_schema.test",
				ImportState:       true,
				ImportStateVerify: true,
				// No attributes to ignore during import
				ImportStateVerifyIgnore: []string{},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "schemaregistry_schema" "test" {
  subject = "test-subject"
  schema  = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"
  schema_type = "avro"
  compatibility_level = "BACKWARD"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated schema attributes
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "subject", "test-subject"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema", "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema_type", "avro"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "schema_id"),
					resource.TestCheckResourceAttrSet("schemaregistry_schema.test", "version"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
