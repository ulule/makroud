package sqlxx

import "sync"

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

	schema, ok := c.schemas[model.TableName()]
	if !ok {
		c.mu.RUnlock()
		return Schema{}, false
	}

	c.mu.RUnlock()

	return schema, true
}
