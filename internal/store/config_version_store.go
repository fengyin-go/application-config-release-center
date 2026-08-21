package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateConfigVersion(c *model.ConfigVersion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	seq := s.configItemVerSeq[c.ConfigItemID] + 1
	c.Version = seq
	s.configItemVerSeq[c.ConfigItemID] = seq
	s.configVersions[c.ID] = c
	return nil
}

func (s *MemoryStore) GetConfigVersion(id string) (*model.ConfigVersion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.configVersions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListConfigVersions() []*model.ConfigVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConfigVersion, 0, len(s.configVersions))
	for _, c := range s.configVersions {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) ListConfigVersionsByConfigItemID(configItemID string) []*model.ConfigVersion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConfigVersion, 0)
	for _, c := range s.configVersions {
		if c.ConfigItemID == configItemID {
			list = append(list, c)
		}
	}
	return list
}
