package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	ID                 types.String `tfsdk:"id"`
	Subject            types.String `tfsdk:"subject"`
	Schema             types.String `tfsdk:"schema"`
	SchemaID           types.Int64  `tfsdk:"schema_id"`
	SchemaType         types.String `tfsdk:"schema_type"`
	Version            types.Int64  `tfsdk:"version"`
	Reference          types.List   `tfsdk:"reference"`
	CompatibilityLevel types.String `tfsdk:"compatibility_level"`
}

// Metadata returns the resource type name.
func (r *schemaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_schema"
}

// Schema defines the schema for the resource.
func (r *schemaResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Schema resource. Manages a schema in the Schema Registry.",
		Description:         "Manages a schema in the Schema Registry.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the schema, which is the subject name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject": schema.StringAttribute{
				Description: "The subject related to the schema.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema string.",
				Required:    true,
			},
			"schema_id": schema.Int64Attribute{
				Description: "The ID of the schema.",
				Computed:    true,
			},
			"schema_type": schema.StringAttribute{
				Description: "The schema type. Default is avro.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Computed:    true,
			},
			"reference": schema.ListNestedAttribute{
				Description: "The referenced schema list.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The referenced schema name.",
							Required:    true,
						},
						"subject": schema.StringAttribute{
							Description: "The referenced schema subject.",
							Required:    true,
						},
						"version": schema.Int64Attribute{
							Description: "The referenced schema version.",
							Required:    true,
						},
					},
				},
			},
			"compatibility_level": schema.StringAttribute{
				Description: "The compatibility level of the schema. Default is FORWARD_TRANSITIVE.",
				Optional:    true,
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

// Configure adds the provider configured client to the resource.
func (r *schemaResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
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
	schemaType := ToSchemaType(plan.SchemaType.ValueString())
	references := ToRegistryReferences(plan.Reference)
	compatibilityLevel := ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())

	// Create new schema resource
	schema, err := r.client.CreateSchema(plan.Subject.ValueString(), schemaString, schemaType, references...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schema",
			"Could not create schema, unexpected error: "+err.Error(),
		)
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

	//convert *srclient.SchemaType to string
	schemaTypeStr := FromSchemaType(schema.SchemaType())

	// Map response body to schema
	plan.ID = plan.Subject
	plan.Schema = types.StringValue(schema.Schema())
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.SchemaType = types.StringValue(schemaTypeStr)
	plan.Version = types.Int64Value(1) // Set the version to 1 for new schema
	plan.Reference = FromRegistryReferences(schema.References())

	// Set state to fully populated data
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

	schemaType := FromSchemaType(schema.SchemaType())

	// Update state with refreshed values
	state.Schema = types.StringValue(schema.Schema())
	state.SchemaID = types.Int64Value(int64(schema.ID()))
	state.SchemaType = types.StringValue(schemaType)
	state.Version = types.Int64Value(int64(schema.Version()))
	state.Reference = FromRegistryReferences(schema.References())

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemaResourceModel

	// Read Terraform plan data into the model
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
		resp.Diagnostics.AddError(
			"Error updating schema",
			"Could not update schema, unexpected error: "+err.Error(),
		)
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
	plan.Schema = types.StringValue(schema.Schema())
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Reference = FromRegistryReferences(schema.References())

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemaResourceModel

	// Get current state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Soft deletes existing schema
	err := r.client.DeleteSubject(state.Subject.ValueString(), false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Schema",
			"Could not delete schema, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *schemaResource) ImportState(ctx context.Context, req resource.ImportStateRequest,
	resp *resource.ImportStateResponse) {
	subject := req.ID

	// Retrieve the latest schema for the subject
	schema, err := r.client.GetLatestSchema(subject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Schema",
			fmt.Sprintf("Could not retrieve schema for subject %s: %s", subject,
				err.Error(),
			),
		)
		return
	}

	// Retrieve the compatibility level for the subject
	compatibilityLevel, err := r.client.GetCompatibilityLevel(subject, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Compatibility Level",
			fmt.Sprintf("Could not retrieve compatibility level for subject %s: %s", subject,
				err.Error(),
			),
		)
		return
	}

	schemaType := FromSchemaType(schema.SchemaType())

	// Create state from retrieved schema
	state := schemaResourceModel{
		ID:                 types.StringValue(subject),
		Subject:            types.StringValue(subject),
		Schema:             types.StringValue(schema.Schema()),
		SchemaID:           types.Int64Value(int64(schema.ID())),
		SchemaType:         types.StringValue(schemaType),
		Version:            types.Int64Value(int64(schema.Version())),
		Reference:          FromRegistryReferences(schema.References()),
		CompatibilityLevel: types.StringValue(FromCompatibilityLevelType(*compatibilityLevel)),
	}

	// Set the state
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
			// Log error and skip this reference
			fmt.Printf("Error converting reference to object value: %v\n", diags)
			continue
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
		// Log error and return null list
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
			// Log error and skip this reference
			fmt.Printf("Error converting reference to object value: %v\n", reference)
			continue
		}
		attributes := r.Attributes()
		name, nameOk := attributes["name"].(types.String)
		subject, subjectOk := attributes["subject"].(types.String)
		version, versionOk := attributes["version"].(types.Int64)

		if !nameOk || !subjectOk || !versionOk {
			// Log error and skip this reference
			fmt.Printf("Error extracting attributes from reference object: %v\n", attributes)
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

func FromSchemaType(schemaType *srclient.SchemaType) string {
	if schemaType == nil {
		return "avro"
	}
	return string(*schemaType)
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

func FromCompatibilityLevelType(compatibilityLevel srclient.CompatibilityLevel) string {
	switch compatibilityLevel {
	case srclient.None:
		return "NONE"
	case srclient.Backward:
		return "BACKWARD"
	case srclient.BackwardTransitive:
		return "BACKWARD_TRANSITIVE"
	case srclient.Forward:
		return "FORWARD"
	case srclient.ForwardTransitive:
		return "FORWARD_TRANSITIVE"
	case srclient.Full:
		return "FULL"
	case srclient.FullTransitive:
		return "FULL_TRANSITIVE"
	default:
		return "FORWARD_TRANSITIVE"
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
	case "FORWARD_TRANSITIVE":
		return srclient.ForwardTransitive
	case "FULL":
		return srclient.Full
	case "FULL_TRANSITIVE":
		return srclient.FullTransitive
	default:
		return srclient.ForwardTransitive
	}
}
