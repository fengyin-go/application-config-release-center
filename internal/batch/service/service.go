package service

import (
	"configcenter/internal/batch/audit"
	"configcenter/internal/batch/repository"
)

type Publisher struct {
	Repo  *repository.Repository
	Audit *audit.Log
}

func (p *Publisher) Publish(items []string) error {
	p.Audit.Success()
	err := p.Repo.Store(items)
	if err != nil {
		p.Audit.Failure()
	}
	return err
}
