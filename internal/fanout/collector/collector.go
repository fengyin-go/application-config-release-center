package collector

import "context"

// Await blocks until either the pipeline reports a result on done or ctx is
// cancelled — whichever happens first. Honoring cancellation is what lets a
// cancelled pipeline wake the collector instead of blocking forever on done.
func Await(ctx context.Context, done <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-done:
		return err
	}
}
