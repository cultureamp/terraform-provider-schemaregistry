package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/riferrei/srclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &schemaDataSource{}
	_ datasource.DataSourceWithConfigure = &schemaDataSource{}
)

// NewSchemaDataSource is a helper function to simplify the provider implementation.
func NewSchemaDataSource() datasource.DataSource {
	return &schemaDataSource{}
}

// schemaDataSource is the data source implementation.
type schemaDataSource struct {
	client *srclient.SchemaRegistryClient
}

// schemaDataSourceModel describes the data source data model.
type schemaDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Subject            types.String `tfsdk:"subject"`
	Schema             types.String `tfsdk:"schema"`
	SchemaID           types.Int64  `tfsdk:"schema_id"`
	SchemaType         types.String `tfsdk:"schema_type"`
	Version            types.Int64  `tfsdk:"version"`
	Reference          types.List   `tfsdk:"reference"`
	CompatibilityLevel types.String `tfsdk:"compatibility_level"`
}

// Metadata returns the data source type name.
func (d *schemaDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

// Schema defines the schema for the data source.
func (d *schemaDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema data source. Fetches a schema from the Schema Registry.",
		Description:         "Fetches a schema from the Schema Registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "UID for the schema, which is the subject name.",
				Computed:    true,
			},
			"subject": schema.StringAttribute{
				Description: "The subject related to the schema.",
				Required:    true,
			},
			"schema": schema.StringAttribute{
				Description: "The schema string.",
				Computed:    true,
			},
			"schema_id": schema.Int64Attribute{
				Description: "The ID of the schema.",
				Computed:    true,
			},
			"schema_type": schema.StringAttribute{
				Description: "The schema format.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"avro",
						"json",
						"protobuf",
					),
				},
			},
			"version": schema.Int64Attribute{
				Description: "The version of the schema.",
				Optional:    true,
			},
			"reference": schema.ListNestedAttribute{
				Description: "The referenced schema list.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The referenced schema name.",
							Computed:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The referenced schema subject.",
							Computed:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The referenced schema version.",
							Computed:    true,
						},
					},
				},
			},
			"compatibility_level": schema.StringAttribute{
				Description: "The compatibility level of the schema. Default is FORWARD_TRANSITIVE.",
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"NONE",
						"BACKWARD",
						"BACKWARD_TRANSITIVE",
						"FORWARD",
						"FORWARD_TRANSITIVE",
						"FULL",
						"FULL_TRANSITIVE",
					),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *schemaDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*srclient.SchemaRegistryClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *srclient.SchemaRegistryClient, got: %T.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Read fetches the schema details from the Schema Registry.
func (d *schemaDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get input values
	var inputs schemaDataSourceModel
	diags := req.Config.Get(ctx, &inputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	subject := inputs.Subject.ValueString()
	version := inputs.Version.ValueInt64()

	// Fetch schema and compatibility level
	schema, err := d.fetchSchema(subject, version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema",
			"Could not read schema: "+err.Error())
		return
	}

	compatibilityLevel, err := d.fetchCompatibilityLevel(subject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Compatibility Level",
			"Could not read compatibility level: "+err.Error(),
		)
		return
	}

	// Map response body to schema data source model
	outputs := mapSchemaToOutputs(subject, schema, compatibilityLevel)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, outputs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// fetchSchema retrieves the schema from the registry, either by specific version or the latest version.
func (d *schemaDataSource) fetchSchema(subject string, version int64) (*srclient.Schema, error) {
	if version > 0 {
		return d.client.GetSchemaByVersion(subject, int(version))
	}
	return d.client.GetLatestSchema(subject)
}

// fetchCompatibilityLevel retrieves the compatibility level for the given subject.
func (d *schemaDataSource) fetchCompatibilityLevel(subject string) (*srclient.CompatibilityLevel, error) {
	return d.client.GetCompatibilityLevel(subject, true)
}

// mapSchemaToOutputs maps the schema and compatibility level to the schema data source model.
func mapSchemaToOutputs(subject string, schema *srclient.Schema, compatibilityLevel *srclient.CompatibilityLevel) schemaDataSourceModel {
	return schemaDataSourceModel{
		ID:                 types.StringValue(subject),
		Subject:            types.StringValue(subject),
		Schema:             types.StringValue(schema.Schema()),
		SchemaID:           types.Int64Value(int64(schema.ID())),
		SchemaType:         types.StringValue(FromSchemaType(schema.SchemaType())),
		Version:            types.Int64Value(int64(schema.Version())),
		Reference:          FromRegistryReferences(schema.References()),
		CompatibilityLevel: types.StringValue(FromCompatibilityLevelType(*compatibilityLevel)),
	}
}
