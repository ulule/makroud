package sqlxx

import (
	"github.com/heetch/sqalx"
)

// wrapClient creates a new Client using given database connection.
func wrapClient(connection sqalx.Node) Driver {
	return &Client{
		node: connection,
	}
}
