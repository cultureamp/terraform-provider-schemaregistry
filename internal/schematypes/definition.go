// Package schematypes provides Terraform values for Schema Registry definitions.
package schematypes

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ basetypes.StringTypable                    = DefinitionType{}
	_ basetypes.StringValuableWithSemanticEquals = Definition{}
)

// DefinitionType accepts both JSON definitions and raw Protobuf source text.
// Its Terraform representation remains a string, matching existing state.
type DefinitionType struct {
	basetypes.StringType
}

func (DefinitionType) String() string { return "schematypes.DefinitionType" }

func (DefinitionType) ValueType(context.Context) attr.Value { return Definition{} }

func (t DefinitionType) Equal(other attr.Type) bool {
	value, ok := other.(DefinitionType)
	return ok && t.StringType.Equal(value.StringType)
}

func (DefinitionType) ValueFromString(_ context.Context, value basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	return Definition{StringValue: value}, nil
}

func (t DefinitionType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	converted, err := t.StringType.ValueFromTerraform(ctx, value)
	if err != nil {
		return nil, err
	}
	stringValue, ok := converted.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected schema value type %T", converted)
	}
	return Definition{StringValue: stringValue}, nil
}

// Definition preserves JSON semantic equality without requiring every schema to
// be JSON. Protobuf syntax and semantic equivalence are checked by the registry.
type Definition struct {
	basetypes.StringValue
}

func (Definition) Type(context.Context) attr.Type { return DefinitionType{} }

func (v Definition) Equal(other attr.Value) bool {
	value, ok := other.(Definition)
	return ok && v.StringValue.Equal(value.StringValue)
}

func (v Definition) StringSemanticEquals(ctx context.Context, other basetypes.StringValuable) (bool, diag.Diagnostics) {
	value, ok := other.(Definition)
	if !ok {
		return false, diag.Diagnostics{diag.NewErrorDiagnostic(
			"Schema Semantic Equality Error", fmt.Sprintf("Expected schematypes.Definition, got %T", other),
		)}
	}
	if v.IsNull() || v.IsUnknown() || value.IsNull() || value.IsUnknown() {
		return v.Equal(value), nil
	}
	if !json.Valid([]byte(v.ValueString())) || !json.Valid([]byte(value.ValueString())) {
		return v.Equal(value), nil
	}
	return jsontypes.NewNormalizedValue(v.ValueString()).StringSemanticEquals(
		ctx, jsontypes.NewNormalizedValue(value.ValueString()),
	)
}

// NewDefinitionValue creates a known schema definition without altering its text.
func NewDefinitionValue(value string) Definition {
	return Definition{StringValue: basetypes.NewStringValue(value)}
}
