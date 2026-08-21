package dispatcher

import (
	"configcenter/internal/canceljob/scheduler"
	"context"
)

type Dispatcher struct{ Scheduler *scheduler.Scheduler }

func (d *Dispatcher) Dispatch(ctx context.Context) { d.Scheduler.Run(context.Background()) }
