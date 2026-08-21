package service_test

import (
	"configcenter/internal/retrystate/cache"
	"configcenter/internal/retrystate/model"
	"configcenter/internal/retrystate/service"
	"configcenter/internal/retrystate/worker"
	"testing"
)

func TestRetryKeepsOneEffectAndMonotonicSuccess(t *testing.T) {
	effects := service.NewEffects()
	effects.Trigger("release-7", 1)
	effects.Trigger("release-7", 2)
	if effects.Total() != 1 {
		t.Errorf("external release executed %d times", effects.Total())
	}
	job := &model.Job{ID: "release-7", State: "failed", Version: 1}
	worker.Complete(job, 2)
	worker.LateCallback(job, 1)
	if job.State != "success" || job.Version != 2 {
		t.Errorf("late callback changed job to state=%s version=%d", job.State, job.Version)
	}
	c := cache.New()
	c.Save(model.Job{ID: job.ID, State: "success", Version: 2})
	c.Save(model.Job{ID: job.ID, State: "running", Version: 1})
	view, _ := c.Get(job.ID)
	if view.State != "success" || view.Version != 2 {
		t.Errorf("cached list regressed to state=%s version=%d", view.State, view.Version)
	}
}
