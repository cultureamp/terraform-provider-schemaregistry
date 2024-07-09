package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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

func TestAccSchemaResource_DuplicateSubject(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.7.0"))),
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccSchemaResourceConfig_DuplicateSubject(subjectName),
				ExpectError: regexp.MustCompile("A schema with subject .* is already managed by another Terraform resource"),
			},
			{
				// Skip deleting resource which intentionally failed to create
				// https://github.com/hashicorp/terraform-plugin-testing/issues/85
				Config: testAccSchemaResourceConfig_DuplicateSubjectSkipDelete(resourceName),
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

func testAccSchemaResourceConfig_DuplicateSubject(subject string) string {
	const duplicateTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "BACKWARD"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}

resource "schemaregistry_schema" "test_02" {
  subject              = "%s"
  schema_type          = "avro"
  compatibility_level  = "BACKWARD"
  schema               = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(duplicateTemplate, subject, subject))
}

func testAccSchemaResourceConfig_DuplicateSubjectSkipDelete(resource string) string {
	const SkipDeleteTemplate = `
removed {
  from = %s

  lifecycle {
    destroy = false
  }
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(SkipDeleteTemplate, resource))
}
