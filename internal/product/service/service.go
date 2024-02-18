package service

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/synapsis-test/internal/product/store/postgresql"
)

// service implements user.Service.
type service struct {
	pgStore     postgresql.PGStore
	redisClient *redis.Client
	timeNow     func() time.Time
}

// New creates a new service.
func New(pgStore postgresql.PGStore, redisClient *redis.Client) (*service, error) {
	s := &service{
		pgStore:     pgStore,
		redisClient: redisClient,
		timeNow:     time.Now,
	}

	return s, nil
}
