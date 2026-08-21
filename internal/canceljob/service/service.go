package service

import (
	"configcenter/internal/canceljob/dispatcher"
	"context"
	"time"
)

type Service struct {
	Dispatcher *dispatcher.Dispatcher
	done       chan struct{}
}

func (s *Service) Start(ctx context.Context) {
	s.done = make(chan struct{})
	go func() { defer close(s.done); s.Dispatcher.Dispatch(ctx) }()
}
func (s *Service) Shutdown(timeout time.Duration) bool {
	select {
	case <-s.done:
		return true
	case <-time.After(timeout):
		return false
	}
}
