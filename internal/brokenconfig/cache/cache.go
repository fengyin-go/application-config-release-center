package cache

import "configcenter/internal/brokenconfig/model"

type Cache struct{ values map[string]*model.Bundle }

func New() *Cache                                     { return &Cache{values: map[string]*model.Bundle{}} }
func (c *Cache) Put(key string, value *model.Bundle)  { c.values[key] = value }
func (c *Cache) Get(key string) (*model.Bundle, bool) { value, ok := c.values[key]; return value, ok }
