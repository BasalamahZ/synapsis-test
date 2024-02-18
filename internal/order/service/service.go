package service

import (
	"time"

	"github.com/synapsis-test/internal/order/store/postgresql"
	"github.com/synapsis-test/internal/product"
)

// service implements user.Service.
type service struct {
	pgStore postgresql.PGStore
	product product.Service
	timeNow func() time.Time
}

// New creates a new service.
func New(pgStore postgresql.PGStore, product product.Service) (*service, error) {
	s := &service{
		pgStore: pgStore,
		product: product,
		timeNow: time.Now,
	}

	return s, nil
}
