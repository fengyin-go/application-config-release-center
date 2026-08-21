package worker

import (
	"configcenter/internal/canceljob/client"
	"context"
)

type Worker struct{ Client *client.Client }

func (w *Worker) Attempt(ctx context.Context) { w.Client.Send() }
