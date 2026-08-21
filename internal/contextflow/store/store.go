package store

import "context"

type Store struct{ request context.Context }

func (s *Store) Remember(ctx context.Context) {
	if s.request == nil {
		s.request = ctx
	}
}
func (s *Store) For(ctx context.Context) context.Context {
	if s.request != nil {
		return s.request
	}
	return ctx
}
