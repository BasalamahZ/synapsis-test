package postgresql

import (
	"context"

	"github.com/synapsis-test/internal/product"
)

type PGStore interface {
	NewClient(useTx bool) (PGStoreClient, error)
}

type PGStoreClient interface {
	// Commit commits the transaction.
	Commit() error
	// Rollback aborts the transaction.
	Rollback() error

	// GetProductByID returns a product with the given product ID.
	GetProductByID(ctx context.Context, id int64) (product.Product, error)

	// GetProducts returns list of products that satisfy the given
	// filter.
	GetProducts(ctx context.Context, filter product.GetProductsFilter) ([]product.Product, error)

	// AddProductCart add a prodcut to cart
	AddProductCart(ctx context.Context, productCart product.ProductCart) error

	// GetCartByUserID returns all product in cart with the given user ID.
	GetCartsByUserID(ctx context.Context, userID int64) ([]product.ProductCart, error)

	// DeleteProductCartByID deletes a product in cart with the given cart ID.
	DeleteProductCartByID(ctx context.Context, productCart product.ProductCart) error
}
