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
func (h Handler) Publish(key string, result *model.Bundle, err error) error {
	h.Cache.Put(key, result)
	return err
}
func (h Handler) Read(key string) (*model.Bundle, bool) { return h.Cache.Get(key) }
