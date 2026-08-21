package service

import (
	"sort"
	"strings"
	"time"

	"configcenter/internal/model"
	"configcenter/pkg/idgen"
)

func maskValue(v string) string {
	if len(v) <= 4 {
		return strings.Repeat("*", len(v))
	}
	return v[:2] + strings.Repeat("*", len(v)-4) + v[len(v)-2:]
}

func (s *Service) CreateConfigItem(input model.ConfigItem) (*model.ConfigItem, error) {
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
	c := &model.ConfigItem{
		ID:          idgen.Hex(),
		AppID:       input.AppID,
		EnvID:       input.EnvID,
		Key:         input.Key,
		Value:       input.Value,
		ValueType:   input.ValueType,
		Description: input.Description,
		Encrypted:   input.Encrypted,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateConfigItem(c); err != nil {
		return nil, err
	}
	v := &model.ConfigVersion{
		ID:           idgen.Hex(),
		ConfigItemID: c.ID,
		Value:        c.Value,
		ChangedBy:    "system",
		Remark:       "init",
		CreatedAt:    now,
	}
	if err := s.store.CreateConfigVersion(v); err != nil {
		return nil, err
	}
	s.logAudit("create", "config_item", c.ID, "")
	return s.maskedConfigItem(c), nil
}

func (s *Service) maskedConfigItem(c *model.ConfigItem) *model.ConfigItem {
	out := *c
	if out.Encrypted {
		out.Value = maskValue(out.Value)
	}
	return &out
}

func (s *Service) GetConfigItem(id string) (*model.ConfigItem, error) {
	c, err := s.store.GetConfigItem(id)
	if err != nil {
		return nil, err
	}
	return s.maskedConfigItem(c), nil
}

func (s *Service) ListConfigItems(filter model.ConfigItemFilter, page, size int) ([]*model.ConfigItem, int, error) {
	all := s.store.ListConfigItems()
	matched := make([]*model.ConfigItem, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, s.maskedConfigItem(c))
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ConfigItem{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateConfigItem(id string, input model.ConfigItem) (*model.ConfigItem, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetConfigItem(id)
	if err != nil {
		return nil, err
	}
	if c.AppID != input.AppID {
		if _, err := s.store.GetApplication(input.AppID); err != nil {
			return nil, err
		}
	}
	if c.EnvID != input.EnvID {
		if _, err := s.store.GetEnvironment(input.EnvID); err != nil {
			return nil, err
		}
	}
	oldValue := c.Value
	c.AppID = input.AppID
	c.EnvID = input.EnvID
	c.Key = input.Key
	c.Value = input.Value
	c.ValueType = input.ValueType
	c.Description = input.Description
	c.Encrypted = input.Encrypted
	c.Status = input.Status
	c.UpdatedAt = time.Now()
	if err := s.store.UpdateConfigItem(c); err != nil {
		return nil, err
	}
	if oldValue != c.Value {
		v := &model.ConfigVersion{
			ID:           idgen.Hex(),
			ConfigItemID: c.ID,
			Value:        c.Value,
			ChangedBy:    "system",
			Remark:       "update",
			CreatedAt:    time.Now(),
		}
		if err := s.store.CreateConfigVersion(v); err != nil {
			return nil, err
		}
	}
	s.logAudit("update", "config_item", c.ID, "")
	return s.maskedConfigItem(c), nil
}

func (s *Service) DeleteConfigItem(id string) error {
	if err := s.store.DeleteConfigItem(id); err != nil {
		return err
	}
	s.logAudit("delete", "config_item", id, "")
	return nil
}
