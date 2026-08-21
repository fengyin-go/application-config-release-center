package service_test

import (
	"configcenter/internal/canceljob/client"
	"configcenter/internal/canceljob/dispatcher"
	"configcenter/internal/canceljob/scheduler"
	"configcenter/internal/canceljob/service"
	"configcenter/internal/canceljob/worker"
	"context"
	"testing"
	"time"
)

func TestCancelledRequestStopsRetriesAndShutdown(t *testing.T) {
	gate := make(chan struct{}, 3)
	downstream := client.New()
	job := &service.Service{Dispatcher: &dispatcher.Dispatcher{Scheduler: &scheduler.Scheduler{Worker: &worker.Worker{Client: downstream}, Gate: gate}}}
	ctx, cancel := context.WithCancel(context.Background())
	job.Start(ctx)
	gate <- struct{}{}
	<-downstream.Called
	cancel()
	if job.Shutdown(20*time.Millisecond) != true {
		t.Error("shutdown waited on a cancelled request")
	}
	gate <- struct{}{}
	gate <- struct{}{}
	time.Sleep(10 * time.Millisecond)
	if downstream.Calls != 1 {
		t.Errorf("downstream calls grew to %d after cancellation", downstream.Calls)
	}
}
