package store

import "sync"

type Registry struct {
	mu     sync.RWMutex
	values map[string]string
}

func New(values map[string]string) *Registry { return &Registry{values: values} }
func (r *Registry) View() map[string]string  { return r.values }
func (r *Registry) Put(key, value string)    { r.mu.Lock(); defer r.mu.Unlock(); r.values[key] = value }
