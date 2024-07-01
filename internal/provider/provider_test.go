package provider

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	tfprotov6 "github.com/hashicorp/terraform-plugin-go/tfprotov6"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

const (
	testAccProviderVersion = "test-version"
	testAccProviderType    = "schemaregistry"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"schemaregistry": providerserver.NewProtocol6WithError(New(testAccProviderVersion)()),
}

var compose tc.ComposeStack

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	compose, err = tc.NewDockerCompose("../../docker-compose.yaml")
	if err != nil {
		fmt.Printf("Testcontainers failed: %v\n", err)
		os.Exit(1)
	}

	// Bring up the stack
	if err := compose.Up(ctx, tc.Wait(true)); err != nil {
		fmt.Printf("Docker compose up failed: %v\n", err)
		os.Exit(1)
	}

	// Wait for services to be ready
	time.Sleep(30 * time.Second)

	os.Setenv("SCHEMA_REGISTRY_URL", "http://localhost:8081")

	code := m.Run()

	// Tear down the stack
	if err := compose.Down(ctx, tc.RemoveOrphans(true), tc.RemoveImagesLocal); err != nil {
		fmt.Printf("Docker compose down failed: %v\n", err)
		os.Exit(1)
	}

	os.Exit(code)
}
