package consumer

import (
	"configcenter/internal/fanout/producer"
	"context"
)

func Drain(ctx context.Context, stream producer.Stream) <-chan error {
	done := make(chan error, 1)
	go func() {
		for range stream.Results {
		}
		done <- nil
	}()
	return done
}
