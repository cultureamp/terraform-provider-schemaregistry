package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/riferrei/srclient"
)

// Ensure provider satisfies various expected interfaces.
var _ provider.Provider = &Provider{}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}

// Provider is the provider implementation.
type Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ProviderModel maps provider schema data to a Go type.
type ProviderModel struct {
	URL      types.String `tfsdk:"schema_registry_url"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

const (
	schemaRegistryURLPattern = `^https?://.*$`
)

var schemaRegistryURLRegex = regexp.MustCompile(schemaRegistryURLPattern)

// Metadata returns the provider type name.
func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest,
	resp *provider.MetadataResponse) {
	resp.TypeName = "schemaregistry"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest,
	resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"schema_registry_url": schema.StringAttribute{
				Description: "URI for Schema Registry API. May use SCHEMA_REGISTRY_URL environment variable.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(schemaRegistryURLRegex,
						"Schema Registry URL must start with http or https",
					),
				},
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
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Schema Registry client")

	// Retrieve provider data from configuration
	var config ProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := getEnvOrDefault("SCHEMA_REGISTRY_URL", config.URL.ValueString())
	username := getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", config.Username.ValueString())
	password := getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", config.Password.ValueString())

	ctx = tflog.SetField(ctx, "schema_registry_url", url)
	ctx = tflog.SetField(ctx, "username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "password")

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

// Resources defines the resources implemented in the provider.
func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSchemaResource,
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSchemaDataSource,
	}
}
