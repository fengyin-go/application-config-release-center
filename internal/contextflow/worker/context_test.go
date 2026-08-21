package worker_test

import (
	"configcenter/internal/contextflow/client"
	"configcenter/internal/contextflow/middleware"
	"configcenter/internal/contextflow/store"
	"configcenter/internal/contextflow/worker"
	"context"
	"testing"
	"time"
)

func TestRequestContextSurvivesEveryLayer(t *testing.T) {
	request, cancel := middleware.WithDeadline(context.Background(), 5*time.Millisecond)
	defer cancel()
	time.Sleep(10 * time.Millisecond)
	if request.Err() != context.DeadlineExceeded {
		t.Errorf("entry deadline was lost: %v", request.Err())
	}
	s := &store.Store{}
	expired, expire := context.WithCancel(context.Background())
	expire()
	s.Remember(expired)
	fresh := context.Background()
	if got := s.For(fresh); got != fresh {
		t.Errorf("next request inherited old error %v", got.Err())
	}
	backend := &client.Backend{}
	stopped, stop := context.WithCancel(context.Background())
	stop()
	worker.Worker{Client: client.Client{Backend: backend}}.Retry(stopped, 3)
	if backend.Attempts != 0 {
		t.Errorf("cancelled worker attempted %d retries", backend.Attempts)
	}
	if backend.Calls != 0 {
		t.Errorf("cancelled request made %d retry calls", backend.Calls)
	}
	probe := &client.Backend{}
	client.Client{Backend: probe}.Request(stopped)
	if probe.Seen == nil || probe.Seen.Err() != context.Canceled {
		t.Errorf("downstream saw context error %v", probe.Seen.Err())
	}
}
