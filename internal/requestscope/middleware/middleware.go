package middleware

import "configcenter/internal/requestscope/pool"

type Binder struct{ Pool *pool.Pool }

func (b *Binder) Begin(tenant, label string) *pool.Scope {
	s := b.Pool.Acquire()
	s.Tenant = tenant
	s.Labels = append(s.Labels, label)
	return s
}
