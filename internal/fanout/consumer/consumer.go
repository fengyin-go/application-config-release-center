package consumer

import (
	"configcenter/internal/fanout/producer"
	"context"
)

// Drain consumes a producer stream to completion and reports the first
// error (if any) on the returned channel.
//
// It selects on both Results and Errors concurrently rather than reading
// Results to completion first: the producer publishes its error on an
// unbuffered-ish path, so a results-first reader would deadlock waiting for
// Results to close while the producer is blocked publishing the error.
//
// Draining always runs to completion (the producer is finite and closes both
// channels), so ctx is not used for early exit — bailing mid-stream would
// leak a producer blocked on an unbuffered send. It is kept on the signature
// for symmetry with the rest of the pipeline.
func Drain(ctx context.Context, stream producer.Stream) <-chan error {
	done := make(chan error, 1)
	go func() {
		var firstErr error
		results, errs := stream.Results, stream.Errors
		for results != nil || errs != nil {
			select {
			case _, ok := <-results:
				if !ok {
					results = nil
				}
			case err, ok := <-errs:
				if !ok {
					errs = nil
					continue
				}
				if firstErr == nil {
					firstErr = err
				}
			}
		}
		done <- firstErr
	}()
	return done
}
