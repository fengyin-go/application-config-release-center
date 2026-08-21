package producer

import "errors"

type Stream struct {
	Results <-chan string
	Errors  <-chan error
}

// Start launches a producer that streams each item on Results.
// A "bad" item is reported as a single error on Errors; the producer
// then stops. Both channels are always closed when the producer exits,
// and Errors is buffered so an error can be published without a reader
// already waiting — this is what keeps the error path from blocking.
func Start(items []string) Stream {
	results := make(chan string)
	failures := make(chan error, 1)
	go func() {
		defer close(results)
		defer close(failures)
		for _, item := range items {
			if item == "bad" {
				failures <- errors.New("config source rejected")
				return
			}
			results <- item
		}
	}()
	return Stream{Results: results, Errors: failures}
}
