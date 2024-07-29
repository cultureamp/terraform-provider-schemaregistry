package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

// NormalizeJSON normalizes a JSON string by unmarshaling it into a Go data structure and then marshaling it back to a
// JSON string. If diagnostics is provided, errors encountered during the process are added to it.
func NormalizeJSON(jsonString string, diagnostics *diag.Diagnostics) (string, error) {
	var jsonData interface{}
	err := json.Unmarshal([]byte(jsonString), &jsonData)
	if err != nil {
		diagnostics.AddError(
			"Invalid JSON",
			fmt.Sprintf("Error unmarshaling JSON: %s", err),
		)
		return "", err
	}

	normalizedBytes, err := json.Marshal(jsonData)
	if err != nil {
		diagnostics.AddError(
			"Normalization Error",
			fmt.Sprintf("Error marshaling JSON: %s", err),
		)
		return "", err
	}

	return string(normalizedBytes), nil
}

// NormalizedJSON is a JSON normalization helper function for tests.
func NormalizedJSON(jsonString string) string {
	normalized, err := NormalizeJSON(jsonString, nil)
	if err != nil {
		return ""
	}
	return normalized
}
