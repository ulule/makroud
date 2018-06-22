package sqlxx

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/heetch/sqalx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ClientDriver define the driver name used in sqlxx.
const ClientDriver = "postgres"

// Client is a wrapper that can interact with the database.
type Client struct {
	sqalx.Node
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

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))

	client := &Client{
		Node: connection,
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

// Ping verify that the database connection is healthy.
func (e *Client) Ping() error {
	row, err := e.Query("SELECT true")
	if row != nil {
		defer e.close(row, map[string]string{
			"query": "SELECT true;",
		})
	}
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot ping database")
	}
	return nil
}

func (e *Client) hasCache() bool {
	return e.store != nil
}

func (e *Client) cache() *cache {
	return e.store
}

func (e *Client) logger() Logger {
	return e.log
}

func (e *Client) entropy() io.Reader {
	return e.rnd
}

func (e *Client) close(closer io.Closer, flags map[string]string) {
	thr := closer.Close()
	if thr != nil {
		thr = errors.Wrapf(thr, "trying to close: %T", closer)
		// TODO (novln): Add an observer to collect this error.
		_ = thr
	}
}
