package collector_test

import (
	"configcenter/internal/fanout/collector"
	"configcenter/internal/fanout/consumer"
	"configcenter/internal/fanout/coordinator"
	"configcenter/internal/fanout/producer"
	"context"
	"strings"
	"testing"
	"time"
)

func TestFanoutErrorClosesEveryLifecycle(t *testing.T) {
	direct := producer.Start([]string{"bad"})
	select {
	case <-direct.Errors:
	case <-time.After(20 * time.Millisecond):
		t.Error("producer did not publish its error")
	}
	select {
	case _, ok := <-direct.Results:
		if ok {
			t.Error("result stream remained open")
		}
	case <-time.After(20 * time.Millisecond):
		t.Error("result stream was not closed")
	}
	select {
	case _, ok := <-direct.Errors:
		if ok {
			t.Error("error stream remained open")
		}
	case <-time.After(20 * time.Millisecond):
		t.Error("error stream was not closed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	pipeline := make(chan error, 1)
	go func() { pipeline <- collector.Await(ctx, consumer.Drain(ctx, producer.Start([]string{"ok", "bad"}))) }()
	select {
	case err := <-pipeline:
		if err == nil || !strings.Contains(err.Error(), "rejected") {
			t.Errorf("pipeline error=%v", err)
		}
	case <-time.After(40 * time.Millisecond):
		t.Error("producer error left consumer and collector blocked")
	}

	gate := make(chan struct{})
	done, started := coordinator.Start(gate)
	select {
	case <-done:
		t.Error("coordinator completed before worker started")
	case <-time.After(10 * time.Millisecond):
	}
	close(gate)
	<-started
	select {
	case <-done:
	case <-time.After(20 * time.Millisecond):
		t.Error("coordinator did not complete")
	}

	cancelled, stop := context.WithCancel(context.Background())
	stop()
	never := make(chan error)
	returned := make(chan error, 1)
	go func() { returned <- collector.Await(cancelled, never) }()
	select {
	case <-returned:
	case <-time.After(20 * time.Millisecond):
		t.Error("collector ignored cancellation")
	}
}
