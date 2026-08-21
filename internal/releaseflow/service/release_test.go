package service_test

import (
	"configcenter/internal/releaseflow/audit"
	"configcenter/internal/releaseflow/publisher"
	"configcenter/internal/releaseflow/repository"
	"configcenter/internal/releaseflow/service"
	"configcenter/internal/releaseflow/transaction"
	"reflect"
	"testing"
)

func TestRetryPublishesOneSuccessEvent(t *testing.T) {
	m := &transaction.Manager{}
	p := &publisher.Publisher{}
	a := &audit.Log{}
	r := service.Releaser{Repository: &repository.Repository{Transactions: m, Publisher: p, Audit: a}}
	if err := r.Release(); err != nil {
		t.Fatalf("release failed: %v", err)
	}
	if m.Commits != 1 || p.Events != 1 {
		t.Errorf("commits=%d success events=%d", m.Commits, p.Events)
	}
	if !reflect.DeepEqual(a.Entries, []string{"success"}) {
		t.Errorf("audit entries=%v", a.Entries)
	}
}
