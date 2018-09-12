package sqlxx

import (
	"fmt"
)

// Sqlxx general errors.
var (
	// ErrNoRows is returned when query doesn't return a row.
	ErrNoRows = fmt.Errorf("no rows in result set")
	// ErrInvalidDriver is returned when given driver is undefined.
	ErrInvalidDriver = fmt.Errorf("a sqlxx driver is required")
	// ErrPointerRequired is returned when given value is not a pointer.
	ErrPointerRequired = fmt.Errorf("a pointer is required")
	// ErrSchemaCreatedKey is returned when we cannot find a created key in given schema.
	ErrSchemaCreatedKey = fmt.Errorf("cannot find created key in schema")
	// ErrSchemaUpdatedKey is returned when we cannot find a updated key in given schema.
	ErrSchemaUpdatedKey = fmt.Errorf("cannot find updated key in schema")
	// ErrSchemaDeletedKey is returned when we cannot find a deleted key in given schema.
	ErrSchemaDeletedKey = fmt.Errorf("cannot find deleted key in schema")
	// ErrPreloadInvalidSchema is returned when preload detect an invalid schema from given model.
	ErrPreloadInvalidSchema = fmt.Errorf("given model has an invalid schema")
	// ErrPreloadInvalidModel is returned when preload detect an invalid model.
	ErrPreloadInvalidModel = fmt.Errorf("given model is invalid")
	// ErrPreloadInvalidPath is returned when preload detect an invalid path.
	ErrPreloadInvalidPath = fmt.Errorf("given path is invalid")
)
