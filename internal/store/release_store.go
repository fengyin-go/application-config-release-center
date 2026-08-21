package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateRelease(r *model.Release) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.releases[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRelease(id string) (*model.Release, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.releases[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListReleases() []*model.Release {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Release, 0, len(s.releases))
	for _, r := range s.releases {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) ListReleasesByAppID(appID string) []*model.Release {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Release, 0)
	for _, r := range s.releases {
		if r.AppID == appID {
			list = append(list, r)
		}
	}
	return list
}

func (s *MemoryStore) ListReleasesByEnvID(envID string) []*model.Release {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Release, 0)
	for _, r := range s.releases {
		if r.EnvID == envID {
			list = append(list, r)
		}
	}
	return list
}

func (s *MemoryStore) UpdateRelease(r *model.Release) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.releases[r.ID]; !ok {
		return ErrNotFound
	}
	s.releases[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRelease(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.releases[id]; !ok {
		return ErrNotFound
	}
	delete(s.releases, id)
	return nil
}
