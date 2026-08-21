package store

import "configcenter/internal/model"

func (s *MemoryStore) CreateAuditLog(a *model.AuditLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auditLogs[a.ID] = a
	return nil
}

func (s *MemoryStore) GetAuditLog(id string) (*model.AuditLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.auditLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListAuditLogs() []*model.AuditLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AuditLog, 0, len(s.auditLogs))
	for _, a := range s.auditLogs {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) ListAuditLogsByAction(action string) []*model.AuditLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AuditLog, 0)
	for _, a := range s.auditLogs {
		if a.Action == action {
			list = append(list, a)
		}
	}
	return list
}

func (s *MemoryStore) ListAuditLogsByResource(resource string) []*model.AuditLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.AuditLog, 0)
	for _, a := range s.auditLogs {
		if a.Resource == resource {
			list = append(list, a)
		}
	}
	return list
}
