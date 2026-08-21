package service

import (
	"strings"
	"testing"

	"configcenter/internal/config"
	"configcenter/internal/model"
	"configcenter/internal/store"
	"configcenter/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateApplication(t *testing.T) {
	s := newTestService()
	a, err := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	if err != nil {
		t.Fatalf("create app: %v", err)
	}
	if a.ID == "" {
		t.Fatalf("app id empty")
	}
	if a.Status != model.AppStatusActive {
		t.Fatalf("default status active")
	}
	_, err = s.CreateApplication(model.Application{Name: "app2", Code: "c1"})
	if err == nil {
		t.Fatalf("expected conflict")
	}
	logs := s.store.ListAuditLogs()
	if len(logs) != 1 || logs[0].Action != "create" {
		t.Fatalf("audit log mismatch")
	}
	_, err = s.CreateApplication(model.Application{Name: "", Code: "c2"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error")
	}
	_, err = s.CreateApplication(model.Application{Name: "app3", Code: "c3", Status: "invalid"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for invalid status")
	}
}

func TestUpdateDeleteApplication(t *testing.T) {
	s := newTestService()
	a, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	_, err := s.UpdateApplication(a.ID, model.Application{Name: "app-up", Code: "c1"})
	if err != nil {
		t.Fatalf("update app: %v", err)
	}
	if err := s.DeleteApplication(a.ID); err != nil {
		t.Fatalf("delete app: %v", err)
	}
	if _, err := s.GetApplication(a.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
	_, err = s.UpdateApplication("noexist", model.Application{Name: "app", Code: "c1"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestApplicationListAndFilter(t *testing.T) {
	s := newTestService()
	s.CreateApplication(model.Application{Name: "alpha", Code: "a1", Status: model.AppStatusActive})
	s.CreateApplication(model.Application{Name: "beta", Code: "b1", Status: model.AppStatusArchived})
	s.CreateApplication(model.Application{Name: "gamma", Code: "g1", Status: model.AppStatusActive})
	items, total, err := s.ListApplications(model.ApplicationFilter{Status: model.AppStatusActive}, 1, 10)
	if err != nil || total != 2 || len(items) != 2 {
		t.Fatalf("filter status: total=%d items=%d", total, len(items))
	}
	items, total, err = s.ListApplications(model.ApplicationFilter{Keyword: "beta"}, 1, 10)
	if err != nil || total != 1 || items[0].Name != "beta" {
		t.Fatalf("filter keyword: total=%d", total)
	}
	items, total, err = s.ListApplications(model.ApplicationFilter{Keyword: "  beta  "}, 1, 10)
	if err != nil || total != 1 {
		t.Fatalf("trim keyword: total=%d", total)
	}
	items, total, err = s.ListApplications(model.ApplicationFilter{}, 1, 1)
	if err != nil || total != 3 || len(items) != 1 {
		t.Fatalf("page size 1: total=%d len=%d", total, len(items))
	}
	items, total, err = s.ListApplications(model.ApplicationFilter{}, 2, 1)
	if err != nil || total != 3 || len(items) != 1 {
		t.Fatalf("page 2 size 1: total=%d len=%d", total, len(items))
	}
	items, total, err = s.ListApplications(model.ApplicationFilter{}, 5, 1)
	if err != nil || total != 3 || len(items) != 0 {
		t.Fatalf("page 5: total=%d len=%d", total, len(items))
	}
}

func TestCreateEnvironment(t *testing.T) {
	s := newTestService()
	_, err := s.CreateEnvironment(model.Environment{AppID: "noexist", Name: "dev", Code: "dev"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing app, got %v", err)
	}
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	e, err := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	if err != nil {
		t.Fatalf("create env: %v", err)
	}
	_, err = s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev2", Code: "dev"})
	if err == nil {
		t.Fatalf("expected conflict for same app code")
	}
	_, err = s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "prod", Code: "prod"})
	if err != nil {
		t.Fatalf("create env2: %v", err)
	}
	_, err = s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "", Code: "x"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for empty name")
	}
	logs := s.store.ListAuditLogs()
	found := false
	for _, al := range logs {
		if al.Resource == "environment" && al.Action == "create" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("audit log for env create missing")
	}
	s.UpdateEnvironment(e.ID, model.Environment{AppID: app.ID, Name: "dev-up", Code: "dev"})
}

func TestEnvironmentListAndFilter(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "prod", Code: "prod"})
	items, total, err := s.ListEnvironments(model.EnvironmentFilter{AppID: app.ID}, 1, 10)
	if err != nil || total != 2 {
		t.Fatalf("env list: total=%d", total)
	}
	items, total, err = s.ListEnvironments(model.EnvironmentFilter{Keyword: "prod"}, 1, 10)
	if err != nil || total != 1 || items[0].Name != "prod" {
		t.Fatalf("env keyword filter: total=%d", total)
	}
	items, total, err = s.ListEnvironments(model.EnvironmentFilter{AppID: "noexist"}, 1, 10)
	if err != nil || total != 0 {
		t.Fatalf("env filter noexist: total=%d", total)
	}
}

func TestCreateConfigItem(t *testing.T) {
	s := newTestService()
	_, err := s.CreateConfigItem(model.ConfigItem{AppID: "noexist", EnvID: "e1", Key: "k1", Value: "v1"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing app, got %v", err)
	}
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, err := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1", Encrypted: true})
	if err != nil {
		t.Fatalf("create config: %v", err)
	}
	if !strings.Contains(c.Value, "*") {
		t.Fatalf("expected masked value, got %s", c.Value)
	}
	_, err = s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v2"})
	if err == nil {
		t.Fatalf("expected conflict for same app+env+key")
	}
	versions, total, err := s.ListConfigVersions(model.ConfigVersionFilter{ConfigItemID: c.ID}, 1, 10)
	if err != nil || total != 1 || versions[0].Version != 1 {
		t.Fatalf("expected initial version, total=%d version=%d", total, versions[0].Version)
	}
	_, err = s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "", Value: "v2"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for empty key")
	}
	_, err = s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: "noenv", Key: "k2", Value: "v2"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing env, got %v", err)
	}
}

func TestUpdateConfigItemCreatesVersion(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	_, err := s.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v2"})
	if err != nil {
		t.Fatalf("update config: %v", err)
	}
	versions, total, err := s.ListConfigVersions(model.ConfigVersionFilter{ConfigItemID: c.ID}, 1, 10)
	if err != nil || total != 2 {
		t.Fatalf("expected 2 versions, got %d", total)
	}
	if versions[0].Version != 2 {
		t.Fatalf("expected latest version 2, got %d", versions[0].Version)
	}
	_, err = s.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v2"})
	if err != nil {
		t.Fatalf("update same value: %v", err)
	}
	versions, total, _ = s.ListConfigVersions(model.ConfigVersionFilter{ConfigItemID: c.ID}, 1, 10)
	if total != 2 {
		t.Fatalf("expected still 2 versions when value unchanged, got %d", total)
	}
}

func TestConfigItemMasking(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "secret", Value: "mysecret", Encrypted: true})
	got, _ := s.GetConfigItem(c.ID)
	if got.Value == "mysecret" {
		t.Fatalf("expected masked value")
	}
	c2, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "plain", Value: "open", Encrypted: false})
	got2, _ := s.GetConfigItem(c2.ID)
	if got2.Value != "open" {
		t.Fatalf("expected plain value")
	}
}

func TestConfigItemListAndFilter(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1", Status: model.ConfigStatusEnabled})
	s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k2", Value: "v2", Status: model.ConfigStatusDisabled})
	s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k3", Value: "v3", Status: model.ConfigStatusEnabled})
	items, total, _ := s.ListConfigItems(model.ConfigItemFilter{Status: model.ConfigStatusEnabled}, 1, 10)
	if total != 2 {
		t.Fatalf("filter status mismatch")
	}
	items, total, _ = s.ListConfigItems(model.ConfigItemFilter{Keyword: "k2"}, 1, 10)
	if total != 1 || items[0].Key != "k2" {
		t.Fatalf("filter keyword mismatch")
	}
	items, total, _ = s.ListConfigItems(model.ConfigItemFilter{AppID: app.ID, EnvID: env.ID}, 1, 10)
	if total != 3 {
		t.Fatalf("filter app+env mismatch")
	}
	items, total, _ = s.ListConfigItems(model.ConfigItemFilter{AppID: "noexist"}, 1, 10)
	if total != 0 {
		t.Fatalf("filter noexist mismatch")
	}
}

func TestReleaseLifecycle(t *testing.T) {
	tests := []struct {
		name        string
		transitions []string
		wantErrIdx  int
	}{
		{
			name:        "draft->review->released->rolled_back",
			transitions: []string{model.ReleaseStatusReview, model.ReleaseStatusReleased, model.ReleaseStatusRolledBack},
			wantErrIdx:  -1,
		},
		{
			name:        "draft->released illegal",
			transitions: []string{model.ReleaseStatusReleased},
			wantErrIdx:  0,
		},
		{
			name:        "released->draft illegal",
			transitions: []string{model.ReleaseStatusReview, model.ReleaseStatusReleased, model.ReleaseStatusRolledBack, model.ReleaseStatusDraft},
			wantErrIdx:  3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService()
			app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
			env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
			r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
			for i, to := range tt.transitions {
				_, err := s.UpdateReleaseStatus(r.ID, to, "u1")
				if tt.wantErrIdx == i {
					if !model.IsValidationError(err) {
						t.Fatalf("expected validation error at step %d, got %v", i, err)
					}
					return
				}
				if err != nil {
					t.Fatalf("step %d: %v", i, err)
				}
			}
			if tt.wantErrIdx != -1 {
				t.Fatalf("expected error at step %d", tt.wantErrIdx)
			}
		})
	}
}

func TestReleaseListAndFilter(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
	s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v2", Status: model.ReleaseStatusReleased})
	s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v3", Status: model.ReleaseStatusReview})
	items, total, _ := s.ListReleases(model.ReleaseFilter{Status: model.ReleaseStatusReleased}, 1, 10)
	if total != 1 || items[0].Version != "v2" {
		t.Fatalf("filter status mismatch")
	}
	items, total, _ = s.ListReleases(model.ReleaseFilter{AppID: app.ID}, 1, 10)
	if total != 3 {
		t.Fatalf("filter app mismatch")
	}
	items, total, _ = s.ListReleases(model.ReleaseFilter{EnvID: env.ID}, 1, 10)
	if total != 3 {
		t.Fatalf("filter env mismatch")
	}
	items, total, _ = s.ListReleases(model.ReleaseFilter{AppID: "noexist"}, 1, 10)
	if total != 0 {
		t.Fatalf("filter noexist mismatch")
	}
}

func TestAuditLogQuery(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	s.DeleteApplication(app.ID)
	logs, total, _ := s.ListAuditLogs(model.AuditLogFilter{Action: "delete"}, 1, 10)
	if total < 1 {
		t.Fatalf("expected audit logs for delete")
	}
	found := false
	for _, l := range logs {
		if l.Action == "delete" && l.Resource == "application" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected delete application audit log")
	}
	logs, total, _ = s.ListAuditLogs(model.AuditLogFilter{Operator: "system"}, 1, 10)
	if total < 1 {
		t.Fatalf("expected audit logs for system operator")
	}
	logs, total, _ = s.ListAuditLogs(model.AuditLogFilter{Resource: "application"}, 1, 10)
	if total < 1 {
		t.Fatalf("expected audit logs for application resource")
	}
	logs, total, _ = s.ListAuditLogs(model.AuditLogFilter{Action: "noexist"}, 1, 10)
	if total != 0 {
		t.Fatalf("expected no audit logs for noexist action")
	}
}

func TestStats(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k2", Value: "v2"})
	s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
	s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v2", Status: model.ReleaseStatusReleased})

	counts := s.StatsConfigItemCountByAppEnv()
	if len(counts) != 1 || counts[0].Count != 2 {
		t.Fatalf("config item count stat mismatch")
	}

	rels := s.StatsReleaseByStatus()
	m := make(map[string]int)
	for _, r := range rels {
		m[r.Status] = r.Count
	}
	if m[model.ReleaseStatusDraft] != 1 || m[model.ReleaseStatusReleased] != 1 {
		t.Fatalf("release status stat mismatch")
	}

	acts := s.StatsAuditByAction()
	if len(acts) == 0 {
		t.Fatalf("expected audit action stats")
	}
}

func TestStatsEmpty(t *testing.T) {
	s := newTestService()
	if len(s.StatsConfigItemCountByAppEnv()) != 0 {
		t.Fatalf("expected empty config item stats")
	}
	if len(s.StatsReleaseByStatus()) != 0 {
		t.Fatalf("expected empty release stats")
	}
	if len(s.StatsAuditByAction()) != 0 {
		t.Fatalf("expected empty audit stats")
	}
}

func TestPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 5; i++ {
		s.CreateApplication(model.Application{Name: "app", Code: "c" + string(rune('0'+i))})
	}
	items, total, _ := s.ListApplications(model.ApplicationFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("page 1 size 2: total=%d len=%d", total, len(items))
	}
	items, total, _ = s.ListApplications(model.ApplicationFilter{}, 3, 2)
	if total != 5 || len(items) != 1 {
		t.Fatalf("page 3 size 2: total=%d len=%d", total, len(items))
	}
	items, total, _ = s.ListApplications(model.ApplicationFilter{}, 10, 2)
	if total != 5 || len(items) != 0 {
		t.Fatalf("page 10 size 2: total=%d len=%d", total, len(items))
	}
}

func TestReleaseInvalidAppEnv(t *testing.T) {
	s := newTestService()
	_, err := s.CreateRelease(model.Release{AppID: "noapp", EnvID: "noenv", Version: "v1"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing app, got %v", err)
	}
}

func TestConfigItemInvalidEnv(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	_, err := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: "noenv", Key: "k1", Value: "v1"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing env, got %v", err)
	}
}

func TestConfigVersionOrdering(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	for i := 2; i <= 5; i++ {
		s.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v" + string(rune('0'+i))})
	}
	versions, total, _ := s.ListConfigVersions(model.ConfigVersionFilter{ConfigItemID: c.ID}, 1, 10)
	if total != 5 {
		t.Fatalf("expected 5 versions, got %d", total)
	}
	for i := 0; i < 4; i++ {
		if versions[i].Version <= versions[i+1].Version {
			t.Fatalf("versions not descending")
		}
	}
}

func TestReleaseAuditLog(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
	before := len(s.store.ListAuditLogs())
	s.UpdateReleaseStatus(r.ID, model.ReleaseStatusReview, "u1")
	after := len(s.store.ListAuditLogs())
	if after != before+1 {
		t.Fatalf("expected audit log for transition")
	}
}

func TestDeleteRelease(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1", Status: model.ReleaseStatusDraft})
	if err := s.DeleteRelease(r.ID); err != nil {
		t.Fatalf("delete release: %v", err)
	}
	if _, err := s.GetRelease(r.ID); err == nil {
		t.Fatalf("expected not found after delete")
	}
}

func TestGetReleaseNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.GetRelease("noexist"); err != store.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestGetApplicationNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.GetApplication("noexist"); err != store.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestGetEnvironmentNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.GetEnvironment("noexist"); err != store.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestGetConfigItemNotFound(t *testing.T) {
	s := newTestService()
	if _, err := s.GetConfigItem("noexist"); err != store.ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUpdateEnvironmentMissingApp(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	e, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	_, err := s.UpdateEnvironment(e.ID, model.Environment{AppID: "noexist", Name: "dev", Code: "dev"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing app in update, got %v", err)
	}
}

func TestUpdateConfigItemMissingEnv(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	_, err := s.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: "noenv", Key: "k1", Value: "v2"})
	if err != store.ErrNotFound {
		t.Fatalf("expected not found for missing env in update, got %v", err)
	}
}

func TestMultipleAppEnvironments(t *testing.T) {
	s := newTestService()
	app1, _ := s.CreateApplication(model.Application{Name: "app1", Code: "c1"})
	app2, _ := s.CreateApplication(model.Application{Name: "app2", Code: "c2"})
	s.CreateEnvironment(model.Environment{AppID: app1.ID, Name: "dev", Code: "dev"})
	s.CreateEnvironment(model.Environment{AppID: app1.ID, Name: "prod", Code: "prod"})
	s.CreateEnvironment(model.Environment{AppID: app2.ID, Name: "dev", Code: "dev"})
	_, total, _ := s.ListEnvironments(model.EnvironmentFilter{AppID: app1.ID}, 1, 10)
	if total != 2 {
		t.Fatalf("expected 2 envs for app1, got %d", total)
	}
	_, total, _ = s.ListEnvironments(model.EnvironmentFilter{AppID: app2.ID}, 1, 10)
	if total != 1 {
		t.Fatalf("expected 1 env for app2, got %d", total)
	}
}

func TestEncryptedShortValue(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k", Value: "ab", Encrypted: true})
	got, _ := s.GetConfigItem(c.ID)
	if got.Value != "**" {
		t.Fatalf("expected ** for short value, got %s", got.Value)
	}
}

func TestApplicationCodeUniqueness(t *testing.T) {
	s := newTestService()
	_, err := s.CreateApplication(model.Application{Name: "app1", Code: "same"})
	if err != nil {
		t.Fatalf("create app1: %v", err)
	}
	_, err = s.CreateApplication(model.Application{Name: "app2", Code: "same"})
	if err == nil {
		t.Fatalf("expected conflict for duplicate code")
	}
}

func TestEnvironmentCodeUniquenessPerApp(t *testing.T) {
	s := newTestService()
	app1, _ := s.CreateApplication(model.Application{Name: "app1", Code: "c1"})
	app2, _ := s.CreateApplication(model.Application{Name: "app2", Code: "c2"})
	_, err := s.CreateEnvironment(model.Environment{AppID: app1.ID, Name: "dev", Code: "dev"})
	if err != nil {
		t.Fatalf("create env for app1: %v", err)
	}
	_, err = s.CreateEnvironment(model.Environment{AppID: app1.ID, Name: "dev2", Code: "dev"})
	if err == nil {
		t.Fatalf("expected conflict for same app same code")
	}
	_, err = s.CreateEnvironment(model.Environment{AppID: app2.ID, Name: "dev", Code: "dev"})
	if err != nil {
		t.Fatalf("same code different app should be ok: %v", err)
	}
}

func TestConfigItemKeyUniquenessPerAppEnv(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env1, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	env2, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "prod", Code: "prod"})
	_, err := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env1.ID, Key: "k1", Value: "v1"})
	if err != nil {
		t.Fatalf("create item in env1: %v", err)
	}
	_, err = s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env1.ID, Key: "k1", Value: "v2"})
	if err == nil {
		t.Fatalf("expected conflict for same app+env+key")
	}
	_, err = s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env2.ID, Key: "k1", Value: "v2"})
	if err != nil {
		t.Fatalf("same key different env should be ok: %v", err)
	}
}

func TestConfigItemValueTypeValidation(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	_, err := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1", ValueType: "invalid"})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for invalid value type, got %v", err)
	}
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1", ValueType: model.ValueTypeInt})
	if c.ValueType != model.ValueTypeInt {
		t.Fatalf("expected int value type")
	}
}

func TestReleaseVersionValidation(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	_, err := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: ""})
	if !model.IsValidationError(err) {
		t.Fatalf("expected validation error for empty version")
	}
}

func TestReleaseStatusDefault(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1"})
	if r.Status != model.ReleaseStatusDraft {
		t.Fatalf("expected default status draft, got %s", r.Status)
	}
}

func TestAuditLogForAllOperations(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	c, _ := s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v1"})
	r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1"})

	s.UpdateApplication(app.ID, model.Application{Name: "app-up", Code: "c1"})
	s.UpdateEnvironment(env.ID, model.Environment{AppID: app.ID, Name: "dev-up", Code: "dev"})
	s.UpdateConfigItem(c.ID, model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k1", Value: "v2"})
	s.UpdateReleaseStatus(r.ID, model.ReleaseStatusReview, "u1")

	logs := s.store.ListAuditLogs()
	actions := make(map[string]int)
	for _, l := range logs {
		actions[l.Action]++
	}
	if actions["create"] < 4 {
		t.Fatalf("expected at least 4 create audit logs, got %d", actions["create"])
	}
	if actions["update"] < 3 {
		t.Fatalf("expected at least 3 update audit logs, got %d", actions["update"])
	}
	if actions["transition"] < 1 {
		t.Fatalf("expected at least 1 transition audit log, got %d", actions["transition"])
	}
}

func TestReleaseTransitionToReleasedSetsReleasedAt(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	r, _ := s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v1"})
	s.UpdateReleaseStatus(r.ID, model.ReleaseStatusReview, "u1")
	rel, _ := s.UpdateReleaseStatus(r.ID, model.ReleaseStatusReleased, "u1")
	if rel.ReleasedAt == nil {
		t.Fatalf("expected released_at to be set")
	}
	if rel.ReleasedBy != "u1" {
		t.Fatalf("expected released_by u1, got %s", rel.ReleasedBy)
	}
}

func TestConfigItemPagination(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	for i := 0; i < 5; i++ {
		s.CreateConfigItem(model.ConfigItem{AppID: app.ID, EnvID: env.ID, Key: "k" + string(rune('0'+i)), Value: "v" + string(rune('0'+i))})
	}
	items, total, _ := s.ListConfigItems(model.ConfigItemFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("page 1 size 2: total=%d len=%d", total, len(items))
	}
}

func TestReleasePagination(t *testing.T) {
	s := newTestService()
	app, _ := s.CreateApplication(model.Application{Name: "app", Code: "c1"})
	env, _ := s.CreateEnvironment(model.Environment{AppID: app.ID, Name: "dev", Code: "dev"})
	for i := 0; i < 5; i++ {
		s.CreateRelease(model.Release{AppID: app.ID, EnvID: env.ID, Version: "v" + string(rune('0'+i))})
	}
	items, total, _ := s.ListReleases(model.ReleaseFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("page 1 size 2: total=%d len=%d", total, len(items))
	}
}

func TestAuditLogPagination(t *testing.T) {
	s := newTestService()
	for i := 0; i < 5; i++ {
		s.CreateApplication(model.Application{Name: "app", Code: "c" + string(rune('0'+i))})
	}
	items, total, _ := s.ListAuditLogs(model.AuditLogFilter{}, 1, 2)
	if total != 5 || len(items) != 2 {
		t.Fatalf("page 1 size 2: total=%d len=%d", total, len(items))
	}
}
