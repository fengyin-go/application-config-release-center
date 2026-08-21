package service

import (
	"configcenter/internal/requestscope/audit"
	"configcenter/internal/requestscope/pool"
)

type Lifecycle struct {
	Pool   *pool.Pool
	Logger *audit.Logger
}

func (l *Lifecycle) End(s *pool.Scope) { l.Logger.Defer(s); l.Pool.Release(s) }
