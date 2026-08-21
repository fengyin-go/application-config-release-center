package service

import (
	"configcenter/internal/releaseflow/publisher"
	"configcenter/internal/releaseflow/repository"
	"errors"
)

type Releaser struct{ Repository *repository.Repository }

func (r *Releaser) Release() error {
	err := r.Repository.Apply()
	var temporary publisher.TemporaryError
	if errors.As(err, &temporary) {
		return r.Repository.Apply()
	}
	return err
}
