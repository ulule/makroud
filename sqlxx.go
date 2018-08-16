package sqlxx

import (
	"context"
	"io"
)

// Driver is a high level abstraction of a database connection or a transaction.
type Driver interface {

	// ----------------------------------------------------------------------------
	// Query
	// ----------------------------------------------------------------------------

	// Exec executes a named statement using given arguments.
	Exec(ctx context.Context, query string, args ...interface{}) error

	// MustExec executes a named statement using given arguments.
	// If an error has occurred, it panics.
	MustExec(ctx context.Context, query string, args ...interface{})

	// Query executes a named statement that returns rows using given arguments.
	Query(ctx context.Context, query string, arg interface{}) (Rows, error)

	// MustQuery executes a named statement that returns rows using given arguments.
	// If an error has occurred, it panics.
	MustQuery(ctx context.Context, query string, arg interface{}) Rows

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the returned statement.
	Prepare(ctx context.Context, query string) (Statement, error)

	// Get using given named statement and arguments.
	// If there is no row, an error is returned.
	// Output must be a pointer to a value.
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Select using given named statement and arguments.
	// Output must be a pointer to a slice of value.
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// ----------------------------------------------------------------------------
	// Connection
	// ----------------------------------------------------------------------------

	// Close closes the underlying connection.
	Close() error

	// Ping verifies that the underlying connection is healthy.
	Ping() error

	// DriverName returns the driver name used by this driver.
	DriverName() string

	// ----------------------------------------------------------------------------
	// Transaction
	// ----------------------------------------------------------------------------

	// Begin a new transaction.
	Begin() (Driver, error)

	// Rollback the associated transaction.
	Rollback() error

	// Commit the associated transaction.
	Commit() error

	// ----------------------------------------------------------------------------
	// System
	// ----------------------------------------------------------------------------

	close(closer io.Closer, flags map[string]string)
	hasCache() bool
	cache() *cache
	logger() Logger
	entropy() io.Reader
}

// Statement from PrepareNamed.
type Statement interface {
	// Close closes the statement.
	Close() error
	// Exec executes a named statement using the struct passed.
	Exec(ctx context.Context, arg interface{}) error
	// QueryRow using this Statement.
	QueryRow(ctx context.Context, arg interface{}) (Row, error)
	// QueryRows using this Statement.
	QueryRows(ctx context.Context, arg interface{}) (Rows, error)
	// Get using this Statement.
	Get(ctx context.Context, dest interface{}, arg interface{}) error
	// Select using this Statement.
	Select(ctx context.Context, dest interface{}, arg interface{}) error
}

type Row interface {
	// Err returns the error, if any, that was encountered during iteration.
	// Err may be called after an explicit or implicit Close.
	Err() error
	// MapScan copies the columns in the current row into the given map.
	MapScan(dest map[string]interface{}) error
}

type Rows interface {
	// Next prepares the next result row for reading with the MapScan method.
	// It returns true on success, or false if there is no next result row or an error
	// happened while preparing it.
	// Err should be consulted to distinguish between the two cases.
	// Every call to MapScan, even the first one, must be preceded by a call to Next.
	Next() bool
	// Close closes the Rows, preventing further enumeration/iteration.
	// If Next is called and returns false and there are no further result sets, the Rows are closed automatically
	// and it will suffice to check the result of Err.
	Close() error
	// Err returns the error, if any, that was encountered during iteration.
	// Err may be called after an explicit or implicit Close.
	Err() error
	// MapScan copies the columns in the current row into the given map.
	MapScan(dest map[string]interface{}) error
}
