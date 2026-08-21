package service_test

import (
	"configcenter/internal/batch/audit"
	"configcenter/internal/batch/repository"
	"configcenter/internal/batch/resource"
	"configcenter/internal/batch/service"
	"configcenter/internal/batch/tx"
	"testing"
)

func TestBatchReleaseCompletesBeforeSuccess(t *testing.T) {
	transaction := &tx.Transaction{}
	log := &audit.Log{}
	p := service.Publisher{Repo: &repository.Repository{Pool: resource.New(2), Tx: transaction}, Audit: log}
	err := p.Publish([]string{"dev", "staging", "prod"})
	if err != nil {
		t.Fatalf("batch returned error: %v", err)
	}
	if p.Repo.Processed != 3 {
		t.Errorf("success response covered only %d of 3 environments", p.Repo.Processed)
	}
	if !transaction.Committed() || log.Status != "success" || p.Repo.Processed != 3 {
		t.Errorf("audit=%q committed=%v processed=%d", log.Status, transaction.Committed(), p.Repo.Processed)
	}
}
