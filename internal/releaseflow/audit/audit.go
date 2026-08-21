package audit

type Log struct{ Entries []string }

func (l *Log) Add(status string) { l.Entries = append(l.Entries, status) }
