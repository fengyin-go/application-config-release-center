package model

import (
	"strings"
	"time"
)

const (
	ValueTypeString = "string"
	ValueTypeInt    = "int"
	ValueTypeBool   = "bool"
	ValueTypeJSON   = "json"

	ConfigStatusEnabled  = "enabled"
	ConfigStatusDisabled = "disabled"
)

type ConfigItem struct {
	ID          string    `json:"id"`
	AppID       string    `json:"app_id"`
	EnvID       string    `json:"env_id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	ValueType   string    `json:"value_type"`
	Description string    `json:"description"`
	Encrypted   bool      `json:"encrypted"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *ConfigItem) Validate() error {
	c.Key = strings.TrimSpace(c.Key)
	c.Value = strings.TrimSpace(c.Value)
	if c.AppID == "" {
		return NewValidationError("app_id", "所属应用不能为空")
	}
	if c.EnvID == "" {
		return NewValidationError("env_id", "所属环境不能为空")
	}
	if c.Key == "" {
		return NewValidationError("key", "配置键不能为空")
	}
	if c.ValueType == "" {
		c.ValueType = ValueTypeString
	}
	if c.ValueType != ValueTypeString && c.ValueType != ValueTypeInt &&
		c.ValueType != ValueTypeBool && c.ValueType != ValueTypeJSON {
		return NewValidationError("value_type", "值类型不合法")
	}
	if c.Status == "" {
		c.Status = ConfigStatusEnabled
	}
	if c.Status != ConfigStatusEnabled && c.Status != ConfigStatusDisabled {
		return NewValidationError("status", "配置项状态不合法")
	}
	return nil
}

type ConfigItemFilter struct {
	AppID   string
	EnvID   string
	Status  string
	Keyword string
}

func (f ConfigItemFilter) Match(c *ConfigItem) bool {
	if f.AppID != "" && c.AppID != f.AppID {
		return false
	}
	if f.EnvID != "" && c.EnvID != f.EnvID {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Key), k) {
			return false
		}
	}
	return true
}
