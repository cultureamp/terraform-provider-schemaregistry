package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// getEnvOrDefault returns the value of the configuration or the environment variable.
func getEnvOrDefault(envVar, defaultValue string) string {
	if value, exists := os.LookupEnv(envVar); exists && value != "" {
		return value
	}
	return defaultValue
}

// ConfigCompose can be called to concatenate multiple strings to build test configurations.
func ConfigCompose(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}

// ValidateSchemaString compares two JSON schema strings for semantic equality.
// It returns an error if the schemas are not semantically equivalent or if
// validation fails due to parsing errors.
func ValidateSchemaString(expected, actual string) error {
	ctx := context.Background()

	expectedNormalized := jsontypes.NewNormalizedValue(expected)
	actualNormalized := jsontypes.NewNormalizedValue(actual)

	equal, diags := expectedNormalized.StringSemanticEquals(ctx, actualNormalized)

	// Check for validation errors first
	if diags.HasError() {
		return fmt.Errorf("schema validation failed: %w", formatDiagnostics(diags))
	}

	// Check for semantic inequality
	if !equal {
		return fmt.Errorf("schemas are not semantically equal:\nexpected: %s\nactual: %s", expected, actual)
	}

	return nil
}

// formatDiagnostics converts diagnostic errors into a single error message.
func formatDiagnostics(diags diag.Diagnostics) error {
	if !diags.HasError() {
		return nil
	}

	var messages []string
	for _, d := range diags.Errors() {
		if detail := d.Detail(); detail != "" {
			messages = append(messages, fmt.Sprintf("%s: %s", d.Summary(), detail))
		} else {
			messages = append(messages, d.Summary())
		}
	}

	return errors.New(strings.Join(messages, "; "))
}
