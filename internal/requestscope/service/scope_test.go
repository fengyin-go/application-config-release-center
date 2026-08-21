package service_test

import (
	"configcenter/internal/requestscope/audit"
	"configcenter/internal/requestscope/middleware"
	"configcenter/internal/requestscope/pool"
	"configcenter/internal/requestscope/service"
	"reflect"
	"testing"
)

func TestReusedRequestScopeStaysIsolated(t *testing.T) {
	p := &pool.Pool{}
	log := &audit.Logger{}
	binder := middleware.Binder{Pool: p}
	life := service.Lifecycle{Pool: p, Logger: log}
	a := binder.Begin("tenant-a", "secret-a")
	life.End(a)
	b := binder.Begin("tenant-b", "clean-b")
	if !reflect.DeepEqual(b.Labels, []string{"clean-b"}) {
		t.Errorf("tenant-b inherited labels %v", b.Labels)
	}
	entry := log.Flush()
	if entry.Tenant != "tenant-a" || !reflect.DeepEqual(entry.Labels, []string{"secret-a"}) {
		t.Errorf("tenant-a delayed audit became tenant=%q labels=%v", entry.Tenant, entry.Labels)
	}
}
