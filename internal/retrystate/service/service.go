package service

// Effects records external side effects (e.g. the external publish call).
// A retried job must produce the side effect at most once regardless of how
// many attempts it took to succeed, so the key is the job alone — not the
// attempt. Triggering the same job again is a no-op.
type Effects struct{ triggered map[string]bool }

func NewEffects() *Effects { return &Effects{triggered: map[string]bool{}} }

// Trigger records that jobID fired its external side effect. attempt is kept
// in the signature for API compatibility but does not gate the effect: the
// second trigger (the retry) is suppressed, so Total counts the job once.
func (e *Effects) Trigger(jobID string, attempt int) {
	_ = attempt // attempt is intentionally ignored; see Trigger doc comment
	if e.triggered[jobID] {
		return
	}
	e.triggered[jobID] = true
}

func (e *Effects) Total() int {
	total := 0
	for range e.triggered {
		total++
	}
	return total
}
