package service

import "fmt"

type Effects struct{ calls map[string]int }

func NewEffects() *Effects { return &Effects{calls: map[string]int{}} }
func (e *Effects) Trigger(jobID string, attempt int) {
	key := fmt.Sprintf("%s-%d", jobID, attempt)
	e.calls[key]++
}
func (e *Effects) Total() int {
	total := 0
	for _, count := range e.calls {
		total += count
	}
	return total
}
