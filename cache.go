package makroud

import (
	"reflect"
	"sync"
)

// DriverCache is driver cache used to store Schema and Schemaless information.
type DriverCache struct {
	schemas    sync.Map
	schemaless sync.Map
}

// NewDriverCache returns new schema cache instance.
func NewDriverCache() *DriverCache {
	return &DriverCache{
		schemas:    sync.Map{},
		schemaless: sync.Map{},
	}
}

// GetSchema returns the schema associated to given model from cache.
// If it does not exists, it returns nil.
func (c *DriverCache) GetSchema(model Model) *Schema {
	schema, ok := c.schemas.Load(model.TableName())
	if !ok {
		return nil
	}

	return schema.(*Schema)
}

// SetSchema caches the given schema.
func (c *DriverCache) SetSchema(schema *Schema) {
	if schema == nil {
		panic("makroud: schema shouldn't be nil")
	}
	c.schemas.Store(schema.TableName(), schema)
}

// GetSchemaless returns the schemaless associated to type from cache.
// If it does not exists, it returns nil.
func (c *DriverCache) GetSchemaless(value reflect.Type) *Schemaless {
	schemaless, ok := c.schemaless.Load(value)
	if !ok {
		return nil
	}

	return schemaless.(*Schemaless)
}

// SetSchemaless caches the given schemaless.
func (c *DriverCache) SetSchemaless(schemaless *Schemaless) {
	if schemaless == nil {
		panic("makroud: schemaless shouldn't be nil")
	}
	c.schemaless.Store(schemaless.Type(), schemaless)
}
