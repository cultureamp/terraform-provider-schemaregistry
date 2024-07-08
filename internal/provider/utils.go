package provider

import (
	"os"
)

// getEnvOrDefault returns the value of the configuration or the environment variable.
func getEnvOrDefault(envVar, defaultValue string) string {
	if value, exists := os.LookupEnv(envVar); exists && value != "" {
		return value
	}
	return defaultValue
}
