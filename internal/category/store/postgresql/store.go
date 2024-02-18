package postgresql

import (
	"context"

	"github.com/synapsis-test/internal/category"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// GetCategoryByID returns a category with the given category ID.
	GetCategoryByID(ctx context.Context, id int64) (category.Category, error)

	// GetCategories returns list of categories.
	GetCategories(ctx context.Context) ([]category.Category, error)
}
