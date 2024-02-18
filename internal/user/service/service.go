package service

import (
	"errors"
	"time"

	"github.com/synapsis-test/internal/user/store/postgresql"
)

// Following constans are config default values.
const (
	defaultTokenExpiration = 10000 * time.Hour
)

// Followings are the known error returned from service.
var (
	errMissingMandatoryConfig = errors.New("missing mandatory config")
)

// service implements user.Service.
type service struct {
	pgStore postgresql.PGStore
	config  Config
	timeNow func() time.Time
}

// Config denotes service configuration
//
// Adding a new field should also add the corresponding default
// value in getDefaultConfig().
type Config struct {
	PasswordSalt    string
	TokenExpiration time.Duration
	TokenSecretKey  string
}

// getDefaultConfig returns service configuration with the
// predefined default values.
func getDefaultConfig() Config {
	return Config{
		TokenExpiration: defaultTokenExpiration,
	}
}

// New creates a new service.
func New(pgStore postgresql.PGStore, options ...Option) (*service, error) {
	s := &service{
		pgStore: pgStore,
		config:  getDefaultConfig(),
		timeNow: time.Now,
	}

	// apply options
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	// verify mandatory config
	if s.config.PasswordSalt == "" || s.config.TokenSecretKey == "" {
		return nil, errMissingMandatoryConfig
	}

	return s, nil
}

// Option controls the behavior of service.
type Option func(*service) error

// WithConfig returns Option to set service configuration.
func WithConfig(config Config) Option {
	return func(s *service) error {
		if config.PasswordSalt != "" {
			s.config.PasswordSalt = config.PasswordSalt
		}
		if config.TokenExpiration > 0 {
			s.config.TokenExpiration = config.TokenExpiration
		}
		if config.TokenSecretKey != "" {
			s.config.TokenSecretKey = config.TokenSecretKey
		}
		return nil
	}
}
