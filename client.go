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
		return err
	}
	_, err := c.node.NamedExecContext(ctx, query, args[0])
	return err
}

// MustExec executes a named statement using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustExec(ctx context.Context, query string, args ...interface{}) {
	err := c.Exec(ctx, query, args...)
	if err != nil {
		panic(fmt.Sprintf("sqlxx: %s", err.Error()))
	}
}

// Query executes a named statement that returns rows using given arguments.
func (c *Client) Query(ctx context.Context, query string, arg interface{}) (Rows, error) {
	return c.node.NamedQueryContext(ctx, query, arg)
}

// MustQuery executes a named statement that returns rows using given arguments.
// If an error has occurred, it panics.
func (c *Client) MustQuery(ctx context.Context, query string, arg interface{}) Rows {
	rows, err := c.Query(ctx, query, arg)
	if err != nil {
		panic(fmt.Sprintf("sqlxx: %s", err.Error()))
	}
	return rows
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the returned statement.
func (c *Client) Prepare(ctx context.Context, query string) (Statement, error) {
	stmt, err := c.node.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return wrapStatement(stmt), nil
}

// Get using given named statement and arguments.
// If there is no row, an error is returned.
// Output must be a pointer to a value.
func (c *Client) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.node.GetContext(ctx, dest, query, args...)
}

// Select using given named statement and arguments.
// Output must be a pointer to a slice of value.
func (c *Client) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return c.node.SelectContext(ctx, dest, query, args...)
}

// Begin a new transaction.
func (c *Client) Begin() (Driver, error) {
	node, err := c.node.Beginx()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return wrapClient(c, node), nil
}

// Rollback the associated transaction.
func (c *Client) Rollback() error {
	return c.node.Rollback()
}

// Commit the associated transaction.
func (c *Client) Commit() error {
	return c.node.Commit()
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
func wrapStatement(stmt *sqlx.NamedStmt) *stmtWrapper {
	return &stmtWrapper{
		stmt: stmt,
	}
}

// Close closes the statement.
func (w *stmtWrapper) Close() error {
	return w.stmt.Close()
}

// Exec executes a named statement using the struct passed.
func (w *stmtWrapper) Exec(ctx context.Context, arg interface{}) error {
	_, err := w.stmt.ExecContext(ctx, arg)
	return err
}

// QueryRow using this Statement.
func (w *stmtWrapper) QueryRow(ctx context.Context, arg interface{}) (Row, error) {
	row := w.stmt.QueryRowxContext(ctx, arg)
	err := row.Err()
	return row, err
}

// QueryRows using this Statement.
func (w *stmtWrapper) QueryRows(ctx context.Context, arg interface{}) (Rows, error) {
	return w.stmt.QueryxContext(ctx, arg)
}

// Get using this Statement.
func (w *stmtWrapper) Get(ctx context.Context, dest interface{}, arg interface{}) error {
	return w.stmt.GetContext(ctx, dest, arg)
}

// Select using this Statement.
func (w *stmtWrapper) Select(ctx context.Context, dest interface{}, arg interface{}) error {
	return w.stmt.SelectContext(ctx, dest, arg)
}
