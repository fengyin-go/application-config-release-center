package scheduler

import (
	"configcenter/internal/canceljob/worker"
	"context"
)

type Scheduler struct {
	Worker *worker.Worker
	Gate   <-chan struct{}
}

func (s *Scheduler) Run(ctx context.Context) {
	for i := 0; i < 3; i++ {
		<-s.Gate
		s.Worker.Attempt(ctx)
	}
}
