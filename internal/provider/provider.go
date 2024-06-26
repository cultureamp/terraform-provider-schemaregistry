package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/riferrei/srclient"
)

// Ensure provider satisfies various expected interfaces.
var _ provider.Provider = &schemaRegistryProvider{}
var _ provider.ProviderWithFunctions = &schemaRegistryProvider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &schemaRegistryProvider{
			version: version,
		}
	}
}

// schemaRegistryProvider is the provider implementation.
type schemaRegistryProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel maps provider schema data to a Go type.
type schemaRegistryProviderModel struct {
	URL      types.String `tfsdk:"schema_registry_url"`
	Username types.String `tfsdk:"schema_registry_username"`
	Password types.String `tfsdk:"schema_registry_password"`
}

// Metadata returns the provider type name.
func (p *schemaRegistryProvider) Metadata(ctx context.Context, req provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "schemaRegistry"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *schemaRegistryProvider) Schema(ctx context.Context, req provider.SchemaRequest,
	resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_registry_url": schema.StringAttribute{
				Description: "URI for Schema Registry API. May use SCHEMA_REGISTRY_URL environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Schema Registry API. May use SCHEMA_REGISTRY_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Schema Registry API. May use SCHEMA_REGISTRY_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares an API client for data sources and resources.
func (p *schemaRegistryProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Schema Registry client")

	// Retrieve provider data from configuration
	var config schemaRegistryProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := getEnvOrDefault(config.URL, "SCHEMA_REGISTRY_URL")
	username := getEnvOrDefault(config.Username, "SCHEMA_REGISTRY_USERNAME")
	password := getEnvOrDefault(config.Password, "SCHEMA_REGISTRY_PASSWORD")

	ctx = tflog.SetField(ctx, "schema_registry_url", url)
	ctx = tflog.SetField(ctx, "schema_registry_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "schema_registry_password")

	tflog.Debug(ctx, "Creating Schema Registry client")

	if url == "" {
		resp.Diagnostics.AddError("Invalid credential parameters", "Schema "+
			"Registry URL must be provided")
		return
	}

	client := srclient.CreateSchemaRegistryClient(url)

	if username != "" && password != "" {
		client.SetCredentials(username, password)
	} else if username != "" || password != "" {
		resp.Diagnostics.AddError("Incomplete credentials", "Valid username and "+
			"password must be provided for basic authentication.")
		return
	}

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Schema Registry client", map[string]any{"success": true})
}

func getEnvOrDefault(configValue types.String, envVar string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	}
	return os.Getenv(envVar)
}

// Resources defines the resources implemented in the provider.
func (p *schemaRegistryProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//TODO SchemaResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *schemaRegistryProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//TODO SchemaDatasource,
	}
}

func (p *schemaRegistryProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		//TODO HelperFunction,
	}
}
