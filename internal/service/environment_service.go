package service

import (
	"sort"
	"time"

	"configcenter/internal/model"
	"configcenter/pkg/idgen"
)

func (s *Service) CreateEnvironment(input model.Environment) (*model.Environment, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetApplication(input.AppID); err != nil {
		return nil, err
	}
	now := time.Now()
	e := &model.Environment{
		ID:          idgen.Hex(),
		AppID:       input.AppID,
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateEnvironment(e); err != nil {
		return nil, err
	}
	s.logAudit("create", "environment", e.ID, "")
	return e, nil
}

func (s *Service) GetEnvironment(id string) (*model.Environment, error) {
	return s.store.GetEnvironment(id)
}

func (s *Service) ListEnvironments(filter model.EnvironmentFilter, page, size int) ([]*model.Environment, int, error) {
	all := s.store.ListEnvironments()
	matched := make([]*model.Environment, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Environment{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEnvironment(id string, input model.Environment) (*model.Environment, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e, err := s.store.GetEnvironment(id)
	if err != nil {
		return nil, err
	}
	if e.AppID != input.AppID {
		if _, err := s.store.GetApplication(input.AppID); err != nil {
			return nil, err
		}
	}
	e.AppID = input.AppID
	e.Name = input.Name
	e.Code = input.Code
	e.Description = input.Description
	e.UpdatedAt = time.Now()
	if err := s.store.UpdateEnvironment(e); err != nil {
		return nil, err
	}
	s.logAudit("update", "environment", e.ID, "")
	return e, nil
}

func (s *Service) DeleteEnvironment(id string) error {
	if err := s.store.DeleteEnvironment(id); err != nil {
		return err
	}
	s.logAudit("delete", "environment", id, "")
	return nil
}

func (s *Service) GetEnvironmentConfigItems(envID string) ([]*model.ConfigItem, error) {
	if _, err := s.store.GetEnvironment(envID); err != nil {
		return nil, err
	}
	all := s.store.ListConfigItemsByEnvID(envID)
	res := make([]*model.ConfigItem, 0, len(all))
	for _, c := range all {
		res = append(res, s.maskedConfigItem(c))
	}
	return res, nil
}

func (s *Service) GetEnvironmentReleases(envID string) ([]*model.Release, error) {
	if _, err := s.store.GetEnvironment(envID); err != nil {
		return nil, err
	}
	return s.store.ListReleasesByEnvID(envID), nil
}
