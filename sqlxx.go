package sqlxx

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

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
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Queryx(query string, args ...interface{}) (*sqlx.Rows, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	PrepareNamed(query string) (Statement, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	Close() error
	Ping() error
	Beginx() (Driver, error)
	Rollback() error
	Commit() error
}

// Statement from PrepareNamed.
type Statement interface {
	// Close closes the statement.
	Close() error
	// Exec executes a named statement using the struct passed.
	Exec(arg interface{}) (sql.Result, error)
	// Query using this Statement.
	Query(dest interface{}) (*sql.Rows, error)
	// Get using this Statement.
	Get(dest interface{}, arg interface{}) error
	// Select using this Statement.
	Select(dest interface{}, arg interface{}) error
}

// Model represents a database table.
type Model interface {
	TableName() string
}

// Cache is sqlxx cache.
type Cache struct {
	mu      sync.RWMutex
	schemas map[string]Schema
}

// NewCache returns new cache instance.
func NewCache() *Cache {
	return &Cache{
		schemas: map[string]Schema{},
	}
}

// SetSchema caches the given schema.
func (c *Cache) SetSchema(schema Schema) {
	c.mu.Lock()
	c.schemas[schema.TableName] = schema
	c.mu.Unlock()
}

// Flush flushs the cache
func (c *Cache) Flush() {
	c.mu.Lock()
	c.schemas = map[string]Schema{}
	c.mu.Unlock()
}

// GetSchema returns the given schema from cache.
// If the given schema does not exists, returns false as bool.
func (c *Cache) GetSchema(model Model) (Schema, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	schema, ok := c.schemas[model.TableName()]
	if !ok {
		return Schema{}, false
	}

	return schema, true
}
