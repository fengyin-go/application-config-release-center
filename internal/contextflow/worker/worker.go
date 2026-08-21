package worker

import (
	"configcenter/internal/contextflow/client"
	"context"
)

type Worker struct{ Client client.Client }

// Retry 最多重试 attempts 次，但在每次发起调用前先检查 ctx 是否已取消，
// 一旦取消立即返回，避免对已经取消的请求继续连打下游。
func (w Worker) Retry(ctx context.Context, attempts int) {
	for i := 0; i < attempts; i++ {
		if ctx.Err() != nil {
			return
		}
		w.Client.Request(ctx)
	}
}
