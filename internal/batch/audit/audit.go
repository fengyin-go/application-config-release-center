package audit

type Log struct{ Status string }

func (l *Log) Success() { l.Status = "success" }
func (l *Log) Failure() { l.Status = "failed" }
