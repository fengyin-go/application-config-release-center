package model

import (
	"strings"
	"time"
)

const (
	AppStatusActive   = "active"
	AppStatusArchived = "archived"
)

type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (a *Application) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Code = strings.TrimSpace(a.Code)
	if a.Name == "" {
		return NewValidationError("name", "应用名称不能为空")
	}
	if a.Code == "" {
		return NewValidationError("code", "应用编码不能为空")
	}
	if a.Status == "" {
		a.Status = AppStatusActive
	}
	if a.Status != AppStatusActive && a.Status != AppStatusArchived {
		return NewValidationError("status", "应用状态不合法")
	}
	return nil
}

type ApplicationFilter struct {
	Status  string
	Keyword string
}

func (f ApplicationFilter) Match(a *Application) bool {
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) &&
			!strings.Contains(strings.ToLower(a.Code), k) {
			return false
		}
	}
	return true
}
