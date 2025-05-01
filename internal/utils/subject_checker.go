package utils

import (
	"fmt"
	"strings"

	"github.com/riferrei/srclient"
)

// isSubjectManaged prevents multiple Terraform resources from managing the same subject.
func IsSubjectManaged(client *srclient.SchemaRegistryClient, subject string) error {
	// Fetch the subject-specific versions from the schema registry:
	//   GET /subjects/{subject}/versions
	versions, err := client.GetSchemaVersions(subject)
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			return nil
		}
		return fmt.Errorf("error checking existence of subject %q: %w", subject, err)
	}
	// If one or more subject versions exist, return an error
	if len(versions) > 0 {
		return fmt.Errorf(
			`subject %q already exists in the schema registry; please import it with:

		terraform import schemaregistry_schema.%[1]s %[1]s`,
			subject,
		)
	}

	return nil
}
