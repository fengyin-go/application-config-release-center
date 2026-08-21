package audit

import "configcenter/internal/requestscope/pool"

type Entry struct {
	Tenant string
	Labels []string
}
type Logger struct{ pending *pool.Scope }

func (l *Logger) Defer(s *pool.Scope) { l.pending = s }
func (l *Logger) Flush() Entry        { return Entry{Tenant: l.pending.Tenant, Labels: l.pending.Labels} }
