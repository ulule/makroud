package makroud

import (
	"sync"
)

// cache is driver cache used to store Schema information.
type cache struct {
	mutex   sync.RWMutex
	schemas map[string]*Schema
}

// newCache returns new cache instance.
func newCache() *cache {
	return &cache{
		schemas: map[string]*Schema{},
	}
}

// GetSchema returns the given schema from cache.
// If the given schema does not exists, returns false as bool.
func (c *cache) GetSchema(model Model) *Schema {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	schema, ok := c.schemas[model.TableName()]
	if !ok {
		return nil
	}

	return schema
}

// SetSchema caches the given schema.
func (c *cache) SetSchema(schema *Schema) {
	if schema == nil {
		panic("makroud: schema shouldn't be nil")
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.schemas[schema.TableName()] = schema
}

// Flush flushs the cache.
func (c *cache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.schemas = map[string]*Schema{}
}
