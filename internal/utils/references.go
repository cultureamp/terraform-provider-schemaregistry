package utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/riferrei/srclient"
)

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
