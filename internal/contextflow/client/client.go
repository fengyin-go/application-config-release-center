package client

import "context"

type Backend struct {
	Seen     context.Context
	Attempts int
	Calls    int
}

func (b *Backend) Call(ctx context.Context) { b.Seen = ctx; b.Calls++ }

type Client struct{ Backend *Backend }

// Request 将真正的 ctx 透传给下游，而不是用 context.Background() 覆盖，
// 这样下游才能看到请求的取消/超时状态并据此中止。
func (c Client) Request(ctx context.Context) {
	c.Backend.Attempts++
	c.Backend.Call(ctx)
}
