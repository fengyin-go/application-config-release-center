package worker

import "configcenter/internal/retrystate/model"

func Complete(job *model.Job, version int)     { job.Version = version; job.State = "success" }
func LateCallback(job *model.Job, version int) { job.Version = version; job.State = "running" }
