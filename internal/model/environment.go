package model

import (
	"strings"
	"time"
)

type Environment struct {
	ID          string    `json:"id"`
	AppID       string    `json:"app_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (e *Environment) Validate() error {
	e.Name = strings.TrimSpace(e.Name)
	e.Code = strings.TrimSpace(e.Code)
	if e.AppID == "" {
		return NewValidationError("app_id", "所属应用不能为空")
	}
	if e.Name == "" {
		return NewValidationError("name", "环境名称不能为空")
	}
	if e.Code == "" {
		return NewValidationError("code", "环境编码不能为空")
	}
	return nil
}

type EnvironmentFilter struct {
	AppID   string
	Keyword string
}

func (f EnvironmentFilter) Match(e *Environment) bool {
	if f.AppID != "" && e.AppID != f.AppID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Name), k) &&
			!strings.Contains(strings.ToLower(e.Code), k) {
			return false
		}
	}
	return true
}
