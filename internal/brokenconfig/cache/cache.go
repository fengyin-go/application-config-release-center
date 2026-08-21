package cache

import "configcenter/internal/brokenconfig/model"

type Cache struct{ values map[string]*model.Bundle }

func New() *Cache { return &Cache{values: map[string]*model.Bundle{}} }

// Put stores value only when it represents a complete, usable bundle.
// Partial results from a failed build are rejected so a later Read can never
// surface a half-finished object.
func (c *Cache) Put(key string, value *model.Bundle) {
	if value == nil || !value.Complete {
		return
	}
	c.values[key] = value
}

func (c *Cache) Get(key string) (*model.Bundle, bool) { value, ok := c.values[key]; return value, ok }
