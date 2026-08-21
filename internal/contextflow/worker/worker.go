package worker

import (
	"configcenter/internal/contextflow/client"
	"context"
)

type Worker struct{ Client client.Client }

func (w Worker) Retry(ctx context.Context, attempts int) {
	for i := 0; i < attempts; i++ {
		w.Client.Request(ctx)
	}
}
