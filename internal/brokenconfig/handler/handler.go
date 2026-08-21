package handler

import (
	"configcenter/internal/brokenconfig/cache"
	"configcenter/internal/brokenconfig/model"
	"configcenter/internal/brokenconfig/recovery"
)

type Handler struct {
	Recoverer recovery.Recoverer
	Cache     *cache.Cache
}

func (h Handler) Load(key string, target *model.Bundle, fail bool) error {
	result, err := h.Recoverer.Capture(target, fail)
	return h.Publish(key, result, err)
}
// Publish caches result only when the build succeeded and produced a complete
// bundle. A failed build (err != nil) must never pollute the cache with a
// half-finished object that a later Read would expose.
func (h Handler) Publish(key string, result *model.Bundle, err error) error {
	if err == nil {
		h.Cache.Put(key, result)
	}
	return err
}
func (h Handler) Read(key string) (*model.Bundle, bool) { return h.Cache.Get(key) }
