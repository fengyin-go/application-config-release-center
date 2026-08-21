package coordinator

import "sync"

func Start(gate <-chan struct{}) (<-chan struct{}, <-chan struct{}) {
	var wg sync.WaitGroup
	done := make(chan struct{})
	started := make(chan struct{})
	go func() { <-gate; wg.Add(1); close(started); wg.Done() }()
	go func() { wg.Wait(); close(done) }()
	return done, started
}
