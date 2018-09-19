package makroud

import (
	"sync"
)

// cache is driver cache used to store Schema information.
type cache struct {
	schemas sync.Map
}

// newCache returns new cache instance.
func newCache() *cache {
	return &cache{
		schemas: sync.Map{},
	}
}

// GetSchema returns the given schema from cache.
// If the given schema does not exists, returns false as bool.
func (c *cache) GetSchema(model Model) *Schema {
	schema, ok := c.schemas.Load(model.TableName())
	if !ok {
		return nil
	}

	return schema.(*Schema)
}

// SetSchema caches the given schema.
func (c *cache) SetSchema(schema *Schema) {
	if schema == nil {
		panic("makroud: schema shouldn't be nil")
	}
	c.schemas.Store(schema.TableName(), schema)
}
