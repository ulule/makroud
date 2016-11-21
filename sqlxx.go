package sqlxx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

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

// Model represents a database table.
type Model interface {
	TableName() string
}

// GetByParams executes a WHERE with params and populates the given model
// instance with related data.
func GetByParams(driver Driver, out Model, params map[string]string) error {
	return nil
}

// Preload preloads related fields.
func Preload(driver Driver, out Model, related ...string) error {
	return nil
}
