package repository

import (
	"configcenter/internal/releaseflow/audit"
	"configcenter/internal/releaseflow/publisher"
	"configcenter/internal/releaseflow/transaction"
)

type Repository struct {
	Transactions *transaction.Manager
	Publisher    *publisher.Publisher
	Audit        *audit.Log
}

func (r *Repository) Apply() error {
	tx := r.Transactions.Begin()
	r.Audit.Add("success")
	if err := r.Publisher.Publish(); err != nil {
		tx.Rollback()
		r.Audit.Add("failed")
		return err
	}
	tx.Commit()
	return nil
}
