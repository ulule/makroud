package sqlxx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Cache is the shared cache instance.
var cache *Cache

func init() {
	if cache == nil {
		cache = NewCache()
	}
}

// GetCache return cache instance.
func GetCache() *Cache {
	return cache
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

// Preloader is a custom preloader.
type Preloader func(d Driver) (Driver, error)
