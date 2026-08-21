package cache

import "configcenter/internal/retrystate/model"

type Cache struct{ jobs map[string]model.Job }

func New() *Cache                                  { return &Cache{jobs: map[string]model.Job{}} }
func (c *Cache) Get(id string) (model.Job, bool)   { job, ok := c.jobs[id]; return job, ok }

// Save persists the job, but only if it is not older than what the cache
// already holds. A late callback from an earlier attempt carries a stale
// version; without this guard it would overwrite the success view already
// cached from the successful (later) attempt. The list view thus moves only
// forward, matching the job's own monotonic transitions.
func (c *Cache) Save(job model.Job) {
	if existing, ok := c.jobs[job.ID]; ok && job.Version < existing.Version {
		return
	}
	c.jobs[job.ID] = job
}
