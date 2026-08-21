package collector

import "context"

func Await(ctx context.Context, done <-chan error) error { return <-done }
