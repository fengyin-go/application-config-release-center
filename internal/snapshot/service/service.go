package service

import (
	"configcenter/internal/snapshot/cache"
	"configcenter/internal/snapshot/store"
)

type Publisher struct {
	registry *store.Registry
	cache    *cache.Cache
}

func New(r *store.Registry, c *cache.Cache) *Publisher { return &Publisher{registry: r, cache: c} }
func (p *Publisher) Publish() map[string]string {
	p.cache.Save(p.registry.View())
	return p.cache.Load()
}
func Checksum(values map[string]string) int {
	n := 0
	for k, v := range values {
		n += len(k) + len(v)
	}
	return n
}
