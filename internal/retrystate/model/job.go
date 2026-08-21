package model

// stateRank orders lifecycle states so transitions can only move forward.
// A lower rank may not overwrite a higher one: "success" is terminal, so a
// late "running" callback from an earlier attempt cannot regress it.
var stateRank = map[string]int{
	"":        0,
	"pending": 1,
	"running": 2,
	"failed":  3,
	"success": 4,
}

type Job struct {
	ID, State string
	Version   int
}

// Apply transitions the job to (version, state) but only forward. A transition
// is rejected (returns false) when either:
//   - the incoming version is older than the current one (a late callback from
//     an earlier attempt), or
//   - the incoming state ranks no higher than the current one (e.g. a "running"
//     callback arriving after "success").
//
// This keeps "success" terminal and the version monotonic, so out-of-order
// callbacks cannot rewind a job that has already progressed.
func (j *Job) Apply(version int, state string) bool {
	if version < j.Version {
		return false
	}
	if stateRank[state] < stateRank[j.State] {
		return false
	}
	j.Version = version
	j.State = state
	return true
}
