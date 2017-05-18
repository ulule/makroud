package sqlxx

import (
	"fmt"

	"github.com/heetch/sqalx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ClientDriver define the driver name used in sqlxx.
const ClientDriver = "postgres"

// Client is a wrapper that can interact with the database.
type Client struct {
	sqalx.Node
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

	dbx, err := sqlx.Connect(ClientDriver, opts.String())
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: cannot connect to %s server", ClientDriver)
	}

	dbx.SetMaxIdleConns(opts.maxIdleConnections)
	dbx.SetMaxOpenConns(opts.maxOpenConnections)

	connection, err := sqalx.New(dbx)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: cannot instantiate %s client driver", ClientDriver)
	}

	client := &Client{
		Node: connection,
	}

	return client, nil
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
