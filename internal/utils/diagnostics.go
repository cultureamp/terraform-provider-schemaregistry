package utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppendSchemaDiagnostics adds known schema context to error diagnostics while
// preserving warnings and attribute paths.
func AppendSchemaDiagnostics(
	target *diag.Diagnostics,
	diagnostics diag.Diagnostics,
	subject types.String,
	version types.Int64,
) {
	contextParts := make([]string, 0, 2)
	if !subject.IsNull() && !subject.IsUnknown() {
		contextParts = append(contextParts, fmt.Sprintf("subject %q", subject.ValueString()))
	}
	if !version.IsNull() && !version.IsUnknown() {
		contextParts = append(contextParts, fmt.Sprintf("version %d", version.ValueInt64()))
	}

	if len(contextParts) == 0 {
		target.Append(diagnostics...)
		return
	}

	contextDetail := strings.Join(contextParts, ", ")
	for _, diagnostic := range diagnostics {
		if diagnostic.Severity() != diag.SeverityError {
			target.Append(diagnostic)
			continue
		}

		detail := fmt.Sprintf("%s\n\nSchema context: %s.", diagnostic.Detail(), contextDetail)
		if diagnosticWithPath, ok := diagnostic.(diag.DiagnosticWithPath); ok {
			target.Append(diag.NewAttributeErrorDiagnostic(
				diagnosticWithPath.Path(),
				diagnostic.Summary(),
				detail,
			))
			continue
		}

		target.Append(diag.NewErrorDiagnostic(diagnostic.Summary(), detail))
	}
}
