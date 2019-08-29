package makroud

import (
	"fmt"
)

// Makroud general errors.
var (
	// ErrNoRows is returned when query doesn't return a row.
	ErrNoRows = fmt.Errorf("no rows in result set")
	// ErrInvalidDriver is returned when given driver is undefined.
	ErrInvalidDriver = fmt.Errorf("a makroud driver is required")
	// ErrPointerRequired is returned when given value is not a pointer.
	ErrPointerRequired = fmt.Errorf("a pointer is required")
	// ErrPointerOrSliceRequired is returned when given value is not a pointer or a slice.
	ErrPointerOrSliceRequired = fmt.Errorf("a pointer or a slice is required")
	// ErrUnknownPreloadRule is returned when given rule is unknown.
	ErrUnknownPreloadRule = fmt.Errorf("unknown rule")
	// ErrStructRequired is returned when given value is not a struct.
	ErrStructRequired = fmt.Errorf("a struct is required")
	// ErrSchemaColumnRequired is returned when we cannot find a column in current schema.
	ErrSchemaColumnRequired = fmt.Errorf("cannot find column in schema")
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
	// ErrSelectorNotFoundConnection is returned when the required connection does not exists in selector connections.
	ErrSelectorNotFoundConnection = fmt.Errorf("cannot find connection in selector")
	// ErrSelectorMissingRetryConnection is returned when the retry mechanism has no connection available from selector.
	ErrSelectorMissingRetryConnection = fmt.Errorf("cannot find a healthy connection in selector")
	// ErrSliceOfScalarMultipleColumns is returned when a query return multiple columns for a slice of scalar.
	ErrSliceOfScalarMultipleColumns = fmt.Errorf("slice of scalar with multiple columns")
	// ErrCommitNotInTransaction is returned when using commit outside of a transaction.
	ErrCommitNotInTransaction = fmt.Errorf("cannot commit outside of a transaction")
)
