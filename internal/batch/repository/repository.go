package repository

import (
	"configcenter/internal/batch/resource"
	"configcenter/internal/batch/tx"
)

type Repository struct {
	Pool      *resource.Pool
	Tx        *tx.Transaction
	Processed int
}

func (r *Repository) Store(items []string) (err error) {
	defer func() { err = r.Tx.Commit() }()
	for range items {
		lease, acquireErr := r.Pool.Acquire()
		if acquireErr != nil {
			return acquireErr
		}
		defer lease.Close()
		r.Processed++
	}
	return nil
}
