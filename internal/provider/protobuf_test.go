package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchemaResource_protobuf(t *testing.T) {
	subject := acctest.RandomWithPrefix("tf-acc-protobuf")
	initial := `syntax = "proto3"; message Test { string value = 1; }`
	updated := `syntax = "proto3"; message Test { string value = 1; int32 count = 2; }`
	resourceName := "schemaregistry_schema.protobuf"
	dataSourceName := "data.schemaregistry_schema.protobuf"
	config := func(definition string) string {
		return ConfigCompose(testAccSchemaResourceConfig_base(), fmt.Sprintf(`
resource "schemaregistry_schema" "protobuf" {
  subject             = %q
  schema_type         = "PROTOBUF"
  compatibility_level = "NONE"
  schema              = %q
}
data "schemaregistry_schema" "protobuf" {
  subject = schemaregistry_schema.protobuf.subject
}
`, subject, definition))
	}
	check := func(definition string) resource.TestCheckFunc {
		return resource.ComposeAggregateTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "schema_type", "PROTOBUF"),
			resource.TestCheckResourceAttr(resourceName, "schema", definition),
			resource.TestCheckResourceAttr(dataSourceName, "schema_type", "PROTOBUF"),
			resource.TestCheckResourceAttr(dataSourceName, "schema", definition),
		)
	}
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(initial), Check: check(initial)},
			{Config: config(initial), PlanOnly: true},
			{Config: config(updated), Check: check(updated)},
			{ResourceName: resourceName, ImportState: true, ImportStateVerify: true},
		},
	})
}
