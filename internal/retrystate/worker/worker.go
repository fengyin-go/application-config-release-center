package worker

import "configcenter/internal/retrystate/model"

// Complete marks the job as succeeded at the given (forward-moving) version.
// Because Apply is monotonic, a successful terminal state can never be
// regressed by a later callback.
func Complete(job *model.Job, version int) bool {
	return job.Apply(version, "success")
}

// LateCallback processes an out-of-order callback from an earlier attempt. It
// routes through the same monotonic Apply as Complete, so if the job has
// already reached "success" at a newer version this becomes a no-op rather
// than rewinding state and version.
func LateCallback(job *model.Job, version int) bool {
	return job.Apply(version, "running")
}
