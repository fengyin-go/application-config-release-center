package cache

import "configcenter/internal/retrystate/model"

type Cache struct{ jobs map[string]model.Job }

func New() *Cache                                { return &Cache{jobs: map[string]model.Job{}} }
func (c *Cache) Save(job model.Job)              { c.jobs[job.ID] = job }
func (c *Cache) Get(id string) (model.Job, bool) { job, ok := c.jobs[id]; return job, ok }
