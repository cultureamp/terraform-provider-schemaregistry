package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/riferrei/srclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &schemaResource{}
	_ resource.ResourceWithConfigure   = &schemaResource{}
	_ resource.ResourceWithImportState = &schemaResource{}
)

// NewSchemaResource is a helper function to simplify the provider implementation.
func NewSchemaResource() resource.Resource {
	return &schemaResource{}
}

// schemaResource is the resource implementation.
type schemaResource struct {
	client *srclient.SchemaRegistryClient
}

// schemaResourceModel describes the resource data model.
type schemaResourceModel struct {
	Subject            types.String `tfsdk:"subject"`
	Schema             types.String `tfsdk:"schema"`
	SchemaID           types.Int64  `tfsdk:"schema_id"`
	Version            types.Int64  `tfsdk:"version"`
	Reference          types.List   `tfsdk:"reference"`
	SchemaType         types.String `tfsdk:"schema_type"`
	CompatibilityLevel types.String `tfsdk:"compatibility_level"`
}

// Metadata returns the resource type name.
func (r *schemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

// Schema defines the schema for the resource.
func (r *schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a schema in the Schema Registry",
		Attributes: map[string]schema.Attribute{
			"subject": schema.StringAttribute{
				Description: "The subject related to the schema",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema string",
				Required:    true,
			},
			"schema_id": schema.Int64Attribute{
				Description: "The ID of the schema",
				Computed:    true,
			},
			"version": schema.Int64Attribute{
				Description: "The version of the schema",
				Computed:    true,
			},
			"reference": schema.ListNestedAttribute{
				Description: "The referenced schema list",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The referenced schema name",
							Required:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The referenced schema subject",
							Required:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The referenced schema version",
							Required:    true,
						},
					},
				},
			},
			"schema_type": schema.StringAttribute{
				Description: "The schema type",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"avro",
						"json",
						"protobuf",
					),
				},
			},
			"compatibility_level": schema.StringAttribute{
				Description: "The compatibility level of the schema",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"NONE",
						"BACKWARD",
						"BACKWARD_TRANSITIVE",
						"FORWARD",
						"FULL",
						"FULL_TRANSITIVE",
					),
				},
			},
		},
	}
}

func (r *schemaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan schemaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	schemaString := plan.Schema.ValueString()
	references := ToRegistryReferences(plan.Reference)
	schemaType := ToSchemaType(plan.SchemaType.ValueString())
	compatibilityLevel := ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())

	// Create new schema
	schema, err := r.client.CreateSchema(plan.Subject.ValueString(), schemaString, schemaType, references...)
	if err != nil {
		resp.Diagnostics.AddError("Error creating schema", "Could not create schema, unexpected error: "+err.Error())
		return
	}

	_, err = r.client.ChangeSubjectCompatibilityLevel(plan.Subject.ValueString(), compatibilityLevel)
	if err != nil {
		resp.Diagnostics.AddError("Error setting compatibility level", "Could not set compatibility level, "+
			" unexpected error: "+err.Error())
		return
	}

	// Update state with refreshed values
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Schema = types.StringValue(schema.Schema())
	plan.Reference = FromRegistryReferences(schema.References())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state schemaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch the latest schema from the registry
	schema, err := r.client.GetLatestSchema(state.Subject.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema",
			"Could not read schema: "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	state.Schema = types.StringValue(schema.Schema())
	state.SchemaID = types.Int64Value(int64(schema.ID()))
	state.Version = types.Int64Value(int64(schema.Version()))
	state.Reference = FromRegistryReferences(schema.References())

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan schemaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	schemaString := plan.Schema.ValueString()
	references := ToRegistryReferences(plan.Reference)
	schemaType := ToSchemaType(plan.SchemaType.ValueString())
	compatibilityLevel := ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())

	// Update existing schema
	schema, err := r.client.CreateSchema(plan.Subject.ValueString(), schemaString, schemaType, references...)
	if err != nil {
		if strings.Contains(err.Error(), "409") {
			resp.Diagnostics.AddError("Error updating schema", `Invalid "schema": incompatible`)
			return
		}
		resp.Diagnostics.AddError("Error updating schema", "Could not update schema, unexpected error: "+err.Error())
		return
	}

	_, err = r.client.ChangeSubjectCompatibilityLevel(plan.Subject.ValueString(), compatibilityLevel)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting compatibility level",
			"Could not set compatibility level, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state with refreshed values
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Schema = types.StringValue(schema.Schema())
	plan.Reference = FromRegistryReferences(schema.References())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state schemaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing schema
	err := r.client.DeleteSubject(state.Subject.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Schema",
			"Could not delete schema, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*srclient.SchemaRegistryClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			"Expected *srclient.SchemaRegistryClient",
		)
		return
	}

	r.client = client
}

func (r *schemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("subject"), req, resp)
}

func FromRegistryReferences(references []srclient.Reference) types.List {
	if len(references) == 0 {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":    types.StringType,
				"subject": types.StringType,
				"version": types.Int64Type,
			},
		})
	}

	var elems []attr.Value
	for _, reference := range references {
		objectValue, diags := types.ObjectValue(
			map[string]attr.Type{
				"name":    types.StringType,
				"subject": types.StringType,
				"version": types.Int64Type,
			},
			map[string]attr.Value{
				"name":    types.StringValue(reference.Name),
				"subject": types.StringValue(reference.Subject),
				"version": types.Int64Value(int64(reference.Version)),
			},
		)
		if diags.HasError() {
			continue // TODO: fix this
		}
		elems = append(elems, objectValue)
	}

	listValue, diags := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":    types.StringType,
				"subject": types.StringType,
				"version": types.Int64Type,
			},
		},
		elems,
	)
	if diags.HasError() {
		return types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":    types.StringType,
				"subject": types.StringType,
				"version": types.Int64Type,
			},
		})
	}

	return listValue
}

func ToRegistryReferences(references types.List) []srclient.Reference {
	if references.IsNull() || references.IsUnknown() {
		return nil
	}

	var refs []srclient.Reference
	for _, reference := range references.Elements() {
		r, ok := reference.(types.Object)
		if !ok {
			// TODO: fix this
			continue
		}
		attributes := r.Attributes()
		name, nameOk := attributes["name"].(types.String)
		subject, subjectOk := attributes["subject"].(types.String)
		version, versionOk := attributes["version"].(types.Int64)

		if !nameOk || !subjectOk || !versionOk {
			// TODO: fix this
			continue
		}

		refs = append(refs, srclient.Reference{
			Name:    name.ValueString(),
			Subject: subject.ValueString(),
			Version: int(version.ValueInt64()),
		})
	}

	return refs
}

func ToSchemaType(schemaType string) srclient.SchemaType {
	switch schemaType {
	case "json":
		return srclient.Json
	case "protobuf":
		return srclient.Protobuf
	default:
		return srclient.Avro
	}
}

func ToCompatibilityLevelType(compatibilityLevel string) srclient.CompatibilityLevel {
	switch compatibilityLevel {
	case "NONE":
		return srclient.None
	case "BACKWARD":
		return srclient.Backward
	case "BACKWARD_TRANSITIVE":
		return srclient.BackwardTransitive
	case "FORWARD":
		return srclient.Forward
	case "FULL":
		return srclient.Full
	case "FULL_TRANSITIVE":
		return srclient.FullTransitive
	default:
		return srclient.ForwardTransitive
	}
}
