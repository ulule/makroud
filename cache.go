package sqlxx

import (
	"sync"
)

// cache is sqlxx cache.
type cache struct {
	mutex   sync.RWMutex
	schemas map[string]Schema
}

// newCache returns new cache instance.
func newCache() *cache {
	return &cache{
		schemas: map[string]Schema{},
	}
}

// SetSchema caches the given schema.
func (c *cache) SetSchema(schema Schema) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.schemas[schema.TableName] = schema
}

// Flush flushs the cache
func (c *cache) Flush() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.schemas = map[string]Schema{}
}

// GetSchema returns the given schema from cache.
// If the given schema does not exists, returns false as bool.
func (c *cache) GetSchema(model Model) (Schema, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	schema, ok := c.schemas[model.TableName()]
	if !ok {
		return Schema{}, false
	}

	return schema, true
}
