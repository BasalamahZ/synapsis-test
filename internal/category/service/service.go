package service

import (
	"time"

	"github.com/synapsis-test/internal/category/store/postgresql"
)

// service implements user.Service.
type service struct {
	pgStore postgresql.PGStore
	timeNow func() time.Time
}

// New creates a new service.
func New(pgStore postgresql.PGStore) (*service, error) {
	s := &service{
		pgStore: pgStore,
		timeNow: time.Now,
	}

	return s, nil
}
