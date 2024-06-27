package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfprotov6 "github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories is a map of provider factories.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"schemaregistry": providerserver.NewProtocol6WithError(New("test")()),
}

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

	// datasourceConfig = `
	// data "schemaregistry_schema" "test" {
	//   name = "test-schema"
	// }
	// ` //TODO: add config.

	// resourceConfig = `
	// resource "schemaregistry_schema" "test" {
	//   name   = "test-schema"
	//   schema = "{\"type\": \"record\", \"name\": \"Test\", \"fields\": [{\"name\": \"test\", \"type\": \"string\"}]}"
	// }
	// ` //TODO: add config.
)

func TestAccSchemaRegistryProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig, //TODO: add + datasourceConfig + resourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testProvider(),
					testProvider_impl(),
					testProviderSchema(),
					testProviderMetadata(),
					// TODO: add checks.
					// resource.TestCheckResourceAttr("data.schemaregistry_schema.test", "schemas.0", "schema-1"),
					// resource.TestCheckResourceAttr("resource.schemaregistry_schema.test", "schemas.0", "schema-1"),
				),
			},
		},
	})
}

func getProvider() (provider.Provider, error) {
	providerFactory := testAccProtoV6ProviderFactories["schemaregistry"]
	server, err := providerFactory()
	if err != nil {
		return nil, err
	}
	p, ok := server.(provider.Provider)
	if !ok {
		return nil, fmt.Errorf("failed to get Provider server")
	}
	return p, nil
}

func testProvider() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		providerFactory := testAccProtoV6ProviderFactories["schemaregistry"]
		server, err := providerFactory()
		if err != nil {
			return err
		}

		p, ok := server.(provider.ProviderWithConfigValidators)
		if !ok {
			return fmt.Errorf("failed to get ProviderWithConfigValidators server")
		}

		// Perform provider validation
		req := provider.ValidateConfigRequest{}
		resp := provider.ValidateConfigResponse{}

		for _, validator := range p.ConfigValidators(ctx) {
			validator.ValidateProvider(ctx, req, &resp)
		}

		if resp.Diagnostics.HasError() {
			return fmt.Errorf("provider validation failed: %v", resp.Diagnostics)
		}
		return nil
	}
}

func testProvider_impl() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var _ provider.Provider = &Provider{}
		return nil
	}
}

func testProviderSchema() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		p, err := getProvider()
		if err != nil {
			return err
		}

		resp := provider.SchemaResponse{}
		p.Schema(ctx, provider.SchemaRequest{}, &resp)
		if len(resp.Schema.Attributes) == 0 {
			return fmt.Errorf("provider schema is empty")
		}
		if _, ok := resp.Schema.Attributes["schema_registry_url"]; !ok {
			return fmt.Errorf("schema_registry_url attribute is missing")
		}
		if _, ok := resp.Schema.Attributes["username"]; !ok {
			return fmt.Errorf("username attribute is missing")
		}
		if _, ok := resp.Schema.Attributes["password"]; !ok {
			return fmt.Errorf("password attribute is missing")
		}
		return nil
	}
}

func testProviderMetadata() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()
		p, err := getProvider()
		if err != nil {
			return err
		}

		resp := provider.MetadataResponse{}
		p.Metadata(ctx, provider.MetadataRequest{}, &resp)
		if resp.TypeName != "schemaregistry" {
			return fmt.Errorf("expected provider type to be %s, got %s", testAccProviderType, resp.TypeName)
		}
		if resp.Version != testAccProviderVersion {
			return fmt.Errorf("expected provider version to be %s, got %s", testAccProviderVersion, resp.Version)
		}
		return nil
	}
}
