package service

import (
	"sort"
	"time"

	"configcenter/internal/model"
	"configcenter/pkg/idgen"
)

func (s *Service) CreateApplication(input model.Application) (*model.Application, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	a := &model.Application{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateApplication(a); err != nil {
		return nil, err
	}
	s.logAudit("create", "application", a.ID, "")
	return a, nil
}

func (s *Service) GetApplication(id string) (*model.Application, error) {
	return s.store.GetApplication(id)
}

func (s *Service) ListApplications(filter model.ApplicationFilter, page, size int) ([]*model.Application, int, error) {
	all := s.store.ListApplications()
	matched := make([]*model.Application, 0, len(all))
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
		return []*model.Application{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateApplication(id string, input model.Application) (*model.Application, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a, err := s.store.GetApplication(id)
	if err != nil {
		return nil, err
	}
	a.Name = input.Name
	a.Code = input.Code
	a.Description = input.Description
	a.Status = input.Status
	a.UpdatedAt = time.Now()
	if err := s.store.UpdateApplication(a); err != nil {
		return nil, err
	}
	s.logAudit("update", "application", a.ID, "")
	return a, nil
}

func (s *Service) DeleteApplication(id string) error {
	if err := s.store.DeleteApplication(id); err != nil {
		return err
	}
	s.logAudit("delete", "application", id, "")
	return nil
}

func (s *Service) GetApplicationEnvironments(appID string) ([]*model.Environment, error) {
	if _, err := s.store.GetApplication(appID); err != nil {
		return nil, err
	}
	return s.store.ListEnvironmentsByAppID(appID), nil
}

func (s *Service) GetApplicationConfigItems(appID string) ([]*model.ConfigItem, error) {
	if _, err := s.store.GetApplication(appID); err != nil {
		return nil, err
	}
	all := s.store.ListConfigItemsByAppID(appID)
	res := make([]*model.ConfigItem, 0, len(all))
	for _, c := range all {
		res = append(res, s.maskedConfigItem(c))
	}
	return res, nil
}

func (s *Service) GetApplicationReleases(appID string) ([]*model.Release, error) {
	if _, err := s.store.GetApplication(appID); err != nil {
		return nil, err
	}
	return s.store.ListReleasesByAppID(appID), nil
}
