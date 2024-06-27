package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfprotov6 "github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testAccProviderVersion = "0.0.1"
	testAccProviderType    = "schema_registry"

	providerConfig = `
	provider "schemaregistry" {
	  schema_registry_url = "http://test-url"
	  username            = "test-user"
	  password            = "test-pass"
	}
	`

	resourceConfig = `
	resource "schemaregistry_schema" "test" {
	  subject 			  = "test-subject"
	  schema  			  = "{\"type\":\"record\",\"name\":\"Test\",\"fields\":[{\"name\":\"f1\",\"type\":\"string\"}]}"
	  schema_type 		  = "avro"
	  compatibility_level = "NONE"
	}
	`
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"schemaregistry": providerserver.NewProtocol6WithError(New(testAccProviderVersion)()),
}

func TestAccSchemaRegistryProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "schema_registry_url", "https://schema-registry.kafka.usw2.dev-us.cultureamp.io"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "username", "test-user"),
					resource.TestCheckResourceAttr("schemaregistry_schema.test", "password", "test-pass"),
				),
			},
		},
	})
}
