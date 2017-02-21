package sqlxx

import (
	"database/sql"

	"os"

	"github.com/jmoiron/sqlx"
)

var (
	// Cache is the shared cache instance.
	cache *Cache
	// cacheDisabled is true if cache has been disabled
	cacheDisabled bool
)

func init() {
	if os.Getenv("SQLXX_DISABLE_CACHE") != "" {
		cacheDisabled = true
		return
	}

	if cache == nil {
		cache = NewCache()
	}
}

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
}
