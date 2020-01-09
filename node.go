package makroud

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

// Forked from github.com/heetch/sqalx

// Node is a components that allows to seamlessly create nested transactions and to avoid thinking about whether
// or not a function is called within a transaction.
// With this, you can easily create reusable and composable functions that can be called within or out of
// transactions and that can create transactions themselves.
type Node interface {

	// ----------------------------------------------------------------------------
	// Query
	// ----------------------------------------------------------------------------

	// Exec executes a statement using given arguments. The query shouldn't return rows.
	Exec(query string, args ...interface{}) (sql.Result, error)
	// ExecContext executes a statement using given arguments. The query shouldn't return rows.
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	// Query executes a statement that returns rows using given arguments.
	Query(query string, args ...interface{}) (*sql.Rows, error)
	// QueryContext executes a statement that returns rows using given arguments.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// QueryRow executes a statement that returns at most one row using given arguments.
	QueryRow(query string, args ...interface{}) *sql.Row
	// QueryRowContext executes a statement that returns at most one row using given arguments.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	// Prepare creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the returned statement.
	Prepare(query string) (*sql.Stmt, error)
	// PrepareContext creates a prepared statement for later queries or executions.
	// Multiple queries or executions may be run concurrently from the returned statement.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// ----------------------------------------------------------------------------
	// Connection
	// ----------------------------------------------------------------------------

	// DriverName returns the driver name used by this connector.
	DriverName() string
	// Ping verifies that the underlying connection is healthy.
	Ping() error
	// PingContext verifies that the underlying connection is healthy.
	PingContext(ctx context.Context) error
	// Close closes the underlying connection.
	Close() error
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	// Expired connections may be closed lazily before reuse.
	SetConnMaxLifetime(duration time.Duration)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	SetMaxIdleConns(number int)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	SetMaxOpenConns(number int)
	// EnableSavepoint activate PostgreSQL savepoints for nested transactions.
	EnableSavepoint(enabled bool)
	// Stats returns database statistics.
	Stats() sql.DBStats

	// ----------------------------------------------------------------------------
	// Transaction
	// ----------------------------------------------------------------------------

	// Begin begins a new transaction.
	Begin() (Node, error)
	// BeginContext begins a new transaction.
	BeginContext(ctx context.Context, opts *sql.TxOptions) (Node, error)
	// Rollback rollbacks the associated transaction.
	Rollback() error
	// Commit commits the associated transaction.
	Commit() error

	// ----------------------------------------------------------------------------
	// System
	// ----------------------------------------------------------------------------

	// Tx returns the underlying transaction.
	Tx() *sql.Tx
	// DB returns the underlying connection.
	DB() *sql.DB
}

type node struct {
	driver           string
	db               *sql.DB
	tx               *sql.Tx
	savePointID      string
	savePointEnabled bool
	nested           bool
}

// Connect connects to a database and verifies connection with a ping.
func Connect(driver string, dsn string) (Node, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		// the connection has been opened within this function, we must close it
		// on error.
		_ = db.Close()
		return nil, err
	}

	node := &node{
		driver: driver,
		db:     db,
	}

	return node, nil
}

// DriverName returns the driver name used by this connector.
func (node *node) DriverName() string {
	return node.driver
}

// Ping verifies that the underlying connection is healthy.
func (node *node) Ping() error {
	return node.db.Ping()
}

// PingContext verifies that the underlying connection is healthy.
func (node *node) PingContext(ctx context.Context) error {
	return node.db.PingContext(ctx)
}

// Close closes the underlying connection.
func (node *node) Close() error {
	return node.db.Close()
}

// Begin begins a new transaction.
func (node *node) Begin() (Node, error) {
	return node.BeginContext(context.Background(), nil)
}

// BeginContext begins a new transaction.
func (node *node) BeginContext(ctx context.Context, opts *sql.TxOptions) (Node, error) {

	clone := node.clone()

	switch {
	case clone.tx == nil:

		// Create new transaction.
		tx, err := clone.db.BeginTx(ctx, opts)
		if err != nil {
			return nil, err
		}

		clone.tx = tx

	case clone.savePointEnabled:

		// Already in a transaction: using savepoints
		clone.nested = true

		// Savepoints name must start with a char and cannot contain dashes (-)
		clone.savePointID = "sp_" + strings.Replace(uuid.Must(uuid.NewV1()).String(), "-", "_", -1)
		_, err := node.tx.Exec("SAVEPOINT " + clone.savePointID)
		if err != nil {
			return nil, err
		}

	default:

		// Already in a transaction: reusing current one.
		clone.nested = true
	}

	return clone, nil
}

// Commit commits the associated transaction.
func (node *node) Commit() (err error) {
	if node.tx == nil {
		return ErrCommitNotInTransaction
	}

	if node.savePointID != "" {

		_, err = node.tx.Exec("RELEASE SAVEPOINT " + node.savePointID)
		if err != nil {
			return err
		}

	} else if !node.nested {
		err = node.tx.Commit()
		if err != nil {
			return err
		}
	}

	node.tx = nil

	return nil
}

// Rollback rollbacks the associated transaction.
func (node *node) Rollback() (err error) {
	if node.tx == nil {
		return nil
	}

	if node.savePointEnabled && node.savePointID != "" {

		_, err = node.tx.Exec("ROLLBACK TO SAVEPOINT " + node.savePointID)
		if err != nil {
			return err
		}

	} else if !node.nested {

		err = node.tx.Rollback()
		if err != nil {
			return err
		}

	}

	node.tx = nil

	return nil
}

// Tx returns the underlying transaction.
func (node *node) Tx() *sql.Tx {
	return node.tx
}

// DB returns the underlying transaction.
func (node *node) DB() *sql.DB {
	return node.db
}

// Exec executes a statement using given arguments. The query shouldn't return rows.
func (node *node) Exec(query string, args ...interface{}) (sql.Result, error) {
	if node.tx == nil {
		return node.db.Exec(query, args...)
	}
	return node.tx.Exec(query, args...)
}

// ExecContext executes a statement using given arguments. The query shouldn't return rows.
func (node *node) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if node.tx == nil {
		return node.db.ExecContext(ctx, query, args...)
	}
	return node.tx.ExecContext(ctx, query, args...)
}

// Query executes a statement that returns rows using given arguments.
func (node *node) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if node.tx == nil {
		return node.db.Query(query, args...)
	}
	return node.tx.Query(query, args...)
}

// QueryContext executes a statement that returns rows using given arguments.
func (node *node) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if node.tx == nil {
		return node.db.QueryContext(ctx, query, args...)
	}
	return node.tx.QueryContext(ctx, query, args...)
}

// QueryRow executes a statement that returns at most one row using given arguments.
func (node *node) QueryRow(query string, args ...interface{}) *sql.Row {
	if node.tx == nil {
		return node.db.QueryRow(query, args...)
	}
	return node.tx.QueryRow(query, args...)
}

// QueryRowContext executes a statement that returns at most one row using given arguments.
func (node *node) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if node.tx == nil {
		return node.db.QueryRowContext(ctx, query, args...)
	}
	return node.tx.QueryRowContext(ctx, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (node *node) Prepare(query string) (*sql.Stmt, error) {
	if node.tx == nil {
		return node.db.Prepare(query)
	}
	return node.tx.Prepare(query)
}

// PrepareContext creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (node *node) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if node.tx == nil {
		return node.db.PrepareContext(ctx, query)
	}
	return node.tx.PrepareContext(ctx, query)
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
// Expired connections may be closed lazily before reuse.
func (node *node) SetConnMaxLifetime(duration time.Duration) {
	node.db.SetConnMaxLifetime(duration)
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
func (node *node) SetMaxIdleConns(number int) {
	node.db.SetMaxIdleConns(number)
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (node *node) SetMaxOpenConns(number int) {
	node.db.SetMaxOpenConns(number)
}

// EnableSavepoint activates PostgreSQL savepoints for nested transactions.
func (node *node) EnableSavepoint(enabled bool) {
	node.savePointEnabled = enabled
}

// Stats returns database statistics.
func (node *node) Stats() sql.DBStats {
	return node.db.Stats()
}

// clone clones current node.
func (node *node) clone() *node {
	clone := *node
	return &clone
}
