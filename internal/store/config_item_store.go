package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateConfigItem(c *model.ConfigItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.configItems {
		if exist.AppID == c.AppID && exist.EnvID == c.EnvID && exist.Key == c.Key {
			return ErrConflict
		}
	}
	s.configItems[c.ID] = c
	return nil
}

func (s *MemoryStore) GetConfigItem(id string) (*model.ConfigItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.configItems[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) ListConfigItems() []*model.ConfigItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConfigItem, 0, len(s.configItems))
	for _, c := range s.configItems {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) GetConfigItemByAppEnvKey(appID, envID, key string) (*model.ConfigItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.configItems {
		if c.AppID == appID && c.EnvID == envID && c.Key == key {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListConfigItemsByAppID(appID string) []*model.ConfigItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConfigItem, 0)
	for _, c := range s.configItems {
		if c.AppID == appID {
			list = append(list, c)
		}
	}
	return list
}

func (s *MemoryStore) ListConfigItemsByEnvID(envID string) []*model.ConfigItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ConfigItem, 0)
	for _, c := range s.configItems {
		if c.EnvID == envID {
			list = append(list, c)
		}
	}
	return list
}

func (s *MemoryStore) UpdateConfigItem(c *model.ConfigItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.configItems[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.configItems {
		if exist.ID != c.ID && exist.AppID == c.AppID && exist.EnvID == c.EnvID && exist.Key == c.Key {
			return ErrConflict
		}
	}
	s.configItems[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteConfigItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.configItems[id]; !ok {
		return ErrNotFound
	}
	delete(s.configItems, id)
	return nil
}
