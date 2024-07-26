package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	initialSchema = `{
    "type": "record",
    "name": "Test",
    "fields": [
        {
            "name": "f1",
            "type": "string"
        }
    ]
}`

	updatedSchema = `{
    "type": "record",
    "name": "TestUpdated",
    "fields": [
        {
            "name": "f1",
            "type": "string"
        },
        {
            "name": "f2",
            "type": "int"
        }
    ]
}`
)

func TestAccSchemaResource_basic(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	ref01 := acctest.RandomWithPrefix("tf-acc-test-ref")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, Read, ImportState testing
			{
				Config: testAccSchemaResourceConfig_single(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema",
						jsontypes.NewNormalizedValue(initialSchema).ValueString()),
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
			// Update testing
			{
				Config: testAccSchemaResourceConfig_singleUpdate(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema",
						jsontypes.NewNormalizedValue(updatedSchema).ValueString()),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "avro"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
				),
			},
		},
	})
}

func TestAccSchemaResource_withReferences(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	ref01 := acctest.RandomWithPrefix("tf-acc-test-ref")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaResourceConfig_single(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify references attributes
					resource.TestCheckResourceAttr(resourceName, "references.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "references.0.name", "TestRef01"),
					resource.TestCheckResourceAttr(resourceName, "references.0.subject", ref01),
					resource.TestCheckResourceAttr(resourceName, "references.0.version", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
			{
				Config: testAccSchemaResourceConfig_singleUpdate(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated references attributes
					resource.TestCheckResourceAttr(resourceName, "references.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "references.0.name", "TestRefUpdated"),
					resource.TestCheckResourceAttr(resourceName, "references.0.subject", ref01),
					resource.TestCheckResourceAttr(resourceName, "references.0.version", "2"),
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

func testAccSchemaResourceConfig_single(subject, ref01 string) string {
	const singleTemplate = `
resource "schemaregistry_schema" "ref_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "NONE"
  schema               = <<EOF
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

resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "NONE"
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
references = [
    {
      name    = "TestRef01"
      subject = schemaregistry_schema.ref_01.subject
      version = schemaregistry_schema.ref_01.version
    },
  ]
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(singleTemplate, ref01, subject))
}

func testAccSchemaResourceConfig_singleUpdate(subject, ref01 string) string {
	const updateTemplate = `
resource "schemaregistry_schema" "ref_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "NONE"
  schema               = "{\"type\":\"record\",\"name\":\"TestRefUpdated\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}

resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "BACKWARD"
  schema               = "{\"type\":\"record\",\"name\":\"TestUpdated\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"},{\"name\":\"f2\",\"type\":\"int\"}]}"
  references = [
    {
      name    = "TestRefUpdated"
      subject = schemaregistry_schema.ref_01.subject
      version = 2
    }
  ]
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(updateTemplate, ref01, subject))
}
