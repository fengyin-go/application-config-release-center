package model

import (
	"strings"
	"time"
)

type AuditLog struct {
	ID        string    `json:"id"`
	Operator  string    `json:"operator"`
	Action    string    `json:"action"`
	Resource  string    `json:"resource"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *AuditLog) Validate() error {
	a.Operator = strings.TrimSpace(a.Operator)
	a.Action = strings.TrimSpace(a.Action)
	a.Resource = strings.TrimSpace(a.Resource)
	if a.Operator == "" {
		return NewValidationError("operator", "操作人不能为空")
	}
	if a.Action == "" {
		return NewValidationError("action", "操作类型不能为空")
	}
	if a.Resource == "" {
		return NewValidationError("resource", "操作对象不能为空")
	}
	return nil
}

type AuditLogFilter struct {
	Operator string
	Action   string
	Resource string
}

func (f AuditLogFilter) Match(a *AuditLog) bool {
	if f.Operator != "" && a.Operator != f.Operator {
		return false
	}
	if f.Action != "" && a.Action != f.Action {
		return false
	}
	if f.Resource != "" && a.Resource != f.Resource {
		return false
	}
	return true
}
