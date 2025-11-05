package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/cultureamp/terraform-provider-schemaregistry/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/riferrei/srclient"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &schemaResource{}
	_ resource.ResourceWithConfigure   = &schemaResource{}
	_ resource.ResourceWithImportState = &schemaResource{}
	_ resource.ResourceWithModifyPlan  = &schemaResource{}
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
	ID                 types.String         `tfsdk:"id"`
	Subject            types.String         `tfsdk:"subject"`
	Schema             jsontypes.Normalized `tfsdk:"schema"`
	SchemaID           types.Int64          `tfsdk:"schema_id"`
	SchemaType         types.String         `tfsdk:"schema_type"`
	Version            types.Int64          `tfsdk:"version"`
	Reference          types.List           `tfsdk:"references"`
	CompatibilityLevel types.String         `tfsdk:"compatibility_level"`
	HardDelete         types.Bool           `tfsdk:"hard_delete"`
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
				Description: "The globally unique ID of the schema.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject": schema.StringAttribute{
				Description: "The subject related to the schema.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(249),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9._-]+$`),
						"May only contain letters, digits, dots ('.'), underscores ('_') or hyphens ('-')",
					)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema definition.",
				Required:    true,
				CustomType:  jsontypes.NormalizedType{},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(2),
				},
			},
			"schema_id": schema.Int64Attribute{
				Description: "The ID of the schema.",
				Computed:    true,
			},
			"schema_type": schema.StringAttribute{
				Description: "The schema format.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"AVRO",
						"JSON",
						"PROTOBUF",
					),
				},
			},
			"version": schema.Int64Attribute{
				Description: "The version of the schema.",
				Computed:    true,
			},
			"references": schema.ListNestedAttribute{
				Description: "The referenced schema list.",
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The referenced schema name.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"subject": schema.StringAttribute{
							Description: "The referenced schema subject.",
							Required:    true,
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"version": schema.Int64Attribute{
							Description: "The referenced schema version.",
							Required:    true,
						},
					},
				},
			},
			"compatibility_level": schema.StringAttribute{
				Description: "The compatibility level of the schema.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
			"hard_delete": schema.BoolAttribute{
				Description: "Controls whether a schema should be soft or hard deleted.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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

func (r *schemaResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the state is null we assume it's a new resource
	// If the plan is null we assume it's a destroy operation, so don't modify plan
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Get current state
	var state schemaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get new plan
	var plan schemaResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	references, diags := utils.ToRegistryReferences(ctx, plan.Reference)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if the schemas are semantically equivalent
	equal, err := utils.IsSemanticallyEqual(
		ctx,
		r.client,
		state.Subject.ValueString(),
		plan.Schema.ValueString(),
		utils.ToSchemaType(plan.SchemaType.ValueString()),
		references,
	)
	if err != nil {
		// If the schema lookup fails (e.g., subject doesn't exist yet,
		// referenced schemas don't exist), we should not suppress the plan.
		// This is normal during creation or when dependencies are changing. We
		// only log this as debug information and continue without modifying the
		// plan.
		tflog.Debug(ctx, "Schema lookup failed during ModifyPlan, continuing without plan modification", map[string]interface{}{
			"subject": state.Subject.ValueString(),
			"error":   err.Error(),
		})
		return
	}

	// Copy only schema-related fields from state to plan to suppress a no-op
	// schema diff while preserving other fields from plan (compatibility_level,
	// hard_delete, etc.) so that changes to these fields are properly detected
	if equal {
		plan.Schema = state.Schema
		plan.SchemaID = state.SchemaID
		plan.Version = state.Version
		plan.Reference = state.Reference
		resp.Plan.Set(ctx, plan) // ignore diags: we only copy known-good values
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

	// Check if the subject is already managed in schema registry
	subject := plan.Subject.ValueString()
	err := utils.IsSubjectManaged(r.client, subject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schema",
			fmt.Sprintf("Error checking if subject is managed: %s", err),
		)
		return
	}

	// Generate API request body from plan
	schemaString := plan.Schema.ValueString()
	schemaType := utils.ToSchemaType(plan.SchemaType.ValueString())
	references, diags := utils.ToRegistryReferences(ctx, plan.Reference)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new schema resource
	schema, err := r.client.CreateSchema(subject, schemaString, schemaType, references...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schema",
			"Could not create schema, unexpected error: "+err.Error(),
		)
		return
	}

	// Set compatibility level if specified
	if !plan.CompatibilityLevel.IsNull() && !plan.CompatibilityLevel.IsUnknown() {
		compatibilityLevel := utils.ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())
		_, err = r.client.ChangeSubjectCompatibilityLevel(subject, compatibilityLevel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting compatibility level",
				"Could not set compatibility level, unexpected error: "+err.Error(),
			)
			return
		}
	} else {
		// Fetch the current compatibility level from the server
		compatibilityLevel, err := r.client.GetCompatibilityLevel(subject, true)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting compatibility level",
				"Could not get compatibility level, unexpected error: "+err.Error(),
			)
			return
		}
		plan.CompatibilityLevel = types.StringValue(utils.FromCompatibilityLevelType(*compatibilityLevel))
	}

	// Convert *srclient.SchemaType to string
	schemaTypeStr := utils.FromSchemaType(schema.SchemaType())

	// Map response body to schema
	plan.ID = types.StringValue(subject)
	plan.Schema = jsontypes.NewNormalizedValue(schema.Schema())
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.SchemaType = types.StringValue(schemaTypeStr)
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Reference = utils.FromRegistryReferences(schema.References())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemaResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure we have a valid subject before proceeding
	if state.Subject.IsNull() || state.Subject.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error Reading Schema",
			"Subject is null or unknown, cannot read schema",
		)
		return
	}

	subject := state.Subject.ValueString()

	// Fetch the latest schema from the registry
	schema, err := r.client.GetLatestSchema(subject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Schema",
			"Could not read schema: "+err.Error(),
		)
		return
	}

	// Fetch the current compatibility level from the server
	compatibilityLevel, err := r.client.GetCompatibilityLevel(subject, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting compatibility level",
			"Could not get compatibility level, unexpected error: "+err.Error(),
		)
		return
	}

	schemaString := schema.Schema()
	schemaID := int64(schema.ID())
	schemaVersion := int64(schema.Version())
	schemaType := utils.FromSchemaType(schema.SchemaType())
	references := utils.FromRegistryReferences(schema.References())
	compatString := utils.FromCompatibilityLevelType(*compatibilityLevel)

	// Update state with refreshed values
	state.Schema = jsontypes.NewNormalizedValue(schemaString)
	state.SchemaID = types.Int64Value(schemaID)
	state.SchemaType = types.StringValue(schemaType)
	state.Version = types.Int64Value(schemaVersion)
	state.Reference = references
	state.CompatibilityLevel = types.StringValue(compatString)

	// Ensure hard delete is set to false if not specified
	if state.HardDelete.IsNull() || state.HardDelete.IsUnknown() {
		state.HardDelete = types.BoolValue(false)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *schemaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemaResourceModel
	var state schemaResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read current state
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	subject := plan.Subject.ValueString()
	references, diags := utils.ToRegistryReferences(ctx, plan.Reference)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update or fetch the schema
	schema, err := r.updateSchema(ctx, subject, plan, state, references)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating schema",
			err.Error(),
		)
		return
	}

	// Update or fetch the compatibility level
	compatibilityLevel, err := r.updateCompatibilityLevel(subject, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating compatibility level",
			err.Error(),
		)
		return
	}

	// Update state with refreshed values
	plan.Schema = jsontypes.NewNormalizedValue(schema.Schema())
	plan.SchemaType = types.StringValue(utils.FromSchemaType(schema.SchemaType()))
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Reference = utils.FromRegistryReferences(schema.References())
	plan.CompatibilityLevel = types.StringValue(compatibilityLevel)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// updateSchema updates the schema if it changed, or fetches the current schema if not.
func (r *schemaResource) updateSchema(ctx context.Context, subject string, plan, state schemaResourceModel, references []srclient.Reference) (*srclient.Schema, error) {
	schemaString := plan.Schema.ValueString()
	schemaType := utils.ToSchemaType(plan.SchemaType.ValueString())

	// Schema hasn't changed, just fetch the current schema
	if plan.Schema.Equal(state.Schema) {
		schema, err := r.client.GetLatestSchema(subject)
		if err != nil {
			return nil, fmt.Errorf("could not fetch current schema: %w", err)
		}
		return schema, nil
	}

	// Schema string has changed, check semantic equivalence. This is a fallback
	// in case ModifyPlan is skipped or fails
	equal, err := utils.IsSemanticallyEqual(
		ctx,
		r.client,
		subject,
		schemaString,
		schemaType,
		references,
	)
	if err != nil {
		// If semantic check fails, assume they differ since we know
		// the strings already differ from above
		equal = false
	}

	// Schema has changed, update it
	if !equal {
		schema, err := r.client.CreateSchema(subject, schemaString, schemaType, references...)
		if err != nil {
			return nil, fmt.Errorf("could not update schema: %w", err)
		}
		return schema, nil
	}

	// Schemas are semantically equivalent, just fetch the current schema
	schema, err := r.client.GetLatestSchema(subject)
	if err != nil {
		return nil, fmt.Errorf("could not fetch current schema: %w", err)
	}
	return schema, nil
}

// updateCompatibilityLevel updates or fetches the compatibility level.
func (r *schemaResource) updateCompatibilityLevel(subject string, plan schemaResourceModel) (string, error) {
	if !plan.CompatibilityLevel.IsNull() && !plan.CompatibilityLevel.IsUnknown() {
		compatibilityLevel := utils.ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())
		_, err := r.client.ChangeSubjectCompatibilityLevel(subject, compatibilityLevel)
		if err != nil {
			return "", fmt.Errorf("could not set compatibility level: %w", err)
		}
		return plan.CompatibilityLevel.ValueString(), nil
	}

	// Fetch the global compatibility level from the server
	cl, err := r.client.GetCompatibilityLevel(subject, true)
	if err != nil {
		return "", fmt.Errorf("could not get compatibility level: %w", err)
	}
	return utils.FromCompatibilityLevelType(*cl), nil
}

func (r *schemaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemaResourceModel

	// Get current state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure we have a valid subject before proceeding
	if state.Subject.IsNull() || state.Subject.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error Deleting Schema",
			"Subject is null or unknown, cannot delete schema",
		)
		return
	}

	hardDelete := false
	if !state.HardDelete.IsNull() && !state.HardDelete.IsUnknown() {
		hardDelete = state.HardDelete.ValueBool()
	}

	// Delete existing schema
	err := r.client.DeleteSubject(state.Subject.ValueString(), hardDelete)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Schema",
			"Could not delete schema, unexpected error: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Info(ctx, fmt.Sprintf("Schema %s deleted (%s delete)", state.Subject.ValueString(),
		map[bool]string{true: "hard", false: "soft"}[hardDelete]))
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
			"Error Importing Schema",
			fmt.Sprintf("Could not retrieve compatibility level for subject %s: %s", subject,
				err.Error(),
			),
		)
		return
	}

	schemaType := utils.FromSchemaType(schema.SchemaType())

	// Create state from retrieved schema
	state := schemaResourceModel{
		ID:                 types.StringValue(subject),
		Subject:            types.StringValue(subject),
		Schema:             jsontypes.NewNormalizedValue(schema.Schema()),
		SchemaID:           types.Int64Value(int64(schema.ID())),
		SchemaType:         types.StringValue(schemaType),
		Version:            types.Int64Value(int64(schema.Version())),
		Reference:          utils.FromRegistryReferences(schema.References()),
		CompatibilityLevel: types.StringValue(utils.FromCompatibilityLevelType(*compatibilityLevel)),
		HardDelete:         types.BoolValue(false), // Default to false for imported resources
	}

	// Set the state
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
