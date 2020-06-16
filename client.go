package makroud

import (
	"context"
	"database/sql"
	"io"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ClientDriver defines the driver name used in makroud.
const ClientDriver = "postgres"

// Client is a wrapper that can interact with the database, it's an implementation of Driver.
type Client struct {
	node  Node
	cache *DriverCache
	log   Logger
	obs   Observer
	rnd   io.Reader
}

// New returns a new Client instance.
func New(options ...Option) (*Client, error) {
	opts := NewClientOptions()

	for _, option := range options {
		err := option(opts)
		if err != nil {
			return nil, err
		}
	}

	return NewWithOptions(opts)
}

// NewWithOptions returns a new Client instance.
func NewWithOptions(options *ClientOptions) (*Client, error) {
	_ = pq.Driver{}

	node, err := Connect(ClientDriver, options.String())
	if err != nil {
		return nil, errors.Wrapf(err, "makroud: cannot connect to %s server", ClientDriver)
	}

	node.SetMaxIdleConns(options.MaxIdleConnections)
	node.SetMaxOpenConns(options.MaxOpenConnections)
	node.EnableSavepoint(options.SavepointEnabled)

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := &Client{
		node: node,
		rnd:  entropy,
	}

	if options.WithCache {
		client.cache = NewDriverCache()
	}

	if options.Logger != nil {
		client.log = options.Logger
	}

	if options.Observer != nil {
		client.obs = options.Observer
	}

	if options.Entropy != nil {
		client.rnd = options.Entropy
	}

	return client, nil
}

// Exec executes a statement using given arguments.
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := c.node.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute query")
	}

	return nil
}

// MustExec executes a statement using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustExec(ctx context.Context, query string, args ...interface{}) {
	err := c.Exec(ctx, query, args...)
	if err != nil {
		panic(err)
	}
}

// Query executes a statement that returns rows using given arguments.
func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := c.node.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot execute query")
	}
	return wrapRows(rows), nil
}

// MustQuery executes a statement that returns rows using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustQuery(ctx context.Context, query string, args ...interface{}) Rows {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		panic(err)
	}
	return rows
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (c *Client) Prepare(ctx context.Context, query string) (Statement, error) {
	stmt, err := c.node.PrepareContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot prepare statement")
	}
	return wrapStatement(stmt), nil
}

// Begin starts a new transaction.
//
// The provided context is used until the transaction is committed or rolled back.
// If the context is canceled, the driver will roll back the transaction.
// Commit will return an error if the context provided to Begin is canceled.
//
// The provided TxOptions is optional.
// If a non-default isolation level is used that the driver doesn't support, an error will be returned.
// If no option is provided, the default isolation level of the driver will be used.
func (c *Client) Begin(ctx context.Context, opts ...*TxOptions) (Driver, error) {
	var txOpts *TxOptions
	if len(opts) > 0 {
		txOpts = opts[0]
	}

	node, err := c.node.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot create a transaction")
	}
	return wrapClient(c, node), nil
}

// Rollback rollbacks the associated transaction.
func (c *Client) Rollback() error {
	err := c.node.Rollback()
	if err != nil {
		return errors.Wrap(err, "makroud: cannot rollback transaction")
	}
	return nil
}

// Commit commits the associated transaction.
func (c *Client) Commit() error {
	err := c.node.Commit()
	if err != nil {
		return errors.Wrap(err, "makroud: cannot commit transaction")
	}
	return nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	err := c.node.Close()
	if err == nil {
		return nil
	}

	err = errors.Wrapf(err, "makroud: trying to close %T", c.node)
	if c.obs == nil {
		return err
	}

	c.obs.OnClose(err, map[string]string{
		"action": "close",
	})

	return err
}

// Ping verifies that the underlying connection is healthy.
func (c *Client) Ping() error {
	timeout := 1 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.PingContext(ctx)
}

// PingContext verifies that the underlying connection is healthy.
func (c *Client) PingContext(ctx context.Context) error {
	row, err := c.node.QueryContext(ctx, "SELECT true")
	if row != nil {
		defer close(c, row, map[string]string{
			"query":  "SELECT true;",
			"action": "ping",
		})
	}
	if err != nil {
		return errors.Wrap(err, "makroud: cannot ping database")
	}

	return nil
}

// DriverName returns the driver name used by this driver.
func (c *Client) DriverName() string {
	return c.node.DriverName()
}

// HasCache returns if current driver has an internal cache.
func (c *Client) HasCache() bool {
	return c.cache != nil
}

// GetCache returns the driver internal cache.
//
// WARNING: Please, do not use this method unless you know what you are doing:
// YOU COULD BREAK YOUR DRIVER.
func (c *Client) GetCache() *DriverCache {
	return c.cache
}

// SetCache replace the driver internal cache by the given one.
//
// WARNING: Please, do not use this method unless you know what you are doing:
// YOU COULD BREAK YOUR DRIVER.
func (c *Client) SetCache(cache *DriverCache) {
	c.cache = cache
}

// HasLogger returns if the driver has a logger.
func (c *Client) HasLogger() bool {
	return c.log != nil
}

// Logger returns the driver logger.
//
// WARNING: Please, do not use this method unless you know what you are doing.
func (c *Client) Logger() Logger {
	return c.log
}

// HasObserver returns if the driver has an observer.
func (c *Client) HasObserver() bool {
	return c.obs != nil
}

// Observer returns the driver observer.
//
// WARNING: Please, do not use this method unless you know what you are doing.
func (c *Client) Observer() Observer {
	return c.obs
}

// Entropy returns an entropy source, used for primary key generation (if required).
//
// WARNING: Please, do not use this method unless you know what you are doing.
func (c *Client) Entropy() io.Reader {
	return c.rnd
}

// wrapClient creates a new Client using given database connection.
func wrapClient(client *Client, connection Node) Driver {
	return &Client{
		node:  connection,
		cache: client.cache,
		log:   client.log,
		rnd:   client.rnd,
	}
}

// A stmtWrapper wraps a statement from sql.
type stmtWrapper struct {
	stmt *sql.Stmt
}

// wrapStatement creates a new Statement using given statement.
func wrapStatement(stmt *sql.Stmt) Statement {
	return &stmtWrapper{
		stmt: stmt,
	}
}

// Close closes the statement.
func (w *stmtWrapper) Close() error {
	err := w.stmt.Close()
	if err != nil {
		return errors.Wrap(err, "makroud: cannot close statement")
	}
	return nil
}

// Exec executes this statement using the struct passed.
func (w *stmtWrapper) Exec(ctx context.Context, args ...interface{}) error {
	_, err := w.stmt.ExecContext(ctx, args...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute statement")
	}
	return nil
}

// QueryRow executes this statement returning a single row.
func (w *stmtWrapper) QueryRow(ctx context.Context, args ...interface{}) (Row, error) {
	rows, err := w.stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot execute statement")
	}
	return wrapRow(rows), nil
}

// QueryRows executes this statement returning a list of rows.
func (w *stmtWrapper) QueryRows(ctx context.Context, args ...interface{}) (Rows, error) {
	rows, err := w.stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot execute statement")
	}
	return wrapRows(rows), nil
}

// A rowWrapper is a reimplementation of sql.Row in order to gain access to the underlying
// Columns() function.
type rowWrapper struct {
	rows *sql.Rows
}

// wrapRow creates a new Row using given rows from sql.
func wrapRow(rows *sql.Rows) Row {
	return &rowWrapper{
		rows: rows,
	}
}

// Write copies the columns in the current row into the given map.
func (r *rowWrapper) Write(dest map[string]interface{}) error {
	err := mapScan(r, dest)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot write row")
	}
	return nil
}

// Columns returns the column names.
func (r *rowWrapper) Columns() ([]string, error) {
	columns, err := r.rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot return row columns")
	}

	return columns, nil
}

// Scan copies the columns in the current row into the values pointed at by dest.
// The number of values in dest must be the same as the number of columns in Rows.
func (r *rowWrapper) Scan(dest ...interface{}) error {
	err := r.scan(dest...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot scan given values")
	}
	return nil
}

func (r *rowWrapper) scan(dest ...interface{}) error {
	// From https://github.com/jmoiron/sqlx source code:
	// Discard sql.RawBytes to avoid weird issues with the SQL driver and memory management.
	defer func() {
		_ = r.rows.Close()
	}()
	for i := range dest {
		_, ok := dest[i].(*sql.RawBytes)
		if ok {
			return errors.New("sql.RawBytes isn't allowed on Row.Scan")
		}
	}

	if !r.rows.Next() {
		err := r.rows.Err()
		if err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	err := r.rows.Scan(dest...)
	if err != nil {
		return err
	}

	// Make sure the query can be processed to completion with no errors.
	return r.rows.Close()
}

// A rowsWrapper wraps a rows from sql.
type rowsWrapper struct {
	rows *sql.Rows
}

// wrapRow creates a new Rows using given rows from sql.
func wrapRows(rows *sql.Rows) Rows {
	return &rowsWrapper{
		rows: rows,
	}
}

// Next prepares the next result row for reading with the Scan method.
// It returns true on success, or false if there is no next result row or an error
// happened while preparing it.
// Err should be consulted to distinguish between the two cases.
// Every call to Scan, even the first one, must be preceded by a call to Next.
func (r *rowsWrapper) Next() bool {
	return r.rows.Next()
}

// Close closes the Rows, preventing further enumeration/iteration.
// If Next is called and returns false and there are no further result sets, the Rows are closed automatically
// and it will suffice to check the result of Err.
func (r *rowsWrapper) Close() error {
	err := r.rows.Close()
	if err != nil {
		return errors.Wrap(err, "makroud: cannot close rows")
	}
	return nil
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (r *rowsWrapper) Err() error {
	err := r.rows.Err()
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

// Write copies the columns in the current row into the given map.
func (r *rowsWrapper) Write(dest map[string]interface{}) error {
	err := mapScan(r, dest)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot write row")
	}

	err = r.Err()
	if err != nil {
		return errors.Wrap(err, "makroud: cannot write row")
	}

	return nil
}

// Columns returns the column names.
func (r *rowsWrapper) Columns() ([]string, error) {
	columns, err := r.rows.Columns()
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot return rows columns")
	}
	return columns, nil
}

// Scan copies the columns in the current row into the values pointed at by dest.
// The number of values in dest must be the same as the number of columns in Rows.
func (r *rowsWrapper) Scan(dest ...interface{}) error {
	err := r.rows.Scan(dest...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot scan given values")
	}
	return nil
}

type mapScanner interface {
	Columns() ([]string, error)
	Scan(...interface{}) error
}

// mapScan scans the current row into the given map.
// Use this for debugging or analysis if the results might not be under your control.
// Please do not use this as a primary interface!
func mapScan(scanner mapScanner, dest map[string]interface{}) error {
	columns, err := scanner.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = scanner.Scan(values...)
	if err != nil {
		return err
	}

	for i, column := range columns {
		dest[column] = *(values[i].(*interface{}))
	}

	return nil
}
