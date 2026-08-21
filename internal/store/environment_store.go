package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateEnvironment(e *model.Environment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.environments {
		if exist.AppID == e.AppID && exist.Code == e.Code {
			return ErrConflict
		}
	}
	s.environments[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEnvironment(id string) (*model.Environment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.environments[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListEnvironments() []*model.Environment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Environment, 0, len(s.environments))
	for _, e := range s.environments {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) ListEnvironmentsByAppID(appID string) []*model.Environment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Environment, 0)
	for _, e := range s.environments {
		if e.AppID == appID {
			list = append(list, e)
		}
	}
	return list
}

func (s *MemoryStore) UpdateEnvironment(e *model.Environment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.environments[e.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.environments {
		if exist.ID != e.ID && exist.AppID == e.AppID && exist.Code == e.Code {
			return ErrConflict
		}
	}
	s.environments[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEnvironment(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.environments[id]; !ok {
		return ErrNotFound
	}
	delete(s.environments, id)
	return nil
}
