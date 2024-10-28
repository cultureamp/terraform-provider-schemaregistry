package provider

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfprotov6 "github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/testcontainers/testcontainers-go/modules/redpanda"
)

const (
	testAccProviderVersion = "test"
	testAccProviderType    = "schemaregistry"
	redpandaContainerImage = "docker.redpanda.com/redpandadata/redpanda:v24.1.15"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"schemaregistry": providerserver.NewProtocol6WithError(New(testAccProviderVersion)()),
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	redpandaContainer, err := redpanda.Run(ctx,
		redpandaContainerImage,
		redpanda.WithEnableSASL(),
		redpanda.WithEnableKafkaAuthorization(),
		redpanda.WithEnableWasmTransform(),
		redpanda.WithNewServiceAccount("superuser-1", "test"),
		redpanda.WithSuperusers("superuser-1"),
		redpanda.WithEnableSchemaRegistryHTTPBasicAuth(),
		redpanda.WithBootstrapConfig("schema_registry_normalize_on_startup", true),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	schemaRegistryURL, err := redpandaContainer.SchemaRegistryAddress(ctx)
	if err != nil {
		log.Fatalf("failed to get schema registry address: %s", err)
	}

	// Set environment variables for the tests to use
	os.Setenv("SCHEMA_REGISTRY_URL", schemaRegistryURL)
	os.Setenv("SCHEMA_REGISTRY_USERNAME", "superuser-1")
	os.Setenv("SCHEMA_REGISTRY_PASSWORD", "test")

	// Run the tests
	tests := m.Run()

	// Clean up the container
	if err := redpandaContainer.Terminate(ctx); err != nil {
		log.Fatalf("failed to terminate container: %s", err)
	}

	os.Exit(tests)
}
