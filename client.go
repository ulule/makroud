package sqlxx

import (
	"fmt"

	"github.com/heetch/sqalx"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ClientDriver define the driver name used in sqlxx.
const ClientDriver = "postgres"

// ErrInvalidClient is returned when given client is undefined.
var ErrInvalidClient = errors.New("a sqlxx client is required")

// Client is a wrapper that can interact with the database.
type Client struct {
	sqalx.Node
	option clientOption
}

type clientOption struct {
	port               int
	host               string
	user               string
	password           string
	dbname             string
	sslMode            string
	timezone           string
	maxOpenConnections int
	maxIdleConnections int
}

func (e *clientOption) String() string {
	return fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s;timezone=%s",
		ClientDriver,
		e.user,
		e.password,
		e.host,
		e.port,
		e.dbname,
		e.sslMode,
		e.timezone,
	)
}

// New returns a new Client instance.
func New(options ...Option) (*Client, error) {

	client := &Client{}
	client.init()

	for _, option := range options {
		err := option.apply(client)
		if err != nil {
			return nil, err
		}
	}

	dbx, err := sqlx.Connect(ClientDriver, client.option.String())
	if err != nil {
		return nil, errors.Wrapf(err, "cannot connect to %s server", ClientDriver)
	}

	dbx.SetMaxIdleConns(client.option.maxIdleConnections)
	dbx.SetMaxOpenConns(client.option.maxOpenConnections)

	connection, err := sqalx.New(dbx)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot instantiate %s client driver", ClientDriver)
	}

	client.Node = connection

	return client, nil
}

// init configures default parameters for a Client.
func (e *Client) init() {
	if e.option.port == 0 {
		e.option.port = 5432
	}

	if e.option.host == "" {
		e.option.host = "localhost"
	}

	if e.option.user == "" {
		e.option.user = "postgres"
	}

	if e.option.dbname == "" {
		e.option.dbname = e.option.user
	}

	if e.option.sslMode == "" {
		e.option.sslMode = "disable"
	}

	if e.option.timezone == "" {
		e.option.timezone = "UTC"
	}

	if e.option.maxOpenConnections == 0 {
		e.option.maxOpenConnections = 5
	}

	if e.option.maxIdleConnections == 0 {
		e.option.maxIdleConnections = 2
	}
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
		return errors.Wrap(err, "cannot ping database")
	}
	return nil
}

// copy will create a client clone with given connection.
func (e *Client) copy(connection sqalx.Node) *Client {
	if connection == nil {
		panic("sqlxx: connection is required")
	}

	return &Client{
		Node:   connection,
		option: e.option,
	}
}
