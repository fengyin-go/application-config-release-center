package model

import (
	"strings"
	"time"
)

const (
	ReleaseStatusDraft      = "draft"
	ReleaseStatusReview     = "review"
	ReleaseStatusReleased   = "released"
	ReleaseStatusRolledBack = "rolled_back"
)

var releaseTransitions = map[string]map[string]bool{
	ReleaseStatusDraft:    {ReleaseStatusReview: true},
	ReleaseStatusReview:   {ReleaseStatusReleased: true},
	ReleaseStatusReleased: {ReleaseStatusRolledBack: true},
}

func CanReleaseTransition(from, to string) bool {
	if m, ok := releaseTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Release struct {
	ID         string    `json:"id"`
	AppID      string    `json:"app_id"`
	EnvID      string    `json:"env_id"`
	Version    string    `json:"version"`
	Remark     string    `json:"remark"`
	Status     string    `json:"status"`
	ReleasedBy string    `json:"released_by"`
	ReleasedAt *time.Time `json:"released_at,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Release) Validate() error {
	r.Version = strings.TrimSpace(r.Version)
	if r.AppID == "" {
		return NewValidationError("app_id", "所属应用不能为空")
	}
	if r.EnvID == "" {
		return NewValidationError("env_id", "所属环境不能为空")
	}
	if r.Version == "" {
		return NewValidationError("version", "版本号不能为空")
	}
	if r.Status == "" {
		r.Status = ReleaseStatusDraft
	}
	if r.Status != ReleaseStatusDraft && r.Status != ReleaseStatusReview &&
		r.Status != ReleaseStatusReleased && r.Status != ReleaseStatusRolledBack {
		return NewValidationError("status", "发布状态不合法")
	}
	return nil
}

type ReleaseFilter struct {
	AppID  string
	EnvID  string
	Status string
}

func (f ReleaseFilter) Match(r *Release) bool {
	if f.AppID != "" && r.AppID != f.AppID {
		return false
	}
	if f.EnvID != "" && r.EnvID != f.EnvID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
