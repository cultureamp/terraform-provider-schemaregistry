package utils

import (
	"context"
	"errors"

	"github.com/riferrei/srclient"
)

// IsSemanticallyEqual checks if a given schema is semantically equivalent to
// any existing schema under the specified subject in the Schema Registry.  It
// returns true when the lookup succeeds, false when the registry replies 40403
// (ErrSemanticSchemaNotFound), and an error for anything else.
//
// The function uses the Schema Registry's lookup functionality with
// normalization enabled to determine semantic equivalence. Two schemas are
// considered semantically equal if they have the same structure and meaning,
// even if they differ in formatting, field ordering, or other non-semantic
// aspects.
func IsSemanticallyEqual(
	ctx context.Context,
	client *srclient.SchemaRegistryClient,
	subject string,
	schemaString string,
	schemaType srclient.SchemaType,
	refs []srclient.Reference,
) (bool, error) {
	req := &srclient.RegisterSchemaRequest{
		Schema:     schemaString,
		SchemaType: schemaType,
		References: refs,
	}
	// Lookup schema with `normalize = true` We only care about the error, but
	// need to handle all return values
	_, _, _, err := client.LookupSchemaUnderSubject(ctx, subject, req, true) //nolint:dogsled

	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, srclient.ErrSemanticSchemaNotFound):
		return false, nil

	default:
		return false, err
	}
}
