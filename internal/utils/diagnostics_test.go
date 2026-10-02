package utils

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAppendSchemaDiagnosticsAddsKnownContext(t *testing.T) {
	t.Parallel()

	diagnosticPath := path.Root("references").AtListIndex(0).AtName("version")
	errorDiagnostic := diag.NewAttributeErrorDiagnostic(
		diagnosticPath,
		"Invalid reference",
		"Reference version could not be decoded.",
	)
	warningDiagnostic := diag.NewWarningDiagnostic(
		"Deprecated field",
		"This warning should remain unchanged.",
	)

	var got diag.Diagnostics
	AppendSchemaDiagnostics(
		&got,
		diag.Diagnostics{errorDiagnostic, warningDiagnostic},
		types.StringValue("orders-value"),
		types.Int64Value(7),
	)

	if len(got) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(got))
	}

	const wantDetail = "Reference version could not be decoded.\n\nSchema context: subject \"orders-value\", version 7."
	if got[0].Detail() != wantDetail {
		t.Errorf("unexpected error detail:\nwant: %q\n got: %q", wantDetail, got[0].Detail())
	}

	gotWithPath, ok := got[0].(diag.DiagnosticWithPath)
	if !ok {
		t.Fatal("expected error diagnostic to retain its attribute path")
	}
	if !gotWithPath.Path().Equal(diagnosticPath) {
		t.Errorf("unexpected diagnostic path: %s", gotWithPath.Path())
	}

	if !got[1].Equal(warningDiagnostic) {
		t.Error("expected warning diagnostic to remain unchanged")
	}
}

func TestAppendSchemaDiagnosticsOmitsUnknownContext(t *testing.T) {
	t.Parallel()

	errorDiagnostic := diag.NewErrorDiagnostic("Invalid plan", "Plan could not be decoded.")

	var got diag.Diagnostics
	AppendSchemaDiagnostics(
		&got,
		diag.Diagnostics{errorDiagnostic},
		types.StringUnknown(),
		types.Int64Null(),
	)

	if len(got) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(got))
	}
	if !got[0].Equal(errorDiagnostic) {
		t.Error("expected diagnostic to remain unchanged without known schema context")
	}
}

func TestAppendSchemaDiagnosticsUsesOnlyKnownContext(t *testing.T) {
	t.Parallel()

	errorDiagnostic := diag.NewErrorDiagnostic("Invalid plan", "Plan could not be decoded.")

	var got diag.Diagnostics
	AppendSchemaDiagnostics(
		&got,
		diag.Diagnostics{errorDiagnostic},
		types.StringValue("orders-value"),
		types.Int64Unknown(),
	)

	if len(got) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(got))
	}

	const wantDetail = "Plan could not be decoded.\n\nSchema context: subject \"orders-value\"."
	if got[0].Detail() != wantDetail {
		t.Errorf("unexpected error detail:\nwant: %q\n got: %q", wantDetail, got[0].Detail())
	}
}
