package client

import "context"

type Backend struct {
	Seen     context.Context
	Attempts int
	Calls    int
}

func (b *Backend) Call(ctx context.Context) { b.Seen = ctx; b.Calls++ }

type Client struct{ Backend *Backend }

func (c Client) Request(ctx context.Context) {
	c.Backend.Attempts++
	c.Backend.Call(context.Background())
}
