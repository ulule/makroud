package sqlxx

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/heetch/sqalx"
	"github.com/jmoiron/sqlx"
)

// ErrInvalidDriver is returned when given driver is undefined.
var ErrInvalidDriver = fmt.Errorf("a sqlxx driver is required")

// Driver can either be a *sqlx.DB or a *sqlx.Tx.
type Driver interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.Preparer
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	Close() error
	Ping() error
	Beginx() (sqalx.Node, error)
	Rollback() error
	Commit() error
	close(closer io.Closer, flags map[string]string)
	hasCache() bool
	cache() *cache
	logger() Logger
	entropy() io.Reader
}
