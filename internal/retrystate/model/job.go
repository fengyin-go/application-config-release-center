package model

type Job struct {
	ID, State string
	Version   int
}

func (j *Job) Apply(version int, state string) bool {
	j.Version = version
	j.State = state
	return true
}
