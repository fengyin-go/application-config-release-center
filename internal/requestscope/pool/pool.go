package pool

type Scope struct {
	Tenant string
	Labels []string
}
type Pool struct{ idle *Scope }

func (p *Pool) Acquire() *Scope {
	if p.idle != nil {
		s := p.idle
		p.idle = nil
		return s
	}
	return &Scope{}
}
func (p *Pool) Release(s *Scope) { p.idle = s }
