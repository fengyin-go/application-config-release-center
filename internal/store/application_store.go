package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateApplication(a *model.Application) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.applications {
		if exist.Code == a.Code {
			return ErrConflict
		}
	}
	s.applications[a.ID] = a
	return nil
}

func (s *MemoryStore) GetApplication(id string) (*model.Application, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.applications[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) GetApplicationByCode(code string) (*model.Application, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.applications {
		if a.Code == code {
			return a, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListApplications() []*model.Application {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Application, 0, len(s.applications))
	for _, a := range s.applications {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateApplication(a *model.Application) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.applications[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.applications {
		if exist.ID != a.ID && exist.Code == a.Code {
			return ErrConflict
		}
	}
	s.applications[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteApplication(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.applications[id]; !ok {
		return ErrNotFound
	}
	delete(s.applications, id)
	return nil
}
