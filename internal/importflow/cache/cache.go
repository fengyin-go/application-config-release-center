package cache

type Cache struct{ batches map[string][]byte }

func New() *Cache                            { return &Cache{batches: map[string][]byte{}} }
func (c *Cache) Save(id string, data []byte) { c.batches[id] = data }
func (c *Cache) Get(id string) []byte        { return c.batches[id] }
