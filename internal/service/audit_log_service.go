package service

import (
	"sort"
	"time"

	"configcenter/internal/model"
	"configcenter/pkg/idgen"
)

func (s *Service) logAudit(action, resource, resourceID, detail string) {
	al := &model.AuditLog{
		ID:        idgen.Hex(),
		Operator:  "system",
		Action:    action,
		Resource:  resource,
		Detail:    detail,
		CreatedAt: time.Now(),
	}
	if resourceID != "" {
		al.Detail = resource + ":" + resourceID
		if detail != "" {
			al.Detail += " -> " + detail
		}
	}
	_ = s.store.CreateAuditLog(al)
}

func (s *Service) ListAuditLogs(filter model.AuditLogFilter, page, size int) ([]*model.AuditLog, int, error) {
	all := s.store.ListAuditLogs()
	matched := make([]*model.AuditLog, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AuditLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
