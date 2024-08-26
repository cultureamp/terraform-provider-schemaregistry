package provider

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
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

func ValidateSchemaString(expected, actual string) error {
	ctx := context.Background()

	expectedSchemaString := jsontypes.NewNormalizedValue(expected)
	actualSchemaString := jsontypes.NewNormalizedValue(actual)

	equal, diags := expectedSchemaString.StringSemanticEquals(ctx, actualSchemaString)

	if !equal {
		return fmt.Errorf("input schema does not match output. Expected: %s, Actual: %s", expected, actual)
	}

	if diags.HasError() {
		var sb strings.Builder
		for _, d := range diags.Errors() {
			sb.WriteString(d.Summary() + "\n")
			sb.WriteString(d.Detail() + "\n")
		}
		return errors.New(sb.String())
	}

	return nil
}
