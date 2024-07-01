package provider

import (
	"os"
	"testing"

	godotenv "github.com/joho/godotenv"
)

// getEnvOrDefault returns the value of the configuration or the environment variable.
func getEnvOrDefault(envVar, defaultValue string) string {
	if value, exists := os.LookupEnv(envVar); exists && value != "" {
		return value
	}
	return defaultValue
}

// testAccPreConfig verifies that the required environment variables are set.
func testAccPreConfig(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Fatalf("Error loading .env file: %s", err)
	}

	if v := getEnvOrDefault("SCHEMA_REGISTRY_URL", ""); v == "" {
		t.Fatal("SCHEMA_REGISTRY_URL must be set for acceptance tests")
	}
	if v := getEnvOrDefault("SCHEMA_REGISTRY_USERNAME", ""); v == "" {
		t.Fatal("SCHEMA_REGISTRY_USERNAME must be set for acceptance tests")
	}
	if v := getEnvOrDefault("SCHEMA_REGISTRY_PASSWORD", ""); v == "" {
		t.Fatal("SCHEMA_REGISTRY_PASSWORD must be set for acceptance tests")
	}
}
