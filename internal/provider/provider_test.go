package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Schema Registry client is properly configured.
	// It is also possible to use the SCHEMA_REGISTRY_ environment variables instead,
	// such as updating the Makefile and running the testing through that tool.
	providerConfig = `
provider "schemaRegistry" {
  schema_registry_url = "http://localhost:8081"
  schema_registry_username = "user"
  schema_registry_password = "password"
}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"schemaRegistry": providerserver.NewProtocol6WithError(New("test")()),
	}
)
