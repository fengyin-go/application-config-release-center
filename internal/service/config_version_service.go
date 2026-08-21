package service

import (
	"sort"

	"configcenter/internal/model"
)

func (s *Service) ListConfigVersions(filter model.ConfigVersionFilter, page, size int) ([]*model.ConfigVersion, int, error) {
	all := s.store.ListConfigVersions()
	matched := make([]*model.ConfigVersion, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Version > matched[j].Version
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ConfigVersion{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
