package sqlxx

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/ulule/sqalx"
	"github.com/ulule/sqlx"
)

// ClientDriver define the driver name used in sqlxx.
const ClientDriver = "postgres"

// Client is a wrapper that can interact with the database.
type Client struct {
	node  sqalx.Node
	store *cache
	log   Logger
	rnd   io.Reader
}

// clientOptions configure a Client instance. clientOptions are set by the Option
// values passed to New.
type clientOptions struct {
	port               int
	host               string
	user               string
	password           string
	dbName             string
	sslMode            string
	timezone           string
	maxOpenConnections int
	maxIdleConnections int
	withCache          bool
	savepointEnabled   bool
	logger             Logger
}

func (e clientOptions) String() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s;timezone=%s",
		ClientDriver,
		e.user,
		e.password,
		e.host,
		e.port,
		e.dbName,
		e.sslMode,
		e.timezone,
	)
}

// New returns a new Client instance.
func New(options ...Option) (*Client, error) {
	opts := &clientOptions{
		host:               "localhost",
		port:               5432,
		user:               "postgres",
		password:           "",
		dbName:             "postgres",
		sslMode:            "disable",
		timezone:           "UTC",
		maxOpenConnections: 5,
		maxIdleConnections: 2,
		withCache:          true,
		savepointEnabled:   false,
	}

	for _, option := range options {
		err := option.apply(opts)
		if err != nil {
			return nil, err
		}
	}

	_ = pq.Driver{}

	dbx, err := sqlx.Connect(ClientDriver, opts.String())
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: cannot connect to %s server", ClientDriver)
	}

	dbx.SetMaxIdleConns(opts.maxIdleConnections)
	dbx.SetMaxOpenConns(opts.maxOpenConnections)

	connection, err := sqalx.New(dbx, sqalx.SavePoint(opts.savepointEnabled))
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: cannot instantiate %s client driver", ClientDriver)
	}

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := &Client{
		node: connection,
		log:  &EmptyLogger{},
		rnd:  entropy,
	}

	if opts.withCache {
		client.store = newCache()
	}

	if opts.logger != nil {
		client.log = opts.logger
	}

	return client, nil
}

// Exec executes a named statement using given arguments.
func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) error {
	if len(args) == 0 {

		_, err := c.node.ExecContext(ctx, query)
		if err != nil {
			return errors.Wrap(err, "sqlxx: cannot execute query")
		}

		return nil
	}

	_, err := c.node.NamedExecContext(ctx, query, args[0])
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// MustExec executes a named statement using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustExec(ctx context.Context, query string, args ...interface{}) {
	err := c.Exec(ctx, query, args...)
	if err != nil {
		panic(err)
	}
}

// Query executes a named statement that returns rows using given arguments.
func (c *Client) Query(ctx context.Context, query string, arg interface{}) (Rows, error) {
	rows, err := c.node.NamedQueryContext(ctx, query, arg)
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot execute query")
	}
	return wrapRows(rows), nil
}

// MustQuery executes a named statement that returns rows using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustQuery(ctx context.Context, query string, arg interface{}) Rows {
	rows, err := c.Query(ctx, query, arg)
	if err != nil {
		panic(err)
	}
	return rows
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (c *Client) Prepare(ctx context.Context, query string) (Statement, error) {
	stmt, err := c.node.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	return wrapStatement(stmt), nil
}

// FindOne executes this named statement to fetch one record.
// If there is no row, an error is returned.
// Output must be a pointer to a value.
func (c *Client) FindOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := c.node.GetContext(ctx, dest, query, args...)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}
	return nil
}

// FindAll executes this named statement to fetch a list of records.
// Output must be a pointer to a slice of value.
func (c *Client) FindAll(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := c.node.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}
	return nil
}

// Begin a new transaction.
func (c *Client) Begin() (Driver, error) {
	node, err := c.node.Beginx()
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: create a transaction")
	}
	return wrapClient(c, node), nil
}

// Rollback the associated transaction.
func (c *Client) Rollback() error {
	err := c.node.Rollback()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot rollback transaction")
	}
	return nil
}

// Commit the associated transaction.
func (c *Client) Commit() error {
	err := c.node.Commit()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot commit transaction")
	}
	return nil
}

// Close closes the underlying connection.
func (c *Client) Close() error {
	return c.node.Close()
}

// Ping verifies that the underlying connection is healthy.
func (c *Client) Ping() error {
	row, err := c.node.Query("SELECT true")
	if row != nil {
		defer c.close(row, map[string]string{
			"query": "SELECT true;",
		})
	}
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot ping database")
	}
	return nil
}

// DriverName returns the driver name used by this driver.
func (c *Client) DriverName() string {
	return c.node.DriverName()
}

func (c *Client) hasCache() bool {
	return c.store != nil
}

func (c *Client) cache() *cache {
	return c.store
}

func (c *Client) logger() Logger {
	return c.log
}

func (c *Client) entropy() io.Reader {
	return c.rnd
}

func (c *Client) close(closer io.Closer, flags map[string]string) {
	thr := closer.Close()
	if thr != nil {
		thr = errors.Wrapf(thr, "trying to close: %T", closer)
		// TODO (novln): Add an observer to collect this error.
		_ = thr
	}
}

// wrapClient creates a new Client using given database connection.
func wrapClient(client *Client, connection sqalx.Node) Driver {
	return &Client{
		node:  connection,
		store: client.store,
		log:   client.log,
		rnd:   client.rnd,
	}
}

// A stmtWrapper wraps a named statement from sqlx.
type stmtWrapper struct {
	stmt *sqlx.NamedStmt
}

// wrapStatement creates a new Statement using given named statement from sqlx.
func wrapStatement(stmt *sqlx.NamedStmt) Statement {
	return &stmtWrapper{
		stmt: stmt,
	}
}

// Close closes the statement.
func (w *stmtWrapper) Close() error {
	err := w.stmt.Close()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot close statement")
	}
	return nil
}

// Exec executes this named statement using the struct passed.
func (w *stmtWrapper) Exec(ctx context.Context, arg interface{}) error {
	_, err := w.stmt.ExecContext(ctx, arg)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute statement")
	}
	return nil
}

// QueryRow executes this named statement returning a single row.
func (w *stmtWrapper) QueryRow(ctx context.Context, arg interface{}) (Row, error) {
	row := w.stmt.QueryRowxContext(ctx, arg)
	err := row.Err()
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot execute statement")
	}
	if row == nil {
		return nil, errors.Wrap(ErrNoRows, "sqlxx: cannot execute statement")
	}
	return wrapRow(row), nil
}

// QueryRows executes this named statement returning a list of rows.
func (w *stmtWrapper) QueryRows(ctx context.Context, arg interface{}) (Rows, error) {
	rows, err := w.stmt.QueryxContext(ctx, arg)
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot execute statement")
	}
	return wrapRows(rows), nil
}

// FindOne executes this named statement to fetch one record.
func (w *stmtWrapper) FindOne(ctx context.Context, dest interface{}, arg interface{}) error {
	err := w.stmt.GetContext(ctx, dest, arg)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute statement")
	}
	return nil
}

// FindAll executes this named statement to fetch a list of records.
func (w *stmtWrapper) FindAll(ctx context.Context, dest interface{}, arg interface{}) error {
	err := w.stmt.SelectContext(ctx, dest, arg)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute statement")
	}
	return nil
}

// A rowWrapper wraps a row from sqlx.
type rowWrapper struct {
	row *sqlx.Row
}

// wrapRow creates a new Row using given row from sqlx.
func wrapRow(row *sqlx.Row) Row {
	return &rowWrapper{
		row: row,
	}
}

// Write copies the columns in the current row into the given map.
func (r *rowWrapper) Write(dest map[string]interface{}) error {
	err := r.row.MapScan(dest)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot write row")
	}
	return nil
}

// A rowsWrapper wraps a rows from sqlx.
type rowsWrapper struct {
	rows *sqlx.Rows
}

// wrapRow creates a new Rows using given rows from sqlx.
func wrapRows(rows *sqlx.Rows) Rows {
	return &rowsWrapper{
		rows: rows,
	}
}

// Next prepares the next result row for reading with the MapScan method.
// It returns true on success, or false if there is no next result row or an error
// happened while preparing it.
// Err should be consulted to distinguish between the two cases.
// Every call to MapScan, even the first one, must be preceded by a call to Next.
func (r *rowsWrapper) Next() bool {
	return r.rows.Next()
}

// Close closes the Rows, preventing further enumeration/iteration.
// If Next is called and returns false and there are no further result sets, the Rows are closed automatically
// and it will suffice to check the result of Err.
func (r *rowsWrapper) Close() error {
	err := r.rows.Close()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot close rows")
	}
	return nil
}

// Err returns the error, if any, that was encountered during iteration.
// Err may be called after an explicit or implicit Close.
func (r *rowsWrapper) Err() error {
	return r.rows.Err()
}

// Write copies the columns in the current row into the given map.
func (r *rowsWrapper) Write(dest map[string]interface{}) error {
	err := r.rows.MapScan(dest)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot write row")
	}
	return nil
}
