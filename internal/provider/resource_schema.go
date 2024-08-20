package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schema": schema.StringAttribute{
				Description: "The schema definition.",
				Optional:    true,
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
				Computed:    false,
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
	err := r.isSubjectManaged(subject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schema",
			fmt.Sprintf("Error checking if subject is managed: %s", err),
		)
		return
	}

	// Normalize the schema string
	schemaString := plan.Schema.ValueString()
	normalizedSchema, err := NormalizeJSON(schemaString, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Schema",
			fmt.Sprintf("Schema validation failed: %s", err),
		)
		return
	}

	// Generate API request body from plan
	schemaType := ToSchemaType(plan.SchemaType.ValueString())
	references := ToRegistryReferences(plan.Reference)

	// Create new schema resource
	schema, err := r.client.CreateSchema(subject, normalizedSchema, schemaType, references...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating schema",
			"Could not create schema, unexpected error: "+err.Error(),
		)
		return
	}

	if !plan.CompatibilityLevel.IsUnknown() && !plan.CompatibilityLevel.IsNull() {
		compatibilityLevel := ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())
		_, err = r.client.ChangeSubjectCompatibilityLevel(subject, compatibilityLevel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting compatibility level",
				"Could not set compatibility level, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Convert *srclient.SchemaType to string
	schemaTypeStr := FromSchemaType(schema.SchemaType())

	// Map response body to schema
	plan.ID = plan.Subject
	plan.Schema = jsontypes.NewNormalizedValue(schema.Schema())
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.SchemaType = types.StringValue(schemaTypeStr)
	plan.Version = types.Int64Value(1) // Set the version to 1 for new schema
	plan.Reference = FromRegistryReferences(schema.References())
	plan.HardDelete = types.BoolValue(plan.HardDelete.ValueBool())

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

	// Normalize the schema string
	schemaString := state.Schema.ValueString()
	normalizedSchema, err := NormalizeJSON(schemaString, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Schema",
			fmt.Sprintf("Schema validation failed: %s", err),
		)
		return
	}

	// Update state with refreshed values
	state.Schema = jsontypes.NewNormalizedValue(normalizedSchema)
	state.SchemaID = types.Int64Value(int64(schema.ID()))
	state.SchemaType = types.StringValue(schemaType)
	state.Version = types.Int64Value(int64(schema.Version()))
	state.Reference = FromRegistryReferences(schema.References())
	state.HardDelete = types.BoolValue(state.HardDelete.ValueBool())

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

	// Normalize the schema string
	schemaString := plan.Schema.ValueString()
	normalizedSchema, err := NormalizeJSON(schemaString, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Schema",
			fmt.Sprintf("Schema validation failed: %s", err),
		)
		return
	}

	// Generate API request body from plan
	subject := plan.Subject.ValueString()
	references := ToRegistryReferences(plan.Reference)
	schemaType := ToSchemaType(plan.SchemaType.ValueString())

	// Update existing schema
	schema, err := r.client.CreateSchema(subject, normalizedSchema, schemaType, references...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating schema",
			"Could not update schema, unexpected error: "+err.Error(),
		)
		return
	}

	if !plan.CompatibilityLevel.IsUnknown() && !plan.CompatibilityLevel.IsNull() {
		compatibilityLevel := ToCompatibilityLevelType(plan.CompatibilityLevel.ValueString())
		_, err = r.client.ChangeSubjectCompatibilityLevel(subject, compatibilityLevel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error setting compatibility level",
				"Could not set compatibility level, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Update state with refreshed values
	plan.Schema = jsontypes.NewNormalizedValue(normalizedSchema)
	plan.SchemaID = types.Int64Value(int64(schema.ID()))
	plan.Version = types.Int64Value(int64(schema.Version()))
	plan.Reference = FromRegistryReferences(schema.References())
	plan.HardDelete = types.BoolValue(plan.HardDelete.ValueBool())

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

	deletionType := state.HardDelete.ValueBool()

	// Delete existing schema
	err := r.client.DeleteSubject(state.Subject.ValueString(), deletionType)
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
			"Error Importing Schema",
			fmt.Sprintf("Could not retrieve compatibility level for subject %s: %s", subject,
				err.Error(),
			),
		)
		return
	}

	schemaType := FromSchemaType(schema.SchemaType())

	// Normalize the schema string
	normalizedSchema, err := NormalizeJSON(schema.Schema(), &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid JSON Schema",
			fmt.Sprintf("Schema validation failed: %s", err),
		)
		return
	}

	// Create state from retrieved schema
	state := schemaResourceModel{
		ID:                 types.StringValue(subject),
		Subject:            types.StringValue(subject),
		Schema:             jsontypes.NewNormalizedValue(normalizedSchema),
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

// isSubjectManaged prevents multiple Terraform resources from managing the same subject.
func (r *schemaResource) isSubjectManaged(subject string) error {
	// Fetch the list of subjects from the schema registry
	subjects, err := r.client.GetSubjects()
	if err != nil {
		return fmt.Errorf(
			"Failed to get subjects from schema registry: %w",
			err,
		)
	}

	// Check if the given subject is already managed in the schema registry
	for _, existingSubject := range subjects {
		if existingSubject == subject {
			return fmt.Errorf(
				"Subject %s already exists in the schema registry."+
					"Please import this resource into Terraform using `terraform import`.",
				subject,
			)
		}
	}

	return nil
}

func FromRegistryReferences(references []srclient.Reference) types.List {
	referenceType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name":    types.StringType,
			"subject": types.StringType,
			"version": types.Int64Type,
		},
	}

	if len(references) == 0 {
		return types.ListNull(referenceType)
	}

	var elems []attr.Value
	for _, reference := range references {
		objectValue, diags := types.ObjectValue(
			referenceType.AttrTypes,
			map[string]attr.Value{
				"name":    types.StringValue(reference.Name),
				"subject": types.StringValue(reference.Subject),
				"version": types.Int64Value(int64(reference.Version)),
			},
		)
		if diags.HasError() {
			fmt.Printf("Error converting reference to object value: %v\n", diags)
			continue
		}
		elems = append(elems, objectValue)
	}

	listValue, diags := types.ListValue(referenceType, elems)
	if diags.HasError() {
		return types.ListNull(referenceType)
	}

	return listValue
}

func ToRegistryReferences(references types.List) []srclient.Reference {
	if references.IsNull() || references.IsUnknown() {
		return nil
	}

	var refs []srclient.Reference
	for _, reference := range references.Elements() {
		if !reference.IsNull() && !reference.IsUnknown() {
			ref, ok := reference.(types.Object)
			if !ok {
				fmt.Printf("Invalid reference object type: %v\n", reference)
				continue
			}

			attributes := ref.Attributes()
			nameAttr, nameOk := attributes["name"].(types.String)
			subjectAttr, subjectOk := attributes["subject"].(types.String)
			versionAttr, versionOk := attributes["version"].(types.Int64)

			if !nameOk || !subjectOk || !versionOk {
				// Log error and skip this reference
				fmt.Printf("Error extracting attributes from reference object: %v\n", attributes)
				continue
			}

			refs = append(refs, srclient.Reference{
				Name:    nameAttr.ValueString(),
				Subject: subjectAttr.ValueString(),
				Version: int(versionAttr.ValueInt64()),
			})
		}
	}

	return refs
}

func FromSchemaType(schemaType *srclient.SchemaType) string {
	if schemaType == nil {
		return "AVRO"
	}
	return string(*schemaType)
}

func ToSchemaType(schemaType string) srclient.SchemaType {
	switch schemaType {
	case "AVRO":
		return srclient.Avro
	case "JSON":
		return srclient.Json
	case "PROTOBUF":
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
		return "undefined"
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
		return "undefined"
	}
}
