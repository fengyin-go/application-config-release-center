package store

import (
	"sync"
	"testing"
	"time"

	"configcenter/internal/model"
)

func TestApplicationStore(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Application{ID: "a1", Name: "app1", Code: "code1", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateApplication(a); err != nil {
		t.Fatalf("create app: %v", err)
	}
	if err := s.CreateApplication(&model.Application{ID: "a2", Name: "app2", Code: "code1", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetApplication("a1")
	if err != nil || got.ID != "a1" {
		t.Fatalf("get app: %v", err)
	}
	if _, err := s.GetApplication("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if _, err := s.GetApplicationByCode("code1"); err != nil || got.ID != "a1" {
		t.Fatalf("get by code: %v", err)
	}
	if _, err := s.GetApplicationByCode("noexist"); err != ErrNotFound {
		t.Fatalf("expected not found by code, got %v", err)
	}
	if len(s.ListApplications()) != 1 {
		t.Fatalf("list count mismatch")
	}
	a.Name = "app1-updated"
	if err := s.UpdateApplication(a); err != nil {
		t.Fatalf("update app: %v", err)
	}
	s.CreateApplication(&model.Application{ID: "a2", Name: "app2", Code: "code2", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err := s.UpdateApplication(&model.Application{ID: "a2", Name: "x", Code: "code1", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict on update, got %v", err)
	}
	if err := s.DeleteApplication("a1"); err != nil {
		t.Fatalf("delete app: %v", err)
	}
	if err := s.DeleteApplication("a1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
	if err := s.UpdateApplication(&model.Application{ID: "noexist", Name: "x", Code: "c", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
}

func TestApplicationStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			a := &model.Application{ID: "ac" + string(rune('0'+idx)), Name: "app", Code: "c" + string(rune('0'+idx)), Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			_ = s.CreateApplication(a)
		}(i)
	}
	wg.Wait()
	if len(s.ListApplications()) != 100 {
		t.Fatalf("expected 100 apps, got %d", len(s.ListApplications()))
	}
}

func TestApplicationStoreUpdatePreservesOtherFields(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Application{ID: "a1", Name: "app1", Code: "code1", Description: "desc", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateApplication(a)
	a.Name = "updated"
	s.UpdateApplication(a)
	got, _ := s.GetApplication("a1")
	if got.Description != "desc" {
		t.Fatalf("description should be preserved")
	}
}

func TestEnvironmentStore(t *testing.T) {
	s := NewMemoryStore()
	e := &model.Environment{ID: "e1", AppID: "a1", Name: "dev", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateEnvironment(e); err != nil {
		t.Fatalf("create env: %v", err)
	}
	if err := s.CreateEnvironment(&model.Environment{ID: "e2", AppID: "a1", Name: "dev2", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if err := s.CreateEnvironment(&model.Environment{ID: "e3", AppID: "a2", Name: "dev", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		t.Fatalf("same code different app should ok: %v", err)
	}
	got, err := s.GetEnvironment("e1")
	if err != nil || got.ID != "e1" {
		t.Fatalf("get env: %v", err)
	}
	if _, err := s.GetEnvironment("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListEnvironments()) != 2 {
		t.Fatalf("list count mismatch")
	}
	e.Name = "prod"
	if err := s.UpdateEnvironment(e); err != nil {
		t.Fatalf("update env: %v", err)
	}
	s.CreateEnvironment(&model.Environment{ID: "e4", AppID: "a1", Name: "dev2", Code: "dev2", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err := s.UpdateEnvironment(&model.Environment{ID: "e4", AppID: "a1", Name: "x", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict on update, got %v", err)
	}
	if err := s.DeleteEnvironment("e1"); err != nil {
		t.Fatalf("delete env: %v", err)
	}
	if err := s.DeleteEnvironment("e1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
	if err := s.UpdateEnvironment(&model.Environment{ID: "noexist", AppID: "a1", Name: "x", Code: "c", CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
}

func TestEnvironmentStoreByAppID(t *testing.T) {
	s := NewMemoryStore()
	s.CreateEnvironment(&model.Environment{ID: "e1", AppID: "a1", Name: "dev", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateEnvironment(&model.Environment{ID: "e2", AppID: "a1", Name: "prod", Code: "prod", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateEnvironment(&model.Environment{ID: "e3", AppID: "a2", Name: "dev", Code: "dev2", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListEnvironmentsByAppID("a1")
	if len(list) != 2 {
		t.Fatalf("expected 2 envs for a1, got %d", len(list))
	}
	list = s.ListEnvironmentsByAppID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 envs for noexist, got %d", len(list))
	}
}

func TestEnvironmentStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			e := &model.Environment{ID: "ec" + string(rune('0'+idx)), AppID: "a1", Name: "dev", Code: "c" + string(rune('0'+idx)), CreatedAt: time.Now(), UpdatedAt: time.Now()}
			_ = s.CreateEnvironment(e)
		}(i)
	}
	wg.Wait()
	if len(s.ListEnvironments()) != 100 {
		t.Fatalf("expected 100 envs, got %d", len(s.ListEnvironments()))
	}
}

func TestConfigItemStore(t *testing.T) {
	s := NewMemoryStore()
	c := &model.ConfigItem{ID: "c1", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v1", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateConfigItem(c); err != nil {
		t.Fatalf("create item: %v", err)
	}
	if err := s.CreateConfigItem(&model.ConfigItem{ID: "c2", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v2", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	if err := s.CreateConfigItem(&model.ConfigItem{ID: "c3", AppID: "a1", EnvID: "e2", Key: "k1", Value: "v2", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != nil {
		t.Fatalf("same key different env should ok: %v", err)
	}
	got, err := s.GetConfigItem("c1")
	if err != nil || got.ID != "c1" {
		t.Fatalf("get item: %v", err)
	}
	if _, err := s.GetConfigItem("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListConfigItems()) != 2 {
		t.Fatalf("list count mismatch")
	}
	c.Value = "v1-updated"
	if err := s.UpdateConfigItem(c); err != nil {
		t.Fatalf("update item: %v", err)
	}
	s.CreateConfigItem(&model.ConfigItem{ID: "c4", AppID: "a1", EnvID: "e1", Key: "k2", Value: "v2", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	if err := s.UpdateConfigItem(&model.ConfigItem{ID: "c4", AppID: "a1", EnvID: "e1", Key: "k1", Value: "x", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrConflict {
		t.Fatalf("expected conflict on update, got %v", err)
	}
	if err := s.DeleteConfigItem("c1"); err != nil {
		t.Fatalf("delete item: %v", err)
	}
	if err := s.DeleteConfigItem("c1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
	if err := s.UpdateConfigItem(&model.ConfigItem{ID: "noexist", AppID: "a1", EnvID: "e1", Key: "k", Value: "x", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
}

func TestConfigItemStoreByAppEnvKey(t *testing.T) {
	s := NewMemoryStore()
	s.CreateConfigItem(&model.ConfigItem{ID: "c1", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v1", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	got, err := s.GetConfigItemByAppEnvKey("a1", "e1", "k1")
	if err != nil || got.ID != "c1" {
		t.Fatalf("get by app env key: %v", err)
	}
	if _, err := s.GetConfigItemByAppEnvKey("a1", "e1", "noexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestConfigItemStoreByAppID(t *testing.T) {
	s := NewMemoryStore()
	s.CreateConfigItem(&model.ConfigItem{ID: "c1", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v1", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateConfigItem(&model.ConfigItem{ID: "c2", AppID: "a1", EnvID: "e2", Key: "k2", Value: "v2", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateConfigItem(&model.ConfigItem{ID: "c3", AppID: "a2", EnvID: "e1", Key: "k3", Value: "v3", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListConfigItemsByAppID("a1")
	if len(list) != 2 {
		t.Fatalf("expected 2 items for a1, got %d", len(list))
	}
	list = s.ListConfigItemsByAppID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 items for noexist, got %d", len(list))
	}
}

func TestConfigItemStoreByEnvID(t *testing.T) {
	s := NewMemoryStore()
	s.CreateConfigItem(&model.ConfigItem{ID: "c1", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v1", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateConfigItem(&model.ConfigItem{ID: "c2", AppID: "a2", EnvID: "e1", Key: "k2", Value: "v2", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListConfigItemsByEnvID("e1")
	if len(list) != 2 {
		t.Fatalf("expected 2 items for e1, got %d", len(list))
	}
	list = s.ListConfigItemsByEnvID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 items for noexist, got %d", len(list))
	}
}

func TestConfigItemStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			c := &model.ConfigItem{ID: "cc" + string(rune('0'+idx)), AppID: "a1", EnvID: "e1", Key: "k" + string(rune('0'+idx)), Value: "v", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			_ = s.CreateConfigItem(c)
		}(i)
	}
	wg.Wait()
	if len(s.ListConfigItems()) != 100 {
		t.Fatalf("expected 100 items, got %d", len(s.ListConfigItems()))
	}
}

func TestConfigVersionStore(t *testing.T) {
	s := NewMemoryStore()
	v1 := &model.ConfigVersion{ID: "v1", ConfigItemID: "c1", Value: "a", ChangedBy: "u1", Remark: "r1", CreatedAt: time.Now()}
	if err := s.CreateConfigVersion(v1); err != nil {
		t.Fatalf("create version: %v", err)
	}
	if v1.Version != 1 {
		t.Fatalf("expected version 1, got %d", v1.Version)
	}
	v2 := &model.ConfigVersion{ID: "v2", ConfigItemID: "c1", Value: "b", ChangedBy: "u1", Remark: "r2", CreatedAt: time.Now()}
	if err := s.CreateConfigVersion(v2); err != nil {
		t.Fatalf("create version2: %v", err)
	}
	if v2.Version != 2 {
		t.Fatalf("expected version 2, got %d", v2.Version)
	}
	v3 := &model.ConfigVersion{ID: "v3", ConfigItemID: "c2", Value: "c", ChangedBy: "u1", Remark: "r3", CreatedAt: time.Now()}
	if err := s.CreateConfigVersion(v3); err != nil {
		t.Fatalf("create version3: %v", err)
	}
	if v3.Version != 1 {
		t.Fatalf("expected version 1 for different item, got %d", v3.Version)
	}
	got, err := s.GetConfigVersion("v1")
	if err != nil || got.ID != "v1" {
		t.Fatalf("get version: %v", err)
	}
	if _, err := s.GetConfigVersion("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListConfigVersions()) != 3 {
		t.Fatalf("list count mismatch")
	}
	list := s.ListConfigVersionsByConfigItemID("c1")
	if len(list) != 2 {
		t.Fatalf("expected 2 versions for c1, got %d", len(list))
	}
	list = s.ListConfigVersionsByConfigItemID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 versions for noexist, got %d", len(list))
	}
}

func TestConfigVersionConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			v := &model.ConfigVersion{ID: "vc" + string(rune('0'+idx)), ConfigItemID: "c1", Value: "v", ChangedBy: "u1", Remark: "r", CreatedAt: time.Now()}
			_ = s.CreateConfigVersion(v)
		}(i)
	}
	wg.Wait()
	if len(s.ListConfigVersions()) != 100 {
		t.Fatalf("expected 100 versions, got %d", len(s.ListConfigVersions()))
	}
}

func TestReleaseStore(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Release{ID: "r1", AppID: "a1", EnvID: "e1", Version: "v1", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateRelease(r); err != nil {
		t.Fatalf("create release: %v", err)
	}
	got, err := s.GetRelease("r1")
	if err != nil || got.ID != "r1" {
		t.Fatalf("get release: %v", err)
	}
	if _, err := s.GetRelease("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListReleases()) != 1 {
		t.Fatalf("list count mismatch")
	}
	r.Status = model.ReleaseStatusReleased
	if err := s.UpdateRelease(r); err != nil {
		t.Fatalf("update release: %v", err)
	}
	if err := s.UpdateRelease(&model.Release{ID: "noexist", AppID: "a1", EnvID: "e1", Version: "v1", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}); err != ErrNotFound {
		t.Fatalf("expected not found on update, got %v", err)
	}
	if err := s.DeleteRelease("r1"); err != nil {
		t.Fatalf("delete release: %v", err)
	}
	if err := s.DeleteRelease("r1"); err != ErrNotFound {
		t.Fatalf("expected not found on delete, got %v", err)
	}
}

func TestReleaseStoreByAppID(t *testing.T) {
	s := NewMemoryStore()
	s.CreateRelease(&model.Release{ID: "r1", AppID: "a1", EnvID: "e1", Version: "v1", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateRelease(&model.Release{ID: "r2", AppID: "a1", EnvID: "e2", Version: "v2", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateRelease(&model.Release{ID: "r3", AppID: "a2", EnvID: "e1", Version: "v3", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListReleasesByAppID("a1")
	if len(list) != 2 {
		t.Fatalf("expected 2 releases for a1, got %d", len(list))
	}
	list = s.ListReleasesByAppID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 releases for noexist, got %d", len(list))
	}
}

func TestReleaseStoreByEnvID(t *testing.T) {
	s := NewMemoryStore()
	s.CreateRelease(&model.Release{ID: "r1", AppID: "a1", EnvID: "e1", Version: "v1", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	s.CreateRelease(&model.Release{ID: "r2", AppID: "a2", EnvID: "e1", Version: "v2", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	list := s.ListReleasesByEnvID("e1")
	if len(list) != 2 {
		t.Fatalf("expected 2 releases for e1, got %d", len(list))
	}
	list = s.ListReleasesByEnvID("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 releases for noexist, got %d", len(list))
	}
}

func TestReleaseStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			r := &model.Release{ID: "rc" + string(rune('0'+idx)), AppID: "a1", EnvID: "e1", Version: "v" + string(rune('0'+idx)), Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
			_ = s.CreateRelease(r)
		}(i)
	}
	wg.Wait()
	if len(s.ListReleases()) != 100 {
		t.Fatalf("expected 100 releases, got %d", len(s.ListReleases()))
	}
}

func TestAuditLogStore(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "al1", Operator: "u1", Action: "create", Resource: "app", Detail: "d1", CreatedAt: time.Now()}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatalf("create audit: %v", err)
	}
	got, err := s.GetAuditLog("al1")
	if err != nil || got.ID != "al1" {
		t.Fatalf("get audit: %v", err)
	}
	if _, err := s.GetAuditLog("notexist"); err != ErrNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
	if len(s.ListAuditLogs()) != 1 {
		t.Fatalf("list count mismatch")
	}
}

func TestAuditLogStoreByAction(t *testing.T) {
	s := NewMemoryStore()
	s.CreateAuditLog(&model.AuditLog{ID: "al1", Operator: "u1", Action: "create", Resource: "app", Detail: "d1", CreatedAt: time.Now()})
	s.CreateAuditLog(&model.AuditLog{ID: "al2", Operator: "u1", Action: "delete", Resource: "app", Detail: "d2", CreatedAt: time.Now()})
	s.CreateAuditLog(&model.AuditLog{ID: "al3", Operator: "u1", Action: "create", Resource: "env", Detail: "d3", CreatedAt: time.Now()})
	list := s.ListAuditLogsByAction("create")
	if len(list) != 2 {
		t.Fatalf("expected 2 create audits, got %d", len(list))
	}
	list = s.ListAuditLogsByAction("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 audits for noexist action, got %d", len(list))
	}
}

func TestAuditLogStoreByResource(t *testing.T) {
	s := NewMemoryStore()
	s.CreateAuditLog(&model.AuditLog{ID: "al1", Operator: "u1", Action: "create", Resource: "app", Detail: "d1", CreatedAt: time.Now()})
	s.CreateAuditLog(&model.AuditLog{ID: "al2", Operator: "u1", Action: "create", Resource: "env", Detail: "d2", CreatedAt: time.Now()})
	list := s.ListAuditLogsByResource("app")
	if len(list) != 1 {
		t.Fatalf("expected 1 app audit, got %d", len(list))
	}
	list = s.ListAuditLogsByResource("noexist")
	if len(list) != 0 {
		t.Fatalf("expected 0 audits for noexist resource, got %d", len(list))
	}
}

func TestAuditLogStoreConcurrent(t *testing.T) {
	s := NewMemoryStore()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			a := &model.AuditLog{ID: "alc" + string(rune('0'+idx)), Operator: "u1", Action: "create", Resource: "app", Detail: "d", CreatedAt: time.Now()}
			_ = s.CreateAuditLog(a)
		}(i)
	}
	wg.Wait()
	if len(s.ListAuditLogs()) != 100 {
		t.Fatalf("expected 100 audit logs, got %d", len(s.ListAuditLogs()))
	}
}

func TestEmptyLists(t *testing.T) {
	s := NewMemoryStore()
	if len(s.ListApplications()) != 0 {
		t.Fatalf("expected empty applications")
	}
	if len(s.ListEnvironments()) != 0 {
		t.Fatalf("expected empty environments")
	}
	if len(s.ListEnvironmentsByAppID("a1")) != 0 {
		t.Fatalf("expected empty environments by app")
	}
	if len(s.ListConfigItems()) != 0 {
		t.Fatalf("expected empty config items")
	}
	if len(s.ListConfigItemsByAppID("a1")) != 0 {
		t.Fatalf("expected empty config items by app")
	}
	if len(s.ListConfigItemsByEnvID("e1")) != 0 {
		t.Fatalf("expected empty config items by env")
	}
	if len(s.ListConfigVersions()) != 0 {
		t.Fatalf("expected empty config versions")
	}
	if len(s.ListConfigVersionsByConfigItemID("c1")) != 0 {
		t.Fatalf("expected empty config versions by item")
	}
	if len(s.ListReleases()) != 0 {
		t.Fatalf("expected empty releases")
	}
	if len(s.ListReleasesByAppID("a1")) != 0 {
		t.Fatalf("expected empty releases by app")
	}
	if len(s.ListReleasesByEnvID("e1")) != 0 {
		t.Fatalf("expected empty releases by env")
	}
	if len(s.ListAuditLogs()) != 0 {
		t.Fatalf("expected empty audit logs")
	}
	if len(s.ListAuditLogsByAction("create")) != 0 {
		t.Fatalf("expected empty audit logs by action")
	}
	if len(s.ListAuditLogsByResource("app")) != 0 {
		t.Fatalf("expected empty audit logs by resource")
	}
}

func TestApplicationStoreGetAfterDelete(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Application{ID: "a1", Name: "app1", Code: "code1", Status: model.AppStatusActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateApplication(a)
	s.DeleteApplication("a1")
	if _, err := s.GetApplication("a1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
	if _, err := s.GetApplicationByCode("code1"); err != ErrNotFound {
		t.Fatalf("expected not found by code after delete")
	}
}

func TestEnvironmentStoreGetAfterDelete(t *testing.T) {
	s := NewMemoryStore()
	e := &model.Environment{ID: "e1", AppID: "a1", Name: "dev", Code: "dev", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateEnvironment(e)
	s.DeleteEnvironment("e1")
	if _, err := s.GetEnvironment("e1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}

func TestConfigItemStoreGetAfterDelete(t *testing.T) {
	s := NewMemoryStore()
	c := &model.ConfigItem{ID: "c1", AppID: "a1", EnvID: "e1", Key: "k1", Value: "v1", ValueType: model.ValueTypeString, Status: model.ConfigStatusEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateConfigItem(c)
	s.DeleteConfigItem("c1")
	if _, err := s.GetConfigItem("c1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
	if _, err := s.GetConfigItemByAppEnvKey("a1", "e1", "k1"); err != ErrNotFound {
		t.Fatalf("expected not found by app env key after delete")
	}
}

func TestReleaseStoreGetAfterDelete(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Release{ID: "r1", AppID: "a1", EnvID: "e1", Version: "v1", Status: model.ReleaseStatusDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	s.CreateRelease(r)
	s.DeleteRelease("r1")
	if _, err := s.GetRelease("r1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}

func TestAuditLogStoreGetAfterDelete(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "al1", Operator: "u1", Action: "create", Resource: "app", Detail: "d1", CreatedAt: time.Now()}
	s.CreateAuditLog(a)
	if _, err := s.GetAuditLog("al1"); err != nil {
		t.Fatalf("get audit: %v", err)
	}
}

func TestConfigVersionOrdering(t *testing.T) {
	s := NewMemoryStore()
	for i := 0; i < 10; i++ {
		v := &model.ConfigVersion{ID: "v" + string(rune('0'+i)), ConfigItemID: "c1", Value: "v", ChangedBy: "u1", Remark: "r", CreatedAt: time.Now()}
		s.CreateConfigVersion(v)
		if v.Version != i+1 {
			t.Fatalf("expected version %d, got %d", i+1, v.Version)
		}
	}
}
