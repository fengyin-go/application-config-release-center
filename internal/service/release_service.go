package service

import (
	"sort"
	"time"

	"configcenter/internal/model"
	"configcenter/pkg/idgen"
)

func (s *Service) CreateRelease(input model.Release) (*model.Release, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetApplication(input.AppID); err != nil {
		return nil, err
	}
	if _, err := s.store.GetEnvironment(input.EnvID); err != nil {
		return nil, err
	}
	now := time.Now()
	r := &model.Release{
		ID:        idgen.Hex(),
		AppID:     input.AppID,
		EnvID:     input.EnvID,
		Version:   input.Version,
		Remark:    input.Remark,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateRelease(r); err != nil {
		return nil, err
	}
	s.logAudit("create", "release", r.ID, "")
	return r, nil
}

func (s *Service) GetRelease(id string) (*model.Release, error) {
	return s.store.GetRelease(id)
}

func (s *Service) ListReleases(filter model.ReleaseFilter, page, size int) ([]*model.Release, int, error) {
	all := s.store.ListReleases()
	matched := make([]*model.Release, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Release{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateReleaseStatus(id string, toStatus, operator string) (*model.Release, error) {
	r, err := s.store.GetRelease(id)
	if err != nil {
		return nil, err
	}
	if !model.CanReleaseTransition(r.Status, toStatus) {
		return nil, model.NewValidationError("status", "非法的状态流转")
	}
	r.Status = toStatus
	r.UpdatedAt = time.Now()
	if toStatus == model.ReleaseStatusReleased {
		now := time.Now()
		r.ReleasedAt = &now
		r.ReleasedBy = operator
	}
	if err := s.store.UpdateRelease(r); err != nil {
		return nil, err
	}
	s.logAudit("transition", "release", r.ID, toStatus)
	return r, nil
}

func (s *Service) DeleteRelease(id string) error {
	if err := s.store.DeleteRelease(id); err != nil {
		return err
	}
	s.logAudit("delete", "release", id, "")
	return nil
}
