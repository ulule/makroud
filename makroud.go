package makroud

import (
	"context"
	"io"
)

// Driver is a high level abstraction of a database connection or a transaction.
type Driver interface {

	// ----------------------------------------------------------------------------
	// Query
	// ----------------------------------------------------------------------------

	// Exec executes a statement using given arguments.
	Exec(ctx context.Context, query string, args ...interface{}) error

	// MustExec executes a statement using given arguments.
	// If an error has occurred, it panics.
	MustExec(ctx context.Context, query string, args ...interface{})

	// Query executes a statement that returns rows using given arguments.
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)

	// QueryRow executes a statement returning a single row.
	QueryRow(ctx context.Context, query string, args ...interface{}) (Row, error)

	// MustQuery executes a statement that returns rows using given arguments.
	// If an error has occurred, it panics.
	MustQuery(ctx context.Context, query string, args ...interface{}) Rows

	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the returned statement.
	Prepare(ctx context.Context, query string) (Statement, error)

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

	// Begin starts a new transaction.
	//
	// The provided context is used until the transaction is committed or rolled back.
	// If the context is canceled, the driver will roll back the transaction.
	// Commit will return an error if the context provided to Begin is canceled.
	//
	// The provided TxOptions is optional.
	// If a non-default isolation level is used that the driver doesn't support, an error will be returned.
	// If no option is provided, the default isolation level of the driver will be used.
	Begin(ctx context.Context, opts ...*TxOptions) (Driver, error)

	// Rollback rollbacks the associated transaction.
	Rollback() error

	// Commit commits the associated transaction.
	Commit() error

	// ----------------------------------------------------------------------------
	// System
	// ----------------------------------------------------------------------------

	// HasCache returns if current driver has an internal cache.
	HasCache() bool

	// GetCache returns the driver internal cache.
	//
	// WARNING: Please, do not use this method unless you know what you are doing:
	// YOU COULD BREAK YOUR DRIVER.
	GetCache() *DriverCache

	// SetCache replace the driver internal cache by the given one.
	//
	// WARNING: Please, do not use this method unless you know what you are doing:
	// YOU COULD BREAK YOUR DRIVER.
	SetCache(cache *DriverCache)

	// HasLogger returns if the driver has a logger.
	HasLogger() bool

	// Logger returns the driver logger.
	//
	// WARNING: Please, do not use this method unless you know what you are doing.
	Logger() Logger

	// HasObserver returns if the driver has an observer.
	HasObserver() bool

	// Observer returns the driver observer.
	//
	// WARNING: Please, do not use this method unless you know what you are doing.
	Observer() Observer

	// Entropy returns an entropy source, used for primary key generation (if required).
	//
	// WARNING: Please, do not use this method unless you know what you are doing.
	Entropy() io.Reader
}

// A Statement from prepare.
type Statement interface {
	// Close closes the statement.
	Close() error
	// Exec executes this named statement using the struct passed.
	Exec(ctx context.Context, args ...interface{}) error
	// QueryRow executes this named statement returning a single row.
	QueryRow(ctx context.Context, args ...interface{}) (Row, error)
	// QueryRows executes this named statement returning a list of rows.
	QueryRows(ctx context.Context, args ...interface{}) (Rows, error)
}

// A Row is a simple row.
type Row interface {
	// Write copies the columns in the current row into the given map.
	// Use this for debugging or analysis if the results might not be under your control.
	// Please do not use this as a primary interface!
	Write(dest map[string]interface{}) error
	// Columns returns the column names.
	Columns() ([]string, error)
	// Scan copies the columns in the current row into the values pointed at by dest.
	// The number of values in dest must be the same as the number of columns in Rows.
	Scan(dest ...interface{}) error
}

// A Rows is an iteratee of a list of records.
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
	// Write copies the columns in the current row into the given map.
	// Use this for debugging or analysis if the results might not be under your control.
	// Please do not use this as a primary interface!
	Write(dest map[string]interface{}) error
	// Columns returns the column names.
	Columns() ([]string, error)
	// Scan copies the columns in the current row into the values pointed at by dest.
	// The number of values in dest must be the same as the number of columns in Rows.
	Scan(dest ...interface{}) error
}
