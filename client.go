package sqlxx

import (
	"database/sql"
	"fmt"

	"github.com/heetch/sqalx"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// ClientDriver define the driver name used in sqlxx.
const ClientDriver = "postgres"

// Client is a wrapper that can interact with the database.
type Client struct {
	node sqalx.Node
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
	savepointEnabled   bool
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

	var sqalxOpts []sqalx.Option
	if opts.savepointEnabled {
		sqalxOpts = append(sqalxOpts, sqalx.SavePoint(true))
	}

	connection, err := sqalx.New(dbx, sqalxOpts...)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: cannot instantiate %s client driver", ClientDriver)
	}

	client := &Client{
		node: connection,
	}

	return client, nil
}

func (c *Client) Exec(query string, args ...interface{}) (sql.Result, error) {
	return c.node.Exec(query, args...)
}

func (c *Client) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return c.node.Query(query, args...)
}

func (c *Client) Queryx(query string, args ...interface{}) (*sqlx.Rows, error) {
	return c.node.Queryx(query, args...)
}

func (c *Client) QueryRowx(query string, args ...interface{}) *sqlx.Row {
	return c.node.QueryRowx(query, args...)
}

func (c *Client) DriverName() string {
	return c.node.DriverName()
}

func (c *Client) Get(dest interface{}, query string, args ...interface{}) error {
	return c.node.Get(dest, query, args...)
}

func (c *Client) MustExec(query string, args ...interface{}) sql.Result {
	return c.node.MustExec(query, args...)
}

func (c *Client) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return c.node.NamedExec(query, arg)
}

func (c *Client) PrepareNamed(query string) (Statement, error) {
	return c.node.PrepareNamed(query)
}

func (c *Client) Rebind(query string) string {
	return c.node.Rebind(query)
}

func (c *Client) Select(dest interface{}, query string, args ...interface{}) error {
	return c.node.Select(dest, query, args...)
}

func (c *Client) Close() error {
	return c.node.Close()
}

// Ping verify that the database connection is healthy.
func (e *Client) Ping() error {
	row, err := e.Query("SELECT true")
	if row != nil {
		defer func() {
			// TODO: Add an observer to collect this error.
			thr := row.Close()
			_ = thr
		}()
	}
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot ping database")
	}
	return nil
}

func (c *Client) Beginx() (Driver, error) {
	node, err := c.node.Beginx()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return wrapClient(node), nil
}

func (c *Client) Rollback() error {
	return c.node.Rollback()
}

func (c *Client) Commit() error {
	return c.node.Commit()
}
