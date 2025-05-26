package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/riferrei/srclient"
)

type refItem struct {
	Name    string `tfsdk:"name"`
	Subject string `tfsdk:"subject"`
	Version int64  `tfsdk:"version"`
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

// ToRegistryReferences turns a Terraform list of {name,subject,version} into SRclient refs.
func ToRegistryReferences(ctx context.Context, client *srclient.SchemaRegistryClient, in types.List) ([]srclient.Reference, diag.Diagnostics) {
	if in.IsNull() || in.IsUnknown() {
		return nil, nil
	}

	var items []refItem
	diags := in.ElementsAs(ctx, &items, false)
	if diags.HasError() {
		return nil, diags
	}

	out := make([]srclient.Reference, 0, len(items))
	for _, it := range items {
		vers := int(it.Version)

		if vers <= 0 {
			schema, err := client.GetLatestSchema(it.Subject)
			if err != nil {
				diags.AddError(
					"Error resolving reference version",
					fmt.Sprintf("Could not get latest schema for subject %s: %s", it.Subject, err.Error()),
				)
				return nil, diags
			}
			vers = schema.Version()
		}

		out = append(out, srclient.Reference{
			Name:    it.Name,
			Subject: it.Subject,
			Version: vers,
		})
	}
	return out, diags
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
		return "BACKWARD"
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
		return srclient.Backward
	}
}
