package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
			"type": "int",
			"default": 0
		}
    ]
}`

	expectedSchema = `{"type":"record","name":"TestUpdated","fields":[{"name":"f1","type":"string"},{"name":"f2","type":"int","default":0}]}`
)

func TestAccSchemaResource_basic(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "schemaregistry_schema.test_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSchemaResourceConfig_basic(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "hard_delete", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(initialSchema, state)
					}),
				),
			},
			// Semantic No-Diff testing (PlanOnly)
			// Our ModifyPlan method should detect semantic equivalence and suppress the plan
			{
				Config:   testAccSchemaResourceConfig_basicNormalized(subjectName),
				PlanOnly: true,
			},
			// Update and Read testing
			{
				Config: testAccSchemaResourceConfig_basicUpdate(subjectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify updated schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceName, "hard_delete", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(updatedSchema, state)
					}),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 state, got %d", len(states))
					}
					state := states[0]
					if state.Attributes["subject"] != subjectName {
						return fmt.Errorf("expected subject %s, got %s", subjectName, state.Attributes["subject"])
					}
					err := ValidateSchemaString(expectedSchema, state.Attributes["schema"])
					if err != nil {
						return fmt.Errorf("schema validation error: %v", err)
					}
					return nil
				},
			},
		},
	})
}

func TestAccSchemaResource_withReferences(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test")
	ref01 := acctest.RandomWithPrefix("tf-acc-test-ref")
	resourceName := "schemaregistry_schema.test_01"
	refResourceName := "schemaregistry_schema.ref_01"

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSchemaResourceConfig_single(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify main schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "NONE"),
					resource.TestCheckResourceAttr(resourceName, "hard_delete", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					// Verify reference schema attributes
					resource.TestCheckResourceAttr(refResourceName, "subject", ref01),
					resource.TestCheckResourceAttrSet(refResourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(refResourceName, "version"),
					// Verify references attributes
					resource.TestCheckResourceAttr(resourceName, "references.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "references.0.name", "TestRef01"),
					resource.TestCheckResourceAttr(resourceName, "references.0.subject", ref01),
					resource.TestCheckResourceAttr(resourceName, "references.0.version", "1"),
				),
			},
			{
				Config: testAccSchemaResourceConfig_singleUpdate(subjectName, ref01),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify main schema attributes
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttr(resourceName, "hard_delete", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					// Verify reference schema updated attributes
					resource.TestCheckResourceAttr(refResourceName, "subject", ref01),
					resource.TestCheckResourceAttrSet(refResourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(refResourceName, "version"),
					// Verify updated references attributes
					resource.TestCheckResourceAttr(resourceName, "references.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "references.0.name", "TestRef01"),
					resource.TestCheckResourceAttr(resourceName, "references.0.subject", ref01),
					resource.TestCheckResourceAttr(resourceName, "references.0.version", "2"),
				),
			},
		},
	})
}

func TestAccSchemaResource_modifyPlan(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test-modifyplan")
	resourceName := "schemaregistry_schema.test_01"

	// Define test schema variations
	formattedSchema := `{
    "type": "record",
    "name": "ModifyPlanTest",
    "fields": [
        {
            "name": "id",
            "type": "string"
        },
        {
            "name": "timestamp",
            "type": "long"
        }
    ]
}`

	// Same schema but with different formatting - should trigger ModifyPlan
	compactSchema := `{"type":"record","name":"ModifyPlanTest","fields":[{"name":"id","type":"string"},{"name":"timestamp","type":"long"}]}`

	// Different schema - should not trigger ModifyPlan
	differentSchema := `{
    "type": "record",
    "name": "ModifyPlanTest",
    "fields": [
        {
            "name": "id",
            "type": "string"
        },
        {
            "name": "timestamp",
            "type": "long"
        },
        {
            "name": "newField",
            "type": "string",
            "default": "test"
        }
    ]
}`

	// Compact version of the different schema
	differentSchemaCompact := `{"type":"record","name":"ModifyPlanTest","fields":[{"name":"id","type":"string"},{"name":"timestamp","type":"long"},{"name":"newField","type":"string","default":"test"}]}`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create initial schema
			{
				Config: testAccSchemaResourceConfig_modifyPlanInitial(subjectName, formattedSchema),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(formattedSchema, state)
					}),
				),
			},
			// Step 2: Test ModifyPlan with semantically identical but differently formatted schema
			// This should result in NO plan changes due to ModifyPlan suppression
			{
				Config:   testAccSchemaResourceConfig_modifyPlanCompact(subjectName, compactSchema),
				PlanOnly: true,
			},
			// Step 3: Test with a truly different schema - should show plan changes
			{
				Config: testAccSchemaResourceConfig_modifyPlanDifferent(subjectName, differentSchema),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttrSet(resourceName, "schema_id"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(differentSchema, state)
					}),
				),
			},
			// Step 4: Test ModifyPlan again with the compact version of the different schema
			// This should result in NO plan changes due to ModifyPlan suppression
			{
				Config:   testAccSchemaResourceConfig_modifyPlanCompactDifferent(subjectName, differentSchemaCompact),
				PlanOnly: true,
			},
		},
	})
}

func TestAccSchemaResource_compatibilityLevelChange(t *testing.T) {
	subjectName := acctest.RandomWithPrefix("tf-acc-test-compatibility")
	resourceName := "schemaregistry_schema.test_01"

	schema := `{
    "type": "record",
    "name": "Test",
    "fields": [
        {
            "name": "f1",
            "type": "string"
        }
    ]
}`

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create schema with BACKWARD compatibility
			{
				Config: testAccSchemaResourceConfig_withCompatibility(subjectName, schema, "BACKWARD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "BACKWARD"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(schema, state)
					}),
				),
			},
			// Step 2: Change only compatibility_level to FULL (schema remains the same)
			// This should be detected as a change and applied
			{
				Config: testAccSchemaResourceConfig_withCompatibility(subjectName, schema, "FULL"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "FULL"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(schema, state)
					}),
				),
			},
			// Step 3: Change compatibility_level to FORWARD_TRANSITIVE
			{
				Config: testAccSchemaResourceConfig_withCompatibility(subjectName, schema, "FORWARD_TRANSITIVE"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "subject", subjectName),
					resource.TestCheckResourceAttr(resourceName, "schema_type", "AVRO"),
					resource.TestCheckResourceAttr(resourceName, "compatibility_level", "FORWARD_TRANSITIVE"),
					resource.TestCheckResourceAttrWith(resourceName, "schema", func(state string) error {
						return ValidateSchemaString(schema, state)
					}),
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
  schema_type          = "AVRO"
  schema               = <<EOF
{
  "type": "record",
  "name": "TestRef01",
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
references = [
    {
      name    = "TestRef01"
      subject = schemaregistry_schema.ref_01.subject
      version = 1 # schemaregistry_schema.ref_01.version
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
  schema_type          = "AVRO"
  schema               = jsonencode({
    "type": "record",
    "name": "TestRef01",
    "fields": [
      {
        "name": "f1",
        "type": "string"
      },
      {
        "name": "f2",
        "type": "int",
        "default": 0
      }
    ]
  })
}

resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "BACKWARD"
  hard_delete          = true
  schema               = jsonencode({
    "type": "record",
    "name": "TestUpdated",
    "fields": [
      {
        "name": "f1",
        "type": "string"
      },
      {
        "name": "f2",
        "type": "int",
		"default": 0
      },
    ]
  })
  references = [
    {
      name    = "TestRef01"
      subject = schemaregistry_schema.ref_01.subject
      version = 2 # schemaregistry_schema.ref_01.version
    }
  ]
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(updateTemplate, ref01, subject))
}

// testAccSchemaResourceConfig_basic creates a basic schema configuration for testing.
func testAccSchemaResourceConfig_basic(subject string) string {
	const basicTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(basicTemplate, subject, initialSchema))
}

// testAccSchemaResourceConfig_basicNormalized creates a semantically identical
// schema with different whitespace formatting.
func testAccSchemaResourceConfig_basicNormalized(subject string) string {
	const NormalizedTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(NormalizedTemplate, subject, NormalizeSchemaString(initialSchema)))
}

// testAccSchemaResourceConfig_basicUpdate creates an updated schema configuration.
func testAccSchemaResourceConfig_basicUpdate(subject string) string {
	const updateTemplate = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "BACKWARD"
  hard_delete          = true
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(updateTemplate, subject, updatedSchema))
}

// testAccSchemaResourceConfig_modifyPlanInitial creates a schema configuration for ModifyPlan testing.
func testAccSchemaResourceConfig_modifyPlanInitial(subject, schema string) string {
	const template = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(template, subject, schema))
}

// testAccSchemaResourceConfig_modifyPlanCompact creates a schema with compact formatting.
func testAccSchemaResourceConfig_modifyPlanCompact(subject, schema string) string {
	const template = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(template, subject, schema))
}

// testAccSchemaResourceConfig_modifyPlanDifferent creates a schema with different content.
func testAccSchemaResourceConfig_modifyPlanDifferent(subject, schema string) string {
	const template = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(template, subject, schema))
}

// testAccSchemaResourceConfig_modifyPlanCompactDifferent creates a compact schema with different content.
func testAccSchemaResourceConfig_modifyPlanCompactDifferent(subject, schema string) string {
	const template = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "NONE"
  hard_delete          = false
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(template, subject, schema))
}

// testAccSchemaResourceConfig_withCompatibility creates a schema configuration with a specific compatibility level.
func testAccSchemaResourceConfig_withCompatibility(subject, schema, compatibilityLevel string) string {
	const template = `
resource "schemaregistry_schema" "test_01" {
  subject              = "%s"
  schema_type          = "AVRO"
  compatibility_level  = "%s"
  schema               = <<EOF
%s
EOF
}
`
	return ConfigCompose(testAccSchemaResourceConfig_base(),
		fmt.Sprintf(template, subject, compatibilityLevel, schema))
}
