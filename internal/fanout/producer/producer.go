package producer

import "errors"

type Stream struct {
	Results <-chan string
	Errors  <-chan error
}

func Start(items []string) Stream {
	results := make(chan string)
	failures := make(chan error)
	go func() {
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
