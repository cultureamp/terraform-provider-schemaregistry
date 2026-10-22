package schematypes_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cultureamp/terraform-provider-schemaregistry/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestProviderAcceptsRawProtobufSchema(t *testing.T) {
	ctx := context.Background()
	server, err := providerserver.NewProtocol6WithError(provider.New("test")())()
	if err != nil {
		t.Fatal(err)
	}
	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	resourceType := schemas.ResourceSchemas["schemaregistry_schema"].ValueType()
	config := dynamicObject(t, resourceType, map[string]tftypes.Value{
		"subject":     tftypes.NewValue(tftypes.String, "protobuf-test"),
		"schema_type": tftypes.NewValue(tftypes.String, "PROTOBUF"),
		"schema":      tftypes.NewValue(tftypes.String, `syntax = "proto3"; message Test { string value = 1; }`),
	})
	response, err := server.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
		TypeName: "schemaregistry_schema",
		Config:   config,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, diagnostic := range response.Diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Errorf("unexpected validation error: %s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
}

func TestProviderReadsAndImportsRawProtobufSchema(t *testing.T) {
	const definition = `syntax = "proto3"; message Test { string value = 1; }`
	server, schemas := configuredProvider(t, startRegistry(t, definition))
	ctx := context.Background()
	sourceType := schemas.DataSourceSchemas["schemaregistry_schema"].ValueType()
	read, err := server.ReadDataSource(ctx, &tfprotov6.ReadDataSourceRequest{
		TypeName: "schemaregistry_schema",
		Config: dynamicObject(t, sourceType, map[string]tftypes.Value{
			"subject": tftypes.NewValue(tftypes.String, "protobuf-test"),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	checkDiagnostics(t, read.Diagnostics)
	checkDefinition(t, read.State, sourceType, definition)
	imported, err := server.ImportResourceState(ctx, &tfprotov6.ImportResourceStateRequest{
		TypeName: "schemaregistry_schema", ID: "protobuf-test",
	})
	if err != nil {
		t.Fatal(err)
	}
	checkDiagnostics(t, imported.Diagnostics)
	if len(imported.ImportedResources) != 1 {
		t.Fatalf("expected one imported resource, got %d", len(imported.ImportedResources))
	}
	resourceType := schemas.ResourceSchemas["schemaregistry_schema"].ValueType()
	state := imported.ImportedResources[0].State
	checkDefinition(t, state, resourceType, definition)
	refreshed, err := server.ReadResource(ctx, &tfprotov6.ReadResourceRequest{
		TypeName: "schemaregistry_schema", CurrentState: state,
	})
	if err != nil {
		t.Fatal(err)
	}
	checkDiagnostics(t, refreshed.Diagnostics)
	checkDefinition(t, refreshed.NewState, resourceType, definition)
}

func startRegistry(t *testing.T, definition string) string {
	t.Helper()
	registry := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var result any
		switch r.URL.Path {
		case "/subjects/protobuf-test/versions/latest":
			result = map[string]any{"id": 1, "version": 1, "subject": "protobuf-test", "schemaType": "PROTOBUF", "schema": definition}
		case "/config/protobuf-test":
			result = map[string]any{"compatibilityLevel": "BACKWARD"}
		default:
			t.Errorf("unexpected registry request: %s %s", r.Method, r.URL)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(registry.Close)
	return registry.URL
}

func configuredProvider(t *testing.T, registryURL string) (tfprotov6.ProviderServer, *tfprotov6.GetProviderSchemaResponse) {
	t.Helper()
	// Provider environment variables take precedence over configuration.
	t.Setenv("SCHEMA_REGISTRY_URL", registryURL)
	t.Setenv("SCHEMA_REGISTRY_USERNAME", "")
	t.Setenv("SCHEMA_REGISTRY_PASSWORD", "")
	ctx := context.Background()
	server, err := providerserver.NewProtocol6WithError(provider.New("test")())()
	if err != nil {
		t.Fatal(err)
	}
	schemas, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		t.Fatal(err)
	}
	configured, err := server.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		Config: dynamicObject(t, schemas.Provider.ValueType(), map[string]tftypes.Value{
			"schema_registry_url": tftypes.NewValue(tftypes.String, registryURL),
		}),
	})
	if err != nil {
		t.Fatal(err)
	}
	checkDiagnostics(t, configured.Diagnostics)
	return server, schemas
}

func checkDiagnostics(t *testing.T, diagnostics []*tfprotov6.Diagnostic) {
	t.Helper()
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity == tfprotov6.DiagnosticSeverityError {
			t.Fatalf("unexpected provider error: %s: %s", diagnostic.Summary, diagnostic.Detail)
		}
	}
}

func checkDefinition(t *testing.T, state *tfprotov6.DynamicValue, stateType tftypes.Type, want string) {
	t.Helper()
	decoded, err := state.Unmarshal(stateType)
	if err != nil {
		t.Fatal(err)
	}
	var attributes map[string]tftypes.Value
	if err := decoded.As(&attributes); err != nil {
		t.Fatal(err)
	}
	var got string
	if err := attributes["schema"].As(&got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("schema = %q, want %q", got, want)
	}
}

func dynamicObject(t *testing.T, objectType tftypes.Type, supplied map[string]tftypes.Value) *tfprotov6.DynamicValue {
	t.Helper()
	object, ok := objectType.(tftypes.Object)
	if !ok {
		t.Fatalf("expected object type, got %T", objectType)
	}
	values := make(map[string]tftypes.Value, len(object.AttributeTypes))
	for name, attributeType := range object.AttributeTypes {
		values[name] = tftypes.NewValue(attributeType, nil)
	}
	for name, value := range supplied {
		values[name] = value
	}
	result, err := tfprotov6.NewDynamicValue(objectType, tftypes.NewValue(objectType, values))
	if err != nil {
		t.Fatal(err)
	}
	return &result
}
