package publisher

type TemporaryError struct{}

func (TemporaryError) Error() string { return "publisher temporarily unavailable" }

type Publisher struct{ Calls, Events int }

func (p *Publisher) Publish() error {
	p.Calls++
	p.Events++
	if p.Calls == 1 {
		return TemporaryError{}
	}
	return nil
}
