package resource

import "errors"

type Pool struct{ limit, open int }
type Lease struct {
	pool   *Pool
	closed bool
}

func New(limit int) *Pool { return &Pool{limit: limit} }
func (p *Pool) Acquire() (*Lease, error) {
	if p.open >= p.limit {
		return nil, errors.New("release resource limit reached")
	}
	p.open++
	return &Lease{pool: p}, nil
}
func (l *Lease) Close() {
	if !l.closed {
		l.closed = true
		l.pool.open--
	}
}
