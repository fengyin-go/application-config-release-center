package cache

type Cache struct{ current map[string]string }

func (c *Cache) Save(values map[string]string) { c.current = values }
func (c *Cache) Load() map[string]string       { return c.current }
