package worker

import "configcenter/internal/snapshot/store"

type Updater struct{ registry *store.Registry }

func New(r *store.Registry) *Updater       { return &Updater{registry: r} }
func (u *Updater) Apply(key, value string) { u.registry.Put(key, value) }
