package postgresql

import (
	"context"

	"github.com/synapsis-test/internal/order"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	// CreateOrder creates a new order and return the
	// created order ID.
	CreateOrder(ctx context.Context, order order.Order) (int64, error)

	// GetOrderByID returns a order with the given order ID.
	GetOrderByID(ctx context.Context, id int64) (order.Order, error)
}
