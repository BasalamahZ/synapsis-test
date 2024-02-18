package postgresql

import (
	"context"

	"github.com/synapsis-test/internal/user"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// CreateUser insert the given user
	//
	// CreateUser return created user ID
	CreateUser(ctx context.Context, users user.User) (int64, error)

	// GetUserByID selects user with the given user ID.
	GetUserByID(ctx context.Context, userID int64) (user.User, error)

	// GetUserByEmail selects a user with the given
	// email.
	GetUserByEmail(ctx context.Context, email string) (user.User, error)

	// UpdateUser updates existing data with the given data
	// for a user specified with the given user ID.
	//
	// UpdateUser do updates on all the fields except ID,
	// Username, Type, CreateBy, and CreateTime. So, make
	// sure to use current values in the given data if do
	// not want to update some specific fields.
	UpdateUser(ctx context.Context, user user.User) error
}
