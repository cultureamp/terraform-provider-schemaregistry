package schematypes_test

import (
	"context"
	"testing"

	"github.com/cultureamp/terraform-provider-schemaregistry/internal/schematypes"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDefinitionSemanticEquality(t *testing.T) {
	tests := []struct {
		name, left, right string
		equal             bool
	}{
		{"protobuf unchanged", `syntax = "proto3";`, `syntax = "proto3";`, true},
		{"protobuf changed", `message A {}`, `message B {}`, false},
		{"protobuf formatting is left to registry", `message A {}`, `message A { }`, false},
		{"JSON formatting", `{"type":"record","name":"A"}`, `{ "name": "A", "type": "record" }`, true},
		{"JSON changed", `{"type":"record","name":"A"}`, `{"type":"record","name":"B"}`, false},
		{"JSON large numbers stay distinct", `{"n":9007199254740992}`, `{"n":9007199254740993}`, false},
		{"JSON and protobuf", `{"type":"string"}`, `message A {}`, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			left := schematypes.NewDefinitionValue(test.left)
			right := schematypes.NewDefinitionValue(test.right)
			equal, diagnostics := left.StringSemanticEquals(context.Background(), right)
			if diagnostics.HasError() || equal != test.equal {
				t.Fatalf("semantic equality = %v, diagnostics = %v; want %v", equal, diagnostics, test.equal)
			}
			if left.Equal(right) != (test.left == test.right) {
				t.Fatal("exact equality must remain distinct from semantic equality")
			}
		})
	}
}

func TestDefinitionReadsExistingStringState(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name  string
		value basetypes.StringValue
	}{
		{"protobuf", basetypes.NewStringValue(`syntax = "proto3"; message Test {}`)},
		{"JSON", basetypes.NewStringValue(`{ "type": "record", "name": "A" }`)},
		{"null", basetypes.NewStringNull()},
		{"unknown", basetypes.NewStringUnknown()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Existing jsontypes state is encoded as a Terraform string. The new
			// type must decode that same representation without a state migration.
			old := jsontypes.Normalized{StringValue: test.value}
			encoded, err := old.ToTerraformValue(ctx)
			if err != nil {
				t.Fatal(err)
			}
			decoded, err := (schematypes.DefinitionType{}).ValueFromTerraform(ctx, encoded)
			if err != nil {
				t.Fatal(err)
			}
			want := schematypes.Definition{StringValue: test.value}
			if !decoded.Equal(want) || !decoded.Type(ctx).Equal(schematypes.DefinitionType{}) {
				t.Fatalf("decoded = %v, want %v", decoded, want)
			}
			roundTrip, err := decoded.ToTerraformValue(ctx)
			if err != nil || !roundTrip.Equal(encoded) {
				t.Fatalf("string round trip failed: %v", err)
			}
			converted, diagnostics := (schematypes.DefinitionType{}).ValueFromString(ctx, test.value)
			if diagnostics.HasError() || !converted.Equal(want) {
				t.Fatalf("String conversion failed: %v", diagnostics)
			}
		})
	}
}

func TestDefinitionRejectsWrongValueType(t *testing.T) {
	ctx := context.Background()
	definitionType := schematypes.DefinitionType{}
	if _, err := definitionType.ValueFromTerraform(ctx, tftypes.NewValue(tftypes.Number, 1)); err == nil {
		t.Fatal("a numeric Terraform value should be rejected")
	}
	value := schematypes.NewDefinitionValue("message A {}")
	if value.Equal(basetypes.NewStringValue("message A {}")) || definitionType.Equal(basetypes.StringType{}) {
		t.Fatal("different custom types must not compare equal")
	}
	if _, diagnostics := value.StringSemanticEquals(ctx, basetypes.NewStringValue("message A {}")); !diagnostics.HasError() {
		t.Fatal("semantic equality must diagnose an unexpected value type")
	}
}
