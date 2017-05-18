package sqlxx

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/heetch/sqalx"
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

// AssociationType is an association type.
type AssociationType uint8

func (a AssociationType) String() string {
	return map[AssociationType]string{
		AssociationTypeUndefined:  "undefined",
		AssociationTypeOne:        "one",
		AssociationTypeMany:       "many",
		AssociationTypeManyToMany: "many-to-many",
	}[a]
}

// Association types
const (
	AssociationTypeUndefined = AssociationType(iota)
	AssociationTypeOne
	AssociationTypeMany
	AssociationTypeManyToMany
)

// Constants
const (
	StructTagName       = "sqlxx"
	SQLXStructTagName   = "db"
	StructTagPrimaryKey = "primary_key"
	StructTagIgnored    = "ignored"
	StructTagDefault    = "default"
	StructTagForeignKey = "fk"
	StructTagSQLXField  = "field"
)

// PrimaryKeyFieldName is the default field name for primary keys
const PrimaryKeyFieldName = "ID"

// SupportedTags are supported tags.
var SupportedTags = []string{
	StructTagName,
	SQLXStructTagName,
}

// TagsMapping is the reflekt.Tags mapping to handle struct tag without key:value format
var TagsMapping = map[string]string{
	"db": "field",
}

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
}

// Model represents a database table.
type Model interface {
	TableName() string
}
