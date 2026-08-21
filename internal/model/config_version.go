package model

import (
	"time"
)

type ConfigVersion struct {
	ID           string    `json:"id"`
	ConfigItemID string    `json:"config_item_id"`
	Value        string    `json:"value"`
	Version      int       `json:"version"`
	ChangedBy    string    `json:"changed_by"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"created_at"`
}

func (c *ConfigVersion) Validate() error {
	if c.ConfigItemID == "" {
		return NewValidationError("config_item_id", "配置项 ID 不能为空")
	}
	if c.Version <= 0 {
		return NewValidationError("version", "版本号必须大于 0")
	}
	return nil
}

type ConfigVersionFilter struct {
	ConfigItemID string
}

func (f ConfigVersionFilter) Match(c *ConfigVersion) bool {
	if f.ConfigItemID != "" && c.ConfigItemID != f.ConfigItemID {
		return false
	}
	return true
}
