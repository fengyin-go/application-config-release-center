package coordinator

// Start waits for gate to close, signals that the worker has started, then
// signals completion.
//
// started is guaranteed to be closed before done: both happen in the same
// goroutine, in order. The previous implementation used a sync.WaitGroup
// whose Add(1) raced with a separate Wait() goroutine — Wait could observe
// a zero counter and close done before the worker's Add(1) ever ran.
func Start(gate <-chan struct{}) (<-chan struct{}, <-chan struct{}) {
	done := make(chan struct{})
	started := make(chan struct{})
	go func() {
		<-gate
		close(started)
		close(done)
	}()
	return done, started
}
