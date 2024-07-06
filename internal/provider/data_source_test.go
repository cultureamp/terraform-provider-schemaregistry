package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccSchemaDataSource_basic tests the basic functionality of the data source.
func TestAccSchemaDataSource_basic(t *testing.T) {
	datasourceName := "data.schemaregistry_schema.test"
	subjectName := acctest.RandomWithPrefix("tf-acc-test-subject")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccPreConfig(t) },
				Config:    testAccSchemaDataSourceConfig_basic(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(datasourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(datasourceName, "schema", `{"type":"record","name":"Test","fields":[{"name":"f1","type":"string"}]}`),
					resource.TestCheckResourceAttr(datasourceName, "schema_type", "avro"),
					resource.TestCheckResourceAttr(datasourceName, "compatibility_level", "NONE"),
				),
			},
		},
	})
}

// testAccSchemaDataSourceConfig_basic returns the configuration for testing the data source.
func testAccSchemaDataSourceConfig_basic(subject string) string {
	return fmt.Sprintf(`
provider "schemaregistry" {
  schema_registry_url = "%s"
  username            = "%s"
  password            = "%s"
}

resource "schemaregistry_schema" "test" {
  subject              = "%s"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
  schema_type          = "avro"
  compatibility_level  = "NONE"
}

data "schemaregistry_schema" "test" {
  subject = schemaregistry_schema.test.subject
}
`, getEnvOrDefault("SCHEMA_REGISTRY_URL", "localhost:9092"),
		getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", "superuser-1"),
		getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", "test"),
		subject)
}
